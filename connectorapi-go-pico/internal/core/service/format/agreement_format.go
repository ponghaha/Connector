package format

import (
	"fmt"
	// "strconv"
	"strings"
	//"bytes"

	"connectorapi-go/internal/core/domain"
	"connectorapi-go/internal/adapter/utils"
)

// Converts UpdateStatusRequest to a fixed-length string.
func FormatUpdateStatusRequest(updateStatusReq domain.UpdateStatusRequest) string {
	aeonID    := utils.PadOrTruncate(updateStatusReq.AeonID, 20)
	agreement := utils.PadOrTruncate(updateStatusReq.Agreement, 12)
	status    := utils.PadOrTruncate(updateStatusReq.Status, 1)
	return aeonID + agreement + status
}

func FormatUpdateStatusResponse(raw string) (domain.UpdateStatusResponse, error) {
	const headerLen = 123
	const dataLen = 32

	if len(raw) <= headerLen {
		return domain.UpdateStatusResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.UpdateStatusResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	aeonID      := strings.TrimSpace(data[:20])
	agreementNo := strings.TrimSpace(data[20:32])
	// agreementNoInt, err := strconv.Atoi(agreementNo)
	// if err != nil {
	// 	fmt.Println("connot convert string to int:", err)
	// }

	return domain.UpdateStatusResponse{
		AeonID:    aeonID,
		Agreement: agreementNo,
	}, nil
} 

// Converts AgreeMentBillingRequest to a fixed-length string.
func FormatAgreeMentBillingRequest(AgreeMentBillingReq domain.AgreeMentBillingRequest) string {
	idCardNo    := utils.PadOrTruncate(AgreeMentBillingReq.IDCardNo, 20)
	agreementNo := utils.PadOrTruncate(AgreeMentBillingReq.AgreementNo, 16)
	cardCode    := utils.PadOrTruncate(AgreeMentBillingReq.CardCode, 2)
	return idCardNo + agreementNo + cardCode
}

func FormatAgreeMentBillingResponse(raw string) (domain.AgreeMentBillingResponse, error) {
	const headerLen = 123
	const dataLen = 168

	if len(raw) <= headerLen {
		return domain.AgreeMentBillingResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.AgreeMentBillingResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	dueDate                          := parser.ReadInt(0,8)
	settlementDate                   := parser.ReadInt(8,8)
	billingAmount                    := parser.ReadFloat100(16,11)
	minPaymentAmount                 := parser.ReadFloat100(27,11)
	fullPaymentAmount                := parser.ReadFloat100(38,11)
	unbilledAmount                   := parser.ReadFloat100(49,11)
	creditShoppingFloorLimit         := parser.ReadFloat100(60,11)
	creditShoppingOutstanding        := parser.ReadFloat100(71,11)
	creditShoppingAvailableLimit     := parser.ReadFloat100(82,11)
	creditCashingFloorLimit          := parser.ReadFloat100(93,11)
	creditCashingOutstanding         := parser.ReadFloat100(104,11)
	creditCashingAvailableLimit      := parser.ReadFloat100(115,11)
	installmentNo                    := parser.ReadInt(126,3)
	installmentCurrent               := parser.ReadInt(129,3)
	paymentHistory                   := parser.ReadString(132,36)

	return domain.AgreeMentBillingResponse{
		DueDate:                      dueDate,
		SettlementDate:               settlementDate,
		BillingAmount:                utils.DecimalString(billingAmount),
		MinPaymentAmount:             utils.DecimalString(minPaymentAmount),
		FullPaymentAmount:            utils.DecimalString(fullPaymentAmount),
		UnbilledAmount:               utils.DecimalString(unbilledAmount),
		CreditShoppingFloorLimit:     utils.DecimalString(creditShoppingFloorLimit),
		CreditShoppingOutstanding:    utils.DecimalString(creditShoppingOutstanding),
		CreditShoppingAvailableLimit: utils.DecimalString(creditShoppingAvailableLimit),
		CreditCashingFloorLimit:      utils.DecimalString(creditCashingFloorLimit),
		CreditCashingOutstanding:     utils.DecimalString(creditCashingOutstanding),
		CreditCashingAvailableLimit:  utils.DecimalString(creditCashingAvailableLimit),
		InstallmentNo:                installmentNo,
		InstallmentCurrent:           installmentCurrent,
		PaymentHistory:               paymentHistory,
	}, nil
} 