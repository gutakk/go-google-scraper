package controllers

import (
	"errors"
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"
)

const (
	loginTitle = "Login"
	loginView  = "login"

	invalidUsernameOrPassword = "Username or password is invalid"
)

type UserSessionController struct{}

type LoginForm struct {
	Email    string `form:"email" binding:"email,required"`
	Password string `form:"password" binding:"required,min=6"`
}

func (us *UserSessionController) applyRoutes(engine *gin.RouterGroup) {
	engine.GET("/login", us.displayLogin)
	engine.POST("/login", us.login)
}

func (us *UserSessionController) displayLogin(c *gin.Context) {
	ginview.HTML(c, http.StatusOK, loginView, gin.H{
		"title":   loginTitle,
		"notices": session.Flashes(c),
	})
}

func (us *UserSessionController) login(c *gin.Context) {
	form := &LoginForm{}

	if err := c.ShouldBind(form); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			renderLoginWithError(c, http.StatusBadRequest, errorHandler.ValidationErrorMessage(fieldErr), form)
			return
		}
	}

	user, err := models.FindOneUserBy(&models.User{Email: form.Email})
	if err != nil {
		renderLoginWithError(c, http.StatusUnauthorized, errors.New(invalidUsernameOrPassword), form)
		return
	}

	if err := models.ValidatePassword(user.Password, form.Password); err != nil {
		renderLoginWithError(c, http.StatusUnauthorized, errors.New(invalidUsernameOrPassword), form)
		return
	}

	session.Set(c, "user_id", user.ID)
	c.Redirect(http.StatusFound, "/")
}

func renderLoginWithError(c *gin.Context, status int, err error, form *LoginForm) {
	ginview.HTML(c, status, loginView, gin.H{
		"title":  loginTitle,
		"errors": err.Error(),
		"email":  form.Email,
	})
}
