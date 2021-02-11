package controllers_test

import (
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	err := os.Chdir(path_test.GetRoot())
	if err != nil {
		log.Fatal(errorconf.ChangeToRootDirFailure, err)
	}

	config.LoadEnv()

	os.Exit(m.Run())
}
