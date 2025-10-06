package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"
	"connectorapi-go/internal/adapter/utils"

	"github.com/gin-gonic/gin"
)

// ---------- API GetCustomerInfo ---------
type GetCustomerInfoRequest struct {
	SNSNo                 string `json:"SNSNo"       validate:"max=20"`
	UserRef 			  string `json:"UserRef"     validate:"max=20"`
	Channel 		      string `json:"Channel"`
	Mode                  string `json:"Mode"        validate:"max=1"`
	AEONID 			      string `json:"AEONID"      validate:"max=20"`
	IDCardNo 		      string `json:"IDCardNo"    validate:"max=20"`
	AgreementNo           string `json:"AgreementNo" validate:"max=16"`
}

type GetCustomerInfoResponse001 struct {
	IDCardNo                      string	                      `json:"IDCardNo"` 
	CustomerNameENG               string 	                      `json:"CustomerNameENG"`
	CustomerNameTH                string	                      `json:"CustomerNameTH"` 
	MobileNo                      string 	                      `json:"MobileNo"`
	
}

type GetCustomerInfoResponse004 struct {
	AEONID                        string	                      `json:"AEONID"` 
	CustomerNameENG               string 	                      `json:"CustomerNameENG"`
	CustomerNameTH                string	                      `json:"CustomerNameTH"`
	Sex                           string	                      `json:"Sex"` 
	MobileNo                      string 	                      `json:"MobileNo"`
	Email                         string 	                      `json:"Email"`
	Nationality                   string	                      `json:"Nationality"`
	Birthdate                     string	                      `json:"Birthdate"` 
	MemberStatus                  string 	                      `json:"MemberStatus"`
}

type GetCustomerInfoResponse003 struct {
	IDCardNo                      string	                      `json:"IDCardNo"` 
	CustomerNameENG               string 	                      `json:"CustomerNameENG"`
	CustomerNameTH                string	                      `json:"CustomerNameTH"`
	MobileNo                      string 	                      `json:"MobileNo"`
	FoundDataFlag                 string 	                      `json:"FoundDataFlag"`
	CustomerGroup                 string	                      `json:"CustomerGroup"`
	NamePreFixEN                  string	                      `json:"NamePreFixEN"` 
	Age                           int 	                          `json:"Age"`
    Birthdate                     int	                          `json:"Birthdate"` 
	Gender                        int 	                          `json:"Gender"`
	MarriageStatus                int	                          `json:"MarriageStatus"`
	EducationCode                 string 	                      `json:"EducationCode"`
	EducationDescription          string 	                      `json:"EducationDescription"`
	HomeStatus                    int	                          `json:"HomeStatus"`
	LivingPeriod                  string	                      `json:"LivingPeriod"` 
	StayWith                      int 	                          `json:"StayWith"`
    HomeAddress                   string	                      `json:"HomeAddress"` 
	HomeZip                       int 	                          `json:"HomeZip"`
	HomePhone                     string	                      `json:"HomePhone"`
	HomePhoneExtension            string 	                      `json:"HomePhoneExtension"`
	OfficeName                    string 	                      `json:"OfficeName"`
	OfficeSection                 string	                      `json:"OfficeSection"`
	OfficeAddress                 string	                      `json:"OfficeAddress"` 
	OfficeZip                     int 	                          `json:"OfficeZip"`
	OfficePhone                   string	                      `json:"OfficePhone"` 
	OfficeExtension               string 	                      `json:"OfficeExtension"`
	BusinessType                  string	                      `json:"BusinessType"`
	BusinessTypeDescription       string 	                      `json:"BusinessTypeDescription"`
	JobTypeCode                   int 	                          `json:"JobTypeCode"`
	JobTypeSubCode                string	                      `json:"JobTypeSubCode"`
	OtherJobDescription           string	                      `json:"OtherJobDescription"` 
	WorkingPeriod                 string 	                      `json:"WorkingPeriod"`
	EmploymentStatus              string	                      `json:"EmploymentStatus"` 
	Salary                        utils.DecimalString 	          `json:"Salary"`
	OtherIncome                   utils.DecimalString	          `json:"OtherIncome"`
	OtherIncomeResource           string 	                      `json:"OtherIncomeResource"`
	OtherIncomeResourceDescription string 	                      `json:"OtherIncomeResourceDescription"`
	SourceOfOtherIncomeCountry    string	                      `json:"SourceOfOtherIncomeCountry"`
	EmailAddress                  string	                      `json:"EmailAddress"` 
	MailTo                        string 	                      `json:"MailTo"`
	TimeToContact                 string	                      `json:"TimeToContact"` 
	SpouseName                    string 	                      `json:"SpouseName"`
	SpousePhone                   string	                      `json:"SpousePhone"`
	SpousePhoneExtension          string 	                      `json:"SpousePhoneExtension"`
	ReferenceName                 string 	                      `json:"ReferenceName"`
	ReferenceRelationship         string	                      `json:"ReferenceRelationship"`
	ReferencePhone                string	                      `json:"ReferencePhone"` 
	ReferenceExtension            string 	                      `json:"ReferenceExtension"`
	HouseRegistrationHome         string 	                      `json:"HouseRegistrationHome"`
	HouseRegistrationHomeZip      int	                          `json:"HouseRegistrationHomeZip"`
	DebtReferenceName             string 	                      `json:"DebtReferenceName"`
	DebtReferenceRelationship     string 	                      `json:"DebtReferenceRelationship"`
	DebtReferencePhone            string	                      `json:"DebtReferencePhone"`
	DebtReferencePhoneExtension   string	                      `json:"DebtReferencePhoneExtension"` 
	DebtReferenceMobilePhone      string 	                      `json:"DebtReferenceMobilePhone"`
	PaymentType                   string 	                      `json:"PaymentType"`
	AutoPayBankName               string	                      `json:"AutoPayBankName"`
	AutoPayAccountNo              string 	                      `json:"AutoPayAccountNo"`
	HomeNo                        string 	                      `json:"HomeNo"`
	HomeVillageBuilding           string	                      `json:"HomeVillageBuilding"`
	HomeRoom                      string	                      `json:"HomeRoom"` 
	HomeFloor                     string 	                      `json:"HomeFloor"`
	HomeMoo                       string 	                      `json:"HomeMoo"`
	HomeSoi                       string	                      `json:"HomeSoi"`
	HomeRoad                      string 	                      `json:"HomeRoad"`
	HomeSubDistrict               string 	                      `json:"HomeSubDistrict"`
	HomeDistrict                  string	                      `json:"HomeDistrict"`
	HomeProvince                  string	                      `json:"HomeProvince"` 
	OfficeNo                      string 	                      `json:"OfficeNo"`
	OfficeVillageBuilding         string 	                      `json:"OfficeVillageBuilding"`
	OfficeRoom                    string	                      `json:"OfficeRoom"`
	OfficeFloor                   string 	                      `json:"OfficeFloor"`
	OfficeMoo                     string 	                      `json:"OfficeMoo"`
	OfficeSoi                     string	                      `json:"OfficeSoi"`
	OfficeRoad                    string	                      `json:"OfficeRoad"` 
	OfficeSubDistrict             string 	                      `json:"OfficeSubDistrict"`
	OfficeDistrict                string 	                      `json:"OfficeDistrict"`
	OfficeProvince                string	                      `json:"OfficeProvince"`
	OfficeMobilePhone             string 	                      `json:"OfficeMobilePhone"`
	HouseRegistrationCode         int 	                          `json:"HouseRegistrationCode"`
}

