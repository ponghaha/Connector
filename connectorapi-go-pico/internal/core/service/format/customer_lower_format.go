package format

import (
	"fmt"
	// "strconv"
	// "strings"
	//"bytes"

	"connectorapi-go/internal/core/domain" 
	"connectorapi-go/internal/adapter/utils"
)

// Converts CheckAeonCustomerRequest to a fixed-length string.
func FormatCheckAeonCustomerRequest(checkAeonCustomerReq domain.CheckAeonCustomerRequest) string {
	CustomerID        := utils.PadOrTruncate(checkAeonCustomerReq.CustomerID, 20)
	return CustomerID
}

func FormatCheckAeonCustomerResponse(raw string) (domain.CheckAeonCustomerResponse, error) {
	const headerLen = 123
	const dataLen = 74

	if len(raw) <= headerLen {
		return domain.CheckAeonCustomerResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.CheckAeonCustomerResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	customerID                 := parser.ReadString(0, 20)
	aeonMember                 := parser.ReadString(20, 1)
	resultCode                 := parser.ReadString(21, 1)
	reasonCode                 := parser.ReadInt(22, 2)
	reasonDescription          := parser.ReadString(24, 50)

	return domain.CheckAeonCustomerResponse{
		CustomerID:             customerID,
		AeonMember:             aeonMember,
		ResultCode:             resultCode,
		ReasonCode:             reasonCode,
		ReasonDescription:      reasonDescription,
	}, nil
} 