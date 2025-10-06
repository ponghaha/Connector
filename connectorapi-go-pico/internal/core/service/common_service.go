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
// type commonTCPSocketClient = client.TCPSocketClient

// commonService implements the business logic for customer-related features
type commonService struct {
	config       *config.Config
	logger       *zap.SugaredLogger
	tcpClient    client.TCPSocketClient
	routes       map[string]config.Route
	destinations map[string]config.Destination
}

// NewCommonService creates a new instance of commonService.
func NewCommonService(
	cfg *config.Config,
	logger *zap.SugaredLogger,
	tcpClient creditCardTCPSocketClient,
	routes map[string]config.Route,
	destinations map[string]config.Destination,
) *commonService {
	return &commonService{
		config:       cfg,
		logger:       logger,
		tcpClient:    tcpClient,
		routes:       routes,
		destinations: destinations,
	}
}

//helper
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func HasValidApplyCardItem(list []domain.ApplyCardListobj) bool {
	for _, item := range list {
		if !(item.CardApplyType == 0 && item.CardCode == "" && item.VirtualCardFlag == "") {
			return true
		}
	}
	return false
}

func HasValidApply2ndCardItem(list []domain.CheckApply2ndCardRqOBJ) bool {
	for _, item := range list {
		if item.CardCode != "" {
			return true
		}
	}
	return false
}

func ValidateCustomerInfo(getCustomerInfoReq domain.GetCustomerInfoRequest, timestamp time.Time) *appError.AppError {
    userRef := strings.TrimSpace(getCustomerInfoReq.UserRef)
    snsNo := strings.TrimSpace(getCustomerInfoReq.SNSNo)
    channel := strings.TrimSpace(getCustomerInfoReq.Channel)
    mode := strings.TrimSpace(getCustomerInfoReq.Mode)
    idCard := strings.TrimSpace(getCustomerInfoReq.IDCardNo)
    aeonID := strings.TrimSpace(getCustomerInfoReq.AEONID)
    agreementNo := strings.TrimSpace(getCustomerInfoReq.AgreementNo)

    if userRef == "" && snsNo == "" && idCard == "" && aeonID == "" && agreementNo == "" {
        return appError.ErrRequiedParam
    }

    if userRef != ""{
        if mode != "S" && mode != "F" && mode != "" {
            return appError.ErrInvMode
        }
    }

    if snsNo != "" {
        if channel == "" {
            return appError.ErrRequiedParam
        }

        if mode != "S" && mode != "F" && mode != "" {
            return appError.ErrInvMode
        }
    }

    if (idCard != "" || aeonID != "" || agreementNo != "") && mode == "" {
        return appError.ErrRequiedParam
    }

	if  mode != "S" && mode != "F" && mode != "" {
        return appError.ErrInvMode
    }

    return nil
}






