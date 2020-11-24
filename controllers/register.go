package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"
)

const (
	registerTitle = "Register"
	registerView  = "register"

	registerSuccessfully = "Register successfully"
)

type RegisterController struct{}

type RegisterForm struct {
	Email           string `form:"email" binding:"email,required"`
	Password        string `form:"password" binding:"required,min=6"`
	ConfirmPassword string `form:"confirm-password" binding:"eqfield=Password,required"`
}

func (r *RegisterController) applyRoutes(engine *gin.RouterGroup) {
	engine.GET("/register", r.displayRegister)
	engine.POST("/register", r.register)
}

func (r *RegisterController) displayRegister(c *gin.Context) {
	c.HTML(http.StatusOK, registerView, gin.H{
		"title": registerTitle,
	})
}

func (r *RegisterController) register(c *gin.Context) {
	form := &RegisterForm{}

	if err := c.ShouldBind(form); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			renderRegisterWithError(c, http.StatusBadRequest, errorHandler.ValidationErrorMessage(fieldErr), form)
			return
		}
	}

	if err := models.SaveUser(form.Email, form.Password); err != nil {
		renderRegisterWithError(c, http.StatusBadRequest, err, form)
		return
	}

	session.AddFlash(c, registerSuccessfully)
	c.Redirect(http.StatusFound, "/login")
}

func renderRegisterWithError(c *gin.Context, status int, error error, form *RegisterForm) {
	c.HTML(status, registerView, gin.H{
		"title":  registerTitle,
		"errors": error,
		"email":  form.Email,
	})
}
