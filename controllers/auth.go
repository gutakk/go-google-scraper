package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
)

type AuthController struct{}

type UserCredential struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (a *AuthController) applyRoutes(engine *gin.Engine) {
	engine.GET("/register", a.displayRegister)
	engine.POST("/register", a.register)
}

func (a *AuthController) displayRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}

func (a *AuthController) register(c *gin.Context) {
	credential := &UserCredential{}

	if err := c.ShouldBind(credential); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": err.Error(),
		})
		return
	}

	if result := db.DB.Create(&models.User{Email: credential.Email, Password: credential.Password}); result.Error != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": result.Error,
		})

		return
	}

	c.Redirect(http.StatusFound, "/")
}
