package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/jackc/pgx/v4"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)

var oauthServer *server.Server
var clientStore *pg.ClientStore

func SetupOAuthServer() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	pgxConn, _ := pgx.Connect(context.TODO(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName))
	manager := manage.NewDefaultManager()

	// use PostgreSQL token store with pgx.Connection adapter
	adapter := pgx4adapter.NewConn(pgxConn)
	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	defer tokenStore.Close()

	store, _ := pg.NewClientStore(adapter)

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
}

func GetOAuthServer() *server.Server {
	return oauthServer
}

func GetClientStore() *pg.ClientStore {
	return clientStore
}
