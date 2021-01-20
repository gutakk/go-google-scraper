package controllers

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	err := os.Chdir("..")
	if err != nil {
		glog.Fatal(err)
	}

	os.Exit(m.Run())
}
