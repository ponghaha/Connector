package format

import (
	"fmt"
	// "strconv"
	"strings"
	//"bytes"

	"connectorapi-go/internal/core/domain" 
	"connectorapi-go/internal/adapter/utils"
)

// Converts GetSpendingControlRequest to a fixed-length string.
func FormatGetSpendingControlRequest(getSpendingControlReq domain.GetSpendingControlRequest) string {
	TransactionType        := utils.PadOrTruncate(getSpendingControlReq.TransactionType, 2)
	AEONID                 := utils.PadOrTruncate(getSpendingControlReq.AEONID, 20)
	CardNo                 := utils.PadOrTruncate(getSpendingControlReq.CardNo, 16)
	BusinessCode           := utils.PadOrTruncate(getSpendingControlReq.BusinessCode, 2)
	Channel                := utils.PadOrTruncate(getSpendingControlReq.Channel, 10)
	return TransactionType + AEONID + CardNo + BusinessCode + Channel
}

func FormatGetSpendingControlResponse(raw string) (domain.GetSpendingControlResponse, error) {
	const headerLen = 123
	const dataLen = 216				// กรอกค่านี้ด้วย

	if len(raw) <= headerLen {
		return domain.GetSpendingControlResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetSpendingControlResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data) 				// กรอกค่านี้ด้วย

	transactionType                    := parser.ReadString(0,5)
	cardNotPresentStatus               := parser.ReadString(5,11)
	cNPLimitAmountPerDay               := parser.ReadInt(16,5)
	cNPLimitAmountPerTransaction       := parser.ReadInt(21,11)
	limitStatus                        := parser.ReadString(32,5)
	limitAmountPerDay                  := parser.ReadInt(37,11)
	limitAmountPerTransaction          := parser.ReadInt(48,8)

	return domain.GetSpendingControlResponse{
		TransactionType:               transactionType,
		CardNotPresentStatus:          cardNotPresentStatus,
		CNPLimitAmountPerDay:          cNPLimitAmountPerDay,
		CNPLimitAmountPerTransaction:  cNPLimitAmountPerTransaction,
		LimitStatus:                   limitStatus,
		LimitAmountPerDay:             limitAmountPerDay,
		LimitAmountPerTransaction:     limitAmountPerTransaction,
	}, nil
} 

// Converts UpdateSpendingControlRequest to a fixed-length string.
func FormatUpdateSpendingControlRequest(updateSpendingControlReq domain.UpdateSpendingControlRequest) string {
	TransactionType                := utils.PadOrTruncate(updateSpendingControlReq.TransactionType, 2)
	AEONID                         := utils.PadOrTruncate(updateSpendingControlReq.AEONID, 20)
	CardNo                         := utils.PadOrTruncate(updateSpendingControlReq.CardNo, 16)
	BusinessCode                   := utils.PadOrTruncate(updateSpendingControlReq.BusinessCode, 2)
	Channel                        := utils.PadOrTruncate(updateSpendingControlReq.Channel, 10)
	CardNotPresentStatus           := utils.PadOrTruncate(updateSpendingControlReq.CardNotPresentStatus, 1)
	CNPLimitAmountPerDay           := utils.PadIntWithZero(updateSpendingControlReq.CNPLimitAmountPerDay, 12)
	CNPLimitAmountPerTransaction   := utils.PadIntWithZero(updateSpendingControlReq.CNPLimitAmountPerTransaction, 12)
	LimitStatus                    := utils.PadOrTruncate(updateSpendingControlReq.LimitStatus, 1)
	LimitAmountPerDay              := utils.PadIntWithZero(updateSpendingControlReq.LimitAmountPerDay, 12)
	LimitAmountPerTransaction      := utils.PadIntWithZero(updateSpendingControlReq.LimitAmountPerTransaction, 12)
	Date                           := utils.PadOrTruncate(updateSpendingControlReq.Date, 8)
	Time                           := utils.PadOrTruncate(updateSpendingControlReq.Time, 6)

	return TransactionType + AEONID + CardNo + BusinessCode + Channel + CardNotPresentStatus + CNPLimitAmountPerDay + CNPLimitAmountPerTransaction + LimitStatus + LimitAmountPerDay + LimitAmountPerTransaction + Date + Time
}

func FormatUpdateSpendingControlResponse(raw string) (domain.UpdateSpendingControlResponse, error) {
	const headerLen = 123
	const dataLen = 318				// กรอกค่านี้ด้วย

	if len(raw) <= headerLen {
		return domain.UpdateSpendingControlResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.UpdateSpendingControlResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)				// กรอกค่านี้ด้วย

	transactionType             := parser.ReadString(0, 8)
	date                        := parser.ReadString(8, 6)
	time                        := parser.ReadString(14, 2)

	return domain.UpdateSpendingControlResponse{
		TransactionType:    transactionType,
		Date:               date,
		Time:               time,
	}, nil
}