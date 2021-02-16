package oauth

import (
	"context"
	"time"

	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/jackc/pgx/v4"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)

var oauthServer *server.Server
var clientStore *pg.ClientStore

func SetupOAuthServer() error {
	pgxConn, connectErr := pgx.Connect(context.TODO(), db.GetDatabaseURL())
	if connectErr != nil {
		return connectErr
	}
	manager := manage.NewDefaultManager()

	// use PostgreSQL token store with pgx.Connection adapter
	adapter := pgx4adapter.NewConn(pgxConn)
	tokenStore, tokenStoreErr := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	if tokenStoreErr != nil {
		return tokenStoreErr
	}
	defer tokenStore.Close()

	store, clientStoreErr := pg.NewClientStore(adapter)
	if clientStoreErr != nil {
		return clientStoreErr
	}

	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(store)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	oauthServer = srv
	clientStore = store

	return nil
}

func GetOAuthServer() *server.Server {
	return oauthServer
}

func GetClientStore() *pg.ClientStore {
	return clientStore
}
