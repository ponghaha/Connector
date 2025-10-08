package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"
	"connectorapi-go/internal/adapter/utils"

	"github.com/gin-gonic/gin"
)


// ---------- API CollectionDetail ---------
type CollectionDetailRequest struct {
	IDCardNo            string `json:"IDCardNo"       validate:"required,max=20"`
	RedCaseNo 			string `json:"RedCaseNo"      validate:"max=15"`
	BlackCaseNo 		string `json:"BlackCaseNo"    validate:"max=15"`
}

type CollectionDetailResponse struct {
	IDCardNo            string	                      `json:"IDCardNo"` 
	NoOfAgreement       int 	                      `json:"NoOfAgreement"`
	AgreementList       []CollectionDetailAgreement   `json:"AgreementList"`
}

type CollectionDetailAgreement struct {
	AgreementNo     			string 	`json:"AgreementNo"`
	SeqOfAgreement      		int 	`json:"SeqOfAgreement"`
	OutsourceID 				string 	`json:"OutsourceID"`
	OutsourceName     			string 	`json:"OutsourceName"`
	BlockCode      				string 	`json:"BlockCode"`
	CurrentSUEOSPrincipalNet 	utils.DecimalString `json:"CurrentSUEOSPrincipalNet"`
	CurrentSUEOSPrincipalVAT    utils.DecimalString `json:"CurrentSUEOSPrincipalVAT"`
	CurrentSUEOSInterestNet     utils.DecimalString `json:"CurrentSUEOSInterestNet"`
	CurrentSUEOSInterestVAT 	utils.DecimalString `json:"CurrentSUEOSInterestVAT"`
	CurrentSUEOSPenalty     	utils.DecimalString `json:"CurrentSUEOSPenalty"`
	CurrentSUEOSHDCharge      	utils.DecimalString `json:"CurrentSUEOSHDCharge"`
	CurrentSUEOSOtherFee 		utils.DecimalString `json:"CurrentSUEOSOtherFee"`
	CurrentSUEOSTotal     		utils.DecimalString `json:"CurrentSUEOSTotal"`
	TotalPaymentAmount      	utils.DecimalString `json:"TotalPaymentAmount"`
	LastPaymentDate 			int 	`json:"LastPaymentDate"`
	SUESeqNo     				int 	`json:"SUESeqNo"`
	BeginSUEOSPrincipalNet      utils.DecimalString `json:"BeginSUEOSPrincipalNet"`
	BeginSUEOSPrincipalVAT 		utils.DecimalString `json:"BeginSUEOSPrincipalVAT"`
	BeginSUEOSInterestNet     	utils.DecimalString `json:"BeginSUEOSInterestNet"`
	BeginSUEOSInterestVAT      	utils.DecimalString `json:"BeginSUEOSInterestVAT"`
	BeginSUEOSPenalty 			utils.DecimalString `json:"BeginSUEOSPenalty"`
	BeginSUEOSHDCharge     		utils.DecimalString `json:"BeginSUEOSHDCharge"`
	BeginSUEOSOtherFee     		utils.DecimalString `json:"BeginSUEOSOtherFee"`
	BeginSUEOSTotal 			utils.DecimalString `json:"BeginSUEOSTotal"`
	SUEStatus     				int		`json:"SUEStatus"`
	SUEStatusDescription      	string 	`json:"SUEStatusDescription"`
	BlackCaseNo 				string 	`json:"BlackCaseNo"`
	BlackCaseDate     			int		`json:"BlackCaseDate"`
	RedCaseNo      				string 	`json:"RedCaseNo"`
	RedCaseDate 				int		`json:"RedCaseDate"`
	CourtCode     				string 	`json:"CourtCode"`
	CourtName      				string 	`json:"CourtName"`
	JudgmentDate 				int 	`json:"JudgmentDate"`
	JudgmentResultCode     		int 	`json:"JudgmentResultCode"`
	JudgmentResultDescription   string 	`json:"JudgmentResultDescription"`
	JudgmentDetail 				string 	`json:"JudgmentDetail"`
	ExpectDate     				int 	`json:"ExpectDate"`
	AssetPrice      			utils.DecimalString `json:"AssetPrice"`
	JudgeAmount 				utils.DecimalString `json:"JudgeAmount"`
	NoOfInstallment     		string 	`json:"NoOfInstallment"`
	InstallmentAmount      		utils.DecimalString `json:"InstallmentAmount"`
	TotalCurrentPerSUESeqNo 	utils.DecimalString `json:"TotalCurrentPerSUESeqNo"`
}

type CollectionDetailResult struct {
	Response       *CollectionDetailResponse
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

// ---------- API CollectionLog ---------
type CollectionLogRequest struct {
	AgreementNo         string `json:"AgreementNo"  validate:"required,max=16"`
	RemarkCode          string `json:"RemarkCode"   validate:"required,max=4"`
	LogRemark1 			string `json:"LogRemark1"   validate:"max=120"`
	LogRemark2     		string `json:"LogRemark2"   validate:"max=120"`
	LogRemark3       	string `json:"LogRemark3"   validate:"max=120"`
	LogRemark4 			string `json:"LogRemark4"   validate:"max=120"`
	LogRemark5     		string `json:"LogRemark5"   validate:"max=120"`
	InputDate           string `json:"InputDate"    validate:"required,max=8"`
	InputTime           string `json:"InputTime"    validate:"required,max=6"`
	OperatorID          string `json:"OperatorID"   validate:"required,max=15"`
}

type CollectionLogResponse struct {
	IDCardNo 			string  `json:"IDCardNo"`
	AgreementNo 		string  `json:"AgreementNo"`	
}

type CollectionLogResult struct {
	Response       *CollectionLogResponse
    AppError       *appError.AppError
    GinCtx         *gin.Context
    Timestamp      time.Time
    ReqBody        interface{}
    RespBody       interface{}
    DomainError    *appError.AppError
    ServiceName    string
    LogLine1       string
}
