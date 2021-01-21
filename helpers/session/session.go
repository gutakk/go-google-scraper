package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	cannotSaveSession   = "Cannot save the session: %s"
	sessionTimeoutInSec = 60 * 60 * 24 // 24 hours
)

func AddFlash(c *gin.Context, value interface{}, key string) {
	session := sessions.Default(c)
	session.AddFlash(value, key)
	err := session.Save()
	if err != nil {
		log.Errorf(cannotSaveSession, err)
	}
}

func Flashes(c *gin.Context, key string) []interface{} {
	session := sessions.Default(c)
	flashes := session.Flashes(key)
	err := session.Save()
	if err != nil {
		log.Errorf(cannotSaveSession, err)
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
		log.Errorf(cannotSaveSession, err)
	}
}

func Delete(c *gin.Context, key string) {
	session := sessions.Default(c)
	session.Delete(key)
	err := session.Save()
	if err != nil {
		log.Errorf(cannotSaveSession, err)
	}
}
