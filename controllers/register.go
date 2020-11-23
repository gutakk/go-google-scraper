package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"
)

type RegisterController struct{}

type RegisterForm struct {
	Email           string `form:"email" binding:"email,required"`
	Password        string `form:"password" binding:"required,min=6"`
	ConfirmPassword string `form:"confirm-password" binding:"eqfield=Password,required"`
}

func (r *RegisterController) applyRoutes(engine *gin.Engine) {
	engine.GET("/register", r.displayRegister)
	engine.POST("/register", r.register)
}

func (r *RegisterController) displayRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register", gin.H{
		"title": "Register",
	})
}

func (r *RegisterController) register(c *gin.Context) {
	form := &RegisterForm{}

	if err := c.ShouldBind(form); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"title": "Register",
				"error": errorHandler.ValidationErrorToText(fieldErr),
				"email": form.Email,
			})
			return
		}
	}

	hashedPassword, hashError := models.HashPassword(form.Password)
	if hashError != nil {
		c.HTML(http.StatusUnprocessableEntity, "register.html", gin.H{
			"title": "Register",
			"error": "Something went wrong, please try again.",
		})
	}

	if err := models.SaveUser(form.Email, hashedPassword); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": errorHandler.DatabaseErrorToText(err),
			"email": form.Email,
		})
		return
	}

	session.AddFlash(c, "Register successfully")
	c.Redirect(http.StatusFound, "/")
}
