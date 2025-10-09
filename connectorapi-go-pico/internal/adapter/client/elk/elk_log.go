package elk

import (
	"encoding/json"
	"fmt"
	"os"
	"net"
	// "net/url"
	"time"
	"strconv"
	"strings"

	appError "connectorapi-go/pkg/error"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LogMainData struct {
	TIMESTAMP			string				`json:"TIMESTAMP"`
	LOGLEVEL			string				`json:"LOGLEVEL"`
	RequestID			string				`json:"RequestID"`
	TraceID				string				`json:"TraceID"`
	SourceIP			string				`json:"SourceIP"`
	DestIP				string				`json:"DestIP"`
	SourceHostname		string				`json:"SourceHostname"`
	DestHostname		string				`json:"DestHostname"`
	Method				string				`json:"Method"`
	ServiceName			string				`json:"ServiceName"`
	Uri					string				`json:"Uri"`
	Path				string				`json:"Path"`
	UserAgent			string				`json:"UserAgent"`
	Status				string				`json:"Status"`
	ErrorCode			string				`json:"ErrorCode,omitempty"`
	ErrorMessage		string				`json:"ErrorMessage,omitempty"`
	UsedTime			string				`json:"UsedTime"`
	ServerName			string				`json:"ServerName"`
	SeqNo				string				`json:"SeqNo"`
	Header				string          	`json:"Header"`
	UserToken			string				`json:"UserToken"`
	UserRef				string				`json:"UserRef"`
	RequestDateTime		string				`json:"RequestDateTime"`
	RequestMessage		interface{}			`json:"RequestMessage"`
	ResponseDateTime	string				`json:"ResponseDateTime"`
	ResponseMessage		interface{}			`json:"ResponseMessage"`
}

type LogLineData struct {
	RequestID			string				`json:"RequestID"`
	TraceID				string				`json:"TraceID"`
	SourceIP			string				`json:"SourceIP"`
	DestIP				string				`json:"DestIP"`
	SourceHostname		string				`json:"SourceHostname"`
	DestHostname		string				`json:"DestHostname"`
	Method				string				`json:"Method"`
	ServiceName			string				`json:"ServiceName"`
	Uri					string				`json:"Uri"`
	Path				string				`json:"Path"`
	UserAgent			string				`json:"UserAgent"`
	Status				string				`json:"Status"`
	ErrorCode			string				`json:"ErrorCode,omitempty"`
	ErrorMessage		string				`json:"ErrorMessage,omitempty"`
	UsedTime			string				`json:"UsedTime"`
	ServerName			string				`json:"ServerName"`
	SeqNo				string				`json:"SeqNo"`
	Header				string	        	`json:"Header"`
	UserToken			string				`json:"UserToken"`
	UserRef				string				`json:"UserRef"`
	RequestDateTime		string				`json:"RequestDateTime"`
	RequestMessage		interface{}			`json:"RequestMessage"`
	ResponseDateTime	string				`json:"ResponseDateTime"`
	ResponseMessage		interface{}			`json:"ResponseMessage"`
}

type ResponseMessageFormat struct {
	ErrorCode    string `json:"ErrorCode,omitempty"`
	ErrorMessage string `json:"ErrorMessage,omitempty"`
}

func GenerateELKLogMain(c *gin.Context, timesRequest time.Time, request interface{}, response interface{}, appErr *appError.AppError, serviceName string, userToken string, userRef string) string {
	timestamp := time.Now()
	formattedCurrentTimestamp := timestamp.Format("2006-01-02 15:04:05.000")
	formattedLogTimestamp := timestamp.Format("02/01/2006 15:04:05")
	formattedRequestTimestamp := formattedCurrentTimestamp

	duration := timestamp.Sub(timesRequest)
	usedTime := fmt.Sprintf("%d", duration.Milliseconds())

	path := c.FullPath()

	if request == nil {
		request = ""
	}
	if response == nil {
		response = ""
	}

	var responseFormat interface{} = response
	if appErr != nil {
		responseFormat = ResponseMessageFormat{
			ErrorCode:    appErr.ErrorCode,
			ErrorMessage: appErr.ErrorMessage,
		}
	} else {
		responseFormat = response
	}

	logData := LogMainData{
		TIMESTAMP:        formattedLogTimestamp,
		LOGLEVEL:         "INFO",
		RequestID:        c.GetHeader("Api-RequestID"),
		TraceID:          "",
		SourceIP:         "0.0.0.0",
		DestIP:           GetLocalIP(),
		SourceHostname:   "Unknown host",
		DestHostname:     "ConnectorAPI",
		Method:           c.Request.Method,
		ServiceName:      serviceName,
		Uri:              "https://connectorapi.aeonth.com"+path,
		Path:             path,
		UserAgent:        "",
		Status:           "200",
		ServerName:       "ConnectorAPI",
		SeqNo:            "0",
		Header:           extractHeader(c, "", "main"),
		UserToken:		  userToken,
		UserRef:          userRef,
		RequestDateTime:  formattedRequestTimestamp,
		RequestMessage:   request,
		ResponseDateTime: formattedCurrentTimestamp,
		ResponseMessage:  responseFormat,
		UsedTime:         usedTime,
	}

	if appErr != nil {
		logData.Status = strconv.Itoa(c.Writer.Status())
		logData.ErrorCode = appErr.ErrorCode
		logData.ErrorMessage = appErr.ErrorMessage
	}

	logJSON, err := json.Marshal(logData)
	if err != nil {
		// fmt.Println("-----marshal error------", err)
		return ""
	}

	finalJSON := formattedCurrentTimestamp + " INFO :" + string(logJSON)
	return finalJSON
}

func GenerateELKLogLine(c *gin.Context, timesRequest time.Time, request interface{}, response interface{}, appErr *appError.AppError, apikey string, endpoint string, serviceNameMain string, serviceNameLine string, userToken string, userRef string) string {
	timestamp := time.Now()
	formattedCurrentTimestamp := timestamp.Format("2006-01-02 15:04:05.000")
	formattedRequestTimestamp := formattedCurrentTimestamp

	duration := timestamp.Sub(timesRequest)
	usedTime := fmt.Sprintf("%d", duration.Milliseconds())

	val, exists := c.Get("seqNoCounter")
	var counter int
	if exists {
		counter = val.(int)
	} else {
		counter = 1
	}
	seqNo := strconv.Itoa(counter)
	c.Set("seqNoCounter", counter+1)

	// parsedURL, _ := url.Parse(endpoint)
	// path := parsedURL.Path

	// domain := parsedURL.Hostname()
    // if net.ParseIP(domain) != nil {
	// 	domain = "No Host Name"
    // }
	// parts := strings.Split(domain, ".")
	// if len(parts) > 0 {
	// 	domain = parts[0]
	// }
	destIP := strings.Split(endpoint, ":")[0]

	if request == nil {
		request = ""
	}
	if response == nil {
		response = ""
	}

	logData := LogLineData{
		RequestID:        c.GetHeader("Api-RequestID"),
		TraceID:          "",
		SourceIP:         GetLocalIP(),
		DestIP:           destIP,
		SourceHostname:   "ConnectorAPI",
		DestHostname:     "",
		Method:           "TCP",
		ServiceName:      serviceNameMain+"-"+serviceNameLine,
		Uri:              endpoint,
		Path:             endpoint,
		UserAgent:        "",
		Status:           "200",
		ServerName:       "ConnectorAPI",
		SeqNo:            seqNo,
		Header:           extractHeader(c, apikey, "line"),
		UserToken:		  userToken,
		UserRef:          userRef,
		RequestDateTime:  formattedRequestTimestamp,
		RequestMessage:   request,
		ResponseDateTime: formattedCurrentTimestamp,
		ResponseMessage:  response,
		UsedTime:         usedTime,
	}

	if appErr != nil {
		if appErr.StatusCode != "" {
			logData.Status = appErr.StatusCode
		} else {
			logData.Status = "400"
		}
		logData.ErrorCode = appErr.Code
		logData.ErrorMessage = appErr.Message
	}

	logJSON, err := json.Marshal(logData)
	if err != nil {
		// fmt.Println("-----marshal error------", err)
		return ""
	}

	finalJSON := formattedCurrentTimestamp + " INFO :" + string(logJSON)
	return finalJSON
}

// func extractHeader(c *gin.Context, apikey string, logType string) map[string]string {
// 	apiLanguage := c.GetHeader("Api-Language")
// 	if language, ok := c.Get("Api-Language"); ok && fmt.Sprintf("%v", language) != "" {
// 		apiLanguage = fmt.Sprintf("%v", language)
// 	}

// 	if logType == "main" {
// 		return map[string]string{
// 			"APIKey":            c.GetHeader("Api-Key"),
// 			"APIChannel":        c.GetHeader("Api-Channel"),
// 			"APILanguage":       apiLanguage,
// 			"APIAuthorizeToken": c.GetHeader("Api-AuthorizationToken"),
// 			"APIDeviceOS":       c.GetHeader("Api-DeviceOS"),
// 			"APIAnalyzeField":   c.GetHeader("Api-AnalyzeField"),
// 		}
// 	}
// 	if logType == "line" {
// 		return map[string]string{
// 			"APIKey":            apikey,
// 			"APIChannel":        "",
// 			"APILanguage":       "",
// 			"APIAuthorizeToken": "",
// 			"APIDeviceOS":       "",
// 		}
// 	}
// 	return map[string]string{}
// }

func extractHeader(c *gin.Context, apikey string, logType string) string {
	apiLanguage := c.GetHeader("Api-Language")
	if language, ok := c.Get("Api-Language"); ok && fmt.Sprintf("%v", language) != "" {
		apiLanguage = fmt.Sprintf("%v", language)
	}

	header := ""
	if logType == "main" {
		header = "[APIKey:"+c.GetHeader("Api-Key")+
				 "|APIChannel:"+c.GetHeader("Api-Channel")+
				 "|APILanguage:"+apiLanguage+
				 "|APIAuthorizeToken:"+c.GetHeader("Api-AuthorizationToken")+
				 "|APIDeviceOS:"+c.GetHeader("Api-DeviceOS")+"]"
	}
	if logType == "line" {
		header= "[APIKey:null|APIChannel:null|APILanguage:EN|APIAuthorizeToken:null|APIDeviceOS:null]"
	}
	return header
}

func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
		// fmt.Println("-----"Error getting network interfaces------", err)
        return ""
    }

    for _, addr := range addrs {
        var ip net.IP
        switch v := addr.(type) {
        case *net.IPNet:
            ip = v.IP
        case *net.IPAddr:
            ip = v.IP
        }
        if ip == nil || ip.IsLoopback() {
            continue
        }
        if ip.To4() != nil { // IPv4
            return ip.String()
        }
    }
    return ""
}

