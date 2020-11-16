package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (a *AuthController) applyRoutes(e *gin.Engine) {
	e.GET("/register", a.displayRegister)
}

func (a *AuthController) displayRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}
