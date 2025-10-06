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

// consentService defines the interface
type consentService interface {
	UpdateConsent(c *gin.Context, reqData domain.UpdateConsentRequest) domain.UpdateConsentResult
	GetConsentList(c *gin.Context, reqData domain.GetConsentListRequest) domain.GetConsentListResult
}

// consentHandler handles all customer-related API requests
type consentHandler struct {
	service   consentService
	validator *validator.Validate
	logger    *zap.SugaredLogger
	apikey    *utils.APIKeyRepository
	config    *config.Config
}

// NewConsentHandler creates a new instance of consentHandler
func NewConsentHandler(s consentService, logger *zap.SugaredLogger, apikey *utils.APIKeyRepository, cfg *config.Config) *consentHandler {
	return &consentHandler{
		service:   s,
		validator: validator.New(),
		logger:    logger,
		apikey:    apikey,
		config:    cfg,
	}
}

// RegisterRoutes registers all routes related to Common to the router group
func (h *consentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	consentRoutes := rg.Group("/Consent")
	{
		consentRoutes.POST("/UpdateConsent", h.UpdateConsent)
		consentRoutes.POST("/GetConsentList", h.GetConsentList)
	}
}

// UpdateConsent godoc
// @Tags         Consent 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        DeviceOS             header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        RequestID            header    string                      false  "RequestID"
// @Param        request              body      domain.UpdateConsentRequest  false  "BodyRequest"
// @Success      200  {object}        domain.UpdateConsentResponse
// @Router       /Api/Consent/UpdateConsent [post]
func (h *consentHandler) UpdateConsent(c *gin.Context) {
	var req domain.UpdateConsentRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "UpdateConsent"

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

	updateConsentResult := h.service.UpdateConsent(c, req)
	if updateConsentResult.AppError != nil {
		handleErrorResponse(c, updateConsentResult.AppError)
		return
	}

	var responseError *appError.AppError
	if updateConsentResult.DomainError != nil {
		responseError = updateConsentResult.DomainError
	}
	if !elkLog.FinalELKLog(updateConsentResult.GinCtx, &logList, updateConsentResult.Timestamp, req, updateConsentResult.Response, updateConsentResult.DomainError, updateConsentResult.ServiceName, updateConsentResult.UserRef, "", []string{updateConsentResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, updateConsentResult.Response)
}

// GetConsentList godoc
// @Tags         Consent 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        DeviceOS             header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        RequestID            header    string                      false  "RequestID"
// @Param        request              body      domain.GetConsentListRequest  false  "BodyRequest"
// @Success      200  {object}        domain.GetConsentListResponse
// @Router       /Api/Consent/GetConsentList [post]
func (h *consentHandler) GetConsentList(c *gin.Context) {
	var req domain.GetConsentListRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetConsentList"

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

	getConsentListResult := h.service.GetConsentList(c, req)
	if getConsentListResult.AppError != nil {
		handleErrorResponse(c, getConsentListResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getConsentListResult.DomainError != nil {
		responseError = getConsentListResult.DomainError
	}
	if !elkLog.FinalELKLog(getConsentListResult.GinCtx, &logList, getConsentListResult.Timestamp, req, getConsentListResult.Response, getConsentListResult.DomainError, getConsentListResult.ServiceName, getConsentListResult.UserRef, "", []string{getConsentListResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getConsentListResult.Response)
}