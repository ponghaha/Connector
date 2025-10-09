package format

import (
	"fmt"
	// "strconv"
	"strings"
	// "time"
	//"bytes"

	"connectorapi-go/internal/core/domain"
	"connectorapi-go/internal/adapter/utils"
)

// Converts CollectionDetailRequest to a fixed-length string.
func FormatCollectionDetailRequest(reqData domain.CollectionDetailRequest) string {
	IDCardNo    := utils.PadOrTruncate(reqData.IDCardNo, 20)
	RedCaseNo   := utils.PadOrTruncate(reqData.RedCaseNo, 15)
	BlackCaseNo := utils.PadOrTruncate(reqData.BlackCaseNo, 15)
	return IDCardNo + RedCaseNo + BlackCaseNo
}

func FormatCollectionDetailResponse(raw string) (domain.CollectionDetailResponse, error) {
	const headerLen = 123
	if len(raw) <= headerLen {
		return domain.CollectionDetailResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}
	body := raw[headerLen:]
	runes := []rune(body)

	parser := utils.NewFixedParser(body)

	idCardNo := parser.ReadString(0, 20)
	noOfAgreement := parser.ReadInt(20, 2)

	const agreementLen = 942
	agreementStart := 22
	agreements := make([]domain.CollectionDetailAgreement, 0, noOfAgreement)

	for i := 0; i < noOfAgreement; i++ {
		start := agreementStart + i*agreementLen
		if start >= len(runes) {
			break
		}
		blockRunes := utils.ReadRunes(runes, start, agreementLen)

		agreements = append(agreements, domain.CollectionDetailAgreement{
			AgreementNo:              utils.ReadBlockStr(blockRunes, 0, 16),
			SeqOfAgreement:           utils.ReadBlockInt(blockRunes, 16, 2),
			OutsourceID:              utils.ReadBlockStr(blockRunes, 18, 4),
			OutsourceName:            utils.ReadBlockStr(blockRunes, 2, 30),
			BlockCode:                utils.ReadBlockStr(blockRunes, 52, 2),
			CurrentSUEOSPrincipalNet: utils.ReadBlockFloat100ToDecimal(blockRunes, 54, 10),
			CurrentSUEOSPrincipalVAT: utils.ReadBlockFloat100ToDecimal(blockRunes, 64, 10),
			CurrentSUEOSInterestNet:  utils.ReadBlockFloat100ToDecimal(blockRunes, 74, 10),
			CurrentSUEOSInterestVAT:  utils.ReadBlockFloat100ToDecimal(blockRunes, 84, 10),
			CurrentSUEOSPenalty:      utils.ReadBlockFloat100ToDecimal(blockRunes, 94, 9),
			CurrentSUEOSHDCharge:     utils.ReadBlockFloat100ToDecimal(blockRunes, 103, 9),
			CurrentSUEOSOtherFee:     utils.ReadBlockFloat100ToDecimal(blockRunes, 112, 9),
			CurrentSUEOSTotal:        utils.ReadBlockFloat100ToDecimal(blockRunes, 121, 10),
			TotalPaymentAmount:       utils.ReadBlockFloat100ToDecimal(blockRunes, 131, 10),
			LastPaymentDate:          utils.ReadBlockInt(blockRunes, 141, 8),
			SUESeqNo:                 utils.ReadBlockInt(blockRunes, 149, 2),
			BeginSUEOSPrincipalNet:   utils.ReadBlockFloat100ToDecimal(blockRunes, 151, 10),
			BeginSUEOSPrincipalVAT:   utils.ReadBlockFloat100ToDecimal(blockRunes, 161, 10),
			BeginSUEOSInterestNet:    utils.ReadBlockFloat100ToDecimal(blockRunes, 171, 10),
			BeginSUEOSInterestVAT:    utils.ReadBlockFloat100ToDecimal(blockRunes, 181, 10),
			BeginSUEOSPenalty:        utils.ReadBlockFloat100ToDecimal(blockRunes, 191, 10),
			BeginSUEOSHDCharge:       utils.ReadBlockFloat100ToDecimal(blockRunes, 201, 9),
			BeginSUEOSOtherFee:       utils.ReadBlockFloat100ToDecimal(blockRunes, 210, 9),
			BeginSUEOSTotal:          utils.ReadBlockFloat100ToDecimal(blockRunes, 219, 10),
			SUEStatus:                utils.ReadBlockInt(blockRunes, 229, 2),
			SUEStatusDescription:     utils.ReadBlockStr(blockRunes, 231, 30),
			BlackCaseNo:              utils.ReadBlockStr(blockRunes, 261, 15),
			BlackCaseDate:            utils.ReadBlockInt(blockRunes, 276, 8),
			RedCaseNo:                utils.ReadBlockStr(blockRunes, 284, 15),
			RedCaseDate:              utils.ReadBlockInt(blockRunes, 299, 8),
			CourtCode:                utils.ReadBlockStr(blockRunes, 307, 4),
			CourtName:                utils.ReadBlockStr(blockRunes, 311, 30),
			JudgmentDate:             utils.ReadBlockInt(blockRunes, 341, 8),
			JudgmentResultCode:       utils.ReadBlockInt(blockRunes, 349, 1),
			JudgmentResultDescription: utils.ReadBlockStr(blockRunes, 350, 40),
			JudgmentDetail:           utils.ReadBlockStr(blockRunes, 390, 500),
			ExpectDate:               utils.ReadBlockInt(blockRunes, 890, 8),
			AssetPrice:               utils.ReadBlockFloat100ToDecimal(blockRunes, 898, 10),
			JudgeAmount:              utils.ReadBlockFloat100ToDecimal(blockRunes, 908, 10),
			NoOfInstallment:          utils.ReadBlockStr(blockRunes, 918, 3),
			InstallmentAmount:        utils.ReadBlockFloat100ToDecimal(blockRunes, 921, 10),
			TotalCurrentPerSUESeqNo:  utils.ReadBlockFloat100ToDecimal(blockRunes, 931, 11),
		})
	}

	return domain.CollectionDetailResponse{
		IDCardNo:      idCardNo,
		NoOfAgreement: noOfAgreement,
		AgreementList: agreements,
	}, nil
}


