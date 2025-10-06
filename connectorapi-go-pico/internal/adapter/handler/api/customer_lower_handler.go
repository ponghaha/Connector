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

// customer_lowerService defines the interface
type customer_lowerService interface {
	GetCustomerInfoMobileNo(c *gin.Context, reqData domain.GetCustomerInfoMobileNoRequest) domain.GetCustomerInfoMobileNoResult
	CheckAeonCustomer(c *gin.Context, reqData domain.CheckAeonCustomerRequest) domain.CheckAeonCustomerResult
}

// customerLowerHandler handles all customer-related API requests
type customerLowerHandler struct {
	service   customer_lowerService
	validator *validator.Validate
	logger    *zap.SugaredLogger
	apikey    *utils.APIKeyRepository
	config    *config.Config
}

// NewCustomerLowerHandler creates a new instance of customerHandler
func NewCustomerLowerHandler(s customer_lowerService, logger *zap.SugaredLogger, apikey *utils.APIKeyRepository, cfg *config.Config) *customerLowerHandler {
	return &customerLowerHandler{
		service:   s,
		validator: validator.New(),
		logger:    logger,
		apikey:    apikey,
		config:    cfg,
	}
}

// RegisterRoutes registers all routes related to customer to the router group
func (h *customerLowerHandler) RegisterRoutes(rg *gin.RouterGroup) {
	customer_lowerRoutes := rg.Group("/customer") 
	{
		customer_lowerRoutes.POST("getcustomerinfo/mobileno", h.GetCustomerInfoMobileNo)
		customer_lowerRoutes.POST("checkaeoncustomer", h.CheckAeonCustomer)
	}
}

// GetCustomerInfoMobileNo godoc
// @Tags         customer 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.GetCustomerInfoMobileNoRequest  false  "BodyRequest"
// @Success      200  {object}        domain.GetCustomerInfoMobileNoResponse
// @Router       /Api/customer/GetCustomerInfoMobileNo [post]
func (h *customerLowerHandler) GetCustomerInfoMobileNo(c *gin.Context) {
	var req domain.GetCustomerInfoMobileNoRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetCustomerInfoMobileNo"

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

	getCustomerInfoMobileNoResult := h.service.GetCustomerInfoMobileNo(c, req)
	if getCustomerInfoMobileNoResult.AppError != nil {
		handleErrorResponse(c, getCustomerInfoMobileNoResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getCustomerInfoMobileNoResult.DomainError != nil {
		responseError = getCustomerInfoMobileNoResult.DomainError
	}
	if !elkLog.FinalELKLog(getCustomerInfoMobileNoResult.GinCtx, &logList, getCustomerInfoMobileNoResult.Timestamp, req, getCustomerInfoMobileNoResult.Response, getCustomerInfoMobileNoResult.DomainError, getCustomerInfoMobileNoResult.ServiceName, getCustomerInfoMobileNoResult.UserRef, "", []string{getCustomerInfoMobileNoResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getCustomerInfoMobileNoResult.Response)
}

// CheckAeonCustomer godoc
// @Tags         customer 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.CheckAeonCustomerRequest  false  "BodyRequest"
// @Success      200  {object}        domain.CheckAeonCustomerResponse
// @Router       /Api/customer/CheckAeonCustomer [post]
func (h *customerLowerHandler) CheckAeonCustomer(c *gin.Context) {
	var req domain.CheckAeonCustomerRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "CheckAeonCustomer"

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, appError.ErrInternalServer)
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

	checkAeonCustomerResult := h.service.CheckAeonCustomer(c, req)
	if checkAeonCustomerResult.AppError != nil {
		handleErrorResponse(c, checkAeonCustomerResult.AppError)
		return
	}

	var responseError *appError.AppError
	if checkAeonCustomerResult.DomainError != nil {
		responseError = checkAeonCustomerResult.DomainError
	}
	if !elkLog.FinalELKLog(checkAeonCustomerResult.GinCtx, &logList, checkAeonCustomerResult.Timestamp, req, checkAeonCustomerResult.Response, checkAeonCustomerResult.DomainError, checkAeonCustomerResult.ServiceName, "", checkAeonCustomerResult.UserRef, []string{checkAeonCustomerResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, checkAeonCustomerResult.Response)
}