type GetCustomerInfoResult struct {
	Response       interface{}
    AppError       *appError.AppError
    GinCtx         *gin.Context
    Timestamp      time.Time
    ReqBody        interface{}
    RespBody       interface{}
    DomainError    *appError.AppError
    ServiceName    string
	UserToken	   string
	UserRef        string
    LogLine1       string
}

// ---------- API CheckApplyCondition ---------
type CheckApplyConditionRequest struct {
	ApplicationNo       string `json:"ApplicationNo"       validate:"required"`
	Channel 			string `json:"Channel"`
	IDCardNo 		    string `json:"IDCardNo"            validate:"required,max=20"`
	Birthdate           int    `json:"Birthdate"           validate:"required,gt=0,min=10000000,max=99999999"`
	SuppIDCardNo        string `json:"SuppIDCardNo"`
	SuppBirthdate 		int    `json:"SuppBirthdate"`
	ApplyChannel        string `json:"ApplyChannel"        validate:"required"`
	ApplicationDate     int    `json:"ApplicationDate"     validate:"gt=0",min=10000000,max=99999999"`
	BranchCode 		    string `json:"BranchCode"`
	SourceCode 			string `json:"SourceCode"`
	StaffCode 		    string `json:"StaffCode"`
	TotalApplyCard      int    `json:"TotalApplyCard"      validate:"gt=0"`
    ApplyCardList       []ApplyCardListobj   `json:"CardList_rq"`
}

type ApplyCardListobj struct {
	CardApplyType     	int 	`json:"CardApplyType"`
	CardCode      		string 	`json:"CardCode"`
	PrimaryCreditCard 	string 	`json:"PrimaryCreditCard"`
	VirtualCardFlag     string 	`json:"VirtualCardFlag"`
}

type CheckApplyConditionResponse struct {
	ApplicationNo         string	                      `json:"ApplicationNo"` 
	Status                string 	                      `json:"Status"`
	ReasonCode            string	                      `json:"ReasonCode"` 
	ReasonDescription     string 	                      `json:"ReasonDescription"`
}

type CheckApplyConditionResult struct {
	Response       *CheckApplyConditionResponse
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

// ---------- API CheckApplyCondition2ndCard ---------
type CheckApplyCondition2ndCardRequest struct {
	IDCardNo                    string                        `json:"IDCardNo"                    validate:"required"`
	Channel 			        string                        `json:"Channel"                     validate:"required"`
	TotalOfApplyCard 		    int                           `json:"TotalOfApplyCard"            validate:"gt=0"`
    CheckApply2ndCardList       []CheckApply2ndCardRqOBJ      `json:"CardList"`
}

type CheckApply2ndCardRqOBJ struct {
	CardCode                    string 	                      `json:"CardCode"`
}

type CheckApplyCondition2ndCardResponse struct {
	IDCardNo                    string	                      `json:"IDCardNo"` 
	MaximumCR                   int 	                      `json:"MaximumCR"`
	HaveCardCR                  int	                          `json:"HaveCardCR"` 
	MaximumYC                   int 	                      `json:"MaximumYC"`
	HaveCardYC                  int 	                      `json:"HaveCardYC"`
	TotalOfApplyCard            int 	                      `json:"TotalOfApplyCard"`
	CheckApply2ndCardList       []CheckApply2ndCardRsOBJ 	  `json:"CheckApply2ndCardList"`
}

type CheckApply2ndCardRsOBJ struct {
	CardCode                    string 	                      `json:"CardCode"`
	ResultCode                  string 	                      `json:"ResultCode"`
	ReasonCode                  string 	                      `json:"ReasonCode"`
	ReasonDescription           string 	                      `json:"ReasonDescription"`
}

type CheckApplyCondition2ndCardResult struct {
	Response       *CheckApplyCondition2ndCardResponse
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