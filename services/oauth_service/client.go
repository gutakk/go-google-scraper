package oauth_service

import (
	"fmt"
	"os"

	"github.com/gutakk/go-google-scraper/oauth"

	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/google/uuid"
)

func GenerateClient() (string, string, error) {
	clientID := uuid.New().String()
	clientSecret := uuid.New().String()

	err := oauth.GetClientStore().Create(&models.Client{
		ID:     clientID,
		Secret: clientSecret,
		Domain: fmt.Sprintf("http://localhost:%s", os.Getenv("APP_PORT")),
	})
	if err != nil {
		return "", "", err
	}

	return clientID, clientSecret, nil
}
