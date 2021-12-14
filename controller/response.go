package controller

import (
	"net/http"

	"github.com/frank-yf/go-web-example/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	MsgOK = "ok"
)

var (
	OK = ResponseEntity{
		Code: http.StatusOK,
		Msg:  MsgOK,
	}
)

type ResponseEntity struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseOK(data interface{}) *ResponseEntity {
	return &ResponseEntity{
		Code: http.StatusOK,
		Msg:  MsgOK,
		Data: data,
	}
}

func ResponseError(errMsg string) *ResponseEntity {
	return &ResponseEntity{
		Code: http.StatusInternalServerError,
		Msg:  errMsg,
	}
}

func renderOK(c *gin.Context) {
	c.JSON(http.StatusOK, OK)
}

func renderData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseOK(data))
}

func renderError(c *gin.Context, errMsg string) {
	c.JSON(http.StatusOK, ResponseError(errMsg))
}

func renderServerError(c *gin.Context, errMsg string) {
	c.JSON(http.StatusInternalServerError, ResponseError(errMsg))
	utils.GetLogger().Warn("response error message", zap.String("msg", errMsg))
}
