package format

import (
	"fmt"
	// "strconv"
	"strings"
	//"bytes"

	"connectorapi-go/internal/core/domain" 
	"connectorapi-go/internal/adapter/utils"
)

// Converts GetCardSalesRequest to a fixed-length string.
func FormatGetCardSalesRequest(getCardSalesReq domain.GetCardSalesRequest) string {
	IDCardNo        := utils.PadOrTruncate(getCardSalesReq.IDCardNo, 20)
	CardType        := utils.PadOrTruncate(getCardSalesReq.CardType, 2)
	CardBINno       := utils.PadOrTruncate(getCardSalesReq.CardBINno, 7)
	// UsingTypeCPCH   := utils.PadOrTruncate(getCardSalesReq.UsingTypeCPCH, 1)
	UsingTypeCPCH   := "Y"
	// UsingTypeCA     := utils.PadOrTruncate(getCardSalesReq.UsingTypeCA, 1)
	UsingTypeCA     := "Y"
	SaleDateFrom    := utils.PadOrTruncate(getCardSalesReq.SaleDateFrom, 8)
	SaleDateTo      := utils.PadOrTruncate(getCardSalesReq.SaleDateTo, 8)
	// MCCCodeCPCH     := utils.PadOrTruncate(getCardSalesReq.MCCCodeCPCH, 4)
	MCCCodeCPCH     := "0000"
	// AgencyCodeCPCH  := utils.PadOrTruncate(getCardSalesReq.AgencyCodeCPCH, 4)
	AgencyCodeCPCH  := "0000"
	// ShopCodeCPCH    := utils.PadOrTruncate(getCardSalesReq.ShopCodeCPCH, 2)
	ShopCodeCPCH    := "00"
	return IDCardNo + CardType + CardBINno + UsingTypeCPCH + UsingTypeCA + SaleDateFrom + SaleDateTo + MCCCodeCPCH + AgencyCodeCPCH + ShopCodeCPCH
}

