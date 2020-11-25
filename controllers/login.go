package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginController struct {
	DB *gorm.DB
}

type LoginForm struct {
	Email    string `form:"email" binding:"email,required"`
	Password string `form:"password" binding:"required,min=6"`
}

func (l *LoginController) applyRoutes(engine *gin.Engine) {
	engine.GET("/login", l.displayLogin)
}

func (l *LoginController) displayLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}
