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

// ---------- API AgreeMentBilling ---------
type AgreeMentBillingRequest struct {
	IDCardNo      string `json:"IDCardNo"      validate:"required,max=20"`
	AgreementNo   string `json:"AgreementNo"   validate:"required,max=16"`
	CardCode      string `json:"CardCode"      validate:"max=2"`
}

type AgreeMentBillingResponse struct {
	DueDate                        int                 `json:"DueDate"`
	SettlementDate                 int                 `json:"SettlementDate"`
	BillingAmount                  utils.DecimalString `json:"BillingAmount"`
	MinPaymentAmount               utils.DecimalString `json:"MinPaymentAmount"`
	FullPaymentAmount              utils.DecimalString `json:"FullPaymentAmount"`
	UnbilledAmount                 utils.DecimalString `json:"UnbilledAmount"`
	CreditShoppingFloorLimit       utils.DecimalString `json:"CreditShoppingFloorLimit"`
	CreditShoppingOutstanding      utils.DecimalString `json:"CreditShoppingOutstanding"`
	CreditShoppingAvailableLimit   utils.DecimalString `json:"CreditShoppingAvailableLimit"`
	CreditCashingFloorLimit        utils.DecimalString `json:"CreditCashingFloorLimit"`
	CreditCashingOutstanding       utils.DecimalString `json:"CreditCashingOutstanding"`
	CreditCashingAvailableLimit    utils.DecimalString `json:"CreditCashingAvailableLimit"`
	InstallmentNo                  int                 `json:"InstallmentNo"`
	InstallmentCurrent             int                 `json:"InstallmentCurrent"`
	PaymentHistory                 string              `json:"PaymentHistory"`
}

type AgreeMentBillingResult struct {
	Response       *AgreeMentBillingResponse
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