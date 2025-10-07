package service

import (
	"fmt"
	"time"
	"strings"

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
type applicationLowerTCPSocketClient = client.TCPSocketClient

// applicationLowerService implements the business logic for customer-related features
type applicationLowerService struct {
	config       *config.Config
	logger       *zap.SugaredLogger
	tcpClient    client.TCPSocketClient
	routes       map[string]config.Route
	destinations map[string]config.Destination
}

// NewApplicationLowerService creates a new instance of applicationLowerService.
func NewApplicationLowerService(
	cfg *config.Config,
	logger *zap.SugaredLogger,
	tcpClient applicationLowerTCPSocketClient,
	routes map[string]config.Route,
	destinations map[string]config.Destination,
) *applicationLowerService {
	return &applicationLowerService{
		config:       cfg,
		logger:       logger,
		tcpClient:    tcpClient,
		routes:       routes,
		destinations: destinations,
	}
}

// It sends a request to the TCP service and returns the response.
func (s *applicationLowerService) SubmitLoanApplication(c *gin.Context, submitLoanApplicationReq domain.SubmitLoanApplicationRequest) domain.SubmitLoanApplicationResult {
	timestamp := time.Now()
	//const routeKey = "POST:/Api/application/SubmitLoanApplication"
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "SubmitLoanApplication"
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
		return domain.SubmitLoanApplicationResult{
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitLoanApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.SubmitLoanApplicationResult{
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitLoanApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.SubmitLoanApplicationResult{
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitLoanApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["SubmitLoanApplication"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.SubmitLoanApplicationResult{
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitLoanApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.SubmitLoanApplicationResult{
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitLoanApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	submitLoanApplicationReq.RequestID = formattedRequestID
	fixedLengthData := format.FormatSubmitLoanApplicationRequest(submitLoanApplicationReq)

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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "SubmitLoanApplication", submitLoanApplicationReq.IDCardNo, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.SubmitLoanApplicationResult{
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     submitLoanApplicationReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.SubmitLoanApplicationResult{
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     submitLoanApplicationReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     submitLoanApplicationReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	// errorCode := strings.TrimSpace(responseStr[67:73])
	// errorMessage := strings.TrimSpace(responseStr[73:123])
	// if errorCode != "" {
	// 	switch errorCode {
	// 	case "SVC105", "SVC163", "SVC267", "SVC272", "SVC274",
	// 		 "SVC275", "SVC276", "SVC278", "SVC279", "SVC280",
	// 		 "SVC281", "SVC291", "SVC292", "SVC295", "SVC296",
	// 		 "SVC300", "SVC301", "SVC302", "SVC303", "SVC305",
	// 		 "SVC307", "SVC312", "SVC313", "SVC314", "SVC316":
	// 		domainErr = appError.ErrRequiedParam
	// 	case "SVC168":
	// 		domainErr = appError.ErrDupAppNo
	// 	case "SVC174":
	// 		domainErr = appError.ErrInvCardAppType
	// 	case "SVC285":
	// 		domainErr = appError.ErrInvGender
	// 	case "SVC315":
	// 		domainErr = appError.ErrInvAgentCode
	// 	case "SVC317":
	// 		domainErr = appError.ErrInvCode
	// 	default:
	// 		if errorCode != "" {
	// 			s.logger.Info("Unknown error code from System I : ", "code", errorCode, "message", errorMessage)
	// 			domainErr = appError.ErrSystemIUnexpect
	// 		}
	// 	}
	// 	temp := *domainErr
	// 	temp.Code = errorCode
	// 	temp.Message = errorMessage
	// 	domainErr = &temp

	// 	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "SubmitLoanApplication", submitLoanApplicationReq.IDCardNo, "")
	// 	if logLine1 == "" {
	// 		s.logger.Errorw("Error generating log: %v", logLine1)
	// 		return domain.SubmitLoanApplicationResult{
	// 			AppError:    appError.ErrInternalServer,
	// 			GinCtx:      nil,
	// 			Timestamp:   timestamp,
	// 			ReqBody:     nil,
	// 			RespBody:    nil,
	// 			DomainError: nil,
	// 			ServiceName: serviceName,
	// 			UserRef:     submitLoanApplicationReq.IDCardNo,
	// 			LogLine1:    "",
	// 		}
	// 	}

	// 	return domain.SubmitLoanApplicationResult{
	// 		AppError:    nil,
	// 		GinCtx:      c,
	// 		Timestamp:   timestamp,
	// 		ReqBody:     formatReq,
	// 		RespBody:    domainErr,
	// 		DomainError: domainErr,
	// 		ServiceName: serviceName,
	// 		UserRef:     submitLoanApplicationReq.IDCardNo,
	// 		LogLine1:    logLine1,
	// 	}
	// }

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	if err != nil {
		s.logger.Errorw("Error map submitLoanApplicationResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "SubmitLoanApplication", submitLoanApplicationReq.IDCardNo, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.SubmitLoanApplicationResult{
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     submitLoanApplicationReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "SubmitLoanApplication", submitLoanApplicationReq.IDCardNo, "")
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.SubmitLoanApplicationResult{
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     submitLoanApplicationReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.SubmitLoanApplicationResult{
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     submitLoanApplicationReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserRef:     submitLoanApplicationReq.IDCardNo,
		LogLine1:    logLine1,
	}
}
