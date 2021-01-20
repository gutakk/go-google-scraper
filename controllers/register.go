package controllers

import (
	"net/http"

	"github.com/gutakk/go-google-scraper/helpers/error_handler"
	"github.com/gutakk/go-google-scraper/helpers/html"
	"github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	registerTitle = "Register"
	registerView  = "register"

	registerSuccessFlash = "Registration completed successfully"
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
	html.RenderWithFlash(c, http.StatusOK, registerView, registerTitle, nil)
}

func (r *RegisterController) register(c *gin.Context) {
	form := &RegisterForm{}

	bindFormErr := c.ShouldBind(form)
	if bindFormErr != nil {
		for _, fieldErr := range bindFormErr.(validator.ValidationErrors) {
			renderRegisterWithError(c, http.StatusBadRequest, error_handler.ValidationErrorMessage(fieldErr), form)
			return
		}
	}

	saveUserErr := models.SaveUser(form.Email, form.Password)
	if saveUserErr != nil {
		renderRegisterWithError(c, http.StatusBadRequest, saveUserErr, form)
		return
	}

	session.AddFlash(c, registerSuccessFlash, "notice")
	c.Redirect(http.StatusFound, "/login")
}

func renderRegisterWithError(c *gin.Context, status int, err error, form *RegisterForm) {
	data := map[string]interface{}{
		"email": form.Email,
	}

	html.RenderWithError(c, status, registerView, registerTitle, err, data)
}
