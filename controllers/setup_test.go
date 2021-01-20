package controllers

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	err := os.Chdir("..")
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
