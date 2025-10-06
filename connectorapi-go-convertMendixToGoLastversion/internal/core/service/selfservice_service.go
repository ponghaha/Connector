package service

import (
	"fmt"
	"strings"
	"time"

	"connectorapi-go/internal/adapter/client"
	"connectorapi-go/internal/adapter/utils"
	"connectorapi-go/internal/core/domain"
	"connectorapi-go/pkg/config"
	"connectorapi-go/internal/core/service/format"
	appError "connectorapi-go/pkg/error"
	elkLog "connectorapi-go/internal/adapter/client/elk"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TCPSocketClient defines the interface for a TCP socket client
type selfServiceTCPSocketClient = client.TCPSocketClient

// selfServiceService implements the business logic for customer-related features
type selfServiceService struct {
	config       *config.Config
	logger       *zap.SugaredLogger
	tcpClient    client.TCPSocketClient
	routes       map[string]config.Route
	destinations map[string]config.Destination
}

// NewSelfServiceService creates a new instance of selfServiceService.
func NewSelfServiceService(
	cfg *config.Config,
	logger *zap.SugaredLogger,
	tcpClient selfServiceTCPSocketClient,
	routes map[string]config.Route,
	destinations map[string]config.Destination,
) *selfServiceService {
	return &selfServiceService{
		config:       cfg,
		logger:       logger,
		tcpClient:    tcpClient,
		routes:       routes,
		destinations: destinations,
	}
}

//helper
func ValidateMyCardRequest(req domain.MyCardRequest) *appError.AppError {
	snsNo := strings.TrimSpace(req.SNSNo)
	userRef := strings.TrimSpace(req.UserRef)
	channel := strings.TrimSpace(req.Channel)
	mode := strings.TrimSpace(req.Mode)

	if channel == "L" || channel == "F" {
		if snsNo == "" || channel == "" || mode == "" {
			return appError.ErrRequiedParam
		}
	} else {
		if userRef == "" || channel == "" || mode == "" {
			return appError.ErrRequiedParam
		}
	}

	allowedChannels := map[string]bool{
		"L": true,
		"F": true,
		"A": true,
		"E": true,
		"O": true,
	}
	if !allowedChannels[channel] {
		return appError.ErrApiChannel
	}

	if mode != "Normal" && mode != "All" {
		return appError.ErrInvMode
	}

	return nil
}

// It sends a request to the TCP service and returns the response.
func (s *selfServiceService) MyCard(c *gin.Context, myCardReq domain.MyCardRequest) domain.MyCardResult {
	timestamp := time.Now()
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "MyCard"
	var domainErr *appError.AppError
	var logLine1 string
	var RequestLength, Service  string
	var fixedLengthData string
	var formatReq interface{}
	var formatResp interface{}
	var formatErr interface{}

	reqID, _ := c.Get("Api-RequestID")
	apiRequestID, ok := reqID.(string)
	if !ok {
		apiRequestID = ""
	}

	route, ok := s.routes[routeKey]
	if !ok {
		s.logger.Errorw("Route configuration not found for TCP service", "routeKey", routeKey)
		return domain.MyCardResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   myCardReq.SNSNo,
			UserRef:     myCardReq.UserRef,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.MyCardResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   myCardReq.SNSNo,
			UserRef:     myCardReq.UserRef,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.MyCardResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   myCardReq.SNSNo,
			UserRef:     myCardReq.UserRef,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["MyCard"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.MyCardResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   myCardReq.SNSNo,
			UserRef:     myCardReq.UserRef,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.MyCardResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   myCardReq.SNSNo,
			UserRef:     myCardReq.UserRef,
			LogLine1:    "",
		}
	}

    validateResult := ValidateMyCardRequest(myCardReq)
	s.logger.Info("error in check", validateResult)
    if validateResult != nil {
	    return domain.MyCardResult{
		    Response:    nil,
		    AppError:    validateResult,
		    GinCtx:      nil,
		    Timestamp:   timestamp,
		    ReqBody:     nil,
		    RespBody:    nil,
		    DomainError: nil,
		    ServiceName: serviceName,
			UserToken:   myCardReq.SNSNo,
		    UserRef:     myCardReq.UserRef,
		    LogLine1:    "",
	    }
    }

	switch {
    case myCardReq.Mode == "Normal":
	Service, RequestLength = "INQ_CUST_CALIST", "00038"
	fixedLengthData = format.FormatMyCardRequestNormal(myCardReq)
    default:
	Service, RequestLength = "INQ_CUST_CARDLS", "00022"
	fixedLengthData = format.FormatMyCardRequestAll(myCardReq)
    }

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)

	header := utils.BuildFixedLengthHeader(
		route.System,
		Service,
		route.Format,
		formattedRequestID,
		RequestLength,
	)

	combinedPayloadString := header + fixedLengthData
	s.logger.Info("Sending TCP request payload : ", combinedPayloadString)

	tcpAddress := fmt.Sprintf("%s:%s", destination.IP, port)
	responseStr, err := s.tcpClient.SendAndReceive(tcpAddress, combinedPayloadString)

	cleanRsponseStr := strings.ReplaceAll(responseStr, "\r", "")
	cleanRsponseStr = strings.ReplaceAll(cleanRsponseStr, "\n", "")

	formatReq = map[string]string{"data": combinedPayloadString}
	formatResp = map[string]string{"data": cleanRsponseStr}

	if err != nil {
		s.logger.Errorw("Downstream TCP service call failed", "error", err, "address", tcpAddress)

		formatErr = map[string]string{"data": err.Error()}

        errMsg := err.Error()

        switch {
        case strings.Contains(errMsg, "ER040"), strings.Contains(errMsg, "ER060"):
            temp := *appError.ErrTimeOut
            temp.StatusCode = "504"
            domainErr = &temp

        case strings.Contains(errMsg, "ER099"):
            temp := *appError.ErrInternalServer
            temp.StatusCode = "500"
            domainErr = &temp

        default:
            domainErr = appError.ErrService
        }

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "MyCard", myCardReq.SNSNo, myCardReq.UserRef)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.MyCardResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   myCardReq.SNSNo,
				UserRef:     myCardReq.UserRef,
				LogLine1:    "",
			}
		}

		return domain.MyCardResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     myCardReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   myCardReq.SNSNo,
			UserRef:     myCardReq.UserRef,
			LogLine1:    logLine1,
		}
	}

    errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])

	if errorCode != "" {
		switch errorCode {
		case "SVC117", "SVC105":
			domainErr = appError.ErrInvIDCardNo

		case "SVC102":
			domainErr = appError.ErrUserRefOrAeonID

		case "SVC902":
			domainErr = appError.ErrSystemI

		default:
			if errorCode != "" {
				s.logger.Info("Unknown error code from System I : ", "code", errorCode, "message", errorMessage,)
				domainErr = appError.ErrSystemIUnexpect
			}
		}
		temp := *domainErr
		temp.Code = errorCode
		temp.Message = errorMessage
		domainErr = &temp

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "MyCard", myCardReq.SNSNo, myCardReq.UserRef)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.MyCardResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   myCardReq.SNSNo,
				UserRef:     myCardReq.UserRef,
				LogLine1:    "",
			}
		}

		return domain.MyCardResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   myCardReq.SNSNo,
			UserRef:     myCardReq.UserRef,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))

	var MyCardResponse interface{}
	serviceFromSysi := strings.TrimSpace(responseStr[10:25])

	s.logger.Info("Received downstream TCP response", "response", serviceFromSysi)

	switch serviceFromSysi {
    case "INQ_CUST_CALIST":
         MyCardResponse, err = format.FormatMyCardResponseNormal(responseStr)
    default:
         MyCardResponse, err = format.FormatMyCardResponseAll(responseStr)
    }
	if err != nil {
		s.logger.Errorw("Error map MyCardResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "MyCard", myCardReq.SNSNo, myCardReq.UserRef)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.MyCardResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   myCardReq.SNSNo,
				UserRef:     myCardReq.UserRef,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "MyCard", myCardReq.SNSNo,  myCardReq.UserRef)
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.MyCardResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   myCardReq.SNSNo,
			UserRef:     myCardReq.UserRef,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.MyCardResult{
		Response:    MyCardResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     myCardReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserToken:   myCardReq.SNSNo,
		UserRef:     myCardReq.UserRef,
		LogLine1:    logLine1,
	}
}

