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