package main

import (
	"go.uber.org/zap"
	"net/http"
)

var logger *zap.Logger
var sugaredLogger *zap.SugaredLogger

func initLogger() {
	logger, _ = zap.NewDevelopment()
	sugaredLogger = logger.Sugar()
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
func main() {
	initLogger()
	defer logger.Sync()
	simpleHttpGet("https://www.baidu.com")
}
