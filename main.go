package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/frank-yf/go-web-example/controller"
	"github.com/frank-yf/go-web-example/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	serverPort            = ":8000"
	serverReadTimeout     = 5 * time.Second
	serverWriteTimeout    = 35 * time.Second // 为了满足 pprof 使用，特意调大写入超时
	serverShutdownTimeout = 30 * time.Second
	serverMaxHeaderBytes  = 1 << 20

	ProductionMode  = "prod"
	DevelopmentMode = "dev"

	AppVersion = "1.0.2"
)

var (
	// 应用启动角色
	appMode string
	// 日志目录
	logHome string
	// 输出版本号
	outputVersion bool
)

func init() {
	flag.StringVar(&appMode, "appMode", DevelopmentMode, "application running mode, must in [prod,dev]")
	flag.StringVar(&logHome, "logHome", "logs", "application log home")
	flag.BoolVar(&outputVersion, "v", false, "print application version")
	flag.Parse()

	if outputVersion {
		fmt.Println("abtest version", AppVersion)
		os.Exit(0)
	}

	logOpts := &utils.LoggerOptions{LogHome: logHome}
	var logLevel zapcore.Level
	var outToFile bool
	switch appMode {
	case ProductionMode:
		gin.SetMode(gin.ReleaseMode)
		logOpts.LogLevel = zapcore.InfoLevel
		logOpts.OutToFile = true
		break
	case DevelopmentMode:
		gin.SetMode(gin.DebugMode)
		logOpts.LogLevel = zapcore.DebugLevel
		break
	default:
		log.Panicln("appMode must in [prod|dev]")
	}

	utils.InitLog(logOpts)

	// 设置 gin 框架的日志写入
	gin.DefaultWriter = utils.GetLogger().GetWriter()
	gin.DefaultErrorWriter = utils.GetLogger().GetErrorWriter()

	utils.GetLogger().Info("application start ~~~")
	utils.GetLogger().Debug("application config",
		zap.String("appMode", appMode),
		zap.String("logHome", logHome),
		zap.String("serverPort", serverPort),
		zap.Bool("outToFile", outToFile),
		zap.Int8("logLevel", int8(logLevel)),
	)
}

func main() {
	r := controller.InitRouter()
	readyClient()
	listenAndServe(r)
}

func readyClient() {
	utils.GetRedisCli()
	utils.GetCacheCli()
}

func listenAndServe(router *gin.Engine) {
	srv := &http.Server{
		Addr:           serverPort,
		Handler:        router,
		ReadTimeout:    serverReadTimeout,
		WriteTimeout:   serverWriteTimeout,
		MaxHeaderBytes: serverMaxHeaderBytes,
	}
	utils.GetLogger().Debug("Listening server")
	utils.GetLogger().Debug("web server config",
		zap.Duration("ReadTimeout", serverReadTimeout),
		zap.Duration("WriteTimeout", serverWriteTimeout),
		zap.Int("MaxHeaderBytes", serverMaxHeaderBytes),
	)

	// 在goroutine中初始化服务器，以便它不会阻止下面的正常关闭处理
	go func() {
		utils.GetLogger().S.Infof("Serving on port %s", serverPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			utils.GetLogger().S.Errorf("listen: %v", err)
		}
	}()

	// 等待中断信号正常关闭服务器
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.GetLogger().Info("Shutting down server...")

	// 上下文用于通知服务器它有一定的时间来完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	utils.CloseRedisCli()

	if err := srv.Shutdown(ctx); err != nil {
		// 被迫关闭
		utils.GetLogger().S.Errorf("Server forced to shutdown: %v", err)
	}

	utils.GetLogger().Info("Server exiting")

	// 清空磁盘缓冲，关闭日志写入对象
	if err := utils.GetLogger().SyncAndClose(); err != nil {
		log.Fatalln("cannot sync and close log file writer : ", err)
	}
}
