package format

import (
	"fmt"
	// "strconv"
	"strings"
	//"bytes"

	"connectorapi-go/internal/core/domain"
	"connectorapi-go/internal/adapter/utils"
)

// Converts UpdateConsentRequest to a fixed-length string.
func FormatUpdateConsentRequest(req domain.UpdateConsentRequest) string {
	var builder strings.Builder

	builder.WriteString(utils.PadOrTruncate(req.IDCardNo, 20))
	builder.WriteString(utils.PadOrTruncate(req.ActionChannel, 3))
	builder.WriteString(utils.PadOrTruncate(req.ActionDateTime, 14))
	builder.WriteString(utils.PadOrTruncate(req.ApplicationNo, 20))
	builder.WriteString(utils.PadOrTruncate(req.ApplicationVersion, 13))
	builder.WriteString(utils.PadOrTruncate(req.IPAddress, 50))
	builder.WriteString(utils.PadOrTruncate(req.ATMNo, 5))
	builder.WriteString(utils.PadOrTruncate(req.BranchCode, 4))
	builder.WriteString(utils.PadOrTruncate(req.VoicePath, 150))
	builder.WriteString(utils.PadIntWithZero(req.TotalOfConsentCode, 2))

	for _, item := range req.ConsentLists {
		builder.WriteString(utils.PadOrTruncate(item.ConsentForm, 3))
		builder.WriteString(utils.PadOrTruncate(item.ConsentCode, 3))
		builder.WriteString(utils.PadOrTruncate(item.ConsentFormVersion, 13))
		builder.WriteString(utils.PadOrTruncate(item.ConsentLanguage, 1))
		builder.WriteString(utils.PadOrTruncate(item.ConsentStatus, 2))
	}

	return builder.String()
}

func FormatUpdateConsentResponse(raw string) (domain.UpdateConsentResponse, error) {
	const headerLen = 123
	const dataLen = 122

	if len(raw) <= headerLen {
		return domain.UpdateConsentResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.UpdateConsentResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	// iDCardNo                 := parser.ReadString(0,20)
	// applicationNo            := parser.ReadString(20,20)
	status                   := parser.ReadString(40,2)
	// filler                   := parser.ReadString(42,80)

	var Finalstatus string
	switch status {
		case "00":
			Finalstatus = "C"
		default:
			Finalstatus = "N"
	}

	return domain.UpdateConsentResponse{
		Status:                    Finalstatus,
	}, nil
}

// Converts GetConsentListRequest to a fixed-length string.
func FormatGetConsentListRequest(getConsentListReq domain.GetConsentListRequest) string {
	IDCardNo        := utils.PadOrTruncate(getConsentListReq.IDCardNo, 20)
	ConsentCode     := utils.PadOrTruncate(getConsentListReq.ConsentCode, 3)
	Filler          := utils.PadOrTruncate(getConsentListReq.Filler, 100)

	return IDCardNo + ConsentCode + Filler
}

func FormatGetConsentListResponse(raw string) (domain.GetConsentListResponse, error) {
	const headerLen = 123
	const recordLen = 88
	const agreementStart = 122

	if len(raw) <= headerLen {
		return domain.GetConsentListResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	body := raw[headerLen:]
	runes := []rune(body)

	parser := utils.NewFixedParser(body)
	idCardNo := parser.ReadString(0, 20)
	numberOfConsent := parser.ReadInt(120, 2)

	formMap := map[string]*domain.ConsentListObject{}
	formOrder := []string{}

	for i := 0; i < numberOfConsent; i++ {
		start := agreementStart + i * recordLen
		end := start + recordLen
		if end > len(runes) {
			break
		}
		block := utils.ReadRunes(runes, start, recordLen)

		form := strings.TrimSpace(utils.ReadBlockStr(block, 0, 3))
		code := strings.TrimSpace(utils.ReadBlockStr(block, 3, 3))
		version := strings.TrimSpace(utils.ReadBlockStr(block, 6, 13))
		appNo := strings.TrimSpace(utils.ReadBlockStr(block, 19, 20))
		appVer := strings.TrimSpace(utils.ReadBlockStr(block, 39, 12))
		lastStatus := strings.TrimSpace(utils.ReadBlockStr(block, 52, 2))
		acceptDT := strings.TrimSpace(utils.ReadBlockStr(block, 55, 12))
		acceptCH := strings.TrimSpace(utils.ReadBlockStr(block, 68, 2))
		cancelDT := strings.TrimSpace(utils.ReadBlockStr(block, 71, 13))
		cancelCH := strings.TrimSpace(utils.ReadBlockStr(block, 85, 2))


		detail := domain.ConsentDetailObject{
			ConsentCode:        code,
			ApplicationNo:      appNo,
			ApplicationVersion: appVer,
			LastConsentStatus:  lastStatus,
			AcceptDateTime:     acceptDT,
			AcceptChannel:      acceptCH,
			CancelDateTime:     cancelDT,
			CancelChannel:      cancelCH,
		}

		if _, exists := formMap[form]; !exists {
			formMap[form] = &domain.ConsentListObject{
				ConsentForm:        form,
				ConsentFormVersion: version,
				ConsentName:        "", 
				URL:                "", 
				ConsentDetails:     []domain.ConsentDetailObject{},
			}
			formOrder = append(formOrder, form)
		}

		formMap[form].ConsentDetails = append(formMap[form].ConsentDetails, detail)
	}

	consentList := make([]domain.ConsentListObject, 0, len(formMap))
	for _, form := range formOrder {
		item := formMap[form]
		item.TotalOfConsentCode = len(item.ConsentDetails)
		consentList = append(consentList, *item)
	}

	return domain.GetConsentListResponse{
		IDCardNo:        idCardNo,
		NumberOfConsent: numberOfConsent,
		ConsentList:     consentList,
	}, nil
}