// Converts CollectionLogRequest to a fixed-length string.
func FormatCollectionLogRequest(reqData domain.CollectionLogRequest) string {
	AgreementNo := utils.PadOrTruncate(reqData.AgreementNo, 16)
	RemarkCode  := utils.PadOrTruncate(reqData.RemarkCode, 4)
	LogRemark1  := utils.PadOrTruncate(reqData.LogRemark1, 120)
	LogRemark2  := utils.PadOrTruncate(reqData.LogRemark2, 120)
	LogRemark3  := utils.PadOrTruncate(reqData.LogRemark3, 120)
	LogRemark4  := utils.PadOrTruncate(reqData.LogRemark4, 120)
	LogRemark5  := utils.PadOrTruncate(reqData.LogRemark5, 120)
	InputDate   := utils.PadOrTruncate(reqData.InputDate, 8)
	InputTime   := utils.PadOrTruncate(reqData.InputTime, 6)
	OperatorID  := utils.PadOrTruncate(reqData.OperatorID, 15)
	return AgreementNo + RemarkCode + LogRemark1 + LogRemark2 + LogRemark3 + LogRemark4 + LogRemark5 + InputDate + InputTime + OperatorID
}

func FormatCollectionLogResponse(raw string) (domain.CollectionLogResponse, error) {
	const headerLen = 123
	const dataLen = 36 

	if len(raw) <= headerLen {
		return domain.CollectionLogResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.CollectionLogResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	idCardNo := strings.TrimSpace(data[:20])
	agreementNo := strings.TrimSpace(data[20:36])

	return domain.CollectionLogResponse{
		IDCardNo:    idCardNo,
		AgreementNo: agreementNo,
	}, nil
}