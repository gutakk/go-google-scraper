package google_search_service

import (
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	testPath "github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	testPath.ChangeToRootDir()

	config.LoadEnv()

	os.Exit(m.Run())
}
