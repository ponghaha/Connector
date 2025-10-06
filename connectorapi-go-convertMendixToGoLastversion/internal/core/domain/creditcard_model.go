package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"

	"github.com/gin-gonic/gin"
) 

// ---------- API GetCardSales ---------
type GetCardSalesRequest struct {
	IDCardNo                 string `json:"IDCardNo"`
	CardType 			     string `json:"CardType"`
	CardBINno 		         string `json:"CardBINno"`
	UsingTypeCPCH            string `json:"UsingTypeCPCH"`
	UsingTypeCA 			 string `json:"UsingTypeCA"`
	SaleDateFrom 		     string `json:"SaleDateFrom"`
	SaleDateTo               string `json:"SaleDateTo"`
	MCCCodeCPCH 			 string `json:"MCCCodeCPCH"`
	AgencyCodeCPCH 		     string `json:"AgencyCodeCPCH"`
	ShopCodeCPCH             string `json:"ShopCodeCPCH"`
}

type GetCardSalesResponse struct {
	TotalSaleCount                      string	                      `json:"TotalSaleCount"` 
	TotalSaleAmount                     string 	                      `json:"TotalSaleAmount"`
	TotalFACount                        string	                      `json:"TotalFACount"` 
	TotalFAAmount                       string 	                      `json:"TotalFAAmount"`
	TotalFRCount                        string	                      `json:"TotalFRCount"` 
	TotalFRAmount                       string 	                      `json:"TotalFRAmount"`
	LastSaleDate                        string	                      `json:"LastSaleDate"` 
	CPSaleCount                         string 	                      `json:"CPSaleCount"`
	CPSaleAmount                        string	                      `json:"CPSaleAmount"` 
	CPLastSaleDate                      string 	                      `json:"CPLastSaleDate"`
	CHSaleCount                         string	                      `json:"CHSaleCount"` 
	CHSaleAmount                        string 	                      `json:"CHSaleAmount"`
	CHLastSaleDate                      string	                      `json:"CHLastSaleDate"` 
	CANormalSaleCount                   string 	                      `json:"CANormalSaleCount"`
	CANormalSaleAmount                  string	                      `json:"CANormalSaleAmount"` 
	CANormalLastSaleDate                string 	                      `json:"CANormalLastSaleDate"`
	CACardlessSaleCount                 string	                      `json:"CACardlessSaleCount"` 
	CACardlessSaleAmount                string 	                      `json:"CACardlessSaleAmount"`
	CACardlessLastSaleDate              string	                      `json:"CACardlessLastSaleDate"` 
	CPSaleReversalCount                 string 	                      `json:"CPSaleReversalCount"`
	CPSaleReversalAmount                string	                      `json:"CPSaleReversalAmount"` 
	CHSaleReversalCount                 string 	                      `json:"CHSaleReversalCount"`
	CHSaleReversalAmount                string	                      `json:"CHSaleReversalAmount"` 
	CANormalSaleReversalCount           string 	                      `json:"CANormalSaleReversalCount"`
	CANormalSaleReversalAmount          string	                      `json:"CANormalSaleReversalAmount"` 
	CACardlessSaleReversalCount         string 	                      `json:"CACardlessSaleReversalCount"`
	CACardlessSaleReversalAmount        string	                      `json:"CACardlessSaleReversalAmount"` 
}

type GetCardSalesResult struct {
	Response       *GetCardSalesResponse
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

// ---------- API GetBigCardInfo ---------
type GetBigCardInfoRequest struct {
	TransactionDate                  string `json:"TransactionDate"     validate:"required,max=8"`
	TransactionTime 			     string `json:"TransactionTime"     validate:"required,max=6"`
	TransactionType 		         string `json:"TransactionType"     validate:"required,max=2"`
	TraceNumber                      string `json:"TraceNumber"         validate:"max=20"`
	AeonID 			                 string `json:"AeonID"              validate:"required,max=44"`
	BusinessCode 		             string `json:"BusinessCode"        validate:"required,max=2"`
	CreditCardNo                     string `json:"CreditCardNo"        validate:"required,max=16"`
	Reserve1 			             string `json:"Reserve1"            validate:"max=20"`
}

type GetBigCardInfoResponse struct {
	TransactionDate                  string	                      `json:"TransactionDate"` 
	TransactionTime                  string 	                  `json:"TransactionTime"`
	TransactionType                  string	                      `json:"TransactionType"` 
	AeonID                           string 	                  `json:"AeonID"`
	BusinessCode                     string	                      `json:"BusinessCode"` 
	CreditCardNo                     string 	                  `json:"CreditCardNo"`
	DataEncrypt                      string	                      `json:"DataEncrypt"` 
}

type GetBigCardInfoResult struct {
	Response       *GetBigCardInfoResponse
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

// ---------- API GetCardDelinquent ---------
type GetCardDelinquentRequest struct {
	IDCardNo                  string `json:"IDCardNo"     validate:"max=20"`
	CardType 			      string `json:"CardType"     validate:"max=2"`
}

type GetCardDelinquentResponse struct {
	DelinquentCountFAFR       string `json:"DelinquentCountFAFR"` 
	DelinquentCountAll        string `json:"DelinquentCountAll"`
}

type GetCardDelinquentResult struct {
	Response       *GetCardDelinquentResponse
    AppError       *appError.AppError
    GinCtx         *gin.Context
    Timestamp      time.Time
    ReqBody        interface{}
    RespBody       interface{}
    DomainError    *appError.AppError
    ServiceName    string
	UserRef         string
    LogLine1       string
}

// ---------- API GetFullpan ---------
type GetFullpanRequest struct {
	IDCardNo                    string                        `json:"IDCardNo"                    validate:"required"`
	Channel 			        string                        `json:"Channel"                     validate:"required"`
	TotalCard 		            int                           `json:"TotalOfApplyCard"            validate:"gt=0"`
    CardList                   []GetFullpanRqOBJ             `json:"CardList_rq"`
}

type GetFullpanRqOBJ struct {
	CardNo                      string 	                      `json:"CardNo"`
	CardCode                    string 	                      `json:"CardCode"`
}

type GetFullpanResponse struct {
	IDCardNo                    string	                      `json:"IDCardNo"` 
	TotalCard                   int 	                      `json:"TotalCard"`
	CardList                    []GetFullpanRsOBJ 	          `json:"Fullpan_CardList_rs"`
}

type GetFullpanRsOBJ struct {
	CardNo                    string 	                      `json:"CardNo"`
	CardCode                  string 	                      `json:"CardCode"`
	CardType                  string 	                      `json:"CardType"`
	ExpireDate                int 	                          `json:"ExpireDate"`
	HoldCode                  string 	                      `json:"HoldCode"`
	SendMode                  string 	                      `json:"SendMode"`
	FirstEmbossDate           int 	                          `json:"FirstEmbossDate"`
	FirstConfirmDate          int 	                          `json:"FirstConfirmDate"`
	DigitalCardFlag           string 	                      `json:"DigitalCardFlag"`
}

type GetFullpanResult struct {
	Response       *GetFullpanResponse
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

// ---------- API GetCardEnroll ---------
type GetCardEnrollRequest struct {
	IDCardNo                  string `json:"IDCardNo"     validate:"max=20"`
	CardType 			      string `json:"CardType"     validate:"max=2"`
}

type GetCardEnrollResponse struct {
	EnrollmentNo       string `json:"EnrollmentNo"` 
	Status             string `json:"Status"`
}

type GetCardEnrollResult struct {
	Response       *GetCardEnrollResponse
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