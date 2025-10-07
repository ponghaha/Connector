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

// mobileService defines the interface
type mobileService interface {
	DashboardSummary(c *gin.Context, reqData domain.DashboardSummaryRequest) domain.DashboardSummaryResult
	DashboardDetail(c *gin.Context, reqData domain.DashboardDetailRequest) domain.DashboardDetailResult
	MobileFullPan(c *gin.Context, reqData domain.MobileFullPanRequest) domain.MobileFullPanResult
}

// mobileHandler handles all customer-related API requests
type mobileHandler struct {
	service   mobileService
	validator *validator.Validate
	logger    *zap.SugaredLogger
	apikey    *utils.APIKeyRepository
	config    *config.Config
}

// NewMobileHandler creates a new instance of mobileHandler
func NewMobileHandler(s mobileService, logger *zap.SugaredLogger, apikey *utils.APIKeyRepository, cfg *config.Config) *mobileHandler {
	return &mobileHandler{
		service:   s,
		validator: validator.New(),
		logger:    logger,
		apikey:    apikey,
		config:    cfg,
	}
}

// RegisterRoutes registers all routes related to Mobile to the router group
func (h *mobileHandler) RegisterRoutes(rg *gin.RouterGroup) {
	mobileRoutes := rg.Group("/Mobile")
	{
		mobileRoutes.POST("/DashboardSummary", h.DashboardSummary)
		mobileRoutes.POST("/DashboardDetail", h.DashboardDetail)
		mobileRoutes.POST("/MobileFullPAN", h.MobileFullPan)
	}
}

// DashboardSummary godoc
// @Tags         Mobile 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.DashboardSummaryRequest  false  "BodyRequest"
// @Success      200  {object}        domain.DashboardSummaryResponse
// @Router       /Api/Mobile/DashboardSummary [post]
func (h *mobileHandler) DashboardSummary(c *gin.Context) {
	var req domain.DashboardSummaryRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "DashboardSummary"

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, appError.ErrService)
		return
	}

	apiKey := c.GetHeader("Api-Key")
	if !h.apikey.Validate(apiKey, c.Request.Method, c.FullPath()) {
		h.logger.Errorw("Authorization failed", "path", c.FullPath(), "apiKey", apiKey)
		handleErrorResponse(c, appError.ErrUnauthorized)
		if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appError.ErrUnauthorized, serviceName, "", "", nil, h.logger, h.config.ELKPath, handleErrorResponse) {
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

	if req.IDCardNo != "" && req.Channel != "" && req.AeonID == "" {
		switch req.Channel {
		case "L", "F", "A", "W", "R", "O", "E":
		default:
			handleErrorResponse(c, appError.ErrInvChannel)
			if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appError.ErrInvChannel, serviceName, req.AeonID, req.IDCardNo, nil, h.logger, h.config.ELKPath, handleErrorResponse) {
				return
			}
			return
		}
	} else if req.AeonID != "" {
		appErr := ValidateHeaders(c, c.Request.Method, c.FullPath(), h.apikey, h.logger)
		if appErr != nil {
			handleErrorResponse(c, appErr)
			if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appErr, serviceName, req.AeonID, req.IDCardNo, nil, h.logger, h.config.ELKPath, handleErrorResponse) {
				return
			}
			return
		}
	} else {
		handleErrorResponse(c, appError.ErrRequiedParam)
		if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appError.ErrRequiedParam, serviceName, req.AeonID, req.IDCardNo, nil, h.logger, h.config.ELKPath, handleErrorResponse) {
			return
		}
		return
	}

	dashboardSummaryResult := h.service.DashboardSummary(c, req)
	if dashboardSummaryResult.AppError != nil {
		handleErrorResponse(c, dashboardSummaryResult.AppError)
		return
	}

	var responseError *appError.AppError
	if dashboardSummaryResult.DomainError != nil {
		responseError = dashboardSummaryResult.DomainError
	}
	if !elkLog.FinalELKLog(dashboardSummaryResult.GinCtx, &logList, dashboardSummaryResult.Timestamp, req, dashboardSummaryResult.Response, dashboardSummaryResult.DomainError, dashboardSummaryResult.ServiceName, dashboardSummaryResult.UserToken, dashboardSummaryResult.UserRef, []string{dashboardSummaryResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, dashboardSummaryResult.Response)
}

