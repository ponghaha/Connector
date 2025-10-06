package format

import (
	"fmt"
	"strconv"
	// "strings"
	"time"
	//"bytes"

	"connectorapi-go/internal/core/domain"
	"connectorapi-go/internal/adapter/utils"
)

// Converts MyCardRequest to a fixed-length string.
func FormatMyCardRequestNormal(myCardReq domain.MyCardRequest) string {
	IDCardNo        		:= utils.PadOrTruncate(myCardReq.UserRef, 20)
	CreditCardNo            := "                "
	BusinessCode            := "  "
	return IDCardNo + CreditCardNo + BusinessCode
}

func FormatMyCardRequestAll(myCardReq domain.MyCardRequest) string {
	IDCardNo        		:= utils.PadOrTruncate(myCardReq.UserRef, 20)
	CustomerNameEN          := "Y"
	CustomerNameTH          := "Y"
	return IDCardNo + CustomerNameEN + CustomerNameTH
}

func FormatMyCardResponseNormal(raw string) (domain.MyCardResponseNormal, error) {
	const headerLen = 123
	currentDate := time.Now().Format("20060102")
	currentDateInt, _ := strconv.Atoi(currentDate)

	if len(raw) <= headerLen {
		return domain.MyCardResponseNormal{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}
	body := raw[headerLen:]
	runes := []rune(body)

	parser := utils.NewFixedParser(body)

	idCardNo                  := parser.ReadString(0, 20)
	totalCreditCard           := parser.ReadInt(20, 4)

	const agreementLen = 61
	agreementStart := 24
	agreements := make([]domain.MyCardListNormal, 0, totalCreditCard)

	for i := 0; i < totalCreditCard; i++ {
		start := agreementStart + i*agreementLen
		if start >= len(runes) {
			break
		}
		blockRunes := utils.ReadRunes(runes, start, agreementLen)

		rawCreditCardNo := utils.ReadBlockStr(blockRunes, 0, 16)
		cardStatus := utils.ReadBlockStr(blockRunes, 50, 2)
		expireDateStr := utils.ReadBlockStr(blockRunes, 52, 8)
		digitalCardFlag := utils.ReadBlockStr(blockRunes, 60, 1)

		maskedCreditCardNo := rawCreditCardNo
		if len(rawCreditCardNo) == 16 {
			maskedCreditCardNo = rawCreditCardNo[0:6] + "XXXXXX" + rawCreditCardNo[12:16]
		}

		status := "HLD"
		if cardStatus == "00" || cardStatus == "II" {
			status = "ACT"
		} else {
			expireDateInt, _ := strconv.Atoi(expireDateStr)
			if expireDateInt < currentDateInt {
				status = "EXP"
			}
		}

		finalDigitalCardFlag := digitalCardFlag
		if digitalCardFlag == "" {
			finalDigitalCardFlag = "N"
		}


		agreements = append(agreements, domain.MyCardListNormal{
			CreditCardNo:        maskedCreditCardNo,
			CardName:            utils.ReadBlockStr(blockRunes, 16, 30),
			ProductType:         utils.ReadBlockStr(blockRunes, 46, 2),
			BusinessCode:        utils.ReadBlockStr(blockRunes, 48, 2),
			CardStatus:          status,
			ExpireDate:          utils.ReadBlockStr(blockRunes, 52, 8),
			// DYCA:                readBlockStr(60, 1),
			DigitalCardFlag:     finalDigitalCardFlag,
		})
	}

	return domain.MyCardResponseNormal{
		IDCardNo:                 idCardNo,
		TotalCreditCard:          totalCreditCard,
		CardList:                 agreements,
	}, nil
}

func FormatMyCardResponseAll(raw string) (domain.MyCardResponseAll, error) {
	const headerLen = 123
	if len(raw) <= headerLen {
		return domain.MyCardResponseAll{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}
	body := raw[headerLen:]
	runes := []rune(body)

	parser := utils.NewFixedParser(body)

	idCardNo               := parser.ReadString(0, 20)
	customerNameEN         := parser.ReadString(20, 30)
	customerNameTH         := parser.ReadString(50, 30)
	totalCreditCard        := parser.ReadInt(80, 3)

	const agreementLen = 68
	agreementStart := 83
	agreements := make([]domain.MyCardListAll, 0, totalCreditCard)

	for i := 0; i < totalCreditCard; i++ {
	start := agreementStart + i*agreementLen
	if start >= len(runes) {
		break
	}
	blockRunes := utils.ReadRunes(runes, start, agreementLen)

	// Raw values
	rawCreditCardNo := utils.ReadBlockStr(blockRunes, 0, 16)

	maskedCreditCardNo := rawCreditCardNo
	if len(rawCreditCardNo) == 16 {
		maskedCreditCardNo = rawCreditCardNo[0:6] + "XXXXXX" + rawCreditCardNo[12:16]
	}

	agreements = append(agreements, domain.MyCardListAll{
		CreditCardNo:        maskedCreditCardNo,
		CardCode:            utils.ReadBlockStr(blockRunes, 16, 2),
		ProductType:         utils.ReadBlockStr(blockRunes, 18, 2),
		CardType:            utils.ReadBlockStr(blockRunes, 20, 1),
		CardStatus:          utils.ReadBlockStr(blockRunes, 21, 1),
		ExpireDate:          utils.ReadBlockInt(blockRunes, 22, 8),
		HoldCode:            utils.ReadBlockStr(blockRunes, 30, 2),
		RetreatCode:         utils.ReadBlockStr(blockRunes, 32, 1),
		SendMode:            utils.ReadBlockStr(blockRunes, 33, 1),
		FirstEmbossDate:     utils.ReadBlockInt(blockRunes, 34, 8),
		FirstConfirmDate:    utils.ReadBlockInt(blockRunes, 42, 8),
		ShoppingLimit:       utils.ReadBlockInt(blockRunes, 50, 9),
		CashingLimit:        utils.ReadBlockInt(blockRunes, 59, 9),
	})
}

	return domain.MyCardResponseAll{
		IDCardNo:                 idCardNo,
		CustomerNameEN:           customerNameEN,
		CustomerNameTH:           customerNameTH,
		TotalCreditCard:          totalCreditCard,
		CardList:                 agreements,
	}, nil
}

// Converts GetAvailableLimitRequest to a fixed-length string.
func FormatGetAvailableLimitRequest(getAvailableLimitReq domain.GetAvailableLimitRequest) string {
	IDCardNo    := utils.PadOrTruncate(getAvailableLimitReq.IDCardNo, 20)
	CreditCardNo := utils.PadOrTruncate(getAvailableLimitReq.CreditCardNo, 16)
	BusinessCode    := utils.PadOrTruncate(getAvailableLimitReq.BusinessCode, 2)
	return IDCardNo + CreditCardNo + BusinessCode
}

func FormatGetAvailableLimitResponse(raw string) (domain.GetAvailableLimitResponse, error) {
	const headerLen = 123
	const dataLen = 44

	if len(raw) <= headerLen {
		return domain.GetAvailableLimitResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetAvailableLimitResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	shoppingLimit                   := parser.ReadFloat100(0,11)
	shoppingAvailable               := parser.ReadFloat100(11,11)
	cashingLimit                    := parser.ReadFloat100(22,11)
	cashingAvailable                := parser.ReadFloat100(33,11)

	return domain.GetAvailableLimitResponse{
		ShoppingLimit:               shoppingLimit,
		ShoppingAvailable:           shoppingAvailable,
		CashingLimit:                cashingLimit,
		CashingAvailable:            cashingAvailable,
	}, nil
}