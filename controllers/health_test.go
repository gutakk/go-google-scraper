package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

func TestHealthStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	health := new(HealthController)

	health.Status(c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Working!", w.Body.String())
}
