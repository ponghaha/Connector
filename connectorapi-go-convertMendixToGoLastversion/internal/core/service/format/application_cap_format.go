package format

import (
	"fmt"
	"strconv"
	"strings"

	"connectorapi-go/internal/core/domain"
	"connectorapi-go/internal/adapter/utils"
)

// Converts GetApplicationNoRequest to a fixed-length string.
func FormatGetApplicationNoRequest(getApplicationNoReq domain.GetApplicationNoRequest) string {
	var cardList string
	var strTotalApplyCard string
	
	IDCardNo       := utils.PadOrTruncate(getApplicationNoReq.IDCardNo, 20)
	ApplyChannel   := utils.PadOrTruncate(getApplicationNoReq.ApplyChannel, 1)
	if getApplicationNoReq.TotalApplyCard >= 10 {
		strTotalApplyCard = strconv.Itoa(getApplicationNoReq.TotalApplyCard)
	} else {
		strTotalApplyCard = "0" + strconv.Itoa(getApplicationNoReq.TotalApplyCard)
	}

	for _, card := range getApplicationNoReq.CardListRq {
		CardCode := utils.PadOrTruncate(card.CardCode, 2)
		VirtualCardFlag := utils.PadOrTruncate(card.VirtualCardFlag, 1)
		cardList += CardCode + VirtualCardFlag
	}

	return IDCardNo + ApplyChannel + strTotalApplyCard + cardList
}

func FormatGetApplicationNoResponse(raw string) (domain.GetApplicationNoResponse, error) {
	const headerLen = 123
	const dataLen = 106

	if len(raw) <= headerLen {
		return domain.GetApplicationNoResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetApplicationNoResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	applicationNo     := strings.TrimSpace(data[:20])
	idCardNo          := strings.TrimSpace(data[20:40])
	applicationDate   := strings.TrimSpace(data[40:48])
	applicationTime   := strings.TrimSpace(data[48:54])
	resultCode        := strings.TrimSpace(data[54:56])
	resultDescription := strings.TrimSpace(data[56:106])

	return domain.GetApplicationNoResponse{
		ApplicationNo:     applicationNo,
		IDCardNo:          idCardNo,
		ApplicationDate:   applicationDate,
		ApplicationTime:   applicationTime,
		ResultCode:        resultCode,
		ResultDescription: resultDescription,
	}, nil
}

// Converts SubmitCardApplicationRequest to a fixed-length string.
func FormatSubmitCardApplicationRequest(submitCardApplicationNoReq domain.SubmitCardApplicationRequest) string {
	var cardList string
	var strTotalApplyCard string

	IDCardNo        := utils.PadOrTruncate(submitCardApplicationNoReq.IDCardNo, 20)
	ApplicationNo   := utils.PadOrTruncate(submitCardApplicationNoReq.ApplicationNo, 20)
	ApplyChannel    := utils.PadOrTruncate(submitCardApplicationNoReq.ApplyChannel, 1)
	ApplicationDate := utils.PadOrTruncate(submitCardApplicationNoReq.ApplicationDate, 8)
	BranchCode      := utils.PadOrTruncate(submitCardApplicationNoReq.BranchCode, 4)
	SourceCode      := utils.PadOrTruncate(submitCardApplicationNoReq.SourceCode, 8)
	StaffCode       := utils.PadOrTruncate(submitCardApplicationNoReq.StaffCode, 7)
	MailTo          := utils.PadOrTruncate(submitCardApplicationNoReq.MailTo, 1)

	if submitCardApplicationNoReq.TotalApplyCard >= 10 {
		strTotalApplyCard = strconv.Itoa(submitCardApplicationNoReq.TotalApplyCard)
	} else {
		strTotalApplyCard = "0" + strconv.Itoa(submitCardApplicationNoReq.TotalApplyCard)
	}

	for _, card := range submitCardApplicationNoReq.SubmitCardListRq {
		CardCode := utils.PadOrTruncate(card.CardCode, 2)
		VirtualCardFlag := utils.PadOrTruncate(card.VirtualCardFlag, 1)
		cardList += CardCode + VirtualCardFlag
	}

	return IDCardNo + ApplicationNo +ApplyChannel + ApplicationDate + BranchCode + SourceCode + StaffCode + MailTo + strTotalApplyCard + cardList
}

func FormatSubmitCardApplicationResponse(raw string) (domain.SubmitCardApplicationResponse, error) {
	const headerLen = 123
	const dataLen = 126

	if len(raw) <= headerLen {
		return domain.SubmitCardApplicationResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.SubmitCardApplicationResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser            := utils.NewFixedParser(data)
	idCardNo          := parser.ReadString(0,20)
	applicationNo     := parser.ReadString(20,20)
	applicationDate   := parser.ReadString(40,8)
	resultDate        := parser.ReadString(48,8)
	resultTime        := parser.ReadString(56,6)
	programID         := parser.ReadString(62,10)
	resultCode        := parser.ReadString(72,2)
	resultDescription := parser.ReadString(74,50)
	totalCard         := parser.ReadInt(124,2)
	cardStart         := 126

	const cardLen = 103
	runes := []rune(data)
	totalRunes := len(runes)
	remainingRunes := totalRunes - cardStart
	txtlen := remainingRunes / cardLen 
	cards := make([]domain.SubmitCardListRs, 0, txtlen)

	for i := 0; i < txtlen; i++ {
		start := cardStart + i*cardLen
		if start >= len(runes) {
			break
		}
		blockRunes := utils.ReadRunes(runes, start, cardLen)
		cards = append(cards, domain.SubmitCardListRs{
			MemberTempNo: utils.ReadBlockStr(blockRunes, 0, 16),
			CardCode:     utils.ReadBlockStr(blockRunes, 16, 2),
			ResultCode:   utils.ReadBlockStr(blockRunes, 18, 1),
			ReasonCode:   utils.ReadBlockStr(blockRunes, 19, 2),
			Remark1:      utils.ReadBlockStr(blockRunes, 21, 30),
			Remark2:      utils.ReadBlockStr(blockRunes, 51, 30),
			MaximumLimit: utils.ReadBlockFloat100ToDecimal(blockRunes, 81, 10),
			PINNumber:    utils.ReadBlockStr(blockRunes, 91, 12),
		})
	}

	return domain.SubmitCardApplicationResponse{
		IDCardNo:          idCardNo,
		ApplicationNo:     applicationNo,
		ApplicationDate:   applicationDate,
		ResultDate:        resultDate,
		ResultTime:        resultTime,
		ProgramID:         programID,
		ResultCode:        resultCode,
		ResultDescription: resultDescription,
		TotalApplyCard:    totalCard,
		SubmitCardListRs:  cards,
	}, nil
}
