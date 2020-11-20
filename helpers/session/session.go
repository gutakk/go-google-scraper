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
