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
type mobileTCPSocketClient = client.TCPSocketClient

// mobileService implements the business logic for customer-related features
type mobileService struct {
	config       *config.Config
	logger       *zap.SugaredLogger
	tcpClient    client.TCPSocketClient
	routes       map[string]config.Route
	destinations map[string]config.Destination
}

// NewMobileService creates a new instance of mobileService.
func NewMobileService(
	cfg *config.Config,
	logger *zap.SugaredLogger,
	tcpClient mobileTCPSocketClient,
	routes map[string]config.Route,
	destinations map[string]config.Destination,
) *mobileService {
	return &mobileService{
		config:       cfg,
		logger:       logger,
		tcpClient:    tcpClient,
		routes:       routes,
		destinations: destinations,
	}
}

func (s *mobileService) DashboardSummary(c *gin.Context, dashboardSummaryReq domain.DashboardSummaryRequest) domain.DashboardSummaryResult {
	timestamp := time.Now()
	//const routeKey = "POST:/Api/Mobile/DashboardSummary"
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "DashboardSummary"
	var domainErr *appError.AppError
	var logLine1 string
	var formatReq interface{}
	var formatResp interface{}
	var formatErr interface{}
	var systemName string
	var formatNumber string
	var flagOldFormatReq bool

	reqID, _ := c.Get("Api-RequestID")
	apiRequestID, ok := reqID.(string)
	if !ok {
		return domain.DashboardSummaryResult{
			Response:    nil,
			AppError:    appError.ErrApiRequestID,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardSummaryReq.AeonID,
			UserRef:     dashboardSummaryReq.IDCardNo,
			LogLine1:    "",
		}
	}
	
	route, ok := s.routes[routeKey]
	if !ok {
		s.logger.Errorw("Route configuration not found for TCP service", "routeKey", routeKey)
		return domain.DashboardSummaryResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardSummaryReq.AeonID,
			UserRef:     dashboardSummaryReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.DashboardSummaryResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardSummaryReq.AeonID,
			UserRef:     dashboardSummaryReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.DashboardSummaryResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardSummaryReq.AeonID,
			UserRef:     dashboardSummaryReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["DashboardSummary"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.DashboardSummaryResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardSummaryReq.AeonID,
			UserRef:     dashboardSummaryReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.DashboardSummaryResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardSummaryReq.AeonID,
			UserRef:     dashboardSummaryReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if dashboardSummaryReq.IDCardNo != "" {
		systemName = route.SystemV1
		formatNumber = route.FormatV1
		flagOldFormatReq = true
	} else {
		systemName = route.SystemV2
		formatNumber = route.FormatV2
		flagOldFormatReq = false
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatDashboardSummaryRequest(flagOldFormatReq, dashboardSummaryReq)

	header := utils.BuildFixedLengthHeader(
		systemName,
		route.Service,
		formatNumber,
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "DashboardSummary", dashboardSummaryReq.AeonID, dashboardSummaryReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.DashboardSummaryResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   dashboardSummaryReq.AeonID,
				UserRef:     dashboardSummaryReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.DashboardSummaryResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     dashboardSummaryReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   dashboardSummaryReq.AeonID,
			UserRef:     dashboardSummaryReq.IDCardNo,
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
		case "SVC269":
			domainErr = appError.ErrAeonID
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "DashboardSummary", dashboardSummaryReq.AeonID, dashboardSummaryReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.DashboardSummaryResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   dashboardSummaryReq.AeonID,
				UserRef:     dashboardSummaryReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.DashboardSummaryResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   dashboardSummaryReq.AeonID,
			UserRef:     dashboardSummaryReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	dashboardSummaryResponse, err := format.FormatDashboardSummaryResponse(responseStr, flagOldFormatReq)
	if err != nil {
		s.logger.Errorw("Error map dashboardSummaryResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "DashboardSummary", dashboardSummaryReq.AeonID, dashboardSummaryReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.DashboardSummaryResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   dashboardSummaryReq.AeonID,
				UserRef:     dashboardSummaryReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "DashboardSummary", dashboardSummaryReq.AeonID, dashboardSummaryReq.IDCardNo)
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.DashboardSummaryResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardSummaryReq.AeonID,
			UserRef:     dashboardSummaryReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.DashboardSummaryResult{
		Response:    &dashboardSummaryResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     dashboardSummaryReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserToken:   dashboardSummaryReq.AeonID,
		UserRef:     dashboardSummaryReq.IDCardNo,
		LogLine1:    logLine1,
	}
}

func (s *mobileService) DashboardDetail(c *gin.Context, dashboardDetailReq domain.DashboardDetailRequest) domain.DashboardDetailResult {
	timestamp := time.Now()
	//const routeKey = "POST:/Api/Mobile/DashboardDetail"
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "DashboardDetail"
	var domainErr *appError.AppError
	var logLine1 string
	var formatReq interface{}
	var formatResp interface{}
	var formatErr interface{}
	var systemName string
	var formatNumber string
	var flagOldFormatReq bool

	reqID, _ := c.Get("Api-RequestID")
	apiRequestID, ok := reqID.(string)
	if !ok {
		return domain.DashboardDetailResult{
			Response:    nil,
			AppError:    appError.ErrApiRequestID,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardDetailReq.AeonID,
			UserRef:     dashboardDetailReq.IDCardNo,
			LogLine1:    "",
		}
	}
	
	route, ok := s.routes[routeKey]
	if !ok {
		s.logger.Errorw("Route configuration not found for TCP service", "routeKey", routeKey)
		return domain.DashboardDetailResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardDetailReq.AeonID,
			UserRef:     dashboardDetailReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.DashboardDetailResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardDetailReq.AeonID,
			UserRef:     dashboardDetailReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.DashboardDetailResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardDetailReq.AeonID,
			UserRef:     dashboardDetailReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["DashboardDetail"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.DashboardDetailResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardDetailReq.AeonID,
			UserRef:     dashboardDetailReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.DashboardDetailResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardDetailReq.AeonID,
			UserRef:     dashboardDetailReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if dashboardDetailReq.IDCardNo != "" {
		systemName = route.SystemV1
		formatNumber = route.FormatV1
		flagOldFormatReq = true
	} else {
		systemName = route.SystemV2
		formatNumber = route.FormatV2
		flagOldFormatReq = false
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatDashboardDetailRequest(flagOldFormatReq, dashboardDetailReq)

	header := utils.BuildFixedLengthHeader(
		systemName,
		route.Service,
		formatNumber,
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "DashboardDetail", dashboardDetailReq.AeonID, dashboardDetailReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.DashboardDetailResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   dashboardDetailReq.AeonID,
				UserRef:     dashboardDetailReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.DashboardDetailResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     dashboardDetailReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   dashboardDetailReq.AeonID,
			UserRef:     dashboardDetailReq.IDCardNo,
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
		case "SVC269":
			domainErr = appError.ErrAeonID
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "DashboardDetail", dashboardDetailReq.AeonID, dashboardDetailReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.DashboardDetailResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   dashboardDetailReq.AeonID,
				UserRef:     dashboardDetailReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.DashboardDetailResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:   dashboardDetailReq.AeonID,
			UserRef:     dashboardDetailReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	dashboardDetailResponse, err := format.FormatDashboardDetailResponse(responseStr, flagOldFormatReq)
	if err != nil {
		s.logger.Errorw("Error map dashboardDetailResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "DashboardDetail", dashboardDetailReq.AeonID, dashboardDetailReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.DashboardDetailResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:   dashboardDetailReq.AeonID,
				UserRef:     dashboardDetailReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "DashboardDetail", dashboardDetailReq.AeonID, dashboardDetailReq.IDCardNo)
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.DashboardDetailResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:   dashboardDetailReq.AeonID,
			UserRef:     dashboardDetailReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.DashboardDetailResult{
		Response:    &dashboardDetailResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     dashboardDetailReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserToken:   dashboardDetailReq.AeonID,
		UserRef:     dashboardDetailReq.IDCardNo,
		LogLine1:    logLine1,
	}
}

func (s *mobileService) MobileFullPan(c *gin.Context, mobileFullPanReq domain.MobileFullPanRequest) domain.MobileFullPanResult {
	timestamp := time.Now()
	//const routeKey = "POST:/Api/Mobile/MobileFullPAN"
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "MobileFullPan"
	var domainErr *appError.AppError
	var logLine1 string
	var formatReq interface{}
	var formatResp interface{}
	var formatErr interface{}

	reqID, _ := c.Get("Api-RequestID")
	apiRequestID, ok := reqID.(string)
	if !ok {
		return domain.MobileFullPanResult{
			Response:    nil,
			AppError:    appError.ErrApiRequestID,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     mobileFullPanReq.IDCardNo,
			LogLine1:    "",
		}
	}
	
	switch mobileFullPanReq.Channel {
	case "L", "F", "A", "W", "R", "B", "V", "E":
	default:
		return domain.MobileFullPanResult{
			Response:    nil,
			AppError:    appError.ErrInvChannel,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     mobileFullPanReq.IDCardNo,
			LogLine1:    "",
		}
	}

	route, ok := s.routes[routeKey]
	if !ok {
		s.logger.Errorw("Route configuration not found for TCP service", "routeKey", routeKey)
		return domain.MobileFullPanResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     mobileFullPanReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.MobileFullPanResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     mobileFullPanReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.MobileFullPanResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     mobileFullPanReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["MobileFullPan"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.MobileFullPanResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     mobileFullPanReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.MobileFullPanResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     mobileFullPanReq.IDCardNo,
			LogLine1:    "",
		}
	}

	mobileFullPanFormatRq := domain.MobileFullPanFormatRequest{
		IDCardNo: mobileFullPanReq.IDCardNo,
		CreditCardNo: mobileFullPanReq.CardListRq[0].CardNo,
		BusinessCode: mobileFullPanReq.CardListRq[0].CardCode,
	}
	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatMobileFullPanRequest(mobileFullPanFormatRq)

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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "MobileFullPan", "", mobileFullPanReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.MobileFullPanResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     mobileFullPanReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.MobileFullPanResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     mobileFullPanReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     mobileFullPanReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])
	if errorCode != "" {
		switch errorCode {
		case "SVC105", "SVC117":
			domainErr = appError.ErrInvIDCardNo
		case "SVC118":
			domainErr = appError.ErrInvCreditCard
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "MobileFullPan", "", mobileFullPanReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.MobileFullPanResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     mobileFullPanReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.MobileFullPanResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     mobileFullPanReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	mobileFullPanResponse, err := format.FormatMobileFullPanResponse(responseStr)
	if err != nil {
		s.logger.Errorw("Error map mobileFullPanResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "MobileFullPan", "", mobileFullPanReq.IDCardNo)
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.MobileFullPanResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     mobileFullPanReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "MobileFullpan", "", mobileFullPanReq.IDCardNo)
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.MobileFullPanResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     mobileFullPanReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.MobileFullPanResult{
		Response:    &mobileFullPanResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     mobileFullPanReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserRef:     mobileFullPanReq.IDCardNo,
		LogLine1:    logLine1,
	}
}
