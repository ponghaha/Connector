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

// selfServiceService defines the interface
type selfServiceService interface {
	MyCard(c *gin.Context, reqData domain.MyCardRequest) domain.MyCardResult
	GetAvailableLimit(c *gin.Context, reqData domain.GetAvailableLimitRequest) domain.GetAvailableLimitResult
}

// selfServiceHandler handles all customer-related API requests
type selfServiceHandler struct {
	service   selfServiceService
	validator *validator.Validate
	logger    *zap.SugaredLogger
	apikey    *utils.APIKeyRepository
	config    *config.Config
}

// NewSelfServiceHandler creates a new instance of selfServiceHandler
func NewSelfServiceHandler(s selfServiceService, logger *zap.SugaredLogger, apikey *utils.APIKeyRepository, cfg *config.Config) *selfServiceHandler {
	return &selfServiceHandler{
		service:   s,
		validator: validator.New(),
		logger:    logger,
		apikey:    apikey,
		config:    cfg,
	}
}

// RegisterRoutes registers all routes related to SelfService to the router group
func (h *selfServiceHandler) RegisterRoutes(rg *gin.RouterGroup) {
	selfServiceRoutes := rg.Group("/SelfService")
	{
		selfServiceRoutes.POST("/MyCard", h.MyCard)
		selfServiceRoutes.POST("/GetAvailableLimit", h.GetAvailableLimit)
	}
}

// MyCard godoc
// @Tags         SelfService 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        RequestID            header    string                      false  "RequestID"
// @Param        request              body      domain.MyCardRequest  false  "BodyRequest"
// @Success      200                  {object}  domain.MyCardResponseNormal "Mode Normal"
// @Success      200                  {object}  domain.MyCardResponseAll    "Mode All"
// @Router       /Api/SelfService/MyCard [post]
func (h *selfServiceHandler) MyCard(c *gin.Context) {
	var req domain.MyCardRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "MyCard"

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

	myCardResult := h.service.MyCard(c, req)
	if myCardResult.AppError != nil {
		handleErrorResponse(c, myCardResult.AppError)
		return
	}

	var responseError *appError.AppError
	if myCardResult.DomainError != nil {
		responseError = myCardResult.DomainError
	}
	if !elkLog.FinalELKLog(myCardResult.GinCtx, &logList, myCardResult.Timestamp, req, myCardResult.Response, myCardResult.DomainError, myCardResult.ServiceName, myCardResult.UserRef, "", []string{myCardResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, myCardResult.Response)
}

// GetAvailableLimit godoc
// @Tags         SelfService 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        RequestID            header    string                      false  "RequestID"
// @Param        request              body      domain.GetAvailableLimitRequest  false  "BodyRequest"
// @Success      200  {object}        domain.GetAvailableLimitResponse
// @Router       /Api/SelfService/GetAvailableLimit [post]
func (h *selfServiceHandler) GetAvailableLimit(c *gin.Context) {
	var req domain.GetAvailableLimitRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "GetAvailableLimit"

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

	getAvailableLimitResult := h.service.GetAvailableLimit(c, req)
	if getAvailableLimitResult.AppError != nil {
		handleErrorResponse(c, getAvailableLimitResult.AppError)
		return
	}

	var responseError *appError.AppError
	if getAvailableLimitResult.DomainError != nil {
		responseError = getAvailableLimitResult.DomainError
	}
	if !elkLog.FinalELKLog(getAvailableLimitResult.GinCtx, &logList, getAvailableLimitResult.Timestamp, req, getAvailableLimitResult.Response, getAvailableLimitResult.DomainError, getAvailableLimitResult.ServiceName, getAvailableLimitResult.UserRef, "", []string{getAvailableLimitResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, getAvailableLimitResult.Response)
}