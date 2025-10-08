package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"

	"github.com/gin-gonic/gin"
) 

// ---------- API GetSpendingControl ---------
type GetSpendingControlRequest struct {
	TransactionType      string `json:"TransactionType"      validate:"required,max=2"`
	AEONID 			     string `json:"AEONID"               validate:"required,max=20"`
	CardNo 		         string `json:"CardNo"               validate:"required,max=16"`
	BusinessCode         string `json:"BusinessCode"         validate:"required,max=2"`
	Channel 			 string `json:"Channel"              validate:"required,max=10"`
}

type GetSpendingControlResponse struct {
	TransactionType                          string	                      `json:"TransactionType"` 
	CardNotPresentStatus                     string 	                  `json:"CardNotPresentStatus"`
	CNPLimitAmountPerDay                     int	                      `json:"CNPLimitAmountPerDay"` 
	CNPLimitAmountPerTransaction             int 	                      `json:"CNPLimitAmountPerTransaction"`
	LimitStatus                              string	                      `json:"LimitStatus"` 
	LimitAmountPerDay                        int 	                      `json:"LimitAmountPerDay"`
	LimitAmountPerTransaction                int	                      `json:"LimitAmountPerTransaction"` 
}

type GetSpendingControlResult struct {
	Response       *GetSpendingControlResponse
    AppError       *appError.AppError
    GinCtx         *gin.Context
    Timestamp      time.Time
    ReqBody        interface{}
    RespBody       interface{}
    DomainError    *appError.AppError
    ServiceName    string
	UserToken        string
    LogLine1       string
}

// ---------- API UpdateSpendingControl ---------
type UpdateSpendingControlRequest struct {
	TransactionType                  string `json:"TransactionDate"     validate:"required,max=2"`
	AEONID 			                 string `json:"TransactionTime"     validate:"required,max=20"`
	CardNo 		                     string `json:"TransactionType"     validate:"required,max=16"`
	BusinessCode                     string `json:"TraceNumber"         validate:"required,max=2"`
	Channel 			             string `json:"AeonID"              validate:"required,max=10"`
	CardNotPresentStatus 		     string `json:"BusinessCode"        validate:"required,max=1"`
	CNPLimitAmountPerDay             int    `json:"CreditCardNo"        validate:"required"`
	CNPLimitAmountPerTransaction     int    `json:"Reserve1"            validate:"required"`
	LimitStatus                      string `json:"TraceNumber"         validate:"required,max=1"`
	LimitAmountPerDay 			     int    `json:"AeonID"              validate:"required"`
	LimitAmountPerTransaction 		 int    `json:"BusinessCode"        validate:"required"`
	Date                             string `json:"CreditCardNo"        validate:"required,max=8"`
	Time 			                 string `json:"Reserve1"            validate:"max=6"`
}

type UpdateSpendingControlResponse struct {
	TransactionType                  string	                      `json:"TransactionType"` 
	Date                             string 	                  `json:"Date"`
	Time                             string	                      `json:"Time"` 
}

type UpdateSpendingControlResult struct {
	Response       *UpdateSpendingControlResponse
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