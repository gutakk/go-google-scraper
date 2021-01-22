package oauth_service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gutakk/go-google-scraper/oauth"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/google/uuid"
)

type OAuthClient struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

func GenerateClient() (OAuthClient, error) {
	clientID := uuid.New().String()
	clientSecret := uuid.New().String()

	err := oauth.GetClientStore().Create(&models.Client{
		ID:     clientID,
		Secret: clientSecret,
		Domain: fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("PORT")),
	})
	if err != nil {
		return OAuthClient{}, err
	}

	return OAuthClient{ClientID: clientID, ClientSecret: clientSecret}, nil
}

func GenerateToken(c *gin.Context) {
	err := oauth.GetOAuthServer().HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
	}
}
