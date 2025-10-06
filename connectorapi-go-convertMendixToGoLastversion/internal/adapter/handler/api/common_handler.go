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

// commonService defines the interface
type commonService interface {
	GetCustomerInfo(c *gin.Context, reqData domain.GetCustomerInfoRequest) domain.GetCustomerInfoResult
	CheckApplyCondition(c *gin.Context, reqData domain.CheckApplyConditionRequest) domain.CheckApplyConditionResult
	CheckApplyCondition2ndCard(c *gin.Context, reqData domain.CheckApplyCondition2ndCardRequest) domain.CheckApplyCondition2ndCardResult
}

// commonHandler handles all customer-related API requests
type commonHandler struct {
	service   commonService
	validator *validator.Validate
	logger    *zap.SugaredLogger
	apikey    *utils.APIKeyRepository
	config    *config.Config
}

// NewCommonHandler creates a new instance of commonHandler
func NewCommonHandler(s commonService, logger *zap.SugaredLogger, apikey *utils.APIKeyRepository, cfg *config.Config) *commonHandler {
	return &commonHandler{
		service:   s,
		validator: validator.New(),
		logger:    logger,
		apikey:    apikey,
		config:    cfg,
	}
}

// RegisterRoutes registers all routes related to Common to the router group
func (h *commonHandler) RegisterRoutes(rg *gin.RouterGroup) {
	commonRoutes := rg.Group("/Common")
	{
		commonRoutes.POST("/GetCustomerInfo", h.GetCustomerInfo)
		commonRoutes.POST("/CheckApplyCondition/ApplyCard", h.CheckApplyCondition)
		commonRoutes.POST("/CheckApplyCondition/SecondCard", h.CheckApplyCondition2ndCard)
	}
}

// GetCustomerInfo godoc
// @Tags         Common 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        DeviceOS             header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        RequestID            header    string                      false  "RequestID"
// @Param        request              body      domain.GetCustomerInfoRequest  false  "BodyRequest"
// @Success      200                  {object}  domain.GetCustomerInfoResponse001  "UserRef and Mode S"
// @Success      200                  {object}  domain.GetCustomerInfoResponse003  "AEONID and Mode S"
// @Success      200                  {object}  domain.GetCustomerInfoResponse004  "Other and Mode F"
// @Router       /Api/Common/GetCustomerInfo [post]
func (h *commonHandler) GetCustomerInfo(c *gin.Context) {
	var req domain.GetCustomerInfoRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetCustomerInfo"

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

	getCustomerInfoResult := h.service.GetCustomerInfo(c, req)
	if getCustomerInfoResult.AppError != nil {
		handleErrorResponse(c, getCustomerInfoResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getCustomerInfoResult.DomainError != nil {
		responseError = getCustomerInfoResult.DomainError
	}
	if !elkLog.FinalELKLog(getCustomerInfoResult.GinCtx, &logList, getCustomerInfoResult.Timestamp, req, getCustomerInfoResult.Response, getCustomerInfoResult.DomainError, getCustomerInfoResult.ServiceName, getCustomerInfoResult.UserRef, "", []string{getCustomerInfoResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getCustomerInfoResult.Response)
}

// CheckApplyCondition godoc
// @Tags         Common 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        DeviceOS             header    string                      false  "DeviceOS"
// @Param        RequestID            header    string                      false  "RequestID"
// @Param        request              body      domain.CheckApplyConditionRequest  false  "BodyRequest"
// @Success      200  {object}        domain.CheckApplyConditionResponse
// @Router       /Api/Common/CheckApplyCondition/ApplyCard [post]
func (h *commonHandler) CheckApplyCondition(c *gin.Context) {
	var req domain.CheckApplyConditionRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "CheckApplyCondition"

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

	checkApplyConditionResult := h.service.CheckApplyCondition(c, req)
	if checkApplyConditionResult.AppError != nil {
		handleErrorResponse(c, checkApplyConditionResult.AppError)
		return
	}

	var responseError *appError.AppError
	if checkApplyConditionResult.DomainError != nil {
		responseError = checkApplyConditionResult.DomainError
	}
	if !elkLog.FinalELKLog(checkApplyConditionResult.GinCtx, &logList, checkApplyConditionResult.Timestamp, req, checkApplyConditionResult.Response, checkApplyConditionResult.DomainError, checkApplyConditionResult.ServiceName, checkApplyConditionResult.UserRef, "", []string{checkApplyConditionResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, checkApplyConditionResult.Response)
}

// CheckApplyCondition2ndCard godoc
// @Tags         Common 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        DeviceOS             header    string                      false  "DeviceOS"
// @Param        RequestID            header    string                      false  "RequestID"
// @Param        request              body      domain.CheckApplyCondition2ndCardRequest  false  "BodyRequest"
// @Success      200  {object}        domain.CheckApplyCondition2ndCardResponse
// @Router       /Api/Common/CheckApplyCondition/SecondCard [post]
func (h *commonHandler) CheckApplyCondition2ndCard(c *gin.Context) {
	var req domain.CheckApplyCondition2ndCardRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "CheckApplyCondition2ndCard"

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

	checkApplyCondition2ndCardResult := h.service.CheckApplyCondition2ndCard(c, req)
	if checkApplyCondition2ndCardResult.AppError != nil {
		handleErrorResponse(c, checkApplyCondition2ndCardResult.AppError)
		return
	}

	var responseError *appError.AppError
	if checkApplyCondition2ndCardResult.DomainError != nil {
		responseError = checkApplyCondition2ndCardResult.DomainError
	}
	if !elkLog.FinalELKLog(checkApplyCondition2ndCardResult.GinCtx, &logList, checkApplyCondition2ndCardResult.Timestamp, req, checkApplyCondition2ndCardResult.Response, checkApplyCondition2ndCardResult.DomainError, checkApplyCondition2ndCardResult.ServiceName, checkApplyCondition2ndCardResult.UserRef, "", []string{checkApplyCondition2ndCardResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, checkApplyCondition2ndCardResult.Response)
}