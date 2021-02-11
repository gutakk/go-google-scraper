package helpers

import (
	"github.com/gutakk/go-google-scraper/helpers/log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	sessionTimeoutInSec = 60 * 60 * 24 // 24 hours

	saveSessionFailure = "Failed to save session: "
)

func AddFlash(c *gin.Context, value interface{}, key string) {
	session := sessions.Default(c)
	session.AddFlash(value, key)
	err := session.Save()
	if err != nil {
		log.Error(saveSessionFailure, err)
	}
}

func Flashes(c *gin.Context, key string) []interface{} {
	session := sessions.Default(c)
	flashes := session.Flashes(key)
	err := session.Save()
	if err != nil {
		log.Error(saveSessionFailure, err)
	}

	return flashes
}

func Get(c *gin.Context, key string) interface{} {
	session := sessions.Default(c)
	return session.Get(key)
}

func Set(c *gin.Context, key string, value interface{}) {
	session := sessions.Default(c)
	session.Options(sessions.Options{
		MaxAge: sessionTimeoutInSec,
	})

	session.Set(key, value)
	err := session.Save()
	if err != nil {
		log.Error(saveSessionFailure, err)
	}
}

func Delete(c *gin.Context, key string) {
	session := sessions.Default(c)
	session.Delete(key)
	err := session.Save()
	if err != nil {
		log.Error(saveSessionFailure, err)
	}
}