func FormatGetCardSalesResponse(raw string) (domain.GetCardSalesResponse, error) {
	const headerLen = 123
	const dataLen = 216

	if len(raw) <= headerLen {
		return domain.GetCardSalesResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetCardSalesResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	totalSaleCount                    := parser.ReadString(0,5)
	totalSaleAmount                   := parser.ReadString(5,11)
	totalFACount                      := parser.ReadString(16,5)
	totalFAAmount                     := parser.ReadString(21,11)
	totalFRCount                      := parser.ReadString(32,5)
	totalFRAmount                     := parser.ReadString(37,11)
	lastSaleDate                      := parser.ReadString(48,8)
	cPSaleCount                       := parser.ReadString(56,5)
	cPSaleAmount                      := parser.ReadString(61,11)
	cPLastSaleDate                    := parser.ReadString(72,8)
	cHSaleCount                       := parser.ReadString(80,5)
	cHSaleAmount                      := parser.ReadString(85,11)
	cHLastSaleDate                    := parser.ReadString(96,8)
	cANormalSaleCount                 := parser.ReadString(104,5)
	cANormalSaleAmount                := parser.ReadString(109,11)
	cANormalLastSaleDate              := parser.ReadString(120,8)
	cACardlessSaleCount               := parser.ReadString(128,5)
	cACardlessSaleAmount              := parser.ReadString(133,11)
	cACardlessLastSaleDate            := parser.ReadString(144,8)
	cPSaleReversalCount               := parser.ReadString(152,5)
	cPSaleReversalAmount              := parser.ReadString(157,11)
	cHSaleReversalCount               := parser.ReadString(168,5)
	cHSaleReversalAmount              := parser.ReadString(173,11)
	cANormalSaleReversalCount         := parser.ReadString(184,5)
	cANormalSaleReversalAmount        := parser.ReadString(189,11)
	cACardlessSaleReversalCount       := parser.ReadString(200,5)
	cACardlessSaleReversalAmount      := parser.ReadString(205,11)

	return domain.GetCardSalesResponse{
		TotalSaleCount:                totalSaleCount,
		TotalSaleAmount:               totalSaleAmount,
		TotalFACount:                  totalFACount,
		TotalFAAmount:                 totalFAAmount,
		TotalFRCount:                  totalFRCount,
		TotalFRAmount:                 totalFRAmount,
		LastSaleDate:                  lastSaleDate,
		CPSaleCount:                   cPSaleCount,
		CPSaleAmount:                  cPSaleAmount,
		CPLastSaleDate:                cPLastSaleDate,
		CHSaleCount:                   cHSaleCount,
		CHSaleAmount:                  cHSaleAmount,
		CHLastSaleDate:                cHLastSaleDate,
		CANormalSaleCount:             cANormalSaleCount,
		CANormalSaleAmount:            cANormalSaleAmount,
		CANormalLastSaleDate:          cANormalLastSaleDate,
		CACardlessSaleCount:           cACardlessSaleCount,
		CACardlessSaleAmount:          cACardlessSaleAmount,
		CACardlessLastSaleDate:        cACardlessLastSaleDate,
		CPSaleReversalCount:           cPSaleReversalCount,
		CPSaleReversalAmount:          cPSaleReversalAmount,
		CHSaleReversalCount:           cHSaleReversalCount,
		CHSaleReversalAmount:          cHSaleReversalAmount,
		CANormalSaleReversalCount:     cANormalSaleReversalCount,
		CANormalSaleReversalAmount:    cANormalSaleReversalAmount,
		CACardlessSaleReversalCount:   cACardlessSaleReversalCount,
		CACardlessSaleReversalAmount:  cACardlessSaleReversalAmount,
	}, nil
} 

// Converts GetBigCardInfoRequest to a fixed-length string.
func FormatGetBigCardInfoRequest(getBigCardInfoReq domain.GetBigCardInfoRequest) string {
	TransactionDate        := utils.PadOrTruncate(getBigCardInfoReq.TransactionDate, 8)
	TransactionTime        := utils.PadOrTruncate(getBigCardInfoReq.TransactionTime, 6)
	TransactionType        := utils.PadOrTruncate(getBigCardInfoReq.TransactionType, 2)
	TraceNumber            := utils.PadOrTruncate(getBigCardInfoReq.TraceNumber, 20)
	AeonID                 := utils.PadOrTruncate(getBigCardInfoReq.AeonID, 44)
	BusinessCode           := utils.PadOrTruncate(getBigCardInfoReq.BusinessCode, 2)
	CreditCardNo           := utils.PadOrTruncate(getBigCardInfoReq.CreditCardNo, 16)
	Reserve1               := utils.PadOrTruncate(getBigCardInfoReq.Reserve1, 20)

	return TransactionDate + TransactionTime + TransactionType + TraceNumber + AeonID + BusinessCode + CreditCardNo + Reserve1
}

func FormatGetBigCardInfoResponse(raw string) (domain.GetBigCardInfoResponse, error) {
	const headerLen = 123
	const dataLen = 318

	if len(raw) <= headerLen {
		return domain.GetBigCardInfoResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetBigCardInfoResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	transactionDate             := parser.ReadString(0, 8)
	transactionTime             := parser.ReadString(8, 6)
	transactionType             := parser.ReadString(14, 2)
	// traceNumber                 := parser.ReadString(16, 20)
	aeonID                      := parser.ReadString(36, 44)
	businessCode                := parser.ReadString(80, 2)
	creditCardNo                := parser.ReadString(82, 16)
	// bigCardNo                   := parser.ReadString(98, 20)
	dataEncrypt                 := parser.ReadString(118, 128)
	// returnCode                  := parser.ReadString(246, 2)
	// responseText                := parser.ReadString(248, 50)
	// reserve1                    := parser.ReadString(298, 20)

	return domain.GetBigCardInfoResponse{
		TransactionDate:               transactionDate,
		TransactionTime:               transactionTime,
		TransactionType:               transactionType,
		AeonID:                        aeonID,
		BusinessCode:                  businessCode,
		CreditCardNo:                  creditCardNo,
		DataEncrypt:                   dataEncrypt,
	}, nil
}

// Converts GetCardDelinquentRequest to a fixed-length string.
func FormatGetCardDelinquentRequest(getCardDelinquentReq domain.GetCardDelinquentRequest) string {
	IDCardNo        := utils.PadOrTruncate(getCardDelinquentReq.IDCardNo, 20)
	CardType        := utils.PadOrTruncate(getCardDelinquentReq.CardType, 2)

	return IDCardNo + CardType
}

func FormatGetCardDelinquentResponse(raw string) (domain.GetCardDelinquentResponse, error) {
	const headerLen = 123
	const dataLen = 6

	if len(raw) <= headerLen {
		return domain.GetCardDelinquentResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetCardDelinquentResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	delinquentCountFAFR             := parser.ReadString(0, 3)
	delinquentCountAll              := parser.ReadString(3, 3)

	return domain.GetCardDelinquentResponse{
		DelinquentCountFAFR:              delinquentCountFAFR,
		DelinquentCountAll:               delinquentCountAll,
	}, nil
}

// Converts GetFullpanRequest to a fixed-length string.
func FormatGetFullpanRequest(req domain.GetFullpanRequest) string {
	var builder strings.Builder

	builder.WriteString(utils.PadOrTruncate(req.IDCardNo, 20))
	builder.WriteString(utils.PadIntWithZero(req.TotalCard, 3))

	for _, item := range req.CardList {
		builder.WriteString(utils.PadOrTruncate(item.CardNo, 16))
		builder.WriteString(utils.PadOrTruncate(item.CardCode, 2))

	}

	return builder.String()
}

func FormatGetFullpanResponse(raw string) (domain.GetFullpanResponse, error) {
	const headerLen = 123
	if len(raw) <= headerLen {
		return domain.GetFullpanResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}
	body := raw[headerLen:]
	runes := []rune(body)

	parser := utils.NewFixedParser(body)

	idCardNo         := parser.ReadString(0, 20)
	totalCard        := parser.ReadInt(20, 3)

	const agreementLen = 56
	agreementStart := 23
	agreements := make([]domain.GetFullpanRsOBJ, 0, totalCard)

	for i := 0; i < totalCard; i++ {
		start := agreementStart + i*agreementLen
		if start >= len(runes) {
			break
		}
		blockRunes := utils.ReadRunes(runes, start, agreementLen)

		agreements = append(agreements, domain.GetFullpanRsOBJ{
		CardNo:             utils.ReadBlockStr(blockRunes, 0, 16),
		CardCode:           utils.ReadBlockStr(blockRunes, 16, 2),
		CardType:           utils.ReadBlockStr(blockRunes, 18, 2),
		ExpireDate:         utils.ReadBlockInt(blockRunes, 20, 8),
		HoldCode:           utils.ReadBlockStr(blockRunes, 28, 2),
		SendMode:           utils.ReadBlockStr(blockRunes, 30, 1),
		FirstEmbossDate:    utils.ReadBlockInt(blockRunes, 31, 8),
		FirstConfirmDate:   utils.ReadBlockInt(blockRunes, 39, 8),
		DigitalCardFlag:    utils.ReadBlockStr(blockRunes, 47, 1),

		})
	}

	return domain.GetFullpanResponse{
		IDCardNo:                  idCardNo,
		TotalCard:                 totalCard,
		CardList:                  agreements,
	}, nil
}

// Converts GetCardEnrollRequest to a fixed-length string.
func FormatGetCardEnrollRequest(getCardEnrollReq domain.GetCardEnrollRequest) string {
	IDCardNo        := utils.PadOrTruncate(getCardEnrollReq.IDCardNo, 20)
	CardType        := utils.PadOrTruncate(getCardEnrollReq.CardType, 2)

	return IDCardNo + CardType
}

func FormatGetCardEnrollResponse(raw string) (domain.GetCardEnrollResponse, error) {
	const headerLen = 123
	const dataLen = 18

	if len(raw) <= headerLen {
		return domain.GetCardEnrollResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetCardEnrollResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	enrollmentNo             := parser.ReadString(0, 16)
	status                   := parser.ReadString(16, 2)

	return domain.GetCardEnrollResponse{
		EnrollmentNo:              enrollmentNo,
		Status:                    status,
	}, nil
}