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

// collectionService defines the interface
type collectionService interface {
	CollectionDetail(c *gin.Context, reqData domain.CollectionDetailRequest) domain.CollectionDetailResult
	CollectionLog(c *gin.Context, reqData domain.CollectionLogRequest) domain.CollectionLogResult
}

// collectionHandler handles all Collection-related API requests
type collectionHandler struct {
	service   collectionService
	validator *validator.Validate
	logger    *zap.SugaredLogger
	apikey    *utils.APIKeyRepository
	config    *config.Config
}

// collectionHandler creates a new instance of collectionHandler
func NewCollectionHandler(s collectionService, logger *zap.SugaredLogger, apikey *utils.APIKeyRepository, cfg *config.Config) *collectionHandler {
	return &collectionHandler{
		service: s,
		validator: validator.New(),
		logger:  logger,
		apikey:  apikey,
		config:    cfg,
	}
}

// RegisterRoutes registers all routes related to Collection to the router group
func (h *collectionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	collectionRoutes := rg.Group("/Collection")
	{
		collectionRoutes.POST("/CollectionDetail", h.CollectionDetail)
		collectionRoutes.POST("/CollectionLog", h.CollectionLog)
	}
}

// CollectionDetail
// @Tags         Collection
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                               false  "API key"
// @Param        Api-DeviceOS         header    string                               false  "DeviceOS"
// @Param        Api-Channel          header    string                               false  "Channel"
// @Param        Api-RequestID        header    string                               false  "RequestID"
// @Param        request              body      domain.CollectionDetailRequest       false  "Body Request"
// @Success      200  {object}        domain.CollectionDetailResponse
// @Router       /Api/Collection/CollectionDetail [post]
func (h *collectionHandler) CollectionDetail(c *gin.Context) {
	var req domain.CollectionDetailRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "CollectionDetail"

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

	collectionDetailResult := h.service.CollectionDetail(c, req)
	if collectionDetailResult.AppError != nil {
		handleErrorResponse(c, collectionDetailResult.AppError)
		return
	}

	var responseError *appError.AppError
	if collectionDetailResult.DomainError != nil {
		responseError = collectionDetailResult.DomainError
	}
	if !elkLog.FinalELKLog(collectionDetailResult.GinCtx, &logList, collectionDetailResult.Timestamp, req, collectionDetailResult.Response, collectionDetailResult.DomainError, collectionDetailResult.ServiceName, "", collectionDetailResult.UserRef, []string{collectionDetailResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, collectionDetailResult.Response)
}

// CollectionDetail
// @Tags         Collection
// @Accept       json
// @Produce      json
// @Param        Api-Key              header    string                               false  "API key"
// @Param        Api-DeviceOS         header    string                               false  "DeviceOS"
// @Param        Api-Channel          header    string                               false  "Channel"
// @Param        Api-RequestID        header    string                               false  "RequestID"
// @Param        request              body      domain.CollectionLogRequest          false  "Body Request"
// @Success      200  {object}        domain.CollectionLogResponse
// @Router       /Api/Collection/CollectionLog [post]
func (h *collectionHandler) CollectionLog(c *gin.Context) {
	var req domain.CollectionLogRequest
	timeNow := time.Now()
	var logList []string
	serviceName := "CollectionLog"

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
	
	collectionLogResult := h.service.CollectionLog(c, req)
	if collectionLogResult.AppError != nil {
		handleErrorResponse(c, collectionLogResult.AppError)
		return
	}

	var responseError *appError.AppError
	if collectionLogResult.DomainError != nil {
		responseError = collectionLogResult.DomainError
	}
	if !elkLog.FinalELKLog(collectionLogResult.GinCtx, &logList, collectionLogResult.Timestamp, req, collectionLogResult.Response, collectionLogResult.DomainError, collectionLogResult.ServiceName, "", "", []string{collectionLogResult.LogLine1}, h.logger, h.config.ELKPath, handleErrorResponse) {
		return
	}
	if responseError != nil {
		handleErrorResponse(c, responseError)
		return
	}

	c.JSON(http.StatusOK, collectionLogResult.Response)
}
