package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestAuth(t *testing.T) {
	router := gin.New()
	router.Use(Authorization)
	router.GET("/testing/authorization", func(c *gin.Context) {
		c.String(http.StatusOK, c.MustGet(gin.AuthUserKey).(string))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/testing/authorization", nil)
	req.Header.Set("Authorization", "Basic eXVlZmVpNzc0NjoxMjMxMjM=")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "yuefei7746", w.Body.String())
}

func TestAuthorizationHeader(t *testing.T) {
	assert.Equal(t, "Basic eXVlZmVpNzc0NjoxMjMxMjM=", authorizationHeader("yuefei7746", "123123"))
}
