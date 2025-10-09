package handler

import (
	"fmt"
	"net/http"
	"strings"

	appError "connectorapi-go/pkg/error"
	"connectorapi-go/internal/adapter/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"
)

// --- Helper struct & function ---
type apiHeaders struct {
	APIKey    string
	RequestID string
	Channel   string
	DeviceOS  string
}

func handleErrorResponse(c *gin.Context, appErr *appError.AppError) {
	statusCode := http.StatusInternalServerError
	switch appErr.ErrorCode {
	case appError.ErrService.ErrorCode:
		statusCode = http.StatusBadRequest
	case appError.ErrInternalServer.ErrorCode:
		statusCode = http.StatusInternalServerError
	case appError.ErrUnauthorized.ErrorCode:
		statusCode = http.StatusUnauthorized
	case appError.ErrTimeOut.ErrorCode:
		statusCode = http.StatusGatewayTimeout
	default:
		statusCode = http.StatusBadRequest
	}

	errResponse := appError.ErrorResponse{
		ErrorCode:    appErr.ErrorCode,
		ErrorMessage: appErr.ErrorMessage,
	}

	c.JSON(statusCode, errResponse)
}

func getAPIHeaders(c *gin.Context) apiHeaders {
	return apiHeaders{
		APIKey:    utils.GetHeader(c, "Api-Key", "X-Key", "APIKey"),
		RequestID: utils.GetHeader(c, "Api-RequestID", "RequestID", "X-Request-ID"),
		Channel:   utils.GetHeader(c, "Api-Channel", "Channel", "X-Channel"),
		DeviceOS:  utils.GetHeader(c, "Api-DeviceOS", "DeviceOS", "X-Device-OS"),
	}
}

// func getAPIHeaders(c *gin.Context) apiHeaders {
// 	return apiHeaders{
// 		APIKey:    c.GetHeader("Api-Key"),
// 		RequestID: c.GetHeader("Api-RequestID"),
// 		Channel:   c.GetHeader("Api-Channel"),
// 		DeviceOS:  c.GetHeader("Api-DeviceOS"),
// 	}
// }

func ValidateHeaders(c *gin.Context, method string, path string, apiKeyRepo *utils.APIKeyRepository, logger *zap.SugaredLogger) *appError.AppError {
	headers := getAPIHeaders(c)

	if !apiKeyRepo.Validate(headers.APIKey, method, path) {
		logger.Warnw("Authorization failed", "path", path, "apiKey", headers.APIKey)
		return appError.ErrUnauthorized
	}

	if headers.RequestID == "" || len(headers.RequestID) > 20 {
		return appError.ErrApiRequestID
	}
	if headers.Channel == "" {
		return appError.ErrApiChannel
	}
	if headers.DeviceOS == "" {
		return appError.ErrApiDeviceOS
	}

	return nil
}

func ValidateHeadersForApiKeyAndApiRequestID(c *gin.Context, method string, path string, apiKeyRepo *utils.APIKeyRepository, logger *zap.SugaredLogger) *appError.AppError {
	headers := getAPIHeaders(c)

	if !apiKeyRepo.Validate(headers.APIKey, method, path) {
		logger.Warnw("Authorization failed", "path", path, "apiKey", headers.APIKey)
		return appError.ErrUnauthorized
	}

	if headers.RequestID == "" || len(headers.RequestID) > 20 {
		return appError.ErrApiRequestID
	}
	return nil
}

func formatValidationErrors(err error) []appError.ValidationErrorDetail {
	var validationErrors []appError.ValidationErrorDetail

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range ve {
			validationErrors = append(validationErrors, appError.ValidationErrorDetail{
				Field:   fieldErr.Field(),
				Tag:     fieldErr.Tag(),
				Message: fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fieldErr.Field(), fieldErr.Tag()),
			})
		}
	}
	return validationErrors
}

func HandleValidationError(err error) *appError.AppError {
	validationErrors := formatValidationErrors(err)

	var missingFields []string
	var lengthExceededFields []string
	var invalidValueFields []string

	for _, ve := range validationErrors {
		switch ve.Tag {
		case "required":
			missingFields = append(missingFields, ve.Field)
		case "max", "lte":
			lengthExceededFields = append(lengthExceededFields, ve.Field)
		case "gt", "ne":
			invalidValueFields = append(invalidValueFields, ve.Field)
		}
	}

	if len(missingFields) > 0 {
		return &appError.AppError{
			ErrorCode:    appError.ErrRequiedParam.ErrorCode,
			ErrorMessage: appError.ErrRequiedParam.ErrorMessage + "(" + strings.Join(missingFields, ", ") + ")",
			Err:     fmt.Errorf("required fields: %v", missingFields),
		}
	}

	if len(lengthExceededFields) > 0 {
		return &appError.AppError{
			ErrorCode:    appError.ErrInternalLength.ErrorCode,
			ErrorMessage: appError.ErrInternalLength.ErrorMessage + " (" + strings.Join(lengthExceededFields, ", ") + ")",
			Err:     fmt.Errorf("max length fields: %v", lengthExceededFields),
		}
	}

	if len(invalidValueFields) > 0 {
		return &appError.AppError{
			ErrorCode:    appError.ErrRequiedParam.ErrorCode,
			ErrorMessage: appError.ErrRequiedParam.ErrorMessage + " (" + strings.Join(invalidValueFields, ", ") + ")",
			Err:          fmt.Errorf("invalid value fields: %v", invalidValueFields),
		}
	}

	return &appError.AppError{
		ErrorCode:    appError.ErrService.ErrorCode,
		ErrorMessage: appError.ErrService.ErrorMessage,
		Err:     fmt.Errorf("%v", validationErrors),
	}
}
