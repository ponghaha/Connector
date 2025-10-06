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

// applicationService defines the interface
type applicationLowerService interface {
	SubmitLoanApplication(c *gin.Context, reqData domain.SubmitLoanApplicationRequest) domain.SubmitLoanApplicationResult
}

// applicationLowerHandler handles all customer-related API requests
type applicationLowerHandler struct {
	service   applicationLowerService
	validator *validator.Validate
	logger    *zap.SugaredLogger
	apikey    *utils.APIKeyRepository
	config    *config.Config
}

// NewApplicationLowerHandler creates a new instance of applicationLowerHandler
func NewApplicationLowerHandler(s applicationLowerService, logger *zap.SugaredLogger, apikey *utils.APIKeyRepository, cfg *config.Config) *applicationLowerHandler {
	return &applicationLowerHandler{
		service:   s,
		validator: validator.New(),
		logger:    logger,
		apikey:    apikey,
		config:    cfg,
	}
}

// RegisterRoutes registers all routes related to application to the router group
func (h *applicationLowerHandler) RegisterRoutes(rg *gin.RouterGroup) {
	applicationRoutes := rg.Group("/application")
	{
		applicationRoutes.POST("/submitloanapplication", h.SubmitLoanApplication)
	}
}

// SubmitLoanApplication godoc
// @Tags         application 
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                      false  "API key"
// @Param        Api-DeviceOS         header    string                      false  "DeviceOS"
// @Param        Api-Channel          header    string                      false  "Channel"
// @Param        Api-RequestID        header    string                      false  "RequestID"
// @Param        request              body      domain.SubmitLoanApplicationRequest  false  "BodyRequest"
// @Success      200  {object}        domain.SubmitLoanApplicationResponse
// @Router       /Api/application/submitloanapplication [post]
func (h *applicationLowerHandler) SubmitLoanApplication(c *gin.Context) {
	var req domain.SubmitLoanApplicationRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "SubmitLoanApplication"

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

	submitLoanApplicationResult := h.service.SubmitLoanApplication(c, req)
	if submitLoanApplicationResult.AppError != nil {
		handleErrorResponse(c, submitLoanApplicationResult.AppError)
		return
	}

	var responseError *appError.AppError
	if submitLoanApplicationResult.DomainError != nil {
		responseError = submitLoanApplicationResult.DomainError
	}
	if !elkLog.FinalELKLog(submitLoanApplicationResult.GinCtx, &logList, submitLoanApplicationResult.Timestamp, req, "", submitLoanApplicationResult.DomainError, submitLoanApplicationResult.ServiceName, "", submitLoanApplicationResult.UserRef, []string{submitLoanApplicationResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	// c.Status(http.StatusNoContent)
	c.Status(http.StatusOK)
}
