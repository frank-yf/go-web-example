package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/arthurkiller/rollingwriter"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger     *LoggerWrapper
	loggerOnce sync.Once
)

// InitLog 初始化日志结构体 Logger 与 SugarLogger
func InitLog(opts *LoggerOptions) {
	logger = newLoggerWrapper(opts)
}

// InitDefaultLog 初始化一个默认的日志实例
func InitDefaultLog() {
	logger = newLoggerWrapper(&LoggerOptions{
		LogLevel: zap.DebugLevel,
	})
}

// GetLogger 获取日志实例
func GetLogger() *LoggerWrapper {
	loggerOnce.Do(func() {
		if logger == nil {
			// 只有当没有被主动调用初始化时，才执行
			InitDefaultLog()
		}
	})
	return logger
}

type LoggerWrapper struct {
	// 日志结构体，可以输出 结构化日志
	*zap.Logger
	// sugarLogger 日志结构体，可以输出 结构化日志、非结构化日志
	S *zap.SugaredLogger

	infoWriter io.WriteCloser
	errWriter  io.WriteCloser

	opts *LoggerOptions
}

func newLoggerWrapper(opts *LoggerOptions) *LoggerWrapper {
	config := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "Logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
		},
	}

	l := &LoggerWrapper{opts: opts}

	var cores []zapcore.Core
	if l.opts.OutToFile {
		l.infoWriter = getLogFileWriter(l.opts.LogHome, "console")
		l.errWriter = getLogFileWriter(l.opts.LogHome, "error")
		cores = l.newFileWriter(config)
	} else {
		l.infoWriter = os.Stdout
		l.errWriter = os.Stderr
		cores = l.newConsoleWriter(config)
	}

	l.Logger = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(zap.WarnLevel),
	)
	l.S = l.Logger.Sugar()

	return l
}

func (l LoggerWrapper) GetWriter() io.Writer {
	return l.infoWriter
}

func (l LoggerWrapper) GetErrorWriter() io.Writer {
	return l.errWriter
}

// SyncAndClose 刷新 logger 缓冲区，并关闭写入对象
func (l *LoggerWrapper) SyncAndClose() (err error) {
	defer func() {
		// 关于是否应该调用写入对象的 close 方法：https://www.joeshaw.org/dont-defer-close-on-writable-files/

		// 避免关闭 os.Stdout 后程序无法正常输出日志信息，此处只关闭自定义的写入文件对象
		if l.opts.OutToFile {
			err = multierr.Combine(err,
				l.infoWriter.Close(),
				l.errWriter.Close(),
			)
		}
	}()

	// zap 的 Sync 方法会在 stdout/stderr 指向控制台时发生，在社区没有提供可用的解决方案前，忽略该函数的错误返回信息
	// 相关 issues：https://github.com/uber-go/zap/issues/328
	_ = l.Sync()
	return
}

func (l *LoggerWrapper) newFileWriter(config zapcore.EncoderConfig) []zapcore.Core {
	//自定义日志级别：自定义Info级别
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= l.opts.LogLevel
	})

	//自定义日志级别：自定义Warn级别
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= l.opts.LogLevel
	})

	return []zapcore.Core{
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(config),
			zapcore.AddSync(l.infoWriter),
			infoLevel,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(config),
			zapcore.AddSync(l.errWriter),
			warnLevel,
		),
	}
}

func (l *LoggerWrapper) newConsoleWriter(config zapcore.EncoderConfig) []zapcore.Core {
	return []zapcore.Core{
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(config),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
			l.opts.LogLevel,
		),
	}
}

type LoggerOptions struct {
	// OutToFile 是否将日志记录到文件
	OutToFile bool
	// LogHome 日志文件设定的目录，当 OutToFile 为true时生效
	LogHome string
	// LogLevel 应用日志等级，低于该日志等级的将不会被输出
	LogLevel zapcore.Level
}

// getLogFileWriter 提供一个根据文件大小拆分日志文件的写入类
func getLogFileWriter(dir, filename string) io.WriteCloser {
	writer, _ := rollingwriter.NewWriterFromConfig(&rollingwriter.Config{
		LogPath:                dir,
		TimeTagFormat:          "060102",
		FileName:               filenameAppendIP(filename),
		MaxRemain:              1,
		RollingPolicy:          rollingwriter.TimeRolling,
		RollingTimePattern:     "0 0 0 * * *",
		RollingVolumeSize:      "1G",
		WriterMode:             "buffer",
		BufferWriterThershould: 1024,
		Compress:               true,
	})

	return writer
}

// filenameAppendIP 文件名拼接IP
// 部署到k8s发现容器组共用一个log磁盘，为了文件不混淆会根据IP区分容器
func filenameAppendIP(name string) (appended string) {
	ip, err := GetInternetAddress()
	if err != nil {
		ip = "localhost"
	}
	ip = strings.ReplaceAll(ip, ".", "-")
	appended = fmt.Sprint(name, "-", ip)
	return
}
