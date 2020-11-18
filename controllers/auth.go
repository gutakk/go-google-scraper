package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

type UserCredentials struct {
	Email           string `form:"email" binding:"email,required"`
	Password        string `form:"password" binding:"required,min=6"`
	ConfirmPassword string `form:"confirm-password" binding:"eqfield=Password,required"`
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
	credentials := &UserCredentials{}

	if err := c.ShouldBind(credentials); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": err.Error(),
		})
		return
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)

	if result := a.DB.Create(&models.User{Email: credentials.Email, Password: string(encryptedPassword)}); result.Error != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": result.Error,
		})
		return
	}

	// session.Set(c, "status", "Register successfully")
	c.Redirect(http.StatusFound, "/")
}
