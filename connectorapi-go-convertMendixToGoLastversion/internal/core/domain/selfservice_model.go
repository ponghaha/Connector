package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"

	"github.com/gin-gonic/gin"
)

// ---------- API MyCard ---------
type MyCardRequest struct {
	SNSNo                 string `json:"SNSNo"`
	UserRef 			  string `json:"UserRef"`
	Channel 		      string `json:"Channel"`
	Mode                  string `json:"Mode"`
}

type MyCardResponseNormal struct {
	IDCardNo                      string	                      `json:"IDCardNo"` 
	TotalCreditCard               int 	                     	  `json:"TotalCreditCard"`
	CardList                      []MyCardListNormal	          `json:"CardList"`
}

type MyCardListNormal struct {
	CreditCardNo     	          string 	`json:"CreditCardNo"`
	CardName      		          string 	`json:"CardName"`
	ProductType 	              string 	`json:"ProductType"`
	BusinessCode                  string 	`json:"BusinessCode"`
	CardStatus     	              string 	`json:"CardStatus"`
	ExpireDate      		      string 	`json:"ExpireDate"`
	// DYCA 	                      string 	`json:"DYCA"`
	DigitalCardFlag               string 	`json:"DigitalCardFlag"`
}

type MyCardResponseAll struct {
	IDCardNo                      string	                      `json:"IDCardNo"` 
	CustomerNameEN                string 	                      `json:"CustomerNameEN"`
	CustomerNameTH                string	                      `json:"CustomerNameTH"` 
	TotalCreditCard               int 	                     	  `json:"TotalCreditCard"`
	CardList                      []MyCardListAll	              `json:"CardList"`
}

type MyCardListAll struct {
	CreditCardNo     	          string 	`json:"CreditCardNo"`
	CardCode      		          string 	`json:"CardCode"`
	ProductType 	              string 	`json:"ProductType"`
	CardType                      string 	`json:"CardType"`
	CardStatus     	              string 	`json:"CardStatus"`
	ExpireDate      		      int 	    `json:"ExpireDate"`
	HoldCode 	                  string 	`json:"HoldCode"`
	RetreatCode                   string 	`json:"RetreatCode"`
	SendMode                      string 	`json:"SendMode"`
	FirstEmbossDate     	      int 	    `json:"FirstEmbossDate"`
	FirstConfirmDate      		  int 	    `json:"FirstConfirmDate"`
	ShoppingLimit 	              int 	    `json:"ShoppingLimit"`
	CashingLimit                  int 	    `json:"CashingLimit"`
}

type MyCardResult struct {
	Response       interface{}
    AppError       *appError.AppError
    GinCtx         *gin.Context
    Timestamp      time.Time
    ReqBody        interface{}
    RespBody       interface{}
    DomainError    *appError.AppError
    ServiceName    string
	UserToken      string
	UserRef        string
    LogLine1       string
}

// ---------- API GetAvailableLimit ---------
type GetAvailableLimitRequest struct {
	IDCardNo              string `json:"IDCardNo"       validate:"required,max=20"`
	Channel 			  string `json:"Channel"`
	CreditCardNo 		  string `json:"CreditCardNo"   validate:"required,max=16"`
	BusinessCode          string `json:"BusinessCode"   validate:"required,max=2"`
}

type GetAvailableLimitResponse struct {
	ShoppingLimit         float64 	          `json:"ShoppingLimit"` 
	ShoppingAvailable     float64  	          `json:"ShoppingAvailable"`
	CashingLimit          float64 	          `json:"CashingLimit"`
	CashingAvailable      float64 	          `json:"CashingAvailable"`
}

type GetAvailableLimitResult struct {
	Response       *GetAvailableLimitResponse
    AppError       *appError.AppError
    GinCtx         *gin.Context
    Timestamp      time.Time
    ReqBody        interface{}
    RespBody       interface{}
    DomainError    *appError.AppError
    ServiceName    string
	UserToken      string
	UserRef        string
    LogLine1       string
}