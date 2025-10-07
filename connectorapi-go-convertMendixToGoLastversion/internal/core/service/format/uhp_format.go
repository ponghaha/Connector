package format

import (
	"fmt"
	"strconv"
	"strings"
	//"bytes"

	"connectorapi-go/internal/core/domain" 
	"connectorapi-go/internal/adapter/utils"
)

// Converts GetRedbookInfoRequest to a fixed-length string.
func FormatGetRedbookInfoRequest(getRedbookInfoReq domain.GetRedbookInfoRequest) string {
	AgentCode            := utils.PadOrTruncate(getRedbookInfoReq.AgentCode, 8)
	MarketingCode        := utils.PadOrTruncate(getRedbookInfoReq.MarketingCode, 10)
	Brand                := utils.PadOrTruncate(getRedbookInfoReq.Brand, 30)
	Model                := utils.PadOrTruncate(getRedbookInfoReq.Model, 30)
	CarYear              := utils.PadIntWithZero(getRedbookInfoReq.CarYear, 4)
	CarMonth             := utils.PadIntWithZero(getRedbookInfoReq.CarMonth, 2)
	SubModel             := utils.PadOrTruncate(getRedbookInfoReq.SubModel, 100)
	EffectiveYear        := utils.PadIntWithZero(getRedbookInfoReq.EffectiveYear, 4)
	EffectiveMonth       := utils.PadIntWithZero(getRedbookInfoReq.EffectiveMonth, 2)

	return AgentCode + MarketingCode + Brand + Model + CarYear + CarMonth + SubModel + EffectiveYear + EffectiveMonth
}