// It sends a request to the TCP service and returns the response.
func (s *selfServiceService) GetAvailableLimit(c *gin.Context, getAvailableLimitReq domain.GetAvailableLimitRequest) domain.GetAvailableLimitResult {
	timestamp := time.Now()
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "GetAvailableLimit"
	var domainErr *appError.AppError
	var logLine1 string
	var formatReq interface{}
	var formatResp interface{}
	var formatErr interface{}

	reqID, _ := c.Get("Api-RequestID")
	apiRequestID, ok := reqID.(string)
	if !ok {
		apiRequestID = ""
	}

	route, ok := s.routes[routeKey]
	if !ok {
		s.logger.Errorw("Route configuration not found for TCP service", "routeKey", routeKey)
		return domain.GetAvailableLimitResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   "",
			UserRef:     getAvailableLimitReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.GetAvailableLimitResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   "",
			UserRef:     getAvailableLimitReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.GetAvailableLimitResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   "",
			UserRef:     getAvailableLimitReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["GetAvailableLimit"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetAvailableLimitResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   "",
			UserRef:     getAvailableLimitReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetAvailableLimitResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   "",
			UserRef:     getAvailableLimitReq.IDCardNo,
			LogLine1:    "",
		}
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatGetAvailableLimitRequest(getAvailableLimitReq)

	header := utils.BuildFixedLengthHeader(
		route.System,
		route.Service,
		route.Format,
		formattedRequestID,
		route.RequestLength,
	)

	combinedPayloadString := header + fixedLengthData
	s.logger.Info("Sending TCP request payload : ", combinedPayloadString)

	tcpAddress := fmt.Sprintf("%s:%s", destination.IP, port)
	responseStr, err := s.tcpClient.SendAndReceive(tcpAddress, combinedPayloadString)

	cleanRsponseStr := strings.ReplaceAll(responseStr, "\r", "")
	cleanRsponseStr = strings.ReplaceAll(cleanRsponseStr, "\n", "")

	formatReq = map[string]string{"data": combinedPayloadString}
	formatResp = map[string]string{"data": cleanRsponseStr}

	if err != nil {
		s.logger.Errorw("Downstream TCP service call failed", "error", err, "address", tcpAddress)

		formatErr = map[string]string{"data": err.Error()}

        errMsg := err.Error()

        switch {
        case strings.Contains(errMsg, "ER040"), strings.Contains(errMsg, "ER060"):
            temp := *appError.ErrTimeOut
            temp.StatusCode = "504"
            domainErr = &temp

        case strings.Contains(errMsg, "ER099"):
            temp := *appError.ErrInternalServer
            temp.StatusCode = "500"
            domainErr = &temp

        default:
            domainErr = appError.ErrService
        }

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "GetAvailableLimit", "", getAvailableLimitReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetAvailableLimitResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   "",
				UserRef:     getAvailableLimitReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.GetAvailableLimitResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     getAvailableLimitReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   "",
			UserRef:     getAvailableLimitReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	// errorCode := ""
	errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])
	if errorCode != "" {
		switch errorCode {
		case "SVC105", "SVC117":
			domainErr = appError.ErrInvIDCardNo
		case "SVC118":
			domainErr = appError.ErrAgreement
		case "SVC102":
			domainErr = appError.ErrInvCardCode
		case "SVC902":
			domainErr = appError.ErrSystemI
		default:
			if errorCode != "" {
				s.logger.Info("Unknown error code from System I : ", "code", errorCode, "message", errorMessage)
				domainErr = appError.ErrSystemIUnexpect
			}
		}
		temp := *domainErr
		temp.Code = errorCode
		temp.Message = errorMessage
		domainErr = &temp

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "GetAvailableLimit", "", getAvailableLimitReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetAvailableLimitResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   "",
				UserRef:     getAvailableLimitReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.GetAvailableLimitResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   "",
			UserRef:     getAvailableLimitReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	getAvailableLimitResponse, err := format.FormatGetAvailableLimitResponse(responseStr)
	if err != nil {
		s.logger.Errorw("Error map GetAvailableLimitResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "GetAvailableLimit", "", getAvailableLimitReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetAvailableLimitResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   "",
				UserRef:     getAvailableLimitReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "GetAvailableLimit", "", getAvailableLimitReq.IDCardNo)
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.GetAvailableLimitResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   "",
			UserRef:     getAvailableLimitReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.GetAvailableLimitResult{
		Response:    &getAvailableLimitResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     getAvailableLimitReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserToken:   "",
		UserRef:     getAvailableLimitReq.IDCardNo,
		LogLine1:    logLine1,
	}
}