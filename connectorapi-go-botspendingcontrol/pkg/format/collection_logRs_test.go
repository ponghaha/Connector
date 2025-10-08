package format

import (
	"strings"
	"testing"
)

// ----------------- Unit Test -----------------

func TestFormatCollectionLogResponse(t *testing.T) {
	header := strings.Repeat("H", 123)
	// body ต้องมีอย่างน้อย 36 ตัว (20 + 16)
	body := "12345678901234567890ABCDEFGHIJKLMNO1" // 36 ตัว
	longBody := body + "EXTRA"                     // มากกว่า 36 ตัว
	trimBody := "  1234567890  1234567890123456  " // มีช่องว่าง 36 ตัว

	tests := []struct {
		name      string
		raw       string
		wantID    string
		wantAg    string
		wantError bool
	}{
		{
			name:      "raw too short for header",
			raw:       "short",
			wantError: true,
		},
		{
			name:      "raw too short for body",
			raw:       header + "shortbody",
			wantError: true,
		},
		{
			name:      "exact length",
			raw:       header + body,
			wantID:    "12345678901234567890",
			wantAg:    "ABCDEFGHIJKLMNO1",
			wantError: false,
		},
		{
			name:      "long body",
			raw:       header + longBody,
			wantID:    "12345678901234567890",
			wantAg:    "ABCDEFGHIJKLMNO1",
			wantError: false,
		},
		{
			name:      "trim spaces",
			raw:       header + trimBody,
			wantID:    "1234567890",
			wantAg:    "1234567890123456",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatCollectionLogResponse(tt.raw)
			if (err != nil) != tt.wantError {
				t.Errorf("unexpected error status: got %v, want error %v", err, tt.wantError)
			}
			if !tt.wantError {
				if got.IDCardNo != tt.wantID {
					t.Errorf("IDCardNo = %q, want %q", got.IDCardNo, tt.wantID)
				}
				if got.AgreementNo != tt.wantAg {
					t.Errorf("AgreementNo = %q, want %q", got.AgreementNo, tt.wantAg)
				}
			}
		})
	}
}