// func getIPFromURL(parsedURL *url.URL) string {
//     // convert string to URL struct
//     host := parsedURL.Hostname()

//     // check host is ip or domain
//     ip := net.ParseIP(host)
//     if ip != nil {
//         return ip.String()
//     }

//     ips, err := net.LookupIP(host)
//     if err != nil {
//         return ""
//     }

//     // get IP (IPv4 or IPv6)
//     for _, ip := range ips {
//         if ip.To4() != nil {
//             return ip.String()
//         }
//     }

//     if len(ips) > 0 {
//         return ips[0].String()
//     }

//     return ""
// }


func WriteLogToFile(logLines []string, timestamp time.Time, elkPath string) error {
	// if len(logLines) == 0 {
	// 	return nil
	// }
	// // err := os.MkdirAll(elkPath, os.ModePerm)
	// // if err != nil {
	// // 	// fmt.Println("-----"Cannot create log dir------", err)
	// // 	Cannot open file
	// // 	return err
	// // }

	filename := "LOG" + timestamp.Format("20060102") + ".txt"
	fullPath := elkPath + filename

	logFile, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// fmt.Println("-----"Cannot open file------", err)
		return err
	}
	defer logFile.Close()

	for _, line := range logLines {
		if line == "" {
			continue
		}
		if _, err := logFile.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}

type HandleErrorResponse func(c *gin.Context, appErr *appError.AppError)

func FinalELKLog(
	c *gin.Context,
	logList *[]string,
	timestamp time.Time,
	reqBody interface{},
	respBody interface{},
	appErr *appError.AppError,
	serviceName string,
	userToken string,
	userRef string,
	additionalLines []string,
	logger *zap.SugaredLogger,
	elkPath string,
	handleErrorResponse HandleErrorResponse,
) bool {
	logMain := GenerateELKLogMain(c, timestamp, reqBody, respBody, appErr, serviceName, userToken, userRef)
	if logMain == "" {
		logger.Errorw("Error generating log: %v", logMain)
		handleErrorResponse(c, appError.ErrInternalServer)
		return false
	}

	allLogs := []string{logMain}

	if logList != nil && *logList != nil {
		allLogs = append(allLogs, *logList...)
	}

	if len(additionalLines) > 0 {
		allLogs = append(allLogs, additionalLines...)
	}

	if err := WriteLogToFile(allLogs, timestamp, elkPath); err != nil {
		logger.Errorw("Error writing log file:", err)
		handleErrorResponse(c, appError.ErrInternalServer)
		return false
	}

	return true
}
