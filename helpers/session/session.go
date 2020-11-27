package helpers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
		MaxAge: 3600 * 24, // 24 hours
	})

	session.Set(key, value)
	_ = session.Save()
}