// It sends a request to the TCP service and returns the response.
func (s *commonService) GetCustomerInfo(c *gin.Context, getCustomerInfoReq domain.GetCustomerInfoRequest) domain.GetCustomerInfoResult {
	timestamp := time.Now()
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "GetCustomerInfo"
	var domainErr *appError.AppError
	var logLine1 string
	var System, Format, RequestLength, Language  string
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
		return domain.GetCustomerInfoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
			UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.GetCustomerInfoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
			UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.GetCustomerInfoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
			UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["GetCustomerInfo"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetCustomerInfoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
			UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.GetCustomerInfoResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
			UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
			LogLine1:    "",
		}
	}

    validateResult := ValidateCustomerInfo(getCustomerInfoReq, timestamp)
	s.logger.Info("error in check", validateResult)
    if validateResult != nil {
	    return domain.GetCustomerInfoResult{
		    Response:    nil,
		    AppError:    validateResult,
		    GinCtx:      nil,
		    Timestamp:   timestamp,
		    ReqBody:     nil,
		    RespBody:    nil,
		    DomainError: nil,
			ServiceName: serviceName,
			UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
			UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
		    LogLine1: "",
	    }
    }

	apiLang, exists := c.Get("APILanguage")
    lang := "E"

    if exists {
	   if str, ok := apiLang.(string); ok {
		str = strings.TrimSpace(str)
		 if str != "" {
			lang = string(str[0])
		 }
	   }
    }

	switch {
    case getCustomerInfoReq.Mode == "S" && getCustomerInfoReq.UserRef !="":
	System, Format, RequestLength, Language = "MOB_APP", "001", "00021", lang
	fixedLengthData = format.FormatGetCustomerInfoRequest001And003(getCustomerInfoReq, Language)
    case getCustomerInfoReq.Mode == "S":
	System, Format, RequestLength = "CTI_CLOUD", "004", "00056"
	fixedLengthData = format.FormatGetCustomerInfoRequest004(getCustomerInfoReq)
    default:
	System, Format, RequestLength, Language = "APP_EKYC", "003", "00021", lang
	fixedLengthData = format.FormatGetCustomerInfoRequest001And003(getCustomerInfoReq, Language)
    }

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)

	header := utils.BuildFixedLengthHeader(
		System,
		route.Service,
		Format,
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "GetCustomerInfo", firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo), firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo))
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetCustomerInfoResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
				UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
				LogLine1:    "",
			}
		}

		return domain.GetCustomerInfoResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     getCustomerInfoReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
			UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
			LogLine1:    logLine1,
		}
	}


	formatfromsysi := strings.TrimSpace(responseStr[25:28])
    errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])

	if errorCode != "" {
		switch errorCode {
		case "SVC105":
			domainErr = appError.ErrRequiedParam

		case "SVC117":
			if formatfromsysi == "001" || formatfromsysi == "003" {
				domainErr = appError.ErrUserRefOrAeonID
			} else if formatfromsysi == "004" {
				domainErr = appError.ErrInvIDCardNo
			}

		case "SVC118":
			domainErr = appError.ErrAgreement

		case "SVC269":
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "GetCustomerInfo",  firstNonEmpty(getCustomerInfoReq.SNSNo, getCustomerInfoReq.AEONID), firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo))
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetCustomerInfoResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
				UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
				LogLine1:    "",
			}
		}

		return domain.GetCustomerInfoResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
			UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))

	var getCustomerInfoResponse interface{}

	s.logger.Info("Received downstream TCP response", "response", formatfromsysi)

	switch formatfromsysi {
    case "001":
         getCustomerInfoResponse, err = format.FormatGetCustomerInfoResponse001(responseStr)
    case "004":
         getCustomerInfoResponse, err = format.FormatGetCustomerInfoResponse004(responseStr)
    default:
         getCustomerInfoResponse, err = format.FormatGetCustomerInfoResponse003(responseStr)
    }
	if err != nil {
		s.logger.Errorw("Error map getCustomerInfoResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "GetCustomerInfo",  firstNonEmpty(getCustomerInfoReq.SNSNo, getCustomerInfoReq.AEONID), firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo))
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.GetCustomerInfoResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
				UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "GetCustomerInfo", firstNonEmpty(getCustomerInfoReq.SNSNo, getCustomerInfoReq.AEONID), firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo))
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.GetCustomerInfoResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
			UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.GetCustomerInfoResult{
		Response:    getCustomerInfoResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     getCustomerInfoReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserToken:	 firstNonEmpty(getCustomerInfoReq.AEONID, getCustomerInfoReq.SNSNo),
		UserRef:     firstNonEmpty(getCustomerInfoReq.UserRef, getCustomerInfoReq.IDCardNo, getCustomerInfoReq.AgreementNo),
		LogLine1:    logLine1,
	}
}

