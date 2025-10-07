package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"
	// "connectorapi-go/internal/adapter/utils"

	"github.com/gin-gonic/gin"
)

// ---------- API UpdateConsent ---------
type UpdateConsentRequest struct {
	IDCardNo                    string `json:"IDCardNo"              validate:"required,max=20"`
	Channel 			        string `json:"Channel"               validate:"required,max=1"`
	ActionChannel 		        string `json:"ActionChannel"         validate:"required,max=3"`
	ActionDateTime              string `json:"ActionDateTime"        validate:"required,max=14"`
	ApplicationNo 			    string `json:"ApplicationNo"         validate:"max=20"`
	ApplicationVersion 		    string `json:"ApplicationVersion"    validate:"max=13"`
	IPAddress                   string `json:"IPAddress"             validate:"required,max=50"`
	ATMNo                       string `json:"ATMNo"                 validate:"max=5"`
	BranchCode 			        string `json:"BranchCode"            validate:"max=4"`
	VoicePath 		            string `json:"VoicePath"             validate:"max=150"`
	TotalOfConsentCode          int    `json:"TotalOfConsentCode"    validate:"required,gt=0,max=2"`
	ConsentLists 			    []ConsentListsobj   `json:"ConsentLists"`
}

type ConsentListsobj struct {
	ConsentForm     	        string 	`json:"ConsentForm"`
	ConsentCode      		    string 	`json:"ConsentCode"`
	ConsentFormVersion 	        string 	`json:"ConsentFormVersion"`
	ConsentLanguage             string 	`json:"ConsentLanguage"`
	ConsentStatus               string 	`json:"ConsentStatus"`
}

type UpdateConsentResponse struct {
	Status                      string	`json:"Status"` 
}

type UpdateConsentResult struct {
	Response       *UpdateConsentResponse
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

// ---------- API GetConsentList ---------
type GetConsentListRequest struct {
	IDCardNo                    string `json:"IDCardNo"`
	ConsentCode 			    string `json:"ConsentCode"`
	Filler 		                string `json:"Filler"`
	Mode						string `json:"Mode"`
}

type GetConsentListResponse struct {
	IDCardNo         string                 `json:"IDCardNo"`
	Status           string                 `json:"Status,omitempty"`
	Filler           string                 `json:"Filler,omitempty"`
	NumberOfConsent  int                    `json:"NumberOfConsent"`
	ConsentList      []ConsentListObject    `json:"ConsentList_Rs"`
}

type ConsentListObject struct {
	ConsentForm        string                `json:"ConsentForm"`
	ConsentFormVersion string                `json:"ConsentFormVersion"`
	ConsentName        string                `json:"ConsentName"`
	TotalOfConsentCode int                   `json:"TotalOfConsentCode"`
	URL                string                `json:"URL"`
	ConsentDetails     []ConsentDetailObject `json:"ConsentDetail_Rs"`
}

type ConsentDetailObject struct {
	ConsentCode        string `json:"ConsentCode"`
	ApplicationNo      string `json:"ApplicationNo"`
	ApplicationVersion string `json:"ApplicationVersion"`
	LastConsentStatus  string `json:"LastConsentStatus"`
	AcceptDateTime     string `json:"AcceptDateTime"`
	AcceptChannel      string `json:"AcceptChannel"`
	CancelDateTime     string `json:"CancelDateTime"`
	CancelChannel      string `json:"CancelChannel"`
}

type GetConsentListResult struct {
	Response       *GetConsentListResponse
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