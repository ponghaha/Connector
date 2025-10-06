package format

import (
	"fmt"

	"connectorapi-go/internal/core/domain"
	"connectorapi-go/internal/adapter/utils"
)

// Converts DashboardSummaryRequest to a fixed-length string.
func FormatDashboardSummaryRequest(flagOldFormatReq bool, dashboardSummaryReq domain.DashboardSummaryRequest) string {
	idCardNo := utils.PadOrTruncate(dashboardSummaryReq.IDCardNo, 20)
	aeonID   := utils.PadOrTruncate(dashboardSummaryReq.AeonID, 20)

	var result string
	if flagOldFormatReq {
		result = idCardNo
	} else {
		result = aeonID
	}
	return result
}

func FormatDashboardSummaryResponse(raw string, flagOldFormatReq bool) (domain.DashboardSummaryResponse, error) {
	const headerLen = 123
	const dataLen = 354
	var idCardNo string
	var aeonID string
	var remainMinimumPayment float64
    var remainFullPayment float64
	var remainMinPtr *utils.DecimalString
	var remainFullPtr *utils.DecimalString
	var counterNo string
	var accountBank string
	var branchBank string
	var referenceAccountNo string
	var termStart int

	if len(raw) <= headerLen {
		return domain.DashboardSummaryResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.DashboardSummaryResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)
	if !flagOldFormatReq {
		aeonID = parser.ReadString(0,20)
	} else {
		idCardNo = parser.ReadString(0,20)
	}
	nameTH                       := parser.ReadString(20,30)
	nameEN                       := parser.ReadString(50,30)
	mobileNo                     := parser.ReadString(80,15)
	dueDate                      := parser.ReadInt(95,8)
	// AS400ResponseCode := int 2 not use แต่ต้องตัด
	creditShoppingFloorLimit     := parser.ReadFloat100(105,11)
	creditShoppingOutstanding    := parser.ReadFloat100(116,11)
	creditShoppingAvailableLimit := parser.ReadFloat100(127,11)
	ceditCashingFloorLimit       := parser.ReadFloat100(138,11)
	creditCashingOutstanding     := parser.ReadFloat100(149,11)
	creditCashingAvailableLimit  := parser.ReadFloat100(160,11)
	yourCashFloorLimit           := parser.ReadFloat100(171,11)
	yourCashOutstanding          := parser.ReadFloat100(182,11)
	yourCashAvailableLimit       := parser.ReadFloat100(193,11)
	ropShoppingFloorLimit        := parser.ReadFloat100(204,11)
	ropShoppingOutstanding       := parser.ReadFloat100(215,11)
	ropShoppingAvailableLimit    := parser.ReadFloat100(226,11)
	ropCashingFloorLimit         := parser.ReadFloat100(237,11)
	ropCashingOutstanding        := parser.ReadFloat100(248,11)
	ropCashingAvailableLimit     := parser.ReadFloat100(259,11)
	totalMinimumPayment          := parser.ReadFloat100(270,11)
	totalFullPayment             := parser.ReadFloat100(281,11)
	totalPaidAmount              := parser.ReadFloat100(292,11)
	pendingPaymentStatus         := parser.ReadString(303,2)
	if !flagOldFormatReq {
		remainMinimumPayment     = parser.ReadFloat100(305,11)
    	remainFullPayment        = parser.ReadFloat100(316,11)
		remainMin                := utils.DecimalString(remainMinimumPayment)
		remainFull               := utils.DecimalString(remainFullPayment)
		remainMinPtr             = &remainMin
		remainFullPtr            = &remainFull
		counterNo                = parser.ReadString(327,4)
		accountBank              = parser.ReadString(331,20)
		branchBank               = parser.ReadString(351,5)
		referenceAccountNo       = parser.ReadString(356,20)
		termStart = 376
	} else {
		remainMinPtr             = nil
   	 	remainFullPtr            = nil
		counterNo                = parser.ReadString(305,4)
		accountBank              = parser.ReadString(309,20)
		branchBank               = parser.ReadString(329,5)
		referenceAccountNo       = parser.ReadString(334,20)
		termStart = 354
	}

	const termLen = 32
	runes := []rune(data)
	totalRunes := len(runes)
	remainingRunes := totalRunes - termStart
	txtlen := remainingRunes / termLen 
	terms := make([]domain.DBDetailTermsListRq, 0, txtlen)

	for i := 0; i < txtlen; i++ {
		start := termStart + i*termLen
		if start >= len(runes) {
			break
		}
		blockRunes := utils.ReadRunes(runes, start, termLen)
		terms = append(terms, domain.DBDetailTermsListRq{
			TermsType:         utils.ReadBlockStr(blockRunes, 0, 20),
			TermsVersion:      utils.ReadBlockStr(blockRunes, 20, 10),
			TermsAcceptStatus: utils.ReadBlockStr(blockRunes, 30, 2),
		})
	}

	banks := make([]domain.DBDetailBankListRq, 0)
	banks = append(banks, domain.DBDetailBankListRq{
		CounterNo:    counterNo,
		AccountNo:    accountBank,
		BranchBank:   branchBank,
		RefAccountNo: referenceAccountNo,
	})

	return domain.DashboardSummaryResponse{
		IDCardNo:                     idCardNo,
		AeonID:                       aeonID,
		NameTH:                       nameTH,
		NameEN:                       nameEN,
		MobileNo:                     mobileNo,
		DueDate:                      dueDate,
		CreditShoppingFloorLimit:     utils.DecimalString(creditShoppingFloorLimit),
		CreditShoppingOutstanding:    utils.DecimalString(creditShoppingOutstanding),
		CreditShoppingAvailableLimit: utils.DecimalString(creditShoppingAvailableLimit),
		CreditCashingFloorLimit:      utils.DecimalString(ceditCashingFloorLimit),
		CreditCashingOutstanding:     utils.DecimalString(creditCashingOutstanding),
		CreditCashingAvailableLimit:  utils.DecimalString(creditCashingAvailableLimit),
		YourCashFloorLimit:           utils.DecimalString(yourCashFloorLimit),
		YourCashOutstanding:          utils.DecimalString(yourCashOutstanding),
		YourCashAvailableLimit:       utils.DecimalString(yourCashAvailableLimit),
		ROPShoppingFloorLimit:        utils.DecimalString(ropShoppingFloorLimit),
		ROPShoppingOutstanding:       utils.DecimalString(ropShoppingOutstanding),
		ROPShoppingAvailableLimit:    utils.DecimalString(ropShoppingAvailableLimit),
		ROPCashingFloorLimit:         utils.DecimalString(ropCashingFloorLimit),
		ROPCashingOutstanding:        utils.DecimalString(ropCashingOutstanding),
		ROPCashingAvailableLimit:     utils.DecimalString(ropCashingAvailableLimit),
		TotalMinimumPayment:          utils.DecimalString(totalMinimumPayment),
		TotalFullPayment:             utils.DecimalString(totalFullPayment),
		TotalPaidAmount:              utils.DecimalString(totalPaidAmount),
		PendingPaymentStatus:         pendingPaymentStatus,
		RemainMinimumPayment:         remainMinPtr,
		RemainFullPayment:            remainFullPtr,
		BankList:                     banks,
		TermsList:                    terms,
	}, nil
}

// Converts DashboardDetailRequest to a fixed-length string.
func FormatDashboardDetailRequest(flagOldFormatReq bool, dashboardDetailReq domain.DashboardDetailRequest) string {
	idCardNo := utils.PadOrTruncate(dashboardDetailReq.IDCardNo, 20)
	aeonID   := utils.PadOrTruncate(dashboardDetailReq.AeonID, 20)

	var result string
	if flagOldFormatReq {
		result = idCardNo
	} else {
		result = aeonID
	}
	return result
}

func FormatDashboardDetailResponse(raw string, flagOldFormatReq bool) (domain.DashboardDetailResponse, error) {
	const headerLen = 123
	const dataLen = 30
	var idCardNo string
	var aeonID string
	var remainMinimumPayment float64
    var remainFullPayment float64
	var remainMinPtr *utils.DecimalString
	var remainFullPtr *utils.DecimalString
	var dbDetailLen int

	if len(raw) <= headerLen {
		return domain.DashboardDetailResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.DashboardDetailResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)
	if !flagOldFormatReq {
		aeonID = parser.ReadString(0,20)
		dbDetailLen = 249
	} else {
		idCardNo = parser.ReadString(0,20)
		dbDetailLen = 227
	}
	dueDate := parser.ReadInt(20,8)
	// AS400ResponseCode := int 2 not use แต่ต้องตัด
	dbDetailStart := 30
	runes := []rune(data)
	totalRunes := len(runes)
	remainingRunes := totalRunes - dbDetailStart
	txtlen := remainingRunes / dbDetailLen 
	dbDetails := make([]domain.DBDetailListRs, 0, txtlen)

	for i := 0; i < txtlen; i++ {
		start := dbDetailStart + i*dbDetailLen
		if start >= len(runes) {
			break
		}
		blockRunes := utils.ReadRunes(runes, start, dbDetailLen)
		detail := domain.DBDetailListRs{
			CreditCardNo:         utils.ReadBlockStr(blockRunes, 0, 16),
			CardName:             utils.ReadBlockStr(blockRunes, 16, 30),
			ProductType:          utils.ReadBlockStr(blockRunes, 46, 2),
			CardCode:             utils.ReadBlockStr(blockRunes, 48, 2),
			ATMauthorize:         utils.ReadBlockStr(blockRunes, 50, 14),
			CardStatus:           utils.ReadBlockStr(blockRunes, 64, 16),
			MinimumPaymentAmount: utils.ReadBlockFloat100ToDecimal(blockRunes, 80, 11),
			FullPaymentAmount:    utils.ReadBlockFloat100ToDecimal(blockRunes, 91, 11),
			PaidAmount:           utils.ReadBlockFloat100ToDecimal(blockRunes, 102, 11),
		}
		if !flagOldFormatReq {
			remainMinimumPayment                = parser.ReadFloat100(113,11)
			remainFullPayment                   = parser.ReadFloat100(124,11)
			remainMin                           := utils.DecimalString(remainMinimumPayment)
			remainFull                          := utils.DecimalString(remainFullPayment)
			remainMinPtr                        = &remainMin
			remainFullPtr                       = &remainFull
			detail.RemainMinimumPayment         = remainMinPtr
			detail.RemainFullPayment            = remainFullPtr
			detail.CreditShoppingFloorLimit     = utils.ReadBlockFloat100ToDecimal(blockRunes, 135, 11)
			detail.CreditShoppingOutstanding    = utils.ReadBlockFloat100ToDecimal(blockRunes, 146, 11)
			detail.CreditShoppingAvailableLimit = utils.ReadBlockFloat100ToDecimal(blockRunes, 157, 11)
			detail.CreditCashingFloorLimit      = utils.ReadBlockFloat100ToDecimal(blockRunes, 168, 11)
			detail.CreditCashingOutstanding     = utils.ReadBlockFloat100ToDecimal(blockRunes, 179, 11)
			detail.CreditCashingAvailableLimit  = utils.ReadBlockFloat100ToDecimal(blockRunes, 190, 11)
			detail.AvailablePoint               = utils.ReadBlockFloat100ToDecimal(blockRunes, 201, 11)
			detail.BillingAmount                = utils.ReadBlockFloat100ToDecimal(blockRunes, 212, 11)
			detail.UnbilledAmount               = utils.ReadBlockFloat100ToDecimal(blockRunes, 223, 11)
			detail.InstallmentNo                = utils.ReadBlockInt(blockRunes, 234, 3)
			detail.InstallmentCurrent           = utils.ReadBlockInt(blockRunes, 237, 3)
			detail.DigitalCardFlag              = utils.ReadBlockStr(blockRunes,240, 1)
			detail.ApplicationDate              = utils.ReadBlockInt(blockRunes, 241, 8)
		} else {
			detail.RemainMinimumPayment         = nil
			detail.RemainFullPayment            = nil
			detail.CreditShoppingFloorLimit     = utils.ReadBlockFloat100ToDecimal(blockRunes, 113, 11)
			detail.CreditShoppingOutstanding    = utils.ReadBlockFloat100ToDecimal(blockRunes, 124, 11)
			detail.CreditShoppingAvailableLimit = utils.ReadBlockFloat100ToDecimal(blockRunes, 135, 11)
			detail.CreditCashingFloorLimit      = utils.ReadBlockFloat100ToDecimal(blockRunes, 146, 11)
			detail.CreditCashingOutstanding     = utils.ReadBlockFloat100ToDecimal(blockRunes, 157, 11)
			detail.CreditCashingAvailableLimit  = utils.ReadBlockFloat100ToDecimal(blockRunes, 168, 11)
			detail.AvailablePoint               = utils.ReadBlockFloat100ToDecimal(blockRunes, 179, 11)
			detail.BillingAmount                = utils.ReadBlockFloat100ToDecimal(blockRunes, 190, 11)
			detail.UnbilledAmount               = utils.ReadBlockFloat100ToDecimal(blockRunes, 201, 11)
			detail.InstallmentNo                = utils.ReadBlockInt(blockRunes, 212, 3)
			detail.InstallmentCurrent           = utils.ReadBlockInt(blockRunes, 215, 3)
			detail.DigitalCardFlag              = utils.ReadBlockStr(blockRunes, 218, 1)
			detail.ApplicationDate              = utils.ReadBlockInt(blockRunes, 219, 8)
		}
		dbDetails = append(dbDetails, detail)
	}

	return domain.DashboardDetailResponse{
		IDCardNo:            idCardNo,
		AeonID:              aeonID,
		DueDate:             dueDate,
		DashboardDetailList: dbDetails,
	}, nil
}

