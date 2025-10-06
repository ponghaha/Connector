package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"
	"connectorapi-go/internal/adapter/utils"

	"github.com/gin-gonic/gin"
)

// ---------- API GetRedbookInfo ---------
type GetRedbookInfoRequest struct {
	AgentCode                 string `json:"AgentCode"         validate:"max=8"`
	MarketingCode 			  string `json:"MarketingCode"     validate:"max=10"`
	Brand 		              string `json:"Brand"             validate:"required,max=30"`
	Model                     string `json:"Model"             validate:"required,max=30"`
	CarYear 			      int    `json:"CarYear"           validate:"gte=1500,lte=2100"`
	CarMonth 		          int    `json:"CarMonth"          validate:"gte=1,lte=12"`
	SubModel                  string `json:"SubModel"          validate:"required,max=100"`
	EffectiveYear 			  int    `json:"EffectiveYear"     validate:"gte=1500,lte=2100"`
	EffectiveMonth 		      int    `json:"EffectiveMonth"    validate:"gte=1,lte=12"`
}

type GetRedbookInfoResponse struct {
	AgentCode                 string	                      `json:"AgentCode"` 
	MarketingCode             string 	                      `json:"MarketingCode"`
	Brand                     string	                      `json:"Brand"` 
	Model                     string 	                      `json:"Model"`
	CarYear                   int 	                          `json:"CarYear"`
	CarMonth                  int	                          `json:"CarMonth"` 
	SubModel                  string 	                      `json:"SubModel"`
	EffectiveYear             int 	                          `json:"EffectiveYear"`
	EffectiveMonth            int	                          `json:"EffectiveMonth"` 
	VehicleCode               string 	                      `json:"VehicleCode"`
	AvgWholesale              utils.DecimalString 	          `json:"AvgWholesale"`
	AvgRetail                 utils.DecimalString	          `json:"AvgRetail"` 
	GoodWholesale             utils.DecimalString 	          `json:"GoodWholesale"`
	GoodRetail                utils.DecimalString 	          `json:"GoodRetail"`
	NewPrice                  utils.DecimalString	          `json:"NewPrice"`
}

type GetRedbookInfoResult struct {
	Response       *GetRedbookInfoResponse
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

// ---------- API GetDealerCommission ---------
type GetDealerCommissionRequest struct {
	AgentCode                 string `json:"AgentCode"         validate:"max=8"`
	MarketingCode 			  string `json:"MarketingCode"     validate:"max=10"`
	AgreementNo 		      string `json:"AgreementNo"       validate:"required,max=12"`
	CommissionCode            string `json:"CommissionCode"    validate:"required,max=8"`
}

type GetDealerCommissionResponse struct {
	AgentCode                 string	                      `json:"AgentCode"` 
	MarketingCode             string 	                      `json:"MarketingCode"`
	AgreementNo               string	                      `json:"AgreementNo"` 
	CommissionCode            string 	                      `json:"CommissionCode"`
	AgentCategory             string 	                      `json:"AgentCategory"`
	TotalCommission           utils.DecimalString	          `json:"TotalCommission"` 
	VATRate                   utils.DecimalString 	          `json:"VATRate"`
	VAT                       utils.DecimalString 	          `json:"VAT"`
	GrandTotalCommission      utils.DecimalString	          `json:"GrandTotalCommission"` 
	WHTRate                   utils.DecimalString 	          `json:"WHTRate"`
	WHTTax                    utils.DecimalString 	          `json:"WHTTax"`
	NetTotalCommission        utils.DecimalString	          `json:"NetTotalCommission"` 
}

type GetDealerCommissionResult struct {
	Response       *GetDealerCommissionResponse
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

// ---------- API GetDealerAgreement ---------
type GetDealerAgreementRequest struct {
	AgentCode                 string `json:"AgentCode"                 validate:"max=8"`
	MarketingCode 			  string `json:"MarketingCode"             validate:"max=10"`
	TransactionDateFrom       int    `json:"TransactionDateFrom"       validate:"gt=0,min=10000000,max=99999999"`
	TransactionDateTo         int    `json:"TransactionDateTo"         validate:"gt=0,min=10000000,max=99999999"`
	AgreementNo               string `json:"AgreementNo"               validate:"max=12"`
}

type GetDealerAgreementResponse struct {
	AgentCode                 string	                      `json:"AgentCode"` 
	MarketingCode             string 	                      `json:"MarketingCode"`
	TransactionDateFrom       int	                          `json:"TransactionDateFrom"` 
	TransactionDateTo         int 	                          `json:"TransactionDateTo"`
	AgreementNo               string 	                      `json:"AgreementNo"`
	TotalAgreement            int	                          `json:"TotalAgreement"` 
	AgreementList             []AgreementListobj 	          `json:"AgreementList"`
}

type AgreementListobj struct {
	AgreementNo               string	                      `json:"AgreementNo"` 
	TransactionDate           int 	                          `json:"TransactionDate"`
	CustomerName              string	                      `json:"CustomerName"` 
	Status                    string 	                      `json:"Status"`
}

type GetDealerAgreementResult struct {
	Response       *GetDealerAgreementResponse
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