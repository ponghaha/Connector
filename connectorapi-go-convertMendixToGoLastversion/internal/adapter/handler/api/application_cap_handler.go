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

// applicationCapService defines the interface
type applicationCapService interface {
	GetApplicationNo(c *gin.Context, reqData domain.GetApplicationNoRequest) domain.GetApplicationNoResult
	SubmitCardApplication(c *gin.Context, reqData domain.SubmitCardApplicationRequest) domain.SubmitCardApplicationResult
}

// applicationCapHandler handles all customer-related API requests
type applicationCapHandler struct {
	service   applicationCapService
	validator *validator.Validate
	logger    *zap.SugaredLogger
	apikey    *utils.APIKeyRepository
	config    *config.Config
}

// NewApplicationCapHandler creates a new instance of applicationCapHandler
func NewApplicationCapHandler(s applicationCapService, logger *zap.SugaredLogger, apikey *utils.APIKeyRepository, cfg *config.Config) *applicationCapHandler {
	return &applicationCapHandler{
		service:   s,
		validator: validator.New(),
		logger:    logger,
		apikey:    apikey,
		config:    cfg,
	}
}

// RegisterRoutes registers all routes related to Application to the router group
func (h *applicationCapHandler) RegisterRoutes(rg *gin.RouterGroup) {
	applicationRoutes := rg.Group("/Application")
	{
		applicationRoutes.POST("/GetApplicationNo", h.GetApplicationNo)
		applicationRoutes.POST("/SubmitCardApplication", h.SubmitCardApplication)
	}
}

// GetApplicationNo godoc
// @Tags         Application 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.GetApplicationNoRequest false  "BodyRequest"
// @Success      200                  {object}  domain.GetApplicationNoResponse
// @Router       /Api/Application/GetApplicationNo [post]
func (h *applicationCapHandler) GetApplicationNo(c *gin.Context) {
	var req domain.GetApplicationNoRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetApplicationNo"

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

	getApplicationNoResult := h.service.GetApplicationNo(c, req)
	if getApplicationNoResult.AppError != nil {
		handleErrorResponse(c, getApplicationNoResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getApplicationNoResult.DomainError != nil {
		responseError = getApplicationNoResult.DomainError
	}
	if !elkLog.FinalELKLog(getApplicationNoResult.GinCtx, &logList, getApplicationNoResult.Timestamp, req, getApplicationNoResult.Response, getApplicationNoResult.DomainError, getApplicationNoResult.ServiceName, "", getApplicationNoResult.UserRef, []string{getApplicationNoResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getApplicationNoResult.Response)
}

// SubmitCardApplication godoc
// @Tags         Application 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.SubmitCardApplicationRequest false  "BodyRequest"
// @Success      200                  {object}  domain.SubmitCardApplicationResponse
// @Router       /Api/Application/SubmitCardApplication [post]
func (h *applicationCapHandler) SubmitCardApplication(c *gin.Context) {
	var req domain.SubmitCardApplicationRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "SubmitCardApplication"

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

	submitCardApplicationResult := h.service.SubmitCardApplication(c, req)
	if submitCardApplicationResult.AppError != nil {
		handleErrorResponse(c, submitCardApplicationResult.AppError)
		return
	}

	var responseError *appError.AppError
	if submitCardApplicationResult.DomainError != nil {
		responseError = submitCardApplicationResult.DomainError
	}
	if !elkLog.FinalELKLog(submitCardApplicationResult.GinCtx, &logList, submitCardApplicationResult.Timestamp, req, submitCardApplicationResult.Response, submitCardApplicationResult.DomainError, submitCardApplicationResult.ServiceName, "", submitCardApplicationResult.UserRef, []string{submitCardApplicationResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, submitCardApplicationResult.Response)
}
