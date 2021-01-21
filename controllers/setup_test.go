package controllers

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	err := os.Chdir("..")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
