package format

import (
	"fmt"
	"strings"
	"testing"
)

func generateMockAgreementBlock() string {
	block := ""
	block += fmt.Sprintf("%016s", "1234567891234567")          // AgreementNo
	block += fmt.Sprintf("%02d", 1)                            // SeqOfAgreement
	block += fmt.Sprintf("%-4s", "")                           // OutsourceID
	block += fmt.Sprintf("%-30s", "")                          // OutsourceName
	block += fmt.Sprintf("%-2s", "B1")                         // BlockCode
	block += fmt.Sprintf("%010.2f", 12345.67)                 // CurrentSUEOSPrincipalNet
	block += fmt.Sprintf("%010.2f", 123.45)                   // CurrentSUEOSPrincipalVAT
	block += fmt.Sprintf("%010.2f", 1000.12)                  // CurrentSUEOSInterestNet
	block += fmt.Sprintf("%010.2f", 50.00)                    // CurrentSUEOSInterestVAT
	block += fmt.Sprintf("%09.2f", 5.50)                      // CurrentSUEOSPenalty
	block += fmt.Sprintf("%09.2f", 2.00)                      // CurrentSUEOSHDCharge
	block += fmt.Sprintf("%09.2f", 1.00)                      // CurrentSUEOSOtherFee
	block += fmt.Sprintf("%010.2f", 57.67)                    // CurrentSUEOSTotal
	block += fmt.Sprintf("%010.2f", 10000.00)                 // TotalPaymentAmount
	block += fmt.Sprintf("%08d", 20250820)                    // LastPaymentDate
	block += fmt.Sprintf("%02d", 1)                            // SUESeqNo
	block += fmt.Sprintf("%010.2f", 5000.00)                  // BeginSUEOSPrincipalNet
	block += fmt.Sprintf("%010.2f", 250.00)                   // BeginSUEOSPrincipalVAT
	block += fmt.Sprintf("%010.2f", 2000.00)                  // BeginSUEOSInterestNet
	block += fmt.Sprintf("%010.2f", 100.00)                   // BeginSUEOSInterestVAT
	block += fmt.Sprintf("%010.2f", 50.00)                    // BeginSUEOSPenalty
	block += fmt.Sprintf("%09.2f", 2.50)                      // BeginSUEOSHDCharge
	block += fmt.Sprintf("%09.2f", 1.00)                      // BeginSUEOSOtherFee
	block += fmt.Sprintf("%010.2f", 7153.50)                  // BeginSUEOSTotal
	block += fmt.Sprintf("%02d", 1)                            // SUEStatus
	block += fmt.Sprintf("%-30s", "SUE Status Description")   // SUEStatusDescription
	block += fmt.Sprintf("%-15s", "BLACK123")                 // BlackCaseNo
	block += fmt.Sprintf("%08d", 20250720)                    // BlackCaseDate
	block += fmt.Sprintf("%-15s", "RED123")                   // RedCaseNo
	block += fmt.Sprintf("%08d", 20250721)                    // RedCaseDate
	block += fmt.Sprintf("%-4s", "C001")                      // CourtCode
	block += fmt.Sprintf("%-30s", "Court Name")               // CourtName
	block += fmt.Sprintf("%08d", 20250722)                    // JudgmentDate
	block += fmt.Sprintf("%01d", 1)                            // JudgmentResultCode
	block += fmt.Sprintf("%-40s", "Judgment Result Desc")     // JudgmentResultDescription
	block += fmt.Sprintf("%-500s", "Judgment Detail")         // JudgmentDetail
	block += fmt.Sprintf("%08d", 20250830)                    // ExpectDate
	block += fmt.Sprintf("%010.2f", 100000.00)               // AssetPrice
	block += fmt.Sprintf("%010.2f", 50000.00)                // JudgeAmount
	block += fmt.Sprintf("%-3s", "12")                        // NoOfInstallment
	block += fmt.Sprintf("%010.2f", 8333.33)                 // InstallmentAmount
	block += fmt.Sprintf("%011.2f", 25000.00)                // TotalCurrentPerSUESeqNo

	if len(block) < 924 {
		block += strings.Repeat(" ", 924-len(block))
	}
	return block
}

func generateMockData(numAgreements int) string {
	header := strings.Repeat(" ", 123)
	idCardNo := fmt.Sprintf("%-20s", "1234567890123456789")
	noOfAgreement := fmt.Sprintf("%02d", numAgreements)
	body := ""
	for i := 0; i < numAgreements; i++ {
		body += generateMockAgreementBlock()
	}
	return header + idCardNo + noOfAgreement + body
}

// ----------------- Unit Test -----------------

func TestFormatCollectionDetailResponse_AllFieldsWithLog(t *testing.T) {
	numAgreements := 4
	data := generateMockData(numAgreements)

	resp, err := FormatCollectionDetailResponse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("IDCardNo: '%s'", resp.IDCardNo)
	t.Logf("NoOfAgreement: %d", resp.NoOfAgreement)
	t.Logf("Number of parsed agreements: %d", len(resp.AgreementList))

	for i, ag := range resp.AgreementList {
		t.Logf("Agreement %d:", i+1)
		t.Logf("  AgreementNo: '%s'", ag.AgreementNo)
		t.Logf("  SeqOfAgreement: %d", ag.SeqOfAgreement)
		t.Logf("  OutsourceID: '%s'", ag.OutsourceID)
		t.Logf("  CurrentSUEOSTotal: %f", ag.CurrentSUEOSTotal)
		t.Logf("  SUEStatusDescription: '%s'", ag.SUEStatusDescription)
	}

	if len(resp.AgreementList) != numAgreements {
		t.Errorf("Expected %d agreements, got %d", numAgreements, len(resp.AgreementList))
	}
}
