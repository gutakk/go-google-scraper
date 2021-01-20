package google_search_service

import (
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	err := os.Chdir(path_test.GetRoot())
	if err != nil {
		glog.Fatal(err)
	}

	config.LoadEnv()

	os.Exit(m.Run())
}