func FormatGetRedbookInfoResponse(raw string) (domain.GetRedbookInfoResponse, error) {
	const headerLen = 123
	const dataLen = 238

	if len(raw) <= headerLen {
		return domain.GetRedbookInfoResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetRedbookInfoResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	agentCode      := parser.ReadString(0, 8)
	marketingCode  := parser.ReadString(8, 10)
	brand          := parser.ReadString(18, 30)
	model          := parser.ReadString(48, 30)
	carYear        := parser.ReadInt(78, 4)
	carMonth       := parser.ReadInt(82, 2)
	subModel       := parser.ReadString(84, 100)
	effectiveYear  := parser.ReadInt(184, 4)
	effectiveMonth := parser.ReadInt(188, 2)
	vehicleCode    := parser.ReadString(190, 8)
	avgWholesale   := parser.ReadFloat100(198, 8)
	avgRetail      := parser.ReadFloat100(206, 8)
	goodWholesale  := parser.ReadFloat100(214, 8)
	goodRetail     := parser.ReadFloat100(222, 8)
	newPrice       := parser.ReadFloat100(230, 8)

	return domain.GetRedbookInfoResponse{
		AgentCode:      agentCode,
		MarketingCode:  marketingCode,
		Brand:          brand,
		Model:          model,
		CarYear:        carYear,
		CarMonth:       carMonth,
		SubModel:       subModel,
		EffectiveYear:  effectiveYear,
		EffectiveMonth: effectiveMonth,
		VehicleCode:    vehicleCode,
		AvgWholesale:   utils.DecimalString(avgWholesale),
		AvgRetail:      utils.DecimalString(avgRetail),
		GoodWholesale:  utils.DecimalString(goodWholesale),
		GoodRetail:     utils.DecimalString(goodRetail),
		NewPrice:       utils.DecimalString(newPrice),
	}, nil
}

// Converts GetDealerCommissionRequest to a fixed-length string.
func FormatGetDealerCommissionRequest(getDealerCommissionReq domain.GetDealerCommissionRequest) string {
	AgentCode            := utils.PadOrTruncate(getDealerCommissionReq.AgentCode, 8)
	MarketingCode        := utils.PadOrTruncate(getDealerCommissionReq.MarketingCode, 10)
	AgreementNo          := utils.PadOrTruncate(getDealerCommissionReq.AgreementNo, 12)
	CommissionCode       := utils.PadOrTruncate(getDealerCommissionReq.CommissionCode, 8)

	return AgentCode + MarketingCode + AgreementNo + CommissionCode
}

func FormatGetDealerCommissionResponse(raw string) (domain.GetDealerCommissionResponse, error) {
	const headerLen = 123
	const dataLen = 93

	if len(raw) <= headerLen {
		return domain.GetDealerCommissionResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetDealerCommissionResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	agentCode            := parser.ReadString(0, 8)
	marketingCode        := parser.ReadString(8, 10)
	agreementNo          := parser.ReadString(18, 12)
	commissionCode       := parser.ReadString(30, 8)
	agentCategory        := parser.ReadString(38, 2)
	totalCommission      := parser.ReadFloat100(40, 9)
	vatRate              := parser.ReadFloat100(49, 4)
	vat                  := parser.ReadFloat100(53, 9)
	grandTotalCommission := parser.ReadFloat100(62, 9)
	whtRate              := parser.ReadFloat100(71, 4)
	whtTax               := parser.ReadFloat100(75, 9)
	netTotalCommission   := parser.ReadFloat100(84, 9)

	return domain.GetDealerCommissionResponse{
		AgentCode:            agentCode,
		MarketingCode:        marketingCode,
		AgreementNo:          agreementNo,
		CommissionCode:       commissionCode,
		AgentCategory:        agentCategory,
		TotalCommission:      utils.DecimalString(totalCommission),
		VATRate:              utils.DecimalString(vatRate),
		VAT:                  utils.DecimalString(vat),
		GrandTotalCommission: utils.DecimalString(grandTotalCommission),
		WHTRate:              utils.DecimalString(whtRate),
		WHTTax:               utils.DecimalString(whtTax),
		NetTotalCommission:   utils.DecimalString(netTotalCommission),
	}, nil
}

// Converts GetDealerAgreementRequest to a fixed-length string.
func FormatGetDealerAgreementRequest(getDealerAgreementReq domain.GetDealerAgreementRequest) string {
	AgentCode            := utils.PadOrTruncate(getDealerAgreementReq.AgentCode, 8)
	MarketingCode        := utils.PadOrTruncate(getDealerAgreementReq.MarketingCode, 10)
	TransactionDateFrom  := utils.PadIntWithZero(getDealerAgreementReq.TransactionDateFrom, 8)
	TransactionDateTo    := utils.PadIntWithZero(getDealerAgreementReq.TransactionDateTo, 8)
	AgreementNo          := utils.PadOrTruncate(getDealerAgreementReq.AgreementNo, 12)

	return AgentCode + MarketingCode + TransactionDateFrom + TransactionDateTo + AgreementNo
}

func FormatGetDealerAgreementResponse(raw string) (domain.GetDealerAgreementResponse, error) {
	const headerLen = 123
	if len(raw) <= headerLen {
		return domain.GetDealerAgreementResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}
	body := raw[headerLen:]
	runes := []rune(body)

	readRunes := func(start, length int) []rune {
		if start >= len(runes) {
			return []rune{}
		}
		end := start + length
		if end > len(runes) {
			end = len(runes)
		}
		return runes[start:end]
	}

	readString := func(start, length int) string {
		rs := readRunes(start, length)
		return strings.TrimSpace(string(rs))
	}

	readInt := func(start, length int) int {
		s := readString(start, length)
		if s == "" {
			return 0
		}
		i, _ := strconv.Atoi(s)
		return i
	}

	// readFloat := func(start, length int) float64 {
	// 	s := readString(start, length)
	// 	if s == "" {
	// 		return 0
	// 	}
	// 	f, _ := strconv.ParseFloat(s, 64)
	// 	return f
	// }

	agentCode            := readString(0, 8)
	marketingCode        := readString(8, 10)
	transactionDateFrom  := readInt(18, 8)
	transactionDateTo    := readInt(26, 8)
	agreementNo          := readString(34, 12)
	totalAgreement       := readInt(46, 3)

	const agreementLen = 76
	agreementStart := 49
	agreements := make([]domain.AgreementListobj, 0, totalAgreement)

	for i := 0; i < totalAgreement; i++ {
		start := agreementStart + i*agreementLen
		if start >= len(runes) {
			break
		}
		blockRunes := readRunes(start, agreementLen)

		readBlockStr := func(startField, length int) string {
			if startField >= len(blockRunes) {
				return ""
			}
			end := startField + length
			if end > len(blockRunes) {
				end = len(blockRunes)
			}
			return strings.TrimSpace(string(blockRunes[startField:end]))
		}

		readBlockInt := func(startField, length int) int {
			s := readBlockStr(startField, length)
			if s == "" {
				return 0
			}
			i, _ := strconv.Atoi(s)
			return i
		}

	// 	readBlockFloat := func(startField, length int) utils.DecimalString {
    //         s := readBlockStr(startField, length)
    //         i, _ := strconv.ParseInt(s, 10, 64)
    //         f := float64(i) / 100.0
    //         return utils.DecimalString(f)
    //    }


		agreements = append(agreements, domain.AgreementListobj{
			AgreementNo:             readBlockStr(0, 12),
			TransactionDate:         readBlockInt(12, 8),
			CustomerName:            readBlockStr(20, 50),
			Status:                  readBlockStr(70, 6),
		})
	}

	return domain.GetDealerAgreementResponse{
		AgentCode:                  agentCode,
		MarketingCode:              marketingCode,
		TransactionDateFrom:        transactionDateFrom,
		TransactionDateTo:          transactionDateTo,
		AgreementNo:                agreementNo,
		TotalAgreement:             totalAgreement,
		AgreementList:              agreements,
	}, nil
}