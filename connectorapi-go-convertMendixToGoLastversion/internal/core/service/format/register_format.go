package format

import (
	"fmt"

	"connectorapi-go/internal/core/domain"
	"connectorapi-go/internal/adapter/utils"
)

// Converts CheckRegisterRequest to a fixed-length string.
func FormatCheckRegisterRequest(checkRegisterReq domain.CheckRegisterRequest) string {
	IDCardNo        := utils.PadOrTruncate(checkRegisterReq.IDCardNo, 20)
	MobileNo        := utils.PadIntWithZero(utils.ConvertStringToInt(checkRegisterReq.MobileNo), 10)
	AgreementNo     := utils.PadIntWithZero(utils.ConvertStringToInt(checkRegisterReq.AgreementNo), 16)
	return IDCardNo + MobileNo + AgreementNo
}

func FormatCheckRegisterResponse(raw string) (domain.CheckRegisterResponse, error) {
	const headerLen = 123
	const dataLen = 133

	if len(raw) <= headerLen {
		return domain.CheckRegisterResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.CheckRegisterResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	iDCardNo                    := parser.ReadString(0,20)
	customerNameTH              := parser.ReadString(20,30)
	customerNameEN              := parser.ReadString(50,30)
	mobileNo                    := parser.ReadString(80,10)
	email                       := parser.ReadString(90,35)
	result                      := parser.ReadString(125,2)
	resultCode                  := parser.ReadString(127,3)
	cRRegisterFlag              := parser.ReadString(130,1)
	dYCRegisterFlag             := parser.ReadString(131,1)
	agreementRegisterFlag       := parser.ReadString(132,1)


	return domain.CheckRegisterResponse{
		IDCardNo:                iDCardNo,
		CustomerNameTH:          customerNameTH,
		CustomerNameEN:          customerNameEN,
		MobileNo:                mobileNo,
		Email:                   email,
		Result:                  result,
		ResultCode:              resultCode,
		CRRegisterFlag:          cRRegisterFlag,
		DYCRegisterFlag:         dYCRegisterFlag,
		AgreementRegisterFlag:   agreementRegisterFlag,
	}, nil
}

// Converts CheckRegisterSocialRequest to a fixed-length string.
func FormatCheckRegisterSocialRequest(checkRegisterSocialReq domain.CheckRegisterSocialRequest) string {
	IDCardNo        := utils.PadOrTruncate(checkRegisterSocialReq.IDCardNo, 20)
	return IDCardNo
}

func FormatCheckRegisterSocialResponse(raw string) (domain.CheckRegisterSocialResponse, error) {
	const headerLen = 123
	const dataLen = 35

	if len(raw) <= headerLen {
		return domain.CheckRegisterSocialResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.CheckRegisterSocialResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	iDCardNo                    := parser.ReadString(0,20)
	mobileNo                    := parser.ReadString(20,15)



	return domain.CheckRegisterSocialResponse{
		IDCardNo:                iDCardNo,
		MobileNo:                mobileNo,
	}, nil
}

// Converts UpdateUserTokenRequest to a fixed-length string.
func FormatUpdateUserTokenRequest(updateUserTokenReq domain.UpdateUserTokenRequest) string {
	IDCardNo            := utils.PadOrTruncate(updateUserTokenReq.IDCardNo, 20)
	UserRef             := utils.PadOrTruncate(updateUserTokenReq.UserRef, 44)
	RegisterDate        := utils.PadOrTruncate(updateUserTokenReq.RegisterDate, 8)
	RegisterTime        := utils.PadOrTruncate(updateUserTokenReq.RegisterTime, 6)
	return IDCardNo + UserRef + RegisterDate + RegisterTime
}

func FormatUpdateUserTokenResponse(raw string) (domain.UpdateUserTokenResponse, error) {
	const headerLen = 123
	const dataLen = 66

	if len(raw) <= headerLen {
		return domain.UpdateUserTokenResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.UpdateUserTokenResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	iDCardNo                    := parser.ReadString(0,20)
	userToken                   := parser.ReadString(20,44)
	result                      := parser.ReadString(64,2)



	return domain.UpdateUserTokenResponse{
		IDCardNo:                iDCardNo,
		UserToken:               userToken,
		Result:                  result,
	}, nil
} 