package error

import (
	"fmt"
)

type AppError struct {
	StatusCode   string
	Code         string
	Message      string
	Err          error
	ErrorCode    string
	ErrorFields  string `json:"ErrorFields,omitempty"`
	ErrorMessage string
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.ErrorMessage, e.Err)
	}
	return e.ErrorMessage
}
func (e *AppError) Unwrap() error {
	return e.Err
}

var (
	ErrService          = &AppError{ErrorCode: "SYS001", ErrorMessage: "System unavailable"}
	ErrUnauthorized     = &AppError{ErrorCode: "SYS002", ErrorMessage: "Unauthorized"}
	ErrTimeOut          = &AppError{ErrorCode: "SYS003", ErrorMessage: "System Time out"}
	ErrMember           = &AppError{ErrorCode: "SYS005", ErrorMessage: "Member Service System Unavailable"}
	ErrSystemI  		= &AppError{ErrorCode: "SYS008", ErrorMessage: "System-I Unavailable"}
	ErrSystemIUnexpect	= &AppError{ErrorCode: "SYS009", ErrorMessage: "System-I Unexpected error occurred"}
	ErrMemberUnexpect	= &AppError{ErrorCode: "SYS012", ErrorMessage: "Member Service System Unexpected Error"}
	ErrInternalServer   = &AppError{ErrorCode: "SYS500", ErrorMessage: "An unexpected internal error occurred"}
	ErrInternalLength   = &AppError{ErrorCode: "SYS500", ErrorMessage: "An unexpected internal error occurred: max length"}

	ErrRequiedParam     = &AppError{ErrorCode: "COM001", ErrorMessage: "Required Parameter"}
	ErrInvChannel       = &AppError{ErrorCode: "COM002", ErrorMessage: "Invalid Channel"}
	ErrApiChannel       = &AppError{ErrorCode: "COM002", ErrorMessage: "Invalid Api-Channel"}
	ErrInvDate          = &AppError{ErrorCode: "COM004", ErrorMessage: "Invalid Date"}
	ErrInvMode          = &AppError{ErrorCode: "COM007", ErrorMessage: "Invalid Mode"}
	ErrAeonID 			= &AppError{ErrorCode: "COM008", ErrorMessage: "Invalid AEON ID."}
	ErrUserRefOrAeonID  = &AppError{ErrorCode: "COM008", ErrorMessage: "Invalid User Reference / Invalid AEON ID."}
	ErrInvDateTime      = &AppError{ErrorCode: "COM009", ErrorMessage: "Invalid Date Time"}
	ErrStatus           = &AppError{ErrorCode: "COM014", ErrorMessage: "Invalid Status"}
	ErrInvTotalOfList  	= &AppError{ErrorCode: "COM016", ErrorMessage: "Invalid Total of List."}
	ErrInvTime  	    = &AppError{ErrorCode: "COM018", ErrorMessage: "Invalid Time"}
	ErrInvServChannel   = &AppError{ErrorCode: "COM025", ErrorMessage: "Invalid Service Channel"}
	ErrDupinSYSI  	    = &AppError{ErrorCode: "COM026", ErrorMessage: "Duplicate in SystemI"}
	ErrApiRequestID     = &AppError{ErrorCode: "COM033", ErrorMessage: "Invalid Api-RequestID"}
	ErrApiDeviceOS      = &AppError{ErrorCode: "COM034", ErrorMessage: "Invalid Api-DeviceOS"}
    ErrConNotPass       = &AppError{ErrorCode: "COM043", ErrorMessage: "Condition not passed"}
	ErrConNotPassCust   = &AppError{ErrorCode: "COM043", ErrorMessage: "Condition not passed(Customer Cannot Register)"}
	ErrNoMatchProduct   = &AppError{ErrorCode: "COM065", ErrorMessage: "No Product Match with The Conditions"}
	ErrInvAgentCode     = &AppError{ErrorCode: "COM065", ErrorMessage: "Invalid Agent Code"}
	ErrInvCode          = &AppError{ErrorCode: "COM065", ErrorMessage: "Invalid Code"}
	ErrIDCardNotFound   = &AppError{ErrorCode: "COM067", ErrorMessage: "ID Card No. Not Found"}

	ErrAgreement        = &AppError{ErrorCode: "AGR001", ErrorMessage: "Invalid Agreement No."}
	ErrAgreementInAct   = &AppError{ErrorCode: "AGR003", ErrorMessage: "Agreement Inactive"}

	ErrSUEInfoNotFound  = &AppError{ErrorCode: "COL001", ErrorMessage: "SUE Information Not Found"}

	ErrAgrNotFound  	= &AppError{ErrorCode: "UHP003", ErrorMessage: "Agreement No. Not Foundd"}

	ErrInvCreditCard  	= &AppError{ErrorCode: "CRC001", ErrorMessage: "Invalid Credit Card"}
	ErrInvBusCode  	    = &AppError{ErrorCode: "CRC002", ErrorMessage: "Invalid Business Code"}
	ErrInvCardCode  	= &AppError{ErrorCode: "CRC003", ErrorMessage: "Invalid Card Code"}
	ErrInvCardNo  	    = &AppError{ErrorCode: "CRC006", ErrorMessage: "Invalid Card No."}
	ErrInvCardStatus  	= &AppError{ErrorCode: "CRC008", ErrorMessage: "Invalid card status"}
	ErrOverCRLimit  	= &AppError{ErrorCode: "CRC009", ErrorMessage: "Over current CR limit"}
	ErrInvCardNotStatus = &AppError{ErrorCode: "CRC010", ErrorMessage: "Invalid Card Not Present Status"}
	ErrInvLimitStatus   = &AppError{ErrorCode: "CRC011", ErrorMessage: "Invalid Limit Status"}
	ErrInvLimitAmount   = &AppError{ErrorCode: "CRC012", ErrorMessage: "Invalid Limit amount per day"}
	ErrBigCardNotFound  = &AppError{ErrorCode: "CRC012", ErrorMessage: "Not found Big Card No."}

	ErrInvIDCardNo  	= &AppError{ErrorCode: "CUS001", ErrorMessage: "Invalid ID Card No."}
	ErrInvMobileNo  	= &AppError{ErrorCode: "CUS002", ErrorMessage: "Invalid Mobile no.'"}
	ErrInvHBDFormat  	= &AppError{ErrorCode: "CUS003", ErrorMessage: "Birthdate invalid format"}
	ErrInvSupHBDFormat  = &AppError{ErrorCode: "CUS004", ErrorMessage: "Supplement Birthdate invalid format"}
	ErrInvMailTo        = &AppError{ErrorCode: "CUS005", ErrorMessage: "Invalid Mail To"}
	ErrInvGender        = &AppError{ErrorCode: "CUS014", ErrorMessage: "Invalid Gender"}

	ErrInvAppNo  	    = &AppError{ErrorCode: "APP001", ErrorMessage: "Invalid Application no."}
	ErrInvAppChannel  	= &AppError{ErrorCode: "APP002", ErrorMessage: "Invalid Apply Channel"}
	ErrInvViCardFlag  	= &AppError{ErrorCode: "APP003", ErrorMessage: "Invalid Virtual Card Flag"}
	ErrInvAppDateFormat = &AppError{ErrorCode: "APP004", ErrorMessage: "Application Date invalid format"}
	ErrInvSourceCode  	= &AppError{ErrorCode: "APP005", ErrorMessage: "Invalid Source Code"}
	ErrInvCardAppType  	= &AppError{ErrorCode: "APP006", ErrorMessage: "Invalid Card Apply Type"}
	ErrInvAppDate  	    = &AppError{ErrorCode: "APP008", ErrorMessage: "Invalid Application Date"}
	ErrDupAppNo  	    = &AppError{ErrorCode: "APP010", ErrorMessage: "Duplication Application No."}
	
	ErrInvBranchCode  	= &AppError{ErrorCode: "BRN002", ErrorMessage: "Invalid Branch Code"}
	ErrInvATMNo  	    = &AppError{ErrorCode: "BRN003", ErrorMessage: "Invalid ATM No."}

	ErrInvOTPType  	    = &AppError{ErrorCode: "SMS001", ErrorMessage: "Invalid OTP Type"}

	ErrInvSNSNo  	    = &AppError{ErrorCode: "SOC001", ErrorMessage: "Invalid SNS no."}
	ErrCardNotAva  	    = &AppError{ErrorCode: "SOC004", ErrorMessage: "Card not available to register"}

	ErrInvConsentFrom   = &AppError{ErrorCode: "CST001", ErrorMessage: "Invalid Consent Form"}
	ErrInvConsentCode   = &AppError{ErrorCode: "CST002", ErrorMessage: "Invalid Consent Code"}
	ErrInvConsentVer    = &AppError{ErrorCode: "CST003", ErrorMessage: "Invalid Consent Version"}
	ErrInvConsentStatus = &AppError{ErrorCode: "CST005", ErrorMessage: "Invalid Consent Status"}
	ErrInvIPAddress     = &AppError{ErrorCode: "CST006", ErrorMessage: "Invalid IP Address"}
	ErrInvActChannel    = &AppError{ErrorCode: "CST007", ErrorMessage: "Invalid Action Channel"}
	ErrInvAppNoCST      = &AppError{ErrorCode: "CST011", ErrorMessage: "Invalid Application No."}
	ErrNotfoundConsent  = &AppError{ErrorCode: "CST013", ErrorMessage: "Not found Consent"}

	ErrDataNotFound     = &AppError{ErrorCode: "MAC061", ErrorMessage: "Data not found"}
	ErrComCodeNotFound  = &AppError{ErrorCode: "MAC062", ErrorMessage: "Commission Code not found"}

	ErrNotAuthor        = &AppError{ErrorCode: "MCM077", ErrorMessage: "Not Authorizied"}

	ErrAgentNotMatch    = &AppError{ErrorCode: "HPS002", ErrorMessage: "Agent code not match"}

	ErrAlready			= &AppError{ErrorCode: "MST004", ErrorMessage: "Already settlement, Cannot use this menu"}
	ErrCheckerNotMatch  = &AppError{ErrorCode: "MST008", ErrorMessage: "Checker is not match"}

	ErrAgreeNotFound    = &AppError{ErrorCode: "MSG113", ErrorMessage: "Agreement not found"}
	ErrAgentNotFound    = &AppError{ErrorCode: "MSG975", ErrorMessage: "Agent Code not found"}
	ErrInvDate          = &AppError{ErrorCode: "MSG902", ErrorMessage: "Invalid Date"}
)

type ErrorResponse struct {
	ErrorCode    string    `json:"ErrorCode"`
	ErrorMessage string    `json:"ErrorMessage"`
}

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}