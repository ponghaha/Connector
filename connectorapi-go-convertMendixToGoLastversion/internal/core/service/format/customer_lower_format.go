package format

import (
	"fmt"
	// "strconv"
	// "strings"
	//"bytes"

	"connectorapi-go/internal/core/domain" 
	"connectorapi-go/internal/adapter/utils"
)

// Converts GetCustomerInfoMobileNoRequest to a fixed-length string.
func FormatGetCustomerInfoMobileNoRequest(getCustomerInfoMobileNoReq domain.GetCustomerInfoMobileNoRequest) string {
	Mobileno        := utils.PadOrTruncate(getCustomerInfoMobileNoReq.Mobileno, 20)
	return Mobileno
}

func FormatGetCustomerInfoMobileNoResponse(raw string) (domain.GetCustomerInfoMobileNoResponse, error) {
	const headerLen = 123
	const dataLen = 26

	if len(raw) <= headerLen {
		return domain.GetCustomerInfoMobileNoResponse{}, fmt.Errorf("raw data too short for header, length=%d", len(raw))
	}

	data := raw[headerLen:]
	if len(data) < dataLen {
		return domain.GetCustomerInfoMobileNoResponse{}, fmt.Errorf("raw data too short for body, length=%d, need %d", len(data), dataLen)
	}

	parser := utils.NewFixedParser(data)

	setmobileno                 := parser.ReadString(0, 20)
	setresultcode               := parser.ReadString(20, 2)
	setmobileappflag            := parser.ReadString(22, 1)
	setvvipflag                 := parser.ReadString(23, 1)
	setsweetheartflag           := parser.ReadString(24, 1)
	setfraudflag                := parser.ReadString(25, 1)

	return domain.GetCustomerInfoMobileNoResponse{
		Mobileno:               setmobileno,
		Resultcode:             setresultcode,
		Mobileappflag:          setmobileappflag,
		Vvipflag:               setvvipflag,
		Sweetheartflag:         setsweetheartflag,
		Fraudflag:              setfraudflag,
	}, nil
} 