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
type creditCardTCPSocketClient = client.TCPSocketClient

// creditCardService implements the business logic for customer-related features
type creditCardService struct {
	config       *config.Config
	logger       *zap.SugaredLogger
	tcpClient    client.TCPSocketClient
	routes       map[string]config.Route
	destinations map[string]config.Destination
}

// NewCreditCardService creates a new instance of creditCardService.
func NewCreditCardService(
	cfg *config.Config,
	logger *zap.SugaredLogger,
	tcpClient creditCardTCPSocketClient,
	routes map[string]config.Route,
	destinations map[string]config.Destination,
) *creditCardService {
	return &creditCardService{
		config:       cfg,
		logger:       logger,
		tcpClient:    tcpClient,
		routes:       routes,
		destinations: destinations,
	}
}

// It sends a request to the TCP service and returns the response.
func (s *creditCardService) GetCardSales(c *gin.Context, getCardSalesReq domain.GetCardSalesRequest) domain.GetCardSalesResult {
	timestamp := time.Now()
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "GetCardSales"
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
		return domain.GetCardSalesResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardSalesReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.GetCardSalesResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardSalesReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.GetCardSalesResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardSalesReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["GetCardSales"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetCardSalesResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardSalesReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetCardSalesResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardSalesReq.IDCardNo,
			LogLine1:    "",
		}
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatGetCardSalesRequest(getCardSalesReq)

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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "GetCardSales", "", getCardSalesReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetCardSalesResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getCardSalesReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.GetCardSalesResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     getCardSalesReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     getCardSalesReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])
	if errorCode != "" {
		switch errorCode {
		case "SVC105":
			domainErr = appError.ErrRequiedParam
		case "SVC117":
			domainErr = appError.ErrInvIDCardNo
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "GetCardSales", "", getCardSalesReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetCardSalesResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getCardSalesReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.GetCardSalesResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     getCardSalesReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	getCardSalesResponse, err := format.FormatGetCardSalesResponse(responseStr)
	if err != nil {
		s.logger.Errorw("Error map getCardSalesResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "GetCardSales", "", getCardSalesReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetCardSalesResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getCardSalesReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "GetCardSales", "", getCardSalesReq.IDCardNo)
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.GetCardSalesResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardSalesReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.GetCardSalesResult{
		Response:    &getCardSalesResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     getCardSalesReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserRef:     getCardSalesReq.IDCardNo,
		LogLine1:    logLine1,
	}
}

// It sends a request to the TCP service and returns the response.
func (s *creditCardService) GetBigCardInfo(c *gin.Context, getBigCardInfoReq domain.GetBigCardInfoRequest) domain.GetBigCardInfoResult {
	timestamp := time.Now()
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "GetBigCardInfo"
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
		return domain.GetBigCardInfoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   getBigCardInfoReq.AeonID,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.GetBigCardInfoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   getBigCardInfoReq.AeonID,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.GetBigCardInfoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   getBigCardInfoReq.AeonID,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["GetBigCardInfo"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetBigCardInfoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   getBigCardInfoReq.AeonID,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetBigCardInfoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   getBigCardInfoReq.AeonID,
			LogLine1:    "",
		}
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatGetBigCardInfoRequest(getBigCardInfoReq)

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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "GetBigCardInfo", getBigCardInfoReq.AeonID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetBigCardInfoResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   getBigCardInfoReq.AeonID,
				LogLine1:    "",
			}
		}

		return domain.GetBigCardInfoResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     getBigCardInfoReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   getBigCardInfoReq.AeonID,
			LogLine1:    logLine1,
		}
	}

	errorCode := strings.TrimSpace(responseStr[246:248])
	errorMessage := strings.TrimSpace(responseStr[248:298])
	if errorCode != "" {
		switch errorCode {
		case "01":
			domainErr = appError.ErrUserRefOrAeonID
		case "02":
			domainErr = appError.ErrInvBusCode
		case "04":
			domainErr = appError.ErrInvCreditCard
		case "05":
			domainErr = appError.ErrBigCardNotFound
		case "03":
			domainErr = appError.ErrInvIDCardNo
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


		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "GetBigCardInfo", getBigCardInfoReq.AeonID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetBigCardInfoResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   getBigCardInfoReq.AeonID,
				LogLine1:    "",
			}
		}

		return domain.GetBigCardInfoResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   getBigCardInfoReq.AeonID,
			LogLine1:    logLine1,
		}
	}

	errorCodeHeader := strings.TrimSpace(responseStr[67:73])
	errorMessageHeader := strings.TrimSpace(responseStr[73:123])
	if errorCodeHeader != "" {
		switch errorCodeHeader {
		case "SVC902":
			domainErr = appError.ErrSystemI
		default:
			if errorCodeHeader != "" {
				s.logger.Info("Unknown error code from System I : ", "code", errorCodeHeader, "message", errorMessageHeader)
				domainErr = appError.ErrSystemIUnexpect
			}
		}

		temps := *domainErr
		temps.Code = errorCode
		temps.Message = errorMessageHeader
		domainErr = &temps
	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "GetBigCardInfo", getBigCardInfoReq.AeonID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetBigCardInfoResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   getBigCardInfoReq.AeonID,
				LogLine1:    "",
			}
		}

		return domain.GetBigCardInfoResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   getBigCardInfoReq.AeonID,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	getBigCardInfoResponse, err := format.FormatGetBigCardInfoResponse(responseStr)
	if err != nil {
		s.logger.Errorw("Error map getBigCardInfoResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "GetBigCardInfo", getBigCardInfoReq.AeonID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetBigCardInfoResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   getBigCardInfoReq.AeonID,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "GetBigCardInfo", getBigCardInfoReq.AeonID, "")
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.GetBigCardInfoResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   getBigCardInfoReq.AeonID,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.GetBigCardInfoResult{
		Response:    &getBigCardInfoResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     getBigCardInfoReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserToken:   getBigCardInfoReq.AeonID,
		LogLine1:    logLine1,
	}
}

// It sends a request to the TCP service and returns the response.
func (s *creditCardService) GetCardDelinquent(c *gin.Context, getCardDelinquentReq domain.GetCardDelinquentRequest) domain.GetCardDelinquentResult {
	timestamp := time.Now()
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "GetCardDelinquent"
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
		return domain.GetCardDelinquentResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardDelinquentReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.GetCardDelinquentResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardDelinquentReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.GetCardDelinquentResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardDelinquentReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["GetCardDelinquent"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetCardDelinquentResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardDelinquentReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetCardDelinquentResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardDelinquentReq.IDCardNo,
			LogLine1:    "",
		}
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatGetCardDelinquentRequest(getCardDelinquentReq)

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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "GetCardDelinquent", "", getCardDelinquentReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetCardDelinquentResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getCardDelinquentReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.GetCardDelinquentResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     getCardDelinquentReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     getCardDelinquentReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])

	if errorCode != "" {
		domainErr := &appError.AppError{
			ErrorCode:    errorCode,
			ErrorMessage: errorMessage,
		}

		logLine1 := elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "GetCardDelinquent", "", getCardDelinquentReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetCardDelinquentResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getCardDelinquentReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.GetCardDelinquentResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     getCardDelinquentReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}


	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	getCardSalesResponse, err := format.FormatGetCardDelinquentResponse(responseStr)
	if err != nil {
		s.logger.Errorw("Error map getCardSalesResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "GetCardDelinquent", "", getCardDelinquentReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetCardDelinquentResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getCardDelinquentReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "GetCardDelinquent", "", getCardDelinquentReq.IDCardNo)
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.GetCardDelinquentResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardDelinquentReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.GetCardDelinquentResult{
		Response:    &getCardSalesResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     getCardDelinquentReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserRef:     getCardDelinquentReq.IDCardNo,
		LogLine1:    logLine1,
	}
}

// It sends a request to the TCP service and returns the response.
func (s *creditCardService) GetFullpan(c *gin.Context, getFullpanReq domain.GetFullpanRequest) domain.GetFullpanResult {
	timestamp := time.Now()
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "GetFullpan"
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
		return domain.GetFullpanResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getFullpanReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.GetFullpanResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getFullpanReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.GetFullpanResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getFullpanReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["GetFullpan"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetFullpanResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getFullpanReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetFullpanResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getFullpanReq.IDCardNo,
			LogLine1:    "",
		}
	}

	hasValid := false
	for _, item := range getFullpanReq.CardList {
		if item.CardCode != "" && item.CardCode != "" {
			hasValid = true
			break
		}

	}

	if !hasValid {
		return domain.GetFullpanResult{
			Response:    nil,
			AppError:    appError.ErrRequiedParam,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getFullpanReq.IDCardNo,
			LogLine1:    "",
		}
	}


	switch getFullpanReq.Channel {
		case "L", "F", "A", "W", "R", "B", "E", "V":
	
		default:
		s.logger.Errorw("Invalid Channel", "Channel", getFullpanReq.Channel)
		return domain.GetFullpanResult{
		Response:    nil,
		AppError:    appError.ErrApiChannel,
		GinCtx:      nil,
		Timestamp:   timestamp,
		ReqBody:     nil,
		RespBody:    nil,
		DomainError: nil,
		ServiceName: serviceName,
		UserRef:     getFullpanReq.IDCardNo,
		LogLine1:    "",
		}
	}

	if len(getFullpanReq.CardList) != getFullpanReq.TotalCard {
	s.logger.Errorw("Mismatch in number of cards", 
		"Expected", getFullpanReq.TotalCard, 
		"Actual", len(getFullpanReq.CardList),
	)
		return domain.GetFullpanResult{
			Response:    nil,
			AppError:    appError.ErrInvTotalOfList,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getFullpanReq.IDCardNo,
			LogLine1:    "",
			}
		}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatGetFullpanRequest(getFullpanReq)

	requestLength := fmt.Sprintf("%05d", len(fixedLengthData))

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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "GetFullpan", getFullpanReq.IDCardNo, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetFullpanResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getFullpanReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.GetFullpanResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     getFullpanReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     getFullpanReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}


	errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])
	if errorCode != "" {
		switch errorCode {
		case "SVC117":
			domainErr = appError.ErrInvIDCardNo
		case "SVC921":
			domainErr = appError.ErrInvTotalOfList
		case "SVC118":
			domainErr = appError.ErrInvCardNo
		case "SVC136":
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "GetFullpan", getFullpanReq.IDCardNo, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetFullpanResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getFullpanReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.GetFullpanResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     getFullpanReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	getFullpanResponse, err := format.FormatGetFullpanResponse(responseStr)
	if err != nil {
		s.logger.Errorw("Error map GetFullpanResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "GetFullpan", getFullpanReq.IDCardNo, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetFullpanResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getFullpanReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "GetFullpan", getFullpanReq.IDCardNo, "")
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.GetFullpanResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getFullpanReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.GetFullpanResult{
		Response:    &getFullpanResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     getFullpanReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserRef:     getFullpanReq.IDCardNo,
		LogLine1:    logLine1,
	}
}

// It sends a request to the TCP service and returns the response.
func (s *creditCardService) GetCardEnroll(c *gin.Context, getCardEnrollReq domain.GetCardEnrollRequest) domain.GetCardEnrollResult {
	timestamp := time.Now()
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "GetCardEnroll"
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
		return domain.GetCardEnrollResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardEnrollReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.GetCardEnrollResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardEnrollReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.GetCardEnrollResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardEnrollReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["GetCardEnroll"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetCardEnrollResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardEnrollReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetCardEnrollResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardEnrollReq.IDCardNo,
			LogLine1:    "",
		}
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatGetCardEnrollRequest(getCardEnrollReq)

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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "GetCardEnroll", "", getCardEnrollReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetCardEnrollResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getCardEnrollReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.GetCardEnrollResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     getCardEnrollReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     getCardEnrollReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])

	if errorCode != "" {
		domainErr := &appError.AppError{
			ErrorCode:    errorCode,
			ErrorMessage: errorMessage,
		}

		logLine1 := elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "GetCardEnroll", "", getCardEnrollReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetCardEnrollResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getCardEnrollReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.GetCardEnrollResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     getCardEnrollReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}


	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	getCardEnrollResponse, err := format.FormatGetCardEnrollResponse(responseStr)
	if err != nil {
		s.logger.Errorw("Error map getCardEnrollResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "GetCardEnroll", "", getCardEnrollReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetCardEnrollResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     getCardEnrollReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "GetCardEnroll", "", getCardEnrollReq.IDCardNo)
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.GetCardEnrollResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     getCardEnrollReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.GetCardEnrollResult{
		Response:    &getCardEnrollResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     getCardEnrollReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserRef:     getCardEnrollReq.IDCardNo,
		LogLine1:    logLine1,
	}
}
