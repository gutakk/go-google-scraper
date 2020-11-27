package helpers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	sessionTimeoutInSec = 60 * 60 * 24 // 24 hours
)

func AddFlash(c *gin.Context, value interface{}) {
	session := sessions.Default(c)
	session.AddFlash(value)
	_ = session.Save()
}

func Flashes(c *gin.Context) []interface{} {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

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
	_ = session.Save()
}

func Delete(c *gin.Context, key string) {
	session := sessions.Default(c)
	session.Delete(key)
	_ = session.Save()
}
