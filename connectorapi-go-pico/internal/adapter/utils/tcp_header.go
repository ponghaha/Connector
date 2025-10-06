package utils

import (
	"strings"
	"time"
	"fmt"
	"math"
	//"bytes"
)

func PadOrTruncate(s string, length int) string {
	runes := []rune(s)
	if len(runes) > length {
		runes = runes[:length]
	}
	padded := string(runes)
	if len(runes) < length {
		padded = padded + strings.Repeat(" ", length-len(runes))
	}
	return padded
}

func PadIntWithZero(n int, length int) string {
	return fmt.Sprintf("%0*d", length, n)
}

func PadFloatWithZero(n float64, length int, decimalDigit int) string {
	multiplier := math.Pow(10, float64(decimalDigit))
	intValue := int(math.Round(n * multiplier))

	return fmt.Sprintf("%0*d", length, intValue)
}

// BuildFixedLengthHeader constructs the fixed-length header.
func BuildFixedLengthHeader(routeSystem, routeService, routeFormat, requestID string,  routeRequestLength string) string {
	now := time.Now()
	requestDate := now.Format("20060102")
	requestTime := now.Format("150405")

	responseCode := PadOrTruncate("", 6)
	responseMessage := PadOrTruncate("", 50)

	// const fixedHeaderLength = 10 + 15 + 3 + 20 + 8 + 6 + 5 + 6 + 50 // 123 characters
	// totalMessageLength := len(fixedLengthData)
	// requestLengthStr := fmt.Sprintf("%05d", totalMessageLength)

	header := PadOrTruncate(routeSystem, 10) +
		PadOrTruncate(routeService, 15) +
		PadOrTruncate(routeFormat, 3) +
		PadOrTruncate(requestID, 20) +
		PadOrTruncate(requestDate, 8) +
		PadOrTruncate(requestTime, 6) +
		PadOrTruncate(routeRequestLength, 5) +
		responseCode +
		responseMessage

	return header
}
