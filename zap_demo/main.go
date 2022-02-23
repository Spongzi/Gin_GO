package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var sugaredLogger *zap.SugaredLogger

func initLogger() {
	writeSync := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSync, zapcore.DebugLevel)
	// AddCaller()添加调用信息！
	sugaredLogger = zap.New(core, zap.AddCaller()).Sugar()
}

func getLogWriter() zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    10,
		MaxAge:     30, // 最大备份天数
		MaxBackups: 5,  // 备份数量
		LocalTime:  false,
		Compress:   false, // 是否压缩
	}
	return zapcore.AddSync(lumberjackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		//logger.Error(
		sugaredLogger.Error(
			"Error fetching url...",
			zap.String("url", url),
			zap.Error(err),
		)
		return
	} else {
		//logger.Info(
		sugaredLogger.Info(
			"Success",
			zap.String("statusCode", resp.Status),
			zap.String("url", url),
		)
		resp.Body.Close()
	}

}

func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("code", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func main() {
	r := gin.New()
	r.Use(GinRecovery(logger, true))
	r.Use(GinLogger(logger))
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, SuXu")
	})
	initLogger()
	defer sugaredLogger.Sync()
	simpleHttpGet("https://www.baidu.com")
	panic(r.Run())
}
