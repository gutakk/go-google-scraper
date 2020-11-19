package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterController struct {
	DB *gorm.DB
}

type UserCredentials struct {
	Email           string `form:"email" binding:"email,required"`
	Password        string `form:"password" binding:"required,min=6"`
	ConfirmPassword string `form:"confirm-password" binding:"eqfield=Password,required"`
}

func (r *RegisterController) applyRoutes(engine *gin.Engine) {
	engine.GET("/register", r.displayRegister)
	engine.POST("/register", r.register)
}

func (r *RegisterController) displayRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}

func (r *RegisterController) register(c *gin.Context) {
	credentials := &UserCredentials{}

	if err := c.ShouldBind(credentials); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"title": "Register",
				"error": errorHandler.ValidationErrorToText(fieldErr),
			})
			return
		}
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)

	if result := r.DB.Create(&models.User{Email: credentials.Email, Password: string(encryptedPassword)}); result.Error != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": errorHandler.DatabaseErrorToText(result.Error),
		})
		return
	}

	session.Set(c, "status", "Register successfully")
	c.Redirect(http.StatusFound, "/")
}
