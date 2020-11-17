package helpers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetAndDelete(c *gin.Context, key string) interface{} {
	session := sessions.Default(c)
	value := session.Get(key)
	session.Delete(key)
	_ = session.Save()

	return value
}

func Set(c *gin.Context, key string, value interface{}) {
	session := sessions.Default(c)

	session.Set(key, value)
	_ = session.Save()
}
