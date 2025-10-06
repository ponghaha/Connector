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

// creditCardService defines the interface
type creditCardService interface {
	GetCardSales(c *gin.Context, reqData domain.GetCardSalesRequest) domain.GetCardSalesResult
	GetBigCardInfo(c *gin.Context, reqData domain.GetBigCardInfoRequest) domain.GetBigCardInfoResult
	GetCardDelinquent(c *gin.Context, reqData domain.GetCardDelinquentRequest) domain.GetCardDelinquentResult
	GetFullpan(c *gin.Context, reqData domain.GetFullpanRequest) domain.GetFullpanResult
	GetCardEnroll(c *gin.Context, reqData domain.GetCardEnrollRequest) domain.GetCardEnrollResult
}

// creditCardHandler handles all customer-related API requests
type creditCardHandler struct {
	service   creditCardService
	validator *validator.Validate
	logger    *zap.SugaredLogger
	apikey    *utils.APIKeyRepository
	config    *config.Config
}

// NewCreditCardHandler creates a new instance of creditCardHandler
func NewCreditCardHandler(s creditCardService, logger *zap.SugaredLogger, apikey *utils.APIKeyRepository, cfg *config.Config) *creditCardHandler {
	return &creditCardHandler{
		service:   s,
		validator: validator.New(),
		logger:    logger,
		apikey:    apikey,
		config:    cfg,
	}
}

// RegisterRoutes registers all routes related to CreditCard to the router group
func (h *creditCardHandler) RegisterRoutes(rg *gin.RouterGroup) {
	creditCardRoutes := rg.Group("/CreditCard")
	{
		creditCardRoutes.POST("/GetCardSales", h.GetCardSales)
		creditCardRoutes.POST("/GetBigCardInfo", h.GetBigCardInfo)
		creditCardRoutes.POST("/GetCardDelinquent", h.GetCardDelinquent)
		creditCardRoutes.POST("/GetFullPAN", h.GetFullpan)
		creditCardRoutes.POST("/GetCardEnroll", h.GetCardEnroll)
	}
}

// GetCardSales godoc
// @Tags         CreditCard 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.GetCardSalesRequest  false  "BodyRequest"
// @Success      200  {object}        domain.GetCardSalesResponse
// @Router       /Api/CreditCard/GetCardSales [post]
func (h *creditCardHandler) GetCardSales(c *gin.Context) {
	var req domain.GetCardSalesRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetCardSales"

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

	getCardSalesResult := h.service.GetCardSales(c, req)
	if getCardSalesResult.AppError != nil {
		handleErrorResponse(c, getCardSalesResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getCardSalesResult.DomainError != nil {
		responseError = getCardSalesResult.DomainError
	}
	if !elkLog.FinalELKLog(getCardSalesResult.GinCtx, &logList, getCardSalesResult.Timestamp, req, getCardSalesResult.Response, getCardSalesResult.DomainError, getCardSalesResult.ServiceName, getCardSalesResult.UserRef, "", []string{getCardSalesResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getCardSalesResult.Response)
}

// GetBigCardInfo godoc
// @Tags         CreditCard 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.GetBigCardInfoRequest  false  "BodyRequest"
// @Success      200  {object}        domain.GetBigCardInfoResponse
// @Router       /Api/CreditCard/GetBigCardInfo [post]
func (h *creditCardHandler) GetBigCardInfo(c *gin.Context) {
	var req domain.GetBigCardInfoRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetBigCardInfo"

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

	getBigCardInfoResult := h.service.GetBigCardInfo(c, req)
	if getBigCardInfoResult.AppError != nil {
		handleErrorResponse(c, getBigCardInfoResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getBigCardInfoResult.DomainError != nil {
		responseError = getBigCardInfoResult.DomainError
	}
	if !elkLog.FinalELKLog(getBigCardInfoResult.GinCtx, &logList, getBigCardInfoResult.Timestamp, req, getBigCardInfoResult.Response, getBigCardInfoResult.DomainError, getBigCardInfoResult.ServiceName, getBigCardInfoResult.UserToken, "", []string{getBigCardInfoResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getBigCardInfoResult.Response)
}

// GetCardDelinquent godoc
// @Tags         CreditCard 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.GetCardDelinquentRequest  false  "BodyRequest"
// @Success      200  {object}        domain.GetCardDelinquentResponse
// @Router       /Api/CreditCard/GetCardDelinquent [post]
func (h *creditCardHandler) GetCardDelinquent(c *gin.Context) {
	var req domain.GetCardDelinquentRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetCardDelinquent"

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

	getCardDelinquentResult := h.service.GetCardDelinquent(c, req)
	if getCardDelinquentResult.AppError != nil {
		handleErrorResponse(c, getCardDelinquentResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getCardDelinquentResult.DomainError != nil {
		responseError = getCardDelinquentResult.DomainError
	}
	if !elkLog.FinalELKLog(getCardDelinquentResult.GinCtx, &logList, getCardDelinquentResult.Timestamp, req, getCardDelinquentResult.Response, getCardDelinquentResult.DomainError, getCardDelinquentResult.ServiceName, getCardDelinquentResult.UserRef, "", []string{getCardDelinquentResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getCardDelinquentResult.Response)
}

// GetFullpan godoc
// @Tags         CreditCard 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.GetFullpanRequest  false  "BodyRequest"
// @Success      200  {object}        domain.GetFullpanResponse
// @Router       /Api/CreditCard/GetFullpan [post]
func (h *creditCardHandler) GetFullpan(c *gin.Context) {
	var req domain.GetFullpanRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetFullpan"

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

	getFullpanResult := h.service.GetFullpan(c, req)
	if getFullpanResult.AppError != nil {
		handleErrorResponse(c, getFullpanResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getFullpanResult.DomainError != nil {
		responseError = getFullpanResult.DomainError
	}
	if !elkLog.FinalELKLog(getFullpanResult.GinCtx, &logList, getFullpanResult.Timestamp, req, getFullpanResult.Response, getFullpanResult.DomainError, getFullpanResult.ServiceName, getFullpanResult.UserRef, "", []string{getFullpanResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getFullpanResult.Response)
}

// GetCardEnroll godoc
// @Tags         CreditCard 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.GetCardEnrollRequest  false  "BodyRequest"
// @Success      200  {object}        domain.GetCardEnrollResponse
// @Router       /Api/CreditCard/GetCardEnroll [post]
func (h *creditCardHandler) GetCardEnroll(c *gin.Context) {
	var req domain.GetCardEnrollRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetCardEnroll"

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

	getCardEnrollResult := h.service.GetCardEnroll(c, req)
	if getCardEnrollResult.AppError != nil {
		handleErrorResponse(c, getCardEnrollResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getCardEnrollResult.DomainError != nil {
		responseError = getCardEnrollResult.DomainError
	}
	if !elkLog.FinalELKLog(getCardEnrollResult.GinCtx, &logList, getCardEnrollResult.Timestamp, req, getCardEnrollResult.Response, getCardEnrollResult.DomainError, getCardEnrollResult.ServiceName, getCardEnrollResult.UserRef, "", []string{getCardEnrollResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getCardEnrollResult.Response)
}