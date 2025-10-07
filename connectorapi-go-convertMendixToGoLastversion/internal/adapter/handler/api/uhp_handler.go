package handler

import (
	"net/http"
	"time"
	
	"connectorapi-go/internal/adapter/utils"
	"connectorapi-go/internal/core/domain"
	"connectorapi-go/pkg/config"
	appError "connectorapi-go/pkg/error"
	elkLog "connectorapi-go/internal/adapter/client/elk"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"
)

// uhpService defines the interface
type uhpService interface {
	GetRedbookInfo(c *gin.Context, reqData domain.GetRedbookInfoRequest) domain.GetRedbookInfoResult
	GetDealerCommission(c *gin.Context, reqData domain.GetDealerCommissionRequest) domain.GetDealerCommissionResult
	GetDealerAgreement(c *gin.Context, reqData domain.GetDealerAgreementRequest) domain.GetDealerAgreementResult
}

// uhpHandler handles all customer-related API requests
type uhpHandler struct {
	service   uhpService
	validator *validator.Validate
	logger    *zap.SugaredLogger
	apikey    *utils.APIKeyRepository
	config    *config.Config
}

// NewUhpHandler creates a new instance of uhpHandler
func NewUhpHandler(s uhpService, logger *zap.SugaredLogger, apikey *utils.APIKeyRepository, cfg *config.Config) *uhpHandler {
	return &uhpHandler{
		service:   s,
		validator: validator.New(),
		logger:    logger,
		apikey:    apikey,
		config:    cfg,
	}
}

// RegisterRoutes registers all routes related to Common to the router group
func (h *uhpHandler) RegisterRoutes(rg *gin.RouterGroup) {
	uhpRoutes := rg.Group("/uhp")
	{
		uhpRoutes.POST("/GetRedbookInfo", h.GetRedbookInfo)
		uhpRoutes.POST("/GetDealerCommission", h.GetDealerCommission)
		uhpRoutes.POST("/GetDealerAgreement", h.GetDealerAgreement)
	}
}

