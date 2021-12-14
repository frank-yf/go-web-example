package controller

import (
	"encoding/base64"

	"github.com/frank-yf/go-web-example/utils/json"
	"github.com/gin-gonic/gin"
)

var (
	accounts = gin.Accounts{
		"yuefei7746": "123123",
	}
)

func Authorization(c *gin.Context) {
	gin.BasicAuth(accounts)(c)
}

func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString(json.StringToBytes(base))
}