// Converts MobileFullPanRequest to a fixed-length string.
func FormatMobileFullPanRequest(mobileFullPanReq domain.MobileFullPanFormatRequest) string {
	IDCardNo     := utils.PadOrTruncate(mobileFullPanReq.IDCardNo, 20)
	CreditCardNo := utils.PadOrTruncate(mobileFullPanReq.CreditCardNo, 16)
	BusinessCode := utils.PadOrTruncate(mobileFullPanReq.BusinessCode, 2)
	return IDCardNo + CreditCardNo + BusinessCode
}

func FormatMobileFullPanResponse(raw string) (domain.MobileFullPanResponse, error) {
	const headerLen = 123
	const dataLen = 24
	// var lastChar string

	if len(raw) <= headerLen {
		return domain.MobileFullPanResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.MobileFullPanResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser    := utils.NewFixedParser(data)
	idCardNo  := parser.ReadString(0,20)
	totalCard := parser.ReadInt(20,4)
	cardStart := 24

	const cardLen = 61
	runes := []rune(data)
	totalRunes := len(runes)
	remainingRunes := totalRunes - cardStart
	txtlen := remainingRunes / cardLen 
	cards := make([]domain.MobileCardListRs, 0, txtlen)

	for i := 0; i < txtlen; i++ {
		start := cardStart + i*cardLen
		if start >= len(runes) {
			break
		}
		blockRunes := utils.ReadRunes(runes, start, cardLen)
		cards = append(cards, domain.MobileCardListRs{
			CardNo:           utils.ReadBlockStr(blockRunes, 0, 16),
			CardType:         utils.ReadBlockStr(blockRunes, 46, 2),
			CardCode:         utils.ReadBlockStr(blockRunes, 48, 2),
			HoldCode:         utils.ReadBlockStr(blockRunes, 50, 2),
			ExpireDate:       utils.ReadBlockInt(blockRunes, 52, 8),
			SendMode:         "",
			FirstEmbossDate:  0,
			FirstConfirmDate: 0,
			DigitalCardFlag:  utils.ReadBlockStr(blockRunes, 60, 1),

		})
	}

	return domain.MobileFullPanResponse{
		IDCardNo:   idCardNo,
		TotalCard:  totalCard,
		CardListRs: cards,
	}, nil
}