// GetRedbookInfo godoc
// @Tags         Uhp 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        DeviceOS             header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        RequestID            header    string                      false  "RequestID"
// @Param        request              body      domain.GetRedbookInfoRequest  false  "BodyRequest"
// @Success      200  {object}        domain.GetRedbookInfoResponse
// @Router       /Api/uhp/GetRedbookInfo [post]
func (h *uhpHandler) GetRedbookInfo(c *gin.Context) {
	var req domain.GetRedbookInfoRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetRedbookInfo"

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, appError.ErrService)
		return
	}

	appErr := ValidateHeadersForApiKeyAndApiRequestID(c, c.Request.Method, c.FullPath(), h.apikey, h.logger)
	if appErr != nil {
		handleErrorResponse(c, appErr)
		if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appErr, serviceName, "", "", nil, h.logger, h.config.ELKPath, handleErrorResponse) {
			return
		}
		return
	}

	if err := h.validator.Struct(req); err != nil {
    	appErr := HandleValidationError(err)
		handleErrorResponse(c, appErr)
		if appErr.ErrorCode == "SYS500" {
			return
		}
		if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appErr, serviceName, "", "", nil, h.logger, h.config.ELKPath, handleErrorResponse) {
			return
		}
    	return
	}

	getRedbookInfoResult := h.service.GetRedbookInfo(c, req)
	if getRedbookInfoResult.AppError != nil {
		handleErrorResponse(c, getRedbookInfoResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getRedbookInfoResult.DomainError != nil {
		responseError = getRedbookInfoResult.DomainError
	}
	if !elkLog.FinalELKLog(getRedbookInfoResult.GinCtx, &logList, getRedbookInfoResult.Timestamp, req, getRedbookInfoResult.Response, getRedbookInfoResult.DomainError, getRedbookInfoResult.ServiceName, getRedbookInfoResult.UserRef, "", []string{getRedbookInfoResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getRedbookInfoResult.Response)
}

// GetDealerCommission godoc
// @Tags         Uhp 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        DeviceOS             header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        RequestID            header    string                      false  "RequestID"
// @Param        request              body      domain.GetDealerCommissionRequest  false  "BodyRequest"
// @Success      200  {object}        domain.GetDealerCommissionResponse
// @Router       /Api/uhp/GetDealerCommission [post]
func (h *uhpHandler) GetDealerCommission(c *gin.Context) {
	var req domain.GetDealerCommissionRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetDealerCommission"

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, appError.ErrService)
		return
	}

	appErr := ValidateHeadersForApiKeyAndApiRequestID(c, c.Request.Method, c.FullPath(), h.apikey, h.logger)
	if appErr != nil {
		handleErrorResponse(c, appErr)
		if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appErr, serviceName, "", "", nil, h.logger, h.config.ELKPath, handleErrorResponse) {
			return
		}
		return
	}

	if err := h.validator.Struct(req); err != nil {
    	appErr := HandleValidationError(err)
		handleErrorResponse(c, appErr)
		if appErr.ErrorCode == "SYS500" {
			return
		}
		if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appErr, serviceName, "", "", nil, h.logger, h.config.ELKPath, handleErrorResponse) {
			return
		}
    	return
	}

	getDealerCommissionResult := h.service.GetDealerCommission(c, req)
	if getDealerCommissionResult.AppError != nil {
		handleErrorResponse(c, getDealerCommissionResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getDealerCommissionResult.DomainError != nil {
		responseError = getDealerCommissionResult.DomainError
	}
	if !elkLog.FinalELKLog(getDealerCommissionResult.GinCtx, &logList, getDealerCommissionResult.Timestamp, req, getDealerCommissionResult.Response, getDealerCommissionResult.DomainError, getDealerCommissionResult.ServiceName, getDealerCommissionResult.UserRef, "", []string{getDealerCommissionResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getDealerCommissionResult.Response)
}

// GetDealerAgreement godoc
// @Tags         Uhp 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        DeviceOS             header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        RequestID            header    string                      false  "RequestID"
// @Param        request              body      domain.GetDealerAgreementRequest  false  "BodyRequest"
// @Success      200  {object}        domain.GetDealerAgreementResponse
// @Router       /Api/uhp/v [post]
func (h *uhpHandler) GetDealerAgreement(c *gin.Context) {
	var req domain.GetDealerAgreementRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetDealerAgreement"

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, appError.ErrService)
		return
	}

	appErr := ValidateHeadersForApiKeyAndApiRequestID(c, c.Request.Method, c.FullPath(), h.apikey, h.logger)
	if appErr != nil {
		handleErrorResponse(c, appErr)
		if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appErr, serviceName, "", "", nil, h.logger, h.config.ELKPath, handleErrorResponse) {
			return
		}
		return
	}

	if err := h.validator.Struct(req); err != nil {
    	appErr := HandleValidationError(err)
		handleErrorResponse(c, appErr)
		if appErr.ErrorCode == "SYS500" {
			return
		}
		if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appErr, serviceName, "", "", nil, h.logger, h.config.ELKPath, handleErrorResponse) {
			return
		}
    	return
	}

	getDealerAgreementResult := h.service.GetDealerAgreement(c, req)
	if getDealerAgreementResult.AppError != nil {
		handleErrorResponse(c, getDealerAgreementResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getDealerAgreementResult.DomainError != nil {
		responseError = getDealerAgreementResult.DomainError
	}
	if !elkLog.FinalELKLog(getDealerAgreementResult.GinCtx, &logList, getDealerAgreementResult.Timestamp, req, getDealerAgreementResult.Response, getDealerAgreementResult.DomainError, getDealerAgreementResult.ServiceName, getDealerAgreementResult.UserRef, "", []string{getDealerAgreementResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getDealerAgreementResult.Response)
}