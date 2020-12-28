package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/controllers/api"
	"github.com/gutakk/go-google-scraper/oauth"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"
	"github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

func init() {
	gin.SetMode(gin.TestMode)

	if err := os.Chdir(path_test.GetRoot()); err != nil {
		panic(err)
	}

	config.LoadEnv()
	oauth.SetupOAuthServer()
}

func TestGenerateClientWithValidBasicAuth(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(api.OAuthController).ApplyRoutes(controllers.BasicAuthAPIGroup(engine))

	headers := http.Header{}
	// Basic auth with username = admin and password = password
	headers.Set("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=")

	resp := testHttp.PerformRequest(engine, "POST", "/api/client", headers, nil)
	data, _ := ioutil.ReadAll(resp.Body)
	var respBody map[string]string
	_ = json.Unmarshal(data, &respBody)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.NotEqual(t, nil, respBody["CLIENT_ID"])
	assert.NotEqual(t, nil, respBody["CLIENT_SECRET"])
}

func TestGenerateClientWithoutBasicAuth(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(api.OAuthController).ApplyRoutes(controllers.BasicAuthAPIGroup(engine))

	resp := testHttp.PerformRequest(engine, "POST", "/api/client", nil, nil)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}
