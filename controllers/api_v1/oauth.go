package api_v1

import (
	"net/http"

	"github.com/gutakk/go-google-scraper/services/oauth_service"

	"github.com/gin-gonic/gin"
)

type OAuthController struct{}

func (oa *OAuthController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.POST("/client", oa.generateClient)
}

func (oa *OAuthController) generateClient(c *gin.Context) {
	clientID, clientSecret, err := oauth_service.GenerateClient()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"CLIENT_ID":     clientID,
			"CLIENT_SECRET": clientSecret,
		})
	}
}
