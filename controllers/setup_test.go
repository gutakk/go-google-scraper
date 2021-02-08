package controllers

import (
	"os"
	"testing"

	errorHelper "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	err := os.Chdir(path_test.GetRoot())
	if err != nil {
		log.Fatal(errorHelper.ChangeToRootDirFailure, err)
	}

	os.Exit(m.Run())
}
