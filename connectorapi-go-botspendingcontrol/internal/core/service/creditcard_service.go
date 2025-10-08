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
func (s *creditCardService) GetSpendingControl(c *gin.Context, getSpendingControlReq domain.GetSpendingControlRequest) domain.GetSpendingControlResult {
	timestamp := time.Now()
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "GetSpendingControl"
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
		return domain.GetSpendingControlResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   getSpendingControlReq.AEONID,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.GetSpendingControlResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   getSpendingControlReq.AEONID,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.GetSpendingControlResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   getSpendingControlReq.AEONID,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["GetSpendingControl"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetSpendingControlResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   getSpendingControlReq.AEONID,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetSpendingControlResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   getSpendingControlReq.AEONID,
			LogLine1:    "",
		}
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatGetSpendingControlRequest(getSpendingControlReq)

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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "GetSpendingControl", getSpendingControlReq.AEONID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetSpendingControlResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   getSpendingControlReq.AEONID,
				LogLine1:    "",
			}
		}

		return domain.GetSpendingControlResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     getSpendingControlReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   getSpendingControlReq.AEONID,
			LogLine1:    logLine1,
		}
	}

	errorCode := strings.TrimSpace(responseStr[246:248])				//check ด้วยว่า response มาที่ตำแหน่งไหนบ้าง
	errorMessage := strings.TrimSpace(responseStr[248:298])
	if errorCode != "" {
		switch errorCode {
		case "01":
			domainErr = appError.ErrInvDate
		case "02":
			domainErr = appError.ErrInvTime
		case "03":
			domainErr = appError.ErrUserRefOrAeonID
		case "04":
			domainErr = appError.ErrInvCardNo
		case "05":
			domainErr = appError.ErrInvBusCode
		case "06":
			domainErr = appError.ErrInvServChannel
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


		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "GetSpendingControl", getSpendingControlReq.AEONID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetSpendingControlResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   getSpendingControlReq.AeonID,
				LogLine1:    "",
			}
		}

		return domain.GetSpendingControlResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   getSpendingControlReq.AeonID,
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
	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "GetSpendingControl", getSpendingControlReq.AEONID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetSpendingControlResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   getSpendingControlReq.AeonID,
				LogLine1:    "",
			}
		}

		return domain.GetSpendingControlResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   getSpendingControlReq.AeonID,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	getSpendingControlResponse, err := format.FormatGetSpendingControlResponse(responseStr)
	if err != nil {
		s.logger.Errorw("Error map getSpendingControlResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "GetSpendingControl", getSpendingControlReq.AEONID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetSpendingControlResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   getSpendingControlReq.AEONID,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "GetSpendingControl", getSpendingControlReq.AEONID, "")
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.GetSpendingControlResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   getSpendingControlReq.AEONID,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.GetSpendingControlResult{
		Response:    &getSpendingControlResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     getSpendingControlReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserToken:   getSpendingControlReq.AEONID,
		LogLine1:    logLine1,
	}
}

// It sends a request to the TCP service and returns the response.
func (s *creditCardService) UpdateSpendingControl(c *gin.Context, updateSpendingControlReq domain.UpdateSpendingControlRequest) domain.UpdateSpendingControlResult {
	timestamp := time.Now()
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "UpdateSpendingControl"
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
		return domain.UpdateSpendingControlResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   updateSpendingControlReq.AeonID,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.UpdateSpendingControlResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   updateSpendingControlReq.AeonID,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.UpdateSpendingControlResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   updateSpendingControlReq.AeonID,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["UpdateSpendingControl"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.UpdateSpendingControlResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   updateSpendingControlReq.AeonID,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.UpdateSpendingControlResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   updateSpendingControlReq.AeonID,
			LogLine1:    "",
		}
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatUpdateSpendingControlRequest(updateSpendingControlReq)

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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "UpdateSpendingControl", updateSpendingControlReq.AeonID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.UpdateSpendingControlResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   updateSpendingControlReq.AeonID,
				LogLine1:    "",
			}
		}

		return domain.UpdateSpendingControlResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     updateSpendingControlReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   updateSpendingControlReq.AeonID,
			LogLine1:    logLine1,
		}
	}

	errorCode := strings.TrimSpace(responseStr[246:248])				//check ด้วยว่า response มาที่ตำแหน่งไหนบ้าง
	errorMessage := strings.TrimSpace(responseStr[248:298])
	if errorCode != "" {
		switch errorCode {
		case "01":
			domainErr = appError.ErrInvDate
		case "02":
			domainErr = appError.ErrInvTime
		case "03":
			domainErr = appError.ErrAeonID
		case "04":
			domainErr = appError.ErrInvCardNo
		case "05":
			domainErr = appError.ErrInvBusCode
		case "06":
			domainErr = appError.ErrInvServChannel
		case "07":
			domainErr = appError.ErrOverCRLimit
		case "08":
			domainErr = appError.ErrInvCardNotStatus
		case "09":
			domainErr = appError.ErrInvLimitStatus
		case "10":
			domainErr = appError.ErrBigCardNotFound
		case "11":
			domainErr = appError.ErrInvLimitAmount
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


		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "UpdateSpendingControl", updateSpendingControlReq.AeonID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.UpdateSpendingControlResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   updateSpendingControlReq.AeonID,
				LogLine1:    "",
			}
		}

		return domain.UpdateSpendingControlResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   updateSpendingControlReq.AeonID,
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
	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "UpdateSpendingControl", updateSpendingControlReq.AeonID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.UpdateSpendingControlResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   updateSpendingControlReq.AeonID,
				LogLine1:    "",
			}
		}

		return domain.UpdateSpendingControlResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   updateSpendingControlReq.AeonID,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	updateSpendingControlResponse, err := format.FormatUpdateSpendingControlResponse(responseStr)
	if err != nil {
		s.logger.Errorw("Error map updateSpendingControlResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "UpdateSpendingControl", updateSpendingControlReq.AeonID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.UpdateSpendingControlResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   updateSpendingControlReq.AeonID,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "UpdateSpendingControl", updateSpendingControlReq.AeonID, "")
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.UpdateSpendingControlResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   updateSpendingControlReq.AeonID,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.UpdateSpendingControlResult{
		Response:    &updateSpendingControlResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     updateSpendingControlReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserToken:   updateSpendingControlReq.AeonID,
		LogLine1:    logLine1,
	}
}