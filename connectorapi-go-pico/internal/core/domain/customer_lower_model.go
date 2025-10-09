package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"

	"github.com/gin-gonic/gin"
) 

// ---------- API CheckAeonCustomer ---------
type CheckAeonCustomerRequest struct {
	CustomerID                 string `json:"CustomerID"    validate:"required,max=20"`
}

type CheckAeonCustomerResponse struct {
	CustomerID                 string	                      `json:"CustomerID"` 
	AeonMember                 string 	                      `json:"AeonMember"`
	ResultCode                 string	                      `json:"ResultCode"` 
	ReasonCode                 int 	                          `json:"ReasonCode"`
	ReasonDescription          string	                      `json:"ReasonDescription"` 
}

type CheckAeonCustomerResult struct {
	Response       *CheckAeonCustomerResponse
    AppError       *appError.AppError
    GinCtx         *gin.Context
    Timestamp      time.Time
    ReqBody        interface{}
    RespBody       interface{}
    DomainError    *appError.AppError
    ServiceName    string
	UserRef        string
    LogLine1       string
}