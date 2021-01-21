package api_v1

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gutakk/go-google-scraper/oauth"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/google/uuid"
)

type OAuthController struct{}

func (oa *OAuthController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.POST("/client", oa.generateClient)
}

func (oa *OAuthController) generateClient(c *gin.Context) {
	clientId := uuid.New().String()
	clientSecret := uuid.New().String()

	err := oauth.GetClientStore().Create(&models.Client{
		ID:     clientId,
		Secret: clientSecret,
		Domain: fmt.Sprintf("http://localhost:%s", os.Getenv("APP_PORT")),
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusCreated, gin.H{
		"CLIENT_ID":     clientId,
		"CLIENT_SECRET": clientSecret,
	})
}
