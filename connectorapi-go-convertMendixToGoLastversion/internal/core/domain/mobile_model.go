package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"
    "connectorapi-go/internal/adapter/utils"

	"github.com/gin-gonic/gin"
)

// ---------- API DashboardSummary ---------
type DashboardSummaryRequest struct {
	IDCardNo string `json:"IDCardNo" validate:"max=20"`
	AeonID   string `json:"AEONID"   validate:"max=20"`
	Channel  string `json:"Channel"  validate:"max=1"`
}

type DashboardSummaryResponse struct {
	IDCardNo                     string                `json:"IDCardNo,omitempty" validate:"max=20"`
	AeonID                       string                `json:"AEONID,omitempty"   validate:"max=20"`
	NameTH                       string                `json:"NameTH"             validate:"max=30"`
	NameEN                       string                `json:"NameEN"             validate:"max=30"`
	MobileNo                     string                `json:"MobileNo"           validate:"max=15"`
	DueDate                      int                   `json:"DueDate"            validate:"lte=99999999"`
    CreditShoppingFloorLimit     utils.DecimalString   `json:"CreditShoppingFloorLimit"`
    CreditShoppingOutstanding    utils.DecimalString   `json:"CreditShoppingOutstanding"`
    CreditShoppingAvailableLimit utils.DecimalString   `json:"CreditShoppingAvailableLimit"`
    CreditCashingFloorLimit      utils.DecimalString   `json:"CreditCashingFloorLimit"`
    CreditCashingOutstanding     utils.DecimalString   `json:"CreditCashingOutstanding"`
    CreditCashingAvailableLimit  utils.DecimalString   `json:"CreditCashingAvailableLimit"`
    YourCashFloorLimit           utils.DecimalString   `json:"YourCashFloorLimit"`
    YourCashOutstanding          utils.DecimalString   `json:"YourCashOutstanding"`
    YourCashAvailableLimit       utils.DecimalString   `json:"YourCashAvailableLimit"`
    ROPShoppingFloorLimit        utils.DecimalString   `json:"ROPShoppingFloorLimit"`
    ROPShoppingOutstanding       utils.DecimalString   `json:"ROPShoppingOutstanding"`
    ROPShoppingAvailableLimit    utils.DecimalString   `json:"ROPShoppingAvailableLimit"`
    ROPCashingFloorLimit         utils.DecimalString   `json:"ROPCashingFloorLimit"`
    ROPCashingOutstanding        utils.DecimalString   `json:"ROPCashingOutstanding"`
    ROPCashingAvailableLimit     utils.DecimalString   `json:"ROPCashingAvailableLimit"`
    TotalMinimumPayment          utils.DecimalString   `json:"TotalMinimumPayment"`
    TotalFullPayment             utils.DecimalString   `json:"TotalFullPayment"`
    TotalPaidAmount              utils.DecimalString   `json:"TotalPaidAmount"`
	PendingPaymentStatus         string                `json:"PendingPaymentStatus" validate:"max=2"`
    RemainMinimumPayment         *utils.DecimalString  `json:"RemainMinimumPayment,omitempty"`
    RemainFullPayment            *utils.DecimalString  `json:"RemainFullPayment,omitempty"`
	BankList                     []DBDetailBankListRq  `json:"BankList_rs"`
	TermsList                    []DBDetailTermsListRq `json:"TermsList_rs"`
}

type DBDetailBankListRq struct {
	CounterNo    string `json:"CounterNo"    validate:"max=4"`
	AccountNo    string `json:"AccountNo"    validate:"max=20"`
	BranchBank   string `json:"BranchBank"   validate:"max=5"`
	RefAccountNo string `json:"RefAccountNo" validate:"max=20"`
}

type DBDetailTermsListRq struct {
	TermsType         string `json:"TermsType"         validate:"max=20"`
	TermsVersion      string `json:"TermsVersion"      validate:"max=10"`
	TermsAcceptStatus string `json:"TermsAcceptStatus" validate:"max=2"`
}

type DashboardSummaryResult struct {
	Response       *DashboardSummaryResponse
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

// ---------- API DashboardDetail ---------
type DashboardDetailRequest struct {
	IDCardNo string `json:"IDCardNo" validate:"max=20"`
	AeonID   string `json:"AEONID"   validate:"max=20"`
	Channel  string `json:"Channel"  validate:"max=1"`
}

type DashboardDetailResponse struct {
	IDCardNo            string           `json:"IDCardNo,omitempty" validate:"max=20"`
	AeonID              string           `json:"AEONID,omitempty"   validate:"max=20"`
    DueDate             int              `json:"DueDate"  validate:"lte=99999999"`
    DashboardDetailList []DBDetailListRs `json:"DashboardDetailList_rs"`
}

type DBDetailListRs struct {
    CreditCardNo                 string               `json:"CreditCardNo"     validate:"max=16"`
    CardName                     string               `json:"CardName"         validate:"max=30"`
    ProductType                  string               `json:"ProductType"      validate:"max=2"`
    CardCode                     string               `json:"CardCode"         validate:"max=2"`
    ATMauthorize                 string               `json:"ATMauthorize"     validate:"max=14"`
    CardStatus                   string               `json:"CardStatus"       validate:"max=16"`
    MinimumPaymentAmount         utils.DecimalString  `json:"MinimumPaymentAmount"`
    FullPaymentAmount            utils.DecimalString  `json:"FullPaymentAmount"`
    PaidAmount                   utils.DecimalString  `json:"PaidAmount"`
    RemainMinimumPayment         *utils.DecimalString `json:"RemainMinimumPayment,omitempty"`
    RemainFullPayment            *utils.DecimalString `json:"RemainFullPayment,omitempty"`
    CreditShoppingFloorLimit     utils.DecimalString  `json:"CreditShoppingFloorLimit"`
    CreditShoppingOutstanding    utils.DecimalString  `json:"CreditShoppingOutstanding"`
    CreditShoppingAvailableLimit utils.DecimalString  `json:"CreditShoppingAvailableLimit"`
    CreditCashingFloorLimit      utils.DecimalString  `json:"CreditCashingFloorLimit"`
    CreditCashingOutstanding     utils.DecimalString  `json:"CreditCashingOutstanding"`
    CreditCashingAvailableLimit  utils.DecimalString  `json:"CreditCashingAvailableLimit"`
    AvailablePoint               utils.DecimalString  `json:"AvailablePoint"`
    BillingAmount                utils.DecimalString  `json:"BillingAmount"`
    UnbilledAmount               utils.DecimalString  `json:"UnbilledAmount"`
    InstallmentNo                int                  `json:"InstallmentNo"      validate:"lte=999"`
    InstallmentCurrent           int                  `json:"InstallmentCurrent" validate:"lte=9"`
    DigitalCardFlag              string               `json:"DigitalCardFlag"    validate:"max=1"`
    ApplicationDate              int                  `json:"ApplicationDate"    validate:"lte=99999999"`
}

type DashboardDetailResult struct {
	Response       *DashboardDetailResponse
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

// ---------- API MobileFullPan ---------
type MobileFullPanRequest struct {
	IDCardNo   string             `json:"IDCardNo"    validate:"required,max=20"`
	Channel    string             `json:"Channel"     validate:"required,max=1"`
	TotalCard  int                `json:"TotalCard"   validate:"lte=9999"`
    CardListRq []MobileCardListRq `json:"CardList_rq" validate:"required,dive"`
}

type MobileCardListRq struct {
	CardNo   string `json:"CardNo"   validate:"required,max=16"`
	CardCode string `json:"CardCode" validate:"required,max=2"`
}

type MobileFullPanFormatRequest struct {
	IDCardNo     string       `json:"IDCardNo"     validate:"max=20"`
	CreditCardNo string       `json:"CreditCardNo" validate:"max=16"`
	BusinessCode string       `json:"TotalCard"    validate:"max=1"`
}

type MobileFullPanResponse struct {
	IDCardNo   string             `json:"IDCardNo"    validate:"max=20"`
	TotalCard  int                `json:"TotalCard"   validate:"lte=9999"`
    CardListRs []MobileCardListRs `json:"CardList_rs"`
}

type MobileCardListRs struct {
	CardNo           string `json:"CardNo"          validate:"max=16"`
	CardCode         string `json:"CardCode"        validate:"max=2"`
    CardType         string `json:"CardType"        validate:"max=2"`
	ExpireDate       int    `json:"ExpireDate"      validate:"max=8"`
	HoldCode         string `json:"HoldCode"        validate:"max=2"`
	SendMode         string `json:"SendMode"`
	FirstEmbossDate  int    `json:"FirstEmbossDate"`
	FirstConfirmDate int    `json:"FirstConfirmDate"`
	DigitalCardFlag  string `json:"DigitalCardFlag" validate:"max=1"`
}

type MobileFullPanResult struct {
	Response       *MobileFullPanResponse
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
