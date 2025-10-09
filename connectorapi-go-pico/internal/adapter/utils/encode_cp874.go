package utils

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// UTF8 to CP874 mapping (only for Thai characters ก - ฮ and common ASCII)
func Utf8ToCP874(input string) ([]byte, error) {
	var result []byte
	for _, r := range input {
		// Handle newline separately
		if r == '\n' {
			result = append(result, 0x0A)
			continue
		}
		if r == '\r' {
			result = append(result, 0x0D)
			continue
		}

		// ASCII range: pass through
		if r <= 0x7F {
			result = append(result, byte(r))
			continue
		}

		// Thai range ก - ฮ (U+0E01 to U+0E5B)
		if r >= 0x0E01 && r <= 0x0E5B {
			cp874Byte := byte(r - 0x0E01 + 0xA1)
			result = append(result, cp874Byte)
			continue
		}

		// Unsupported character
		return nil, fmt.Errorf("cannot encode rune '%c' (U+%04X) to CP874", r, r)
	}
	return result, nil
}

func DecodeCP874(input []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(input), charmap.Windows874.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
