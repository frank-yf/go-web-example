package controller

import (
	"net/http"
	"net/http/pprof"

	"github.com/frank-yf/go-web-example/utils"
	"github.com/gin-gonic/gin"
)

func PprofIndex(c *gin.Context) {
	pprof.Index(c.Writer, c.Request)
}

func PprofCmdline(c *gin.Context) {
	pprof.Cmdline(c.Writer, c.Request)
}

func PprofProfile(c *gin.Context) {
	pprof.Profile(c.Writer, c.Request)
}

func PprofSymbol(c *gin.Context) {
	pprof.Symbol(c.Writer, c.Request)
}

func PprofTrace(c *gin.Context) {
	pprof.Trace(c.Writer, c.Request)
}

var pprofHandlers = utils.NewSet("heap", "goroutine", "allocs", "block", "threadcreate", "mutex")

func PprofHandler(c *gin.Context) {
	handlerName := c.Param("handler")
	if !pprofHandlers.Has(handlerName) {
		c.JSON(http.StatusNotFound, ResponseEntity{
			Code: http.StatusNotFound,
			Msg:  "unknown pprof router",
		})
		return
	}
	pprof.Handler(handlerName).ServeHTTP(c.Writer, c.Request)
}
