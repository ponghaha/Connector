package service

import (
	"fmt"
	"strings"
	"time"
	"strconv"

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
type applicationCapTCPSocketClient = client.TCPSocketClient

// applicationCapService implements the business logic for customer-related features
type applicationCapService struct {
	config       *config.Config
	logger       *zap.SugaredLogger
	tcpClient    client.TCPSocketClient
	routes       map[string]config.Route
	destinations map[string]config.Destination
}

// NewApplicationCapService creates a new instance of applicationCapService.
func NewApplicationCapService(
	cfg *config.Config,
	logger *zap.SugaredLogger,
	tcpClient applicationCapTCPSocketClient,
	routes map[string]config.Route,
	destinations map[string]config.Destination,
) *applicationCapService {
	return &applicationCapService{
		config:       cfg,
		logger:       logger,
		tcpClient:    tcpClient,
		routes:       routes,
		destinations: destinations,
	}
}

func (s *applicationCapService) GetApplicationNo(c *gin.Context, getApplicationNoReq domain.GetApplicationNoRequest) domain.GetApplicationNoResult {
	timestamp := time.Now()
	//const routeKey = "POST:/Api/Application/GetApplicationNo"
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "GetApplicationNo"
	var domainErr *appError.AppError
	var logLine1 string
	var formatReq interface{}
	var formatResp interface{}
	var formatErr interface{}

	reqID, _ := c.Get("Api-RequestID")
	apiRequestID, ok := reqID.(string)
	if !ok {
		return domain.GetApplicationNoResult{
			Response:    nil,
			AppError:    appError.ErrApiRequestID,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getApplicationNoReq.IDCardNo,
			LogLine1:    "",
		}
	}

	countCardListRq := len(getApplicationNoReq.CardListRq)
	if getApplicationNoReq.TotalApplyCard != countCardListRq {
		return domain.GetApplicationNoResult{
			Response:    nil,
			AppError:    appError.ErrInvTotalOfList,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getApplicationNoReq.IDCardNo,
			LogLine1:    "",
		}
	}

	for _, card := range getApplicationNoReq.CardListRq {
		if card.CardCode == "" {
			return domain.GetApplicationNoResult{
				Response:    nil,
				AppError:    appError.ErrRequiedParam,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getApplicationNoReq.IDCardNo,
				LogLine1:    "",
			}
		}
		if card.VirtualCardFlag == "" {
			return domain.GetApplicationNoResult{
				Response:    nil,
				AppError:    appError.ErrRequiedParam,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getApplicationNoReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}
	
	switch getApplicationNoReq.Channel {
	case "L", "F", "A", "W", "R", "O", "E":
	default:
		return domain.GetApplicationNoResult{
			Response:    nil,
			AppError:    appError.ErrInvChannel,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getApplicationNoReq.IDCardNo,
			LogLine1:    "",
		}
	}

	route, ok := s.routes[routeKey]
	if !ok {
		s.logger.Errorw("Route configuration not found for TCP service", "routeKey", routeKey)
		return domain.GetApplicationNoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getApplicationNoReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.GetApplicationNoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getApplicationNoReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.GetApplicationNoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getApplicationNoReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["GetApplicationNo"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetApplicationNoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getApplicationNoReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetApplicationNoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getApplicationNoReq.IDCardNo,
			LogLine1:    "",
		}
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatGetApplicationNoRequest(getApplicationNoReq)

	var requestLength string
	lenData := len(fixedLengthData)
	strLenData := strconv.Itoa(lenData)
	if len(fixedLengthData) >= 10000 {
		requestLength = strLenData
	} else if len(fixedLengthData) >= 10000 {
		requestLength = "0" + strLenData
	} else if len(fixedLengthData) >= 10000 {
		requestLength = "00" + strLenData
	} else if len(fixedLengthData) >= 10000 {
		requestLength = "000" + strLenData
	} else {
		requestLength = "0000" + strLenData
	}

	header := utils.BuildFixedLengthHeader(
		route.System,
		route.Service,
		route.Format,
		formattedRequestID,
		requestLength,
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "GetApplicationNoFormatRq", "", getApplicationNoReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetApplicationNoResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getApplicationNoReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.GetApplicationNoResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     getApplicationNoReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     getApplicationNoReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])
	if errorCode != "" {
		switch errorCode {
		case "SVC105":
			domainErr = appError.ErrRequiedParam
		case "SVC157":
			domainErr = appError.ErrInvAppChannel
		case "SVC158":
			domainErr = appError.ErrInvTotalOfList
		case "SVC102":
			domainErr = appError.ErrInvCardCode
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "GetApplicationNo", "", getApplicationNoReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetApplicationNoResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getApplicationNoReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.GetApplicationNoResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     getApplicationNoReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	getApplicationNoResponse, err := format.FormatGetApplicationNoResponse(responseStr)

	if err != nil {
		s.logger.Errorw("Error map getApplicationNoResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "GetApplicationNo", "", getApplicationNoReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetApplicationNoResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getApplicationNoReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "GetApplicationNo", "", getApplicationNoReq.IDCardNo)
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.GetApplicationNoResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getApplicationNoReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.GetApplicationNoResult{
		Response:    &getApplicationNoResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     getApplicationNoReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserRef:     getApplicationNoReq.IDCardNo,
		LogLine1:    logLine1,
	}
}

func (s *applicationCapService) SubmitCardApplication(c *gin.Context, submitCardApplicationReq domain.SubmitCardApplicationRequest) domain.SubmitCardApplicationResult {
	timestamp := time.Now()
	//const routeKey = "POST:/Api/Application/SubmitCardApplication"
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "SubmitCardApplication"
	var domainErr *appError.AppError
	var logLine1 string
	var formatReq interface{}
	var formatResp interface{}
	var formatErr interface{}

	reqID, _ := c.Get("Api-RequestID")
	apiRequestID, ok := reqID.(string)
	if !ok {
		return domain.SubmitCardApplicationResult{
			Response:    nil,
			AppError:    appError.ErrApiRequestID,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitCardApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}

	countCardListRq := len(submitCardApplicationReq.SubmitCardListRq)
	if submitCardApplicationReq.TotalApplyCard != countCardListRq {
		return domain.SubmitCardApplicationResult{
			Response:    nil,
			AppError:    appError.ErrInvTotalOfList,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitCardApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}

	for _, card := range submitCardApplicationReq.SubmitCardListRq {
		if card.CardCode == "" {
			return domain.SubmitCardApplicationResult{
				Response:    nil,
				AppError:    appError.ErrRequiedParam,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     submitCardApplicationReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}
	
	switch submitCardApplicationReq.Channel {
	case "L", "F", "A", "W", "R", "O", "E":
	default:
		return domain.SubmitCardApplicationResult{
			Response:    nil,
			AppError:    appError.ErrInvChannel,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitCardApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}

	route, ok := s.routes[routeKey]
	if !ok {
		s.logger.Errorw("Route configuration not found for TCP service", "routeKey", routeKey)
		return domain.SubmitCardApplicationResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitCardApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.SubmitCardApplicationResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitCardApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.SubmitCardApplicationResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitCardApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["SubmitCardApplication"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.SubmitCardApplicationResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitCardApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.SubmitCardApplicationResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitCardApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatSubmitCardApplicationRequest(submitCardApplicationReq)

	var requestLength string
	lenData := len(fixedLengthData)
	strLenData := strconv.Itoa(lenData)
	if len(fixedLengthData) >= 10000 {
		requestLength = strLenData
	} else if len(fixedLengthData) >= 10000 {
		requestLength = "0" + strLenData
	} else if len(fixedLengthData) >= 10000 {
		requestLength = "00" + strLenData
	} else if len(fixedLengthData) >= 10000 {
		requestLength = "000" + strLenData
	} else {
		requestLength = "0000" + strLenData
	}

	header := utils.BuildFixedLengthHeader(
		route.System,
		route.Service,
		route.Format,
		formattedRequestID,
		requestLength,
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "SubmitCardApplication", "", submitCardApplicationReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.SubmitCardApplicationResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     submitCardApplicationReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.SubmitCardApplicationResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     submitCardApplicationReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     submitCardApplicationReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])
	if errorCode != "" {
		switch errorCode {
		case "SVC159", "SVC163", "SVC164":
			domainErr = appError.ErrRequiedParam
		case "SVC117":
			domainErr = appError.ErrInvIDCardNo
		case "SVC157":
			domainErr = appError.ErrInvAppChannel
		case "SVC165":
			domainErr = appError.ErrInvAppDate
		case "SVC127":
			domainErr = appError.ErrInvBranchCode
		case "SVC161", "SVC166", "SVC167", "SVC178", "SVC179", "SVC180":
			domainErr = appError.ErrInvSourceCode
		case "SVC160":
			domainErr = appError.ErrInvMailTo
		case "SVC158":
			domainErr = appError.ErrInvTotalOfList
		case "SVC102", "SVC162":
			domainErr = appError.ErrInvCardCode
		case "SVC168":
			domainErr = appError.ErrInvAppNo
		case "SVC170":
			domainErr = appError.ErrNotfoundConsent
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "SubmitCardApplication", "", submitCardApplicationReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.SubmitCardApplicationResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     submitCardApplicationReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.SubmitCardApplicationResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     submitCardApplicationReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	submitCardApplicationResponse, err := format.FormatSubmitCardApplicationResponse(responseStr)

	if err != nil {
		s.logger.Errorw("Error map submitCardApplicationResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "SubmitCardApplication", "", submitCardApplicationReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.SubmitCardApplicationResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     submitCardApplicationReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "SubmitCardApplication", "", submitCardApplicationReq.IDCardNo)
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.SubmitCardApplicationResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitCardApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.SubmitCardApplicationResult{
		Response:    &submitCardApplicationResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     submitCardApplicationReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserRef:     submitCardApplicationReq.IDCardNo,
		LogLine1:    logLine1,
	}
}
