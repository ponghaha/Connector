package domain

import (
	"time"

	appError "connectorapi-go/pkg/error"
	// "connectorapi-go/internal/adapter/utils"

	"github.com/gin-gonic/gin"
)


// ---------- API SubmitLoanApplication ---------
type SubmitLoanApplicationRequest struct {
	RequestID                 string  `json:"requestid"                     validate:"max=20"`
	ApplicationNo             string  `json:"applicationno"                 validate:"required,max=20"` //**
	KeptBoxNo                 string  `json:"keptboxno"                     validate:"max=15"`
	CustomerGroup             string  `json:"customergroup"                 validate:"required,max=1"` //**
	NonmemberType             string  `json:"nonmembertype"                 validate:"required,max=1"` //**
	ApplicationDate           string  `json:"applicationdate"               validate:"required,max=14"` //**
	NCBToken                  string  `json:"ncbtoken"                      validate:"required,max=25"` //**
	IDCardNo                  string  `json:"idcardno"                      validate:"required,max=20"` //**
	TitleNameEN               string  `json:"titlenameen"                   validate:"required,max=20"` //**
	NameEN                    string  `json:"nameen"                        validate:"required,max=30"` //**
	NameTH                    string  `json:"nameth"                        validate:"required,max=30"` //**
	Age                       int     `json:"age"                           validate:"lte=999"`
	Birthdate                 string  `json:"birthdate"                     validate:"required,max=8"` //**
	Gender                    int     `json:"gender"                        validate:"required,lte=9"` //**
	MarriageStatus            int     `json:"marriagestatus"                validate:"required,lte=9"` //**
	HomeAddressNo             string  `json:"homeaddressno"                 validate:"required,max=20"` //**
	HomeMoo                   string  `json:"homemoo"                       validate:"max=3"`
	HomeVillage               string  `json:"homevillage"                   validate:"max=45"`
	HomeRoom                  string  `json:"homeroom"                      validate:"max=10"`
	HomeFloor                 string  `json:"homefloor"                     validate:"max=10"`
	HomeSoi                   string  `json:"homesoi"                       validate:"max=25"`
	HomeRoad                  string  `json:"homeroad"                      validate:"max=25"`
	HomeSubDistrict           string  `json:"homesubdistrict"               validate:"required,max=25"` //**
	HomeDistrict              string  `json:"homedistrict"                  validate:"required,max=25"` //**
	HomeProvince              string  `json:"homeprovince"                  validate:"required,max=20"` //**
	HomeZipCode               string  `json:"homezipcode"                   validate:"required,max=5"` //**
	HomePhone                 string  `json:"homephone"                     validate:"required,max=10"` //**
	HomePhoneExt              string  `json:"homephoneext"                  validate:"max=5"`
	MobileNo                  string  `json:"mobileno"                      validate:"required,max=15"` //**
	HomeStatus                string  `json:"homestatus"                    validate:"required,max=1"` //**
	Email                     string  `json:"email"                         validate:"required,max=35"` //**
	LivingPeriod              float64 `json:"livingperiod"                  validate:"required,lte=9999"` //**
	StayWith                  string  `json:"staywith"                      validate:"required,max=3"` //**
	EducationCode             string  `json:"educationcode"                 validate:"required,max=2"` //**
	EducationDescription      string  `json:"educationdescription"          validate:"max=20"`
	OfficeName                string  `json:"officename"                    validate:"required,max=50"` //**
	OfficeSection             string  `json:"officesection"                 validate:"required,max=30"` //**
	OfficeNo                  string  `json:"officeno"                      validate:"required,max=20"` //**
	OfficeMoo                 string  `json:"officemoo"                     validate:"max=2"`
	OfficeBuildingName        string  `json:"officebuildingname"            validate:"max=45"`
	OfficeRoom                string  `json:"officeroom"                    validate:"max=10"`
	OfficeFloor               string  `json:"officefloor"                   validate:"max=10"`
	OfficeSoi                 string  `json:"officesoi"                     validate:"max=25"`
	OfficeRoad                string  `json:"officeroad"                    validate:"max=25"`
	OfficeSubDistrict         string  `json:"officesubdistrict"             validate:"required,max=25"` //**
	OfficeDistrict            string  `json:"officedistrict"                validate:"required,max=25"` //**
	OfficeProvince            string  `json:"officeprovince"                validate:"required,max=20"` //**
	OfficeZipCode             string  `json:"officezipcode"                 validate:"required,max=5"` //**
	OfficePhone               string  `json:"officephone"                   validate:"required,max=10"` //**
	OfficePhoneExt            string  `json:"officephoneext"                validate:"max=5"`
	JobTypeCode               string  `json:"jobtypecode"                   validate:"required,max=2"` //**
	JobTypeDetailCode         string  `json:"jobtypedetailcode"             validate:"max=2"`
	WorkingPeriod             float64 `json:"workingperiod"                 validate:"required,lte=9999"` //**
	EmploymentStatus          string  `json:"employmentstatus"              validate:"required,max=2"` //**
	BusinessType              string  `json:"businesstype"                  validate:"required,max=2"` //**
	BusinessDescription       string  `json:"businessdescription"           validate:"max=20"`
	TimeToContact             string  `json:"timetocontact"                 validate:"required,max=30"` //**
	SpouseName                string  `json:"spousename"                    validate:"max=30"`
	SpousePhone               string  `json:"spousephone"                   validate:"max=10"`
	SpousePhoneExt            string  `json:"spousephoneext"                validate:"max=5"`
	DebtReferenceName         string  `json:"debtreferencename"             validate:"max=30"`
	DebtReferenceRelationship string  `json:"debtreferencerelationship"     validate:"max=15"`
	DebtReferencePhone        string  `json:"debtreferencephone"            validate:"max=11"`
	DebtReferencePhoneExt     string  `json:"debtreferencephoneext"         validate:"max=5"`
	DebtReferenceMobile       string  `json:"debtreferencemobile"           validate:"max=11"`
	Salary                    float64 `json:"salary"                        validate:"required,lte=99999999999"` //**
	OtherIncome               float64 `json:"otherincome                    validate:"lte=99999999999"`
	OtherIncomResource        string  `json:"otherincomresource"            validate:"max=1"`
	OtherincomResourceDesc    string  `json:"otherincomresourcedescription" validate:"max=20"`
	PaymentType               string  `json:"paymenttype"                   validate:"max=2"`
	AutopayBank               string  `json:"autopaybank"                   validate:"max=30"`
	AutopayAccountNo          string  `json:"autopayaccountno"              validate:"max=10"`
	MailTo                    string  `json:"mailto"                        validate:"required,max=1"` //**
	IDCardIssueDate           string  `json:"idcardissuedate"               validate:"required,max=8"` //**
	IDCardExpiryDate          string  `json:"idcardexpirydate"              validate:"required,max=8"` //**
	IDCardHouseNo             string  `json:"idcardhouseno"                 validate:"required,max=20"` //**
	IDCardMoo                 string  `json:"idcardmoo"                     validate:"max=2"`
	IDCardRoom                string  `json:"idcardroom"                    validate:"max=10"`
	IDCardFloor               string  `json:"idcardfloor"                   validate:"max=10"`
	IDCardSoi                 string  `json:"idcardsoi"                     validate:"max=25"`
	IDCardRoad                string  `json:"idcardroad"                    validate:"max=25"`
	IDCardSubDistrict         string  `json:"idcardsubdistrict"             validate:"required,max=25"` //**
	IDCardDistrict            string  `json:"idcarddistrict"                validate:"required,max=25"` //**
	IDCardProvince            string  `json:"idcardprovince"                validate:"required,max=20"` //**
	AgentCode                 string  `json:"agentcode"                     validate:"required,max=8"` //**
	ApplicationPurposeCode    string  `json:"applicationpurposecode"        validate:"required,max=3"` //**
	OtherPurposeDescription   string  `json:"otherpurposedescription"       validate:"max=100"`
	ApplyType                 string  `json:"applytype"                     validate:"required,max=3"` //**
	ProductCode               string  `json:"productcode"                   validate:"required,max=4"` //**
	BrandCode                 string  `json:"brandcode"                     validate:"required,max=4"` //**
	ModelCode                 string  `json:"modelcode"                     validate:"required,max=10"` //**
	Color                     string  `json:"color"                         validate:"required,max=20"` //**
	Cc                        int     `json:"cc"                            validate:"lte=99999"`
	PowerTransitionType       string  `json:"powertransitiontype"           validate:"max=1"`
	CarusedType               string  `json:"carusedtype"                   validate:"max=1"`
	CarmileNo                 int     `json:"carmileno"                     validate:"lte=9999999999"`
	EngineNo                  string  `json:"engineno"                      validate:"max=30"`
	ChassisNo                 string  `json:"chassisno"                     validate:"max=20"`
	CarRegistrationDate       string  `json:"carregistrationdate"           validate:"max=8"`
	LicensePlateNo            string  `json:"licenseplateno"                validate:"max=30"`
	LicensePlateProvince      string  `json:"licenseplateprovince"          validate:"required,max=3"` //**
	LoanContractFlag          string  `json:"loancontractflag"              validate:"required,max=1"` //**
	EstimationPrice           float64 `json:"estimationprice"               validate:"lte=99999999999"`      
	CashPrice                 float64 `json:"cashprice"                     validate:"required,lte=99999999999"` //**              
	PromotionCode             string  `json:"promotioncode"                 validate:"max=10"`
	InterestRate              float64 `json:"interestrate"                  validate:"required,lte=99999"` //**
	DownPayment               string  `json:"downpayment"                   validate:"required,max=1"` //**
	FinancePrice              float64 `json:"financeprice"                  validate:"required,lte=99999999999"` //**
	DownPaymentPrice          float64 `json:"downpaymentprice"              validate:"lte=9999999999"`
	InstallmentPeriod         string  `json:"installmentperiod"             validate:"required,max=3"` //**
	CarType                   string  `json:"cartype"                       validate:"max=2"`
	Note                      string  `json:"note"                          validate:"max=50"`
	ScanedBy                  string  `json:"scanedby"                      validate:"max=30"`
	MarketingCode             string  `json:"marketingcode"                 validate:"required,max=30"` //**
	ManufactureYear           string  `json:"manufactureyear"               validate:"max=4"`
	ApplicationReceivedDate   string  `json:"applicationreceiveddate"       validate:"required,max=14"` //**
	BankCode                  string  `json:"bankcode"                      validate:"max=3"`
	AccountNo                 string  `json:"accountno"                     validate:"max=20"`
	AccountHolderName         string  `json:"accountholdername"             validate:"max=30"`
}

type SubmitLoanApplicationResult struct {
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