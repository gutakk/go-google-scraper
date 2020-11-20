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

type UserSessionController struct {
	DB *gorm.DB
}

type LoginForm struct {
	Email    string `form:"email" binding:"email,required"`
	Password string `form:"password" binding:"required,min=6"`
}

func (us *UserSessionController) applyRoutes(engine *gin.Engine) {
	engine.GET("/login", us.displayLogin)
	engine.POST("/login", us.login)
}

func (us *UserSessionController) displayLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

func (us *UserSessionController) login(c *gin.Context) {
	form := &LoginForm{}

	if err := c.ShouldBind(form); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"title": "Login",
				"error": errorHandler.ValidationErrorToText(fieldErr),
			})
			return
		}
	}

	user := &models.User{Email: form.Email}
	if result := us.DB.Where(user).First(user); result.Error != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title": "Login",
			"error": "Username or password is invalid",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title": "Login",
			"error": "Username or password is invalid",
		})
		return
	}

	session.Set(c, "user_id", user.ID)
	c.Redirect(http.StatusFound, "/")
}
