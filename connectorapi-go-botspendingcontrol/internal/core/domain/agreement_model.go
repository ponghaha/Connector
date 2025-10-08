package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"
	"connectorapi-go/internal/adapter/utils"

	"github.com/gin-gonic/gin"
)


// ---------- API UpdateStatus ---------
type UpdateStatusRequest struct {
	AeonID      string `json:"AEONID"      validate:"required,max=20"`
	Agreement   string `json:"Agreement"   validate:"required,max=12"`
	Status      string `json:"Status"      validate:"required,max=1"`
}

type UpdateStatusResponse struct {
	AeonID      string `json:"AEONID"      validate:"max=20"`
	Agreement   string `json:"Agreement"   validate:"max=12"`
}

type UpdateStatusResult struct {
	Response       *UpdateStatusResponse
    AppError       *appError.AppError
    GinCtx         *gin.Context
    Timestamp      time.Time
    ReqBody        interface{}
    RespBody       interface{}
    DomainError    *appError.AppError
    ServiceName    string
    UserToken      string
    LogLine1       string
}