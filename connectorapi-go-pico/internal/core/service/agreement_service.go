package service

import (
	"fmt"
	"strings"
	"time"
	// "strconv"

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
type agreementTCPSocketClient = client.TCPSocketClient

// agreementService implements the business logic for customer-related features
type agreementService struct {
	config       *config.Config
	logger       *zap.SugaredLogger
	tcpClient    client.TCPSocketClient
	routes       map[string]config.Route
	destinations map[string]config.Destination
}

// NewAgreementService creates a new instance of agreementService.
func NewAgreementService(
	cfg *config.Config,
	logger *zap.SugaredLogger,
	tcpClient agreementTCPSocketClient,
	routes map[string]config.Route,
	destinations map[string]config.Destination,
) *agreementService {
	return &agreementService{
		config:       cfg,
		logger:       logger,
		tcpClient:    tcpClient,
		routes:       routes,
		destinations: destinations,
	}
}

// It sends a request to the TCP service and returns the response.
func (s *agreementService) UpdateStatus(c *gin.Context, updateStatusReq domain.UpdateStatusRequest) domain.UpdateStatusResult {
	timestamp := time.Now()
	//const routeKey = "POST:/Api/Agreement/UpdateStatus"
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "UpdateAgreementStatus"
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
		return domain.UpdateStatusResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   updateStatusReq.AeonID,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.UpdateStatusResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   updateStatusReq.AeonID,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.UpdateStatusResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   updateStatusReq.AeonID,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["UpdateAgreementStatus"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.UpdateStatusResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   updateStatusReq.AeonID,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.UpdateStatusResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   updateStatusReq.AeonID,
			LogLine1:    "",
		}
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatUpdateStatusRequest(updateStatusReq)

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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "UpdateStatus", updateStatusReq.AeonID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.UpdateStatusResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   updateStatusReq.AeonID,
				LogLine1:    "",
			}
		}

		return domain.UpdateStatusResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     updateStatusReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   updateStatusReq.AeonID,
			LogLine1:    logLine1,
		}
	}

	errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])
	if errorCode != "" {
		switch errorCode {
		case "01":
			domainErr = appError.ErrAeonID
		case "02":
			domainErr = appError.ErrStatus
		case "03":
			domainErr = appError.ErrAgreement
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "UpdateStatus", updateStatusReq.AeonID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.UpdateStatusResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   updateStatusReq.AeonID,
				LogLine1:    "",
			}
		}

		return domain.UpdateStatusResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   updateStatusReq.AeonID,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	updateStatusResponse, err := format.FormatUpdateStatusResponse(responseStr)
	if err != nil {
		s.logger.Errorw("Error map updateStatusResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "UpdateStatus", updateStatusReq.AeonID, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.UpdateStatusResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   updateStatusReq.AeonID,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "UpdateStatus", updateStatusReq.AeonID, "")
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.UpdateStatusResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   updateStatusReq.AeonID,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.UpdateStatusResult{
		Response:    &updateStatusResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     updateStatusReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserToken:   updateStatusReq.AeonID,
		LogLine1:    logLine1,
	}
}