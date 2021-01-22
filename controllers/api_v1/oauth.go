package api_v1

import (
	"net/http"

	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	"github.com/gutakk/go-google-scraper/services/oauth_service"

	"github.com/gin-gonic/gin"
)

type OAuthController struct{}

func (oa *OAuthController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.POST("/client", oa.generateClient)
}

func (oa *OAuthController) generateClient(c *gin.Context) {
	oauthClient, err := oauth_service.GenerateClient()
	if err != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Detail: err.Error(),
			Status: http.StatusUnprocessableEntity,
		}
		c.JSON(errorResponse.Status, errorResponse.NewErrorResponse())
	} else {
		dataResponse := &api_helper.DataResponseObject{
			ID:         oauthClient.ClientID,
			Type:       "client",
			Attributes: oauthClient,
		}
		c.JSON(http.StatusCreated, dataResponse.ConstructDataResponse())
	}
}