// It sends a request to the TCP service and returns the response.
func (s *commonService) CheckApplyCondition(c *gin.Context, checkApplyConditionReq domain.CheckApplyConditionRequest) domain.CheckApplyConditionResult {
	timestamp := time.Now()
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "CheckApplyCondition"
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
		return domain.CheckApplyConditionResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.CheckApplyConditionResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.CheckApplyConditionResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["CheckApplyCondition"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.CheckApplyConditionResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.CheckApplyConditionResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if !HasValidApplyCardItem(checkApplyConditionReq.ApplyCardList) {
		return domain.CheckApplyConditionResult{
		Response:    nil,
		AppError:    appError.ErrRequiedParam,
		GinCtx:      nil,
		Timestamp:   timestamp,
		ReqBody:     nil,
		RespBody:    nil,
		DomainError: nil,
		ServiceName: serviceName,
		UserRef:     checkApplyConditionReq.IDCardNo,
		LogLine1:    "",
		}
	}

	switch checkApplyConditionReq.Channel {
		case "L", "F", "A", "W", "R", "O", "E":
	
		default:
		s.logger.Errorw("Invalid Channel", "Channel", checkApplyConditionReq.Channel)
		return domain.CheckApplyConditionResult{
		Response:    nil,
		AppError:    appError.ErrRequiedParam,
		GinCtx:      nil,
		Timestamp:   timestamp,
		ReqBody:     nil,
		RespBody:    nil,
		DomainError: nil,
		ServiceName: serviceName,
		UserRef:     checkApplyConditionReq.IDCardNo,
		LogLine1:    "",
		}
	}

	if len(checkApplyConditionReq.ApplyCardList) != checkApplyConditionReq.TotalApplyCard {
		s.logger.Errorw("Mismatch in number of cards", "Expected", checkApplyConditionReq.TotalApplyCard, "Actual", len(checkApplyConditionReq.ApplyCardList),)
		return domain.CheckApplyConditionResult{
			Response:    nil,
			AppError:    appError.ErrInvTotalOfList,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionReq.IDCardNo,
			LogLine1:    "",
			}
	}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatCheckApplyConditionRequest(checkApplyConditionReq)

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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "CheckApplyCondition", checkApplyConditionReq.IDCardNo, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.CheckApplyConditionResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     checkApplyConditionReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.CheckApplyConditionResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     checkApplyConditionReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}


	errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])
	if errorCode != "" {
		switch errorCode {
		case "SVC173":
			domainErr = appError.ErrInvAppNo
		case "SVC157" :
			domainErr = appError.ErrInvAppChannel
		case "SVC171" :
			domainErr = appError.ErrInvHBDFormat
		case "SVC172" : 
			domainErr = appError.ErrInvSupHBDFormat
		case "SVC165" : 
			domainErr = appError.ErrInvAppDateFormat
		case "SVC127" : 
			domainErr = appError.ErrInvBranchCode
		case "SVC161","SVC166","SVC167","SVC178","SVC179","SVC180","SVC181": 
			domainErr = appError.ErrInvSourceCode
		case "SVC158" : 
			domainErr = appError.ErrInvTotalOfList
		case "SVC102","SVC183":
			domainErr = appError.ErrInvCardCode
		case "SVC174" : 
			domainErr = appError.ErrInvCardAppType
		case "SVC185" : 
			domainErr = appError.ErrInvViCardFlag
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "CheckApplyCondition", checkApplyConditionReq.IDCardNo, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.CheckApplyConditionResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     checkApplyConditionReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.CheckApplyConditionResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	checkApplyConditionResponse, err := format.FormatCheckApplyConditionResponse(responseStr)
	if err != nil {
		s.logger.Errorw("Error map CheckApplyConditionResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "CheckApplyCondition", checkApplyConditionReq.IDCardNo, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.CheckApplyConditionResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     checkApplyConditionReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "CheckApplyCondition", checkApplyConditionReq.IDCardNo, "")
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.CheckApplyConditionResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.CheckApplyConditionResult{
		Response:    &checkApplyConditionResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     checkApplyConditionReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserRef:     checkApplyConditionReq.IDCardNo,
		LogLine1:    logLine1,
	}
}

// It sends a request to the TCP service and returns the response.
func (s *commonService) CheckApplyCondition2ndCard(c *gin.Context, checkApplyConditionCondition2ndCardReq domain.CheckApplyCondition2ndCardRequest) domain.CheckApplyCondition2ndCardResult {
	timestamp := time.Now()
	routeKey := utils.GetRouteKey(c)
	const destinationName = "systemI"
	serviceName := "CheckApplyCondition2ndCard"
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
		return domain.CheckApplyCondition2ndCardResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
			LogLine1:    "",
		}
	}

	destination, ok := s.destinations[destinationName]
	if !ok {
		s.logger.Errorw("TCP Destination configuration not found", "destinationName", destinationName)
		return domain.CheckApplyCondition2ndCardResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
			LogLine1:    "",
		}
	}
	if destination.Type != "tcp" {
		s.logger.Errorw("Destination type is not TCP", "destinationName", destinationName, "type", destination.Type)
		return domain.CheckApplyCondition2ndCardResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
			LogLine1:    "",
		}
	}

	portList, ok := destination.Ports["CheckApplyCondition2ndCard"]
	if !ok || len(portList) == 0 {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.CheckApplyCondition2ndCardResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
			LogLine1:    "",
		}
	}
	port := utils.RandomPortFromList(portList)
	if port == "" {
		s.logger.Errorw("Invalid port configuration", "port", portList)
		return domain.CheckApplyCondition2ndCardResult{
			Response:    nil,
			AppError:    appError.ErrService,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if !HasValidApply2ndCardItem(checkApplyConditionCondition2ndCardReq.CheckApply2ndCardList) {
		return domain.CheckApplyCondition2ndCardResult{
		Response:    nil,
		AppError:    appError.ErrRequiedParam,
		GinCtx:      nil,
		Timestamp:   timestamp,
		ReqBody:     nil,
		RespBody:    nil,
		DomainError: nil,
		ServiceName: serviceName,
		UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
		LogLine1:    "",
		}
	}

	switch checkApplyConditionCondition2ndCardReq.Channel {
		case "L", "F", "A", "W", "R", "O", "E", "MobileApp", "EKYC", "Lounge", "Branch", "Web":
	
		default:
		s.logger.Errorw("Invalid Channel", "Channel", checkApplyConditionCondition2ndCardReq.Channel)
		return domain.CheckApplyCondition2ndCardResult{
		Response:    nil,
		AppError:    appError.ErrApiChannel,
		GinCtx:      nil,
		Timestamp:   timestamp,
		ReqBody:     nil,
		RespBody:    nil,
		DomainError: nil,
		ServiceName: serviceName,
		UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
		LogLine1:    "",
		}
	}

	if len(checkApplyConditionCondition2ndCardReq.CheckApply2ndCardList) != checkApplyConditionCondition2ndCardReq.TotalOfApplyCard {
	s.logger.Errorw("Mismatch in number of cards", 
		"Expected", checkApplyConditionCondition2ndCardReq.TotalOfApplyCard, 
		"Actual", len(checkApplyConditionCondition2ndCardReq.CheckApply2ndCardList),
	)
		return domain.CheckApplyCondition2ndCardResult{
			Response:    nil,
			AppError:    appError.ErrInvTotalOfList,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
			LogLine1:    "",
			}
		}

	formattedRequestID := utils.PadOrTruncate(apiRequestID, 20)
	fixedLengthData := format.FormatCheckApplyCondition2ndCardRequest(checkApplyConditionCondition2ndCardReq)

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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, domainErr, "", destination.IP+":"+port, serviceName, "CheckApplyCondition2ndCard", checkApplyConditionCondition2ndCardReq.IDCardNo, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.CheckApplyCondition2ndCardResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.CheckApplyCondition2ndCardResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     checkApplyConditionCondition2ndCardReq,
			RespBody:    formatErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}


	errorCode := strings.TrimSpace(responseStr[67:73])
	errorMessage := strings.TrimSpace(responseStr[73:123])
	if errorCode != "" {
		switch errorCode {
		case "SVC101":
			domainErr = appError.ErrInvCardCode
		case "SVC105","SVC164":
			domainErr = appError.ErrRequiedParam
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

		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, domainErr, "", destination.IP+":"+port, serviceName, "CheckApplyCondition2ndCard", checkApplyConditionCondition2ndCardReq.IDCardNo, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.CheckApplyCondition2ndCardResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
				LogLine1:    "",
			}
		}

		return domain.CheckApplyCondition2ndCardResult{
			Response:    nil,
			AppError:    nil,
			GinCtx:      c,
			Timestamp:   timestamp,
			ReqBody:     formatReq,
			RespBody:    domainErr,
			DomainError: domainErr,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
			LogLine1:    logLine1,
		}
	}

	s.logger.Info("Received downstream TCP response", "response", string(responseStr))
	checkApplyCondition2ndCardResponse, err := format.FormatCheckApplyCondition2ndCardResponse(responseStr)
	if err != nil {
		s.logger.Errorw("Error map CheckApplyCondition2ndCardResponse:", err)

		formatErr = map[string]string{"data": err.Error()}
		logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatErr, nil, "", destination.IP+":"+port, serviceName, "CheckApplyCondition2ndCard", checkApplyConditionCondition2ndCardReq.IDCardNo, "")
		if logLine1 == "" {
			s.logger.Errorw("Error generating log: %v", logLine1)
			return domain.CheckApplyCondition2ndCardResult{
				Response:    nil,
				AppError:    appError.ErrInternalServer,
				GinCtx:      nil,
				Timestamp:   timestamp,
				ReqBody:     nil,
				RespBody:    nil,
				DomainError: nil,
				ServiceName: serviceName,
				UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
				LogLine1:    "",
			}
		}
	}

	logLine1 = elkLog.GenerateELKLogLine(c, timestamp, formatReq, formatResp, nil, "", destination.IP+":"+port, serviceName, "CheckApplyCondition2ndCard", checkApplyConditionCondition2ndCardReq.IDCardNo, "")
	if logLine1 == "" {
		s.logger.Errorw("Error generating log: %v", logLine1)
		return domain.CheckApplyCondition2ndCardResult{
			Response:    nil,
			AppError:    appError.ErrInternalServer,
			GinCtx:      nil,
			Timestamp:   timestamp,
			ReqBody:     nil,
			RespBody:    nil,
			DomainError: nil,
			ServiceName: serviceName,
			UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
			LogLine1:    "",
		}
	}

	if domainErr != nil && domainErr.ErrorCode == "" {
		domainErr = nil
	}
	return domain.CheckApplyCondition2ndCardResult{
		Response:    &checkApplyCondition2ndCardResponse,
		AppError:    nil,
		GinCtx:      c,
		Timestamp:   timestamp,
		ReqBody:     checkApplyConditionCondition2ndCardReq,
		RespBody:    formatResp,
		DomainError: domainErr,
		ServiceName: serviceName,
		UserRef:     checkApplyConditionCondition2ndCardReq.IDCardNo,
		LogLine1:    logLine1,
	}
}