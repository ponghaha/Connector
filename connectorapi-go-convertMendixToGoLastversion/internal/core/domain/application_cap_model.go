package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"
	"connectorapi-go/internal/adapter/utils"

	"github.com/gin-gonic/gin"
)

// ---------- API GetApplicationNo ---------
type GetApplicationNoRequest struct {
	IDCardNo       string          `json:"IDCardNo"       validate:"required,max=20"`
	Channel        string          `json:"Channel"        validate:"required,max=1"`
	ApplyChannel   string          `json:"ApplyChannel"   validate:"required,max=1"`
	TotalApplyCard int             `json:"TotalApplyCard" validate:"required,lte=99"`
    CardListRq     []AppCardListRq `json:"CardList_rq"    validate:"required,dive"`
}

type AppCardListRq struct {
	CardCode        string `json:"CardCode"        validate:"max=2"`
	VirtualCardFlag string `json:"VirtualCardFlag" validate:"max=1"`
}

type GetApplicationNoResponse struct {
	IDCardNo          string `json:"IDCardNo"          validate:"max=20"`
	ApplicationNo     string `json:"ApplicationNo"     validate:"max=20"`
	ApplicationDate   string `json:"ApplicationDate"   validate:"max=8"`
	ApplicationTime   string `json:"ApplicationTime"   validate:"max=6"`
	ResultCode        string `json:"ResultCode"        validate:"max=2"`
	ResultDescription string `json:"ResultDescription" validate:"max=50"`
}

type GetApplicationNoResult struct {
	Response       *GetApplicationNoResponse
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

// ---------- API SubmitCardApplication ---------
type SubmitCardApplicationRequest struct {
	IDCardNo         string             `json:"IDCardNo"          validate:"required,max=20"`
	Channel          string             `json:"Channel"           validate:"required,max=1"`
	ApplicationNo    string             `json:"ApplicationNo"     validate:"required,max=20"`
	ApplyChannel     string             `json:"ApplyChannel"      validate:"required,max=1"`
	ApplicationDate  string             `json:"ApplicationDate"   validate:"required,max=8"`
	BranchCode       string             `json:"BranchCode"        validate:"max=4"`
	SourceCode       string             `json:"SourceCode"        validate:"required,max=8"`
	StaffCode        string             `json:"StaffCode"         validate:"max=7"`
	MailTo           string             `json:"MailTo"            validate:"required,max=1"`
	TotalApplyCard   int                `json:"TotalApplyCard"    validate:"required,lte=99"`
    SubmitCardListRq []SubmitCardListRq `json:"CardList_rq"       validate:"required,dive"`
}

type SubmitCardListRq struct {
	CardCode        string `json:"CardCode"        validate:"max=2"`
	VirtualCardFlag string `json:"VirtualCardFlag" validate:"max=1"`
}

type SubmitCardApplicationResponse struct {
	IDCardNo          string             `json:"IDCardNo"          validate:"max=20"`
	ApplicationNo     string             `json:"ApplicationNo"     validate:"max=20"`
	ApplicationDate   string             `json:"ApplicationDate"   validate:"max=8"`
	ResultDate        string             `json:"ResultDate"        validate:"max=8"`
	ResultTime        string             `json:"ResultTime"        validate:"max=6"`
	ProgramID         string             `json:"ProgramID"         validate:"max=10"`
	ResultCode        string             `json:"ResultCode"        validate:"max=2"`
	ResultDescription string             `json:"ResultDescription" validate:"max=50"`
	TotalApplyCard    int                `json:"TotalApplyCard"    validate:"lte=99"`
	SubmitCardListRs  []SubmitCardListRs `json:"CardList_rs"`
}

type SubmitCardListRs struct {
	MemberTempNo string              `json:"MemberTempNo" validate:"max=16"`
	CardCode     string              `json:"CardCode"     validate:"max=2"`
	ResultCode   string              `json:"ResultCode"   validate:"max=1"`
	ReasonCode   string              `json:"ReasonCode"   validate:"max=2"`
	Remark1      string              `json:"Remark1"      validate:"max=30"`
	Remark2      string              `json:"Remark2"      validate:"max=30"`
	MaximumLimit utils.DecimalString `json:"MaximumLimit"`
	PINNumber    string              `json:"PINNumber"    validate:"max=12"`
}

type SubmitCardApplicationResult struct {
	Response       *SubmitCardApplicationResponse
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
