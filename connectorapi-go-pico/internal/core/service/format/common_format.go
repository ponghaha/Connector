package format

import (
	"fmt"
	"strconv"
	"strings"
	//"bytes"

	"connectorapi-go/internal/core/domain"
	"connectorapi-go/internal/adapter/utils"
)

// Converts GetCustomerInfoRequest to a fixed-length string.
func FormatGetCustomerInfoRequest001And003(getCustomerInfoReq domain.GetCustomerInfoRequest,Language string) string {
	IDCardNo        := utils.PadOrTruncate(getCustomerInfoReq.UserRef, 20)
	Lang            := utils.PadOrTruncate(Language, 1)
	return IDCardNo + Lang
}

func FormatGetCustomerInfoRequest004(getCustomerInfoReq domain.GetCustomerInfoRequest) string {
	IDCardNo        := utils.PadOrTruncate(getCustomerInfoReq.IDCardNo, 20)
	AEONID        := utils.PadOrTruncate(getCustomerInfoReq.AEONID, 20)
	AgreementNo        := utils.PadOrTruncate(getCustomerInfoReq.AgreementNo, 16)
	return IDCardNo + AEONID + AgreementNo
}

func FormatGetCustomerInfoResponse001(raw string) (domain.GetCustomerInfoResponse001, error) {
	const headerLen = 123
	const dataLen = 132

	if len(raw) <= headerLen {
		return domain.GetCustomerInfoResponse001{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetCustomerInfoResponse001{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	iDCardNo                     := parser.ReadString(0,20)
	customerNameENG              := parser.ReadString(20,30)
	customerNameTH               := parser.ReadString(50,30)
	// sex                          := parser.ReadString(80,1)
	mobileNo                     := parser.ReadString(81,15)
	// email                        := parser.ReadString(96,35)
	// nationality                  := parser.ReadString(131,1)


	return domain.GetCustomerInfoResponse001{
		IDCardNo:                iDCardNo,
		CustomerNameENG:         customerNameENG,
		CustomerNameTH:          customerNameTH,
		MobileNo:                mobileNo,
	}, nil
}

func FormatGetCustomerInfoResponse004(raw string) (domain.GetCustomerInfoResponse004, error) {
	const headerLen = 123
	const dataLen = 142

	if len(raw) <= headerLen {
		return domain.GetCustomerInfoResponse004{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetCustomerInfoResponse004{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	aEONID                       := parser.ReadString(0,20)
	customerNameENG              := parser.ReadString(20,30)
	customerNameTH               := parser.ReadString(50,30)
	sex                          := parser.ReadString(80,1)
	mobileNo                     := parser.ReadString(81,15)
	email                        := parser.ReadString(96,35)
	nationality                  := parser.ReadString(131,1)
	birthdate                    := parser.ReadString(132,8)
	memberStatus                 := parser.ReadString(140,2)

	return domain.GetCustomerInfoResponse004{
		AEONID:                  aEONID,
		CustomerNameENG:         customerNameENG,
		CustomerNameTH:          customerNameTH,
		Sex:                     sex,
		MobileNo:                mobileNo,
		Email:                   email,
		Nationality:             nationality,
		Birthdate:               birthdate,
		MemberStatus:            memberStatus,
	}, nil
}

func FormatGetCustomerInfoResponse003(raw string) (domain.GetCustomerInfoResponse003, error) {
	const headerLen = 123
	const dataLen = 1481

	if len(raw) <= headerLen {
		return domain.GetCustomerInfoResponse003{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetCustomerInfoResponse003{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	idCardNo                         := parser.ReadString(0, 20)
	foundDataFlag                    := parser.ReadString(20, 1)
	customerGroup                    := parser.ReadString(21, 1)
	namePreFixEN                     := parser.ReadString(22, 20)
	customerNameENG                  := parser.ReadString(42, 30)
	customerNameTH                   := parser.ReadString(72, 30)
	age                              := parser.ReadInt(102, 3)
	birthdate                        := parser.ReadInt(105, 8)
	gender                           := parser.ReadInt(113, 1)
	marriageStatus                   := parser.ReadInt(114, 1)
	educationCode                    := parser.ReadString(115, 2)
	educationDescription             := parser.ReadString(117, 50)
	homeStatus                       := parser.ReadInt(167, 1)
	livingPeriod                     := parser.ReadString(168, 4)
	stayWith                         := parser.ReadInt(172, 3)
	homeAddress                      := parser.ReadString(175, 100)
	homeZip                          := parser.ReadInt(275, 5)
	homePhone                        := parser.ReadString(280, 10)
	homePhoneExtension               := parser.ReadString(290, 5)
	officeName                       := parser.ReadString(295, 50)
	officeSection                    := parser.ReadString(345, 30)
	officeAddress                    := parser.ReadString(375, 100)
	officeZip                        := parser.ReadInt(475, 5)
	officePhone                      := parser.ReadString(480, 10)
	officeExtension                  := parser.ReadString(490, 5)
	businessType                     := parser.ReadString(495, 2)
	businessTypeDescription          := parser.ReadString(497, 50)
	jobTypeCode                      := parser.ReadInt(547, 2)
	jobTypeSubCode                   := parser.ReadString(549, 2)
	otherJobDescription              := parser.ReadString(551, 50)
	workingPeriod                    := parser.ReadString(601, 4)
	employmentStatus                 := parser.ReadString(605, 2)
	salary                           := parser.ReadFloat100(607, 11)
	otherIncome                      := parser.ReadFloat100(618, 11)
	otherIncomeResource              := parser.ReadString(629, 1)
	otherIncomeResourceDescription   := parser.ReadString(630, 20)
	sourceOfOtherIncomeCountry       := parser.ReadString(650, 3)
	mobileNo                         := parser.ReadString(653, 15)
	emailAddress                     := parser.ReadString(668, 35)
	mailTo                           := parser.ReadString(703, 1)
	timeToContact                    := parser.ReadString(704, 30)
	spouseName                       := parser.ReadString(734, 30)
	spousePhone                      := parser.ReadString(764, 10)
	spousePhoneExtension             := parser.ReadString(774, 4)
	referenceName                    := parser.ReadString(778, 30)
	referenceRelationship            := parser.ReadString(808, 15)
	referencePhone                   := parser.ReadString(823, 10)
	referenceExtension               := parser.ReadString(833, 4)
	houseRegistrationHome            := parser.ReadString(837, 100)
	houseRegistrationHomeZip         := parser.ReadInt(937, 5)
	debtReferenceName                := parser.ReadString(942, 30)
	debtReferenceRelationship        := parser.ReadString(972, 15)
	debtReferencePhone               := parser.ReadInt(987, 11)
	debtReferencePhoneExtension      := parser.ReadString(998, 5)
	debtReferenceMobilePhone         := parser.ReadInt(1003, 11)
	paymentType                      := parser.ReadString(1014, 2)
	autoPayBankName                  := parser.ReadString(1016, 30)
	autoPayAccountNo                 := parser.ReadString(1046, 10)
	homeNo                           := parser.ReadString(1056, 20)
	homeVillageBuilding              := parser.ReadString(1076, 45)
	homeRoom                         := parser.ReadString(1121, 10)
	homeFloor                        := parser.ReadString(1131, 10)
	homeMoo                          := parser.ReadString(1141, 2)
	homeSoi                          := parser.ReadString(1143, 25)
	homeRoad                         := parser.ReadString(1168, 25)
	homeSubDistrict                  := parser.ReadString(1193, 25)
	homeDistrict                     := parser.ReadString(1218, 25)
	homeProvince                     := parser.ReadString(1243, 20)
	officeNo                         := parser.ReadString(1263, 20)
	officeVillageBuilding            := parser.ReadString(1283, 45)
	officeRoom                       := parser.ReadString(1328, 10)
	officeFloor                      := parser.ReadString(1338, 10)
	officeMoo                        := parser.ReadString(1348, 2)
	officeSoi                        := parser.ReadString(1350, 25)
	officeRoad                       := parser.ReadString(1375, 25)
	officeSubDistrict                := parser.ReadString(1400, 25)
	officeDistrict                   := parser.ReadString(1425, 25)
	officeProvince                   := parser.ReadString(1450, 20)
	officeMobilePhone                := parser.ReadString(1470, 10)
	houseRegistrationCode            := parser.ReadInt(1480, 1)

	return domain.GetCustomerInfoResponse003{
		IDCardNo:                         idCardNo,
		FoundDataFlag:                    foundDataFlag,
		CustomerGroup:                    customerGroup,
		NamePreFixEN:                     namePreFixEN,
		CustomerNameENG:                  customerNameENG,
		CustomerNameTH:                   customerNameTH,
		Age:                              age,
		Birthdate:                        birthdate,
		Gender:                           gender,
		MarriageStatus:                   marriageStatus,
		EducationCode:                    educationCode,
		EducationDescription:             educationDescription,
		HomeStatus:                       homeStatus,
		LivingPeriod:                     livingPeriod,
		StayWith:                         stayWith,
		HomeAddress:                      homeAddress,
		HomeZip:                          homeZip,
		HomePhone:                        homePhone,
		HomePhoneExtension:               homePhoneExtension,
		OfficeName:                       officeName,
		OfficeSection:                    officeSection,
		OfficeAddress:                    officeAddress,
		OfficeZip:                        officeZip,
		OfficePhone:                      officePhone,
		OfficeExtension:                  officeExtension,
		BusinessType:                     businessType,
		BusinessTypeDescription:          businessTypeDescription,
		JobTypeCode:                      jobTypeCode,
		JobTypeSubCode:                   jobTypeSubCode,
		OtherJobDescription:              otherJobDescription,
		WorkingPeriod:                    workingPeriod,
		EmploymentStatus:                 employmentStatus,
		Salary:                           utils.DecimalString(salary),
		OtherIncome:                      utils.DecimalString(otherIncome),
		OtherIncomeResource:              otherIncomeResource,
		OtherIncomeResourceDescription:   otherIncomeResourceDescription,
		SourceOfOtherIncomeCountry:       sourceOfOtherIncomeCountry,
		MobileNo:                         mobileNo,
		EmailAddress:                     emailAddress,
		MailTo:                           mailTo,
		TimeToContact:                    timeToContact,
		SpouseName:                       spouseName,
		SpousePhone:                      spousePhone,
		SpousePhoneExtension:             spousePhoneExtension,
		ReferenceName:                    referenceName,
		ReferenceRelationship:            referenceRelationship,
		ReferencePhone:                   referencePhone,
		ReferenceExtension:               referenceExtension,
		HouseRegistrationHome:            houseRegistrationHome,
		HouseRegistrationHomeZip:         houseRegistrationHomeZip,
		DebtReferenceName:                debtReferenceName,
		DebtReferenceRelationship:        debtReferenceRelationship,
		DebtReferencePhone:               strconv.Itoa(debtReferencePhone),
		DebtReferencePhoneExtension:      debtReferencePhoneExtension,
		DebtReferenceMobilePhone:         strconv.Itoa(debtReferenceMobilePhone),
		PaymentType:                      paymentType,
		AutoPayBankName:                  autoPayBankName,
		AutoPayAccountNo:                 autoPayAccountNo,
		HomeNo:                           homeNo,
		HomeVillageBuilding:              homeVillageBuilding,
		HomeRoom:                         homeRoom,
		HomeFloor:                        homeFloor,
		HomeMoo:                          homeMoo,
		HomeSoi:                          homeSoi,
		HomeRoad:                         homeRoad,
		HomeSubDistrict:                  homeSubDistrict,
		HomeDistrict:                     homeDistrict,
		HomeProvince:                     homeProvince,
		OfficeNo:                         officeNo,
		OfficeVillageBuilding:            officeVillageBuilding,
		OfficeRoom:                       officeRoom,
		OfficeFloor:                      officeFloor,
		OfficeMoo:                        officeMoo,
		OfficeSoi:                        officeSoi,
		OfficeRoad:                       officeRoad,
		OfficeSubDistrict:                officeSubDistrict,
		OfficeDistrict:                   officeDistrict,
		OfficeProvince:                   officeProvince,
		OfficeMobilePhone:                officeMobilePhone,
		HouseRegistrationCode:            houseRegistrationCode,
	}, nil
}

// Converts CheckApplyConditionRequest to a fixed-length string.
func FormatCheckApplyConditionRequest(req domain.CheckApplyConditionRequest) string {
	var builder strings.Builder

	builder.WriteString(utils.PadOrTruncate(req.ApplicationNo, 20))
	builder.WriteString(utils.PadOrTruncate(req.IDCardNo, 20))
	builder.WriteString(utils.PadIntWithZero(req.Birthdate, 8))
	builder.WriteString(utils.PadOrTruncate(req.SuppIDCardNo, 20))
	builder.WriteString(utils.PadIntWithZero(req.SuppBirthdate, 8))
	builder.WriteString(utils.PadOrTruncate(req.ApplyChannel, 1))
	builder.WriteString(utils.PadIntWithZero(req.ApplicationDate, 8))
	builder.WriteString(utils.PadOrTruncate(req.BranchCode, 4))
	builder.WriteString(utils.PadOrTruncate(req.SourceCode, 8))
	builder.WriteString(utils.PadOrTruncate(req.StaffCode, 7))
	builder.WriteString(utils.PadIntWithZero(req.TotalApplyCard, 2))

	for _, item := range req.ApplyCardList {
		builder.WriteString(utils.PadIntWithZero(item.CardApplyType, 1))
		builder.WriteString(utils.PadOrTruncate(item.CardCode, 2))
		builder.WriteString(utils.PadOrTruncate(item.PrimaryCreditCard, 16))
		builder.WriteString(utils.PadOrTruncate(item.VirtualCardFlag, 1))
	}

	return builder.String()
}

func FormatCheckApplyConditionResponse(raw string) (domain.CheckApplyConditionResponse, error) {
	const headerLen = 123
	const dataLen = 73

	if len(raw) <= headerLen {
		return domain.CheckApplyConditionResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.CheckApplyConditionResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	applicationNo                 := parser.ReadString(0,20)
	status                        := parser.ReadString(20,1)
	seasonCode                    := parser.ReadString(21,2)
	reasonDescription             := parser.ReadString(23,50)

	return domain.CheckApplyConditionResponse{
		ApplicationNo:             applicationNo,
		Status:                    status,
		ReasonCode:                seasonCode,
		ReasonDescription:         reasonDescription,
	}, nil
}

// Converts CheckApplyCondition2ndCardRequest to a fixed-length string.
func FormatCheckApplyCondition2ndCardRequest(req domain.CheckApplyCondition2ndCardRequest) string {
	var builder strings.Builder

	builder.WriteString(utils.PadOrTruncate(req.IDCardNo, 20))
	builder.WriteString(utils.PadIntWithZero(req.TotalOfApplyCard, 2))

	for _, item := range req.CheckApply2ndCardList {
		builder.WriteString(utils.PadOrTruncate(item.CardCode, 2))

	}

	return builder.String()
}

func FormatCheckApplyCondition2ndCardResponse(raw string) (domain.CheckApplyCondition2ndCardResponse, error) {
	const headerLen = 123
	if len(raw) <= headerLen {
		return domain.CheckApplyCondition2ndCardResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}
	body := raw[headerLen:]
	runes := []rune(body)

	parser := utils.NewFixedParser(body)

	idCardNo         := parser.ReadString(0, 20)
	maximumCR        := parser.ReadInt(20, 2)
	haveCardCR       := parser.ReadInt(22, 2)
	maximumYC        := parser.ReadInt(24, 2)
	haveCardYC       := parser.ReadInt(26, 2)
	totalOfApplyCard := parser.ReadInt(28, 2)

	const agreementLen = 56
	agreementStart := 30
	agreements := make([]domain.CheckApply2ndCardRsOBJ, 0, totalOfApplyCard)

	for i := 0; i < totalOfApplyCard; i++ {
		start := agreementStart + i*agreementLen
		if start >= len(runes) {
			break
		}
		blockRunes := utils.ReadRunes(runes, start, agreementLen)

		agreements = append(agreements, domain.CheckApply2ndCardRsOBJ{
			CardCode:              utils.ReadBlockStr(blockRunes, 0, 2),
			ResultCode:            utils.ReadBlockStr(blockRunes, 2, 2),
			ReasonCode:            utils.ReadBlockStr(blockRunes, 4, 2),
			ReasonDescription:     utils.ReadBlockStr(blockRunes, 6, 50),
		})
	}

	return domain.CheckApplyCondition2ndCardResponse{
		IDCardNo:                  idCardNo,
		MaximumCR:                 maximumCR,
		HaveCardCR:                haveCardCR,
		MaximumYC:                 maximumYC,
		HaveCardYC:                haveCardYC,
		TotalOfApplyCard:          totalOfApplyCard,
		CheckApply2ndCardList:     agreements,
	}, nil
}