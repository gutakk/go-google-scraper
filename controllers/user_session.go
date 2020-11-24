package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	render "github.com/gutakk/go-google-scraper/helpers/render"
	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	loginTitle = "Login"
	loginView  = "login.html"

	invalidUsernameOrPassword = "Username or password is invalid"
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
	render.HtmlWithNotice(c, loginTitle, loginView, http.StatusOK, session.GetAndDelete(c, "notice"))
}

func (us *UserSessionController) login(c *gin.Context) {
	form := &LoginForm{}

	if err := c.ShouldBind(form); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			renderLoginWithError(c, http.StatusBadRequest, errorHandler.ValidationErrorToText(fieldErr))
			return
		}
	}

	user := &models.User{Email: form.Email}
	if result := us.DB.Where(user).First(user); result.Error != nil {
		renderLoginWithError(c, http.StatusUnauthorized, invalidUsernameOrPassword)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		renderLoginWithError(c, http.StatusUnauthorized, invalidUsernameOrPassword)
		return
	}

	session.Set(c, "user_id", user.ID)
	c.Redirect(http.StatusFound, "/")
}

func renderLoginWithError(c *gin.Context, status int, errorMsg string) {
	render.HtmlWithError(c, loginTitle, loginView, status, errorMsg)
}
