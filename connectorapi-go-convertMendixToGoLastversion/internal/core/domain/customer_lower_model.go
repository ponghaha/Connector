package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"

	"github.com/gin-gonic/gin"
) 

// ---------- API GetCustomerInfoMobileNo ---------
type GetCustomerInfoMobileNoRequest struct {
	Mobileno                 string `json:"mobileno"    validate:"required,max=20"`
}

type GetCustomerInfoMobileNoResponse struct {
	Mobileno                       string	                      `json:"mobileno"` 
	Resultcode                     string 	                      `json:"resultcode"`
	Mobileappflag                  string	                      `json:"mobileappflag"` 
	Vvipflag                       string 	                      `json:"vvipflag"`
	Sweetheartflag                 string	                      `json:"sweetheartflag"` 
	Fraudflag                      string 	                      `json:"fraudflag"` 
}

type GetCustomerInfoMobileNoResult struct {
	Response       *GetCustomerInfoMobileNoResponse
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