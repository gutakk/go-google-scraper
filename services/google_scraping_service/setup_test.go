package google_scraping_service

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/tests/path_test"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	if err := os.Chdir(path_test.GetRoot()); err != nil {
		panic(err)
	}

	config.LoadEnv()

	os.Exit(m.Run())
}
