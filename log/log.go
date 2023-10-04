package log

/********************************************************************************
* Temancode Example Log Package                                                 *
*                                                                               *
* Version: 1.0.0                                                                *
* Date:    2023-01-05                                                           *
* Author:  Waluyo Ade Prasetio                                                  *
* Github:  https://github.com/abdullahPrasetio                                  *
********************************************************************************/

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/abdullahPrasetio/prasegateway/config"
	"github.com/abdullahPrasetio/prasegateway/constants"
	"github.com/abdullahPrasetio/prasegateway/entity"
	"github.com/abdullahPrasetio/prasegateway/utils"

	log "github.com/sirupsen/logrus"
)

var (
	serviceName string
	folder      = folderDev
	logText     *log.Logger
	logJSON     *log.Logger
	Logger      *log.Logger
)

func init() {
	setText()
	setJSON()
	config := config.DataConfig
	fmt.Println("debug", config.Debug)
	// fmt.Println("config", config)
	if config.Debug {
		logText.SetLevel(log.DebugLevel)
		logJSON.SetLevel(log.DebugLevel)
	} else {
		logText.SetLevel(log.InfoLevel)
		logJSON.SetLevel(log.InfoLevel)
	}

	if config.AppENV != "local" {
		folder = folderOcp
	}

	serviceName = constants.ServerName
	LoadLogger()
}

func LogDebug(refnum, msg string) {
	timestamp := setLogFile()
	// fmt.Println("refnumm", refnum)
	// fmt.Println("msg", msg)
	// fmt.Println("timestamp", timestamp)
	logText.Debug(fmt.Sprintf("%s [%s] %s", timestamp, refnum, msg))
}

func LoadLogger() {
	Logger = log.New()

	switch os.Getenv("LOG_LEVEL") {
	case "trace":
		Logger.SetLevel(log.TraceLevel)
	case "debug":
		Logger.SetLevel(log.DebugLevel)
	case "info":
		Logger.SetLevel(log.InfoLevel)
	case "warn":
		Logger.SetLevel(log.WarnLevel)
	case "error":
		Logger.SetLevel(log.ErrorLevel)
	case "fatal":
		Logger.SetLevel(log.FatalLevel)
	case "panic":
		Logger.SetLevel(log.PanicLevel)
	}

	Logger.SetFormatter(&log.JSONFormatter{})
	Logger.Info("logger Service transaction successfully configured")
}

const (
	httpRequest  = "REQUEST"
	httpResponse = "RESPONSE"
	timeformat   = "2006-01-02T15:04:05-0700"
	nameformat   = "log-2006-01-02.log"
	folderDev    = "log/"
	folderOcp    = "logs/"
)

func LogRequest(httpMethod, trx_type string, request interface{}, header http.Header) {
	timestamp := setLogFile()
	logJSON.WithFields(log.Fields{
		"service":        serviceName,
		"http_type":      httpRequest,
		"http_method":    httpMethod,
		"request_header": header,
		"request_body":   request,
		"trx_type":       trx_type,
		"timestamp":      timestamp,
	}).Info(httpRequest)

}

func setLogFile() string {
	currentTime := time.Now()
	timestamp := currentTime.Format(timeformat)
	filename := folder + currentTime.Format(nameformat)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	} else {
		logText.SetOutput(file)
		logJSON.SetOutput(file)
	}
	return timestamp
}

func setJSON() {
	logJSON = log.New()
	formatter := new(log.JSONFormatter)
	formatter.DisableTimestamp = true
	logJSON.SetFormatter(formatter)
}

func LogResponse(param entity.LogResponseParam) {
	timestamp := setLogFile()
	logJSON.WithFields(log.Fields{
		"third_party":     param.ThirdParty,
		"service":         serviceName,
		"http_type":       httpResponse,
		"response_header": param.ResponseHeader,
		"response_body":   utils.MinifyJson(param.ResponseBody),
		"response_code":   param.ResponseCode,
		"timestamp":       timestamp,
	}).Info(httpResponse)
}

func setText() {
	logText = log.New()
	formatter := new(log.TextFormatter)
	formatter.DisableTimestamp = true
	formatter.DisableQuote = true
	logText.SetFormatter(formatter)
}
