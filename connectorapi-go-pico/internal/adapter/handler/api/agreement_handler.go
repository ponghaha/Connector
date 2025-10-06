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

// agreementService defines the interface
type agreementService interface {
	UpdateStatus(c *gin.Context, reqData domain.UpdateStatusRequest) domain.UpdateStatusResult
	AgreeMentBilling(c *gin.Context, reqData domain.AgreeMentBillingRequest) domain.AgreeMentBillingResult
}

// agreementHandler handles all customer-related API requests
type agreementHandler struct {
	service   agreementService
	validator *validator.Validate
	logger    *zap.SugaredLogger
	apikey    *utils.APIKeyRepository
	config    *config.Config
}

// agreementHandler creates a new instance of agreementHandler
func NewAgreementHandler(s agreementService, logger *zap.SugaredLogger, apikey *utils.APIKeyRepository, cfg *config.Config) *agreementHandler {
	return &agreementHandler{
		service:   s,
		validator: validator.New(),
		logger:    logger,
		apikey:    apikey,
		config:    cfg,
	}
}

// RegisterRoutes registers all routes related to Agreement to the router group
func (h *agreementHandler) RegisterRoutes(rg *gin.RouterGroup) {
	agreementRoutes := rg.Group("/Agreement")
	{
		agreementRoutes.POST("/UpdateStatus", h.UpdateStatus)
		agreementRoutes.POST("/GetBilling", h.AgreeMentBilling)
	}
}

// UpdateStatus godoc
// @Tags         Agreement 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.UpdateStatusRequest  false  "BodyRequest"
// @Success      200  {object}        domain.UpdateStatusResponse
// @Router       /Api/Agreement/UpdateStatus [post]
func (h *agreementHandler) UpdateStatus(c *gin.Context) {
	var req domain.UpdateStatusRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "UpdateAgreementStatus"

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, appError.ErrService)
		return
	}

	appErr := ValidateHeaders(c, c.Request.Method, c.FullPath(), h.apikey, h.logger)
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

	updateStatusResult := h.service.UpdateStatus(c, req)
	if updateStatusResult.AppError != nil {
		handleErrorResponse(c, updateStatusResult.AppError)
		return
	}

	var responseError *appError.AppError
	if updateStatusResult.DomainError != nil {
		responseError = updateStatusResult.DomainError
	}
	if !elkLog.FinalELKLog(updateStatusResult.GinCtx, &logList, updateStatusResult.Timestamp, req, updateStatusResult.Response, updateStatusResult.DomainError, updateStatusResult.ServiceName, updateStatusResult.UserToken, "", []string{updateStatusResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, updateStatusResult.Response)
}

// AgreeMentBilling godoc
// @Tags         Agreement 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Channel              header    string                      false  "Channel"
// @Param        RequestID            header    string                      false  "RequestID"
// @Param        request              body      domain.AgreeMentBillingRequest  false  "BodyRequest"
// @Success      200  {object}        domain.AgreeMentBillingResponse
// @Router       /Api/Agreement/AgreeMentBilling [post]
func (h *agreementHandler) AgreeMentBilling(c *gin.Context) {
	var req domain.AgreeMentBillingRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "AgreeMentBilling"

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, appError.ErrService)
		return
	}

	appErr := ValidateHeaders(c, c.Request.Method, c.FullPath(), h.apikey, h.logger)
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

	agreementBillingResult := h.service.AgreeMentBilling(c, req)
	if agreementBillingResult.AppError != nil {
		handleErrorResponse(c, agreementBillingResult.AppError)
		return
	}

	var responseError *appError.AppError
	if agreementBillingResult.DomainError != nil {
		responseError = agreementBillingResult.DomainError
	}
	if !elkLog.FinalELKLog(agreementBillingResult.GinCtx, &logList, agreementBillingResult.Timestamp, req, agreementBillingResult.Response, agreementBillingResult.DomainError, agreementBillingResult.ServiceName, agreementBillingResult.UserToken, "", []string{agreementBillingResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, agreementBillingResult.Response)
}