// DashboardDetail godoc
// @Tags         Mobile 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.DashboardDetailRequest  false  "BodyRequest"
// @Success      200  {object}        domain.DashboardDetailResponse
// @Router       /Api/Mobile/DashboardDetail [post]
func (h *mobileHandler) DashboardDetail(c *gin.Context) {
	var req domain.DashboardDetailRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "DashboardDetail"

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorResponse(c, appError.ErrService)
		return
	}

	apiKey := c.GetHeader("Api-Key")
	if !h.apikey.Validate(apiKey, c.Request.Method, c.FullPath()) {
		h.logger.Errorw("Authorization failed", "path", c.FullPath(), "apiKey", apiKey)
		handleErrorResponse(c, appError.ErrUnauthorized)
		if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appError.ErrUnauthorized, serviceName, "", "", nil, h.logger, h.config.ELKPath, handleErrorResponse) {
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

	if req.IDCardNo != "" && req.Channel != "" && req.AeonID == "" {
		switch req.Channel {
		case "L", "F", "A", "W", "R", "O", "E":
		default:
			handleErrorResponse(c, appError.ErrInvChannel)
			if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appError.ErrInvChannel, serviceName, req.AeonID, req.IDCardNo, nil, h.logger, h.config.ELKPath, handleErrorResponse) {
				return
			}
			return
		}
	} else if req.AeonID != "" {
		appErr := ValidateHeaders(c, c.Request.Method, c.FullPath(), h.apikey, h.logger)
		if appErr != nil {
			handleErrorResponse(c, appErr)
			if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appErr, serviceName, req.AeonID, req.IDCardNo, nil, h.logger, h.config.ELKPath, handleErrorResponse) {
				return
			}
			return
		}
	} else {
		handleErrorResponse(c, appError.ErrRequiedParam)
		if !elkLog.FinalELKLog(c, &logList, timeNow, &req, "", appError.ErrRequiedParam, serviceName, req.AeonID, req.IDCardNo, nil, h.logger, h.config.ELKPath, handleErrorResponse) {
			return
		}
		return
	}

	dashboardDetailResult := h.service.DashboardDetail(c, req)
	if dashboardDetailResult.AppError != nil {
		handleErrorResponse(c, dashboardDetailResult.AppError)
		return
	}

	var responseError *appError.AppError
	if dashboardDetailResult.DomainError != nil {
		responseError = dashboardDetailResult.DomainError
	}
	if !elkLog.FinalELKLog(dashboardDetailResult.GinCtx, &logList, dashboardDetailResult.Timestamp, req, dashboardDetailResult.Response, dashboardDetailResult.DomainError, dashboardDetailResult.ServiceName, dashboardDetailResult.UserToken, dashboardDetailResult.UserRef, []string{dashboardDetailResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, dashboardDetailResult.Response)
}

// MobileFullPan godoc
// @Tags         Mobile 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.MobileFullPanRequest false  "BodyRequest"
// @Success      200  {object}        domain.MobileFullPanResponse
// @Router       /Api/Mobile/MobileFullPAN [post]
func (h *mobileHandler) MobileFullPan(c *gin.Context) {
	var req domain.MobileFullPanRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "MobileFullPan"

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

	mobileFullPanResult := h.service.MobileFullPan(c, req)
	if mobileFullPanResult.AppError != nil {
		handleErrorResponse(c, mobileFullPanResult.AppError)
		return
	}

	var responseError *appError.AppError
	if mobileFullPanResult.DomainError != nil {
		responseError = mobileFullPanResult.DomainError
	}
	if !elkLog.FinalELKLog(mobileFullPanResult.GinCtx, &logList, mobileFullPanResult.Timestamp, req, mobileFullPanResult.Response, mobileFullPanResult.DomainError, mobileFullPanResult.ServiceName, "", mobileFullPanResult.UserRef, []string{mobileFullPanResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, mobileFullPanResult.Response)
}
