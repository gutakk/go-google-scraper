package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	html "github.com/gutakk/go-google-scraper/helpers/html"
	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"
)

const (
	loginTitle = "Login"
	loginView  = "login"

	invalidUsernameOrPassword = "Username or password is invalid"
)

type LoginController struct{}

type LoginForm struct {
	Email    string `form:"email" binding:"email,required"`
	Password string `form:"password" binding:"required,min=6"`
}

func (l *LoginController) applyRoutes(engine *gin.RouterGroup) {
	engine.GET("/login", l.displayLogin)
	engine.POST("/login", l.login)
}

func (l *LoginController) displayLogin(c *gin.Context) {
	html.RenderWithFlash(c, http.StatusOK, loginView, loginTitle, nil)
}

func (l *LoginController) login(c *gin.Context) {
	form := &LoginForm{}

	if err := c.ShouldBind(form); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			renderLoginWithError(c, http.StatusBadRequest, errorHandler.ValidationErrorMessage(fieldErr), form)
			return
		}
	}

	user, err := models.FindUserBy(&models.User{Email: form.Email})
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
	data := map[string]interface{}{
		"email": form.Email,
	}

	html.RenderWithError(c, status, loginView, loginTitle, err, data)
}
