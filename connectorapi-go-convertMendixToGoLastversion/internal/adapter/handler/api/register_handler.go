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

// registerService defines the interface
type registerService interface {
	CheckRegister(c *gin.Context, reqData domain.CheckRegisterRequest) domain.CheckRegisterResult
	CheckRegisterSocial(c *gin.Context, reqData domain.CheckRegisterSocialRequest) domain.CheckRegisterSocialResult
	UpdateUserToken(c *gin.Context, reqData domain.UpdateUserTokenRequest) domain.UpdateUserTokenResult
}

// registerHandler handles all customer-related API requests
type registerHandler struct {
	service   registerService
	validator *validator.Validate
	logger    *zap.SugaredLogger
	apikey    *utils.APIKeyRepository
	config    *config.Config
}

// NewRegisterHandler creates a new instance of registerHandler
func NewRegisterHandler(s registerService, logger *zap.SugaredLogger, apikey *utils.APIKeyRepository, cfg *config.Config) *registerHandler {
	return &registerHandler{
		service:   s,
		validator: validator.New(),
		logger:    logger,
		apikey:    apikey,
		config:    cfg,
	}
}

// RegisterRoutes registers all routes related to Register to the router group
func (h *registerHandler) RegisterRoutes(rg *gin.RouterGroup) {
	registerRoutes := rg.Group("/Register")
	{
		registerRoutes.POST("/CheckRegister", h.CheckRegister)
		registerRoutes.POST("/CheckRegisterSocial", h.CheckRegisterSocial)
		registerRoutes.POST("/UpdateUserToken", h.UpdateUserToken)
	}
}

// CheckRegister godoc
// @Tags         Register 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        DeviceOS             header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.CheckRegisterRequest false  "BodyRequest"
// @Success      200  {object}        domain.CheckRegisterResponse
// @Router       /Api/Register/CheckRegister [post]
func (h *registerHandler) CheckRegister(c *gin.Context) {
	var req domain.CheckRegisterRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "CheckRegister"

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

	checkRegisterResult := h.service.CheckRegister(c, req)
	if checkRegisterResult.AppError != nil {
		handleErrorResponse(c, checkRegisterResult.AppError)
		return
	}

	var responseError *appError.AppError
	if checkRegisterResult.DomainError != nil {
		responseError = checkRegisterResult.DomainError
	}
	if !elkLog.FinalELKLog(checkRegisterResult.GinCtx, &logList, checkRegisterResult.Timestamp, req, checkRegisterResult.Response, checkRegisterResult.DomainError, checkRegisterResult.ServiceName, checkRegisterResult.UserRef, "", []string{checkRegisterResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, checkRegisterResult.Response)
}

// CheckRegisterSocial godoc
// @Tags         Register 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.CheckRegisterSocialRequest false  "BodyRequest"
// @Success      200  {object}        domain.CheckRegisterSocialResponse
// @Router       /Api/Register/CheckRegisterSocial [post]
func (h *registerHandler) CheckRegisterSocial(c *gin.Context) {
	var req domain.CheckRegisterSocialRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "CheckRegisterSocial"

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

	checkRegisterSocialResult := h.service.CheckRegisterSocial(c, req) 
	if checkRegisterSocialResult.AppError != nil {
		handleErrorResponse(c, checkRegisterSocialResult.AppError)
		return
	}

	var responseError *appError.AppError
	if checkRegisterSocialResult.DomainError != nil {
		responseError = checkRegisterSocialResult.DomainError
	}
	if !elkLog.FinalELKLog(checkRegisterSocialResult.GinCtx, &logList, checkRegisterSocialResult.Timestamp, req, checkRegisterSocialResult.Response, checkRegisterSocialResult.DomainError, checkRegisterSocialResult.ServiceName, checkRegisterSocialResult.UserRef, "", []string{checkRegisterSocialResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, checkRegisterSocialResult.Response)
}

// UpdateUserToken godoc
// @Tags         Register 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.UpdateUserTokenRequest false  "BodyRequest"
// @Success      200  {object}        domain.UpdateUserTokenResponse
// @Router       /Api/Register/UpdateUserToken [post]
func (h *registerHandler) UpdateUserToken(c *gin.Context) {
	var req domain.UpdateUserTokenRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "UpdateUserToken"

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

	updateUserTokenResult := h.service.UpdateUserToken(c, req) 
	if updateUserTokenResult.AppError != nil {
		handleErrorResponse(c, updateUserTokenResult.AppError)
		return
	}

	var responseError *appError.AppError
	if updateUserTokenResult.DomainError != nil {
		responseError = updateUserTokenResult.DomainError
	}
	if !elkLog.FinalELKLog(updateUserTokenResult.GinCtx, &logList, updateUserTokenResult.Timestamp, req, updateUserTokenResult.Response, updateUserTokenResult.DomainError, updateUserTokenResult.ServiceName, updateUserTokenResult.UserRef, "", []string{updateUserTokenResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, updateUserTokenResult.Response)
}