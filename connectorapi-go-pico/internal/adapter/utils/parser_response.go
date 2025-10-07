package utils

import (
	"strconv"
	"strings"
	"fmt"
)

type FixedParser struct {
	data []rune
}

func NewFixedParser(data string) FixedParser {
	return FixedParser{
		data: []rune(data),
	}
}

type DecimalString float64

func (d DecimalString) MarshalJSON() ([]byte, error) {
    s := fmt.Sprintf("%.2f", d)    
    return []byte(s), nil           
}


func (p FixedParser) ReadString(start, length int) string {
	if start >= len(p.data) {
		return ""
	}
	end := start + length
	if end > len(p.data) {
		end = len(p.data)
	}
	return strings.TrimSpace(string(p.data[start:end]))
}

func (p FixedParser) ReadInt(start, length int) int {
	s := p.ReadString(start, length)
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}

func (p FixedParser) ReadFloat100(start, length int) float64 {
	s := p.ReadString(start, length)
	if s == "" {
		return 0.0
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return float64(i) / 100.0
}

func ReadRunes(runes []rune, start int, length int) []rune {
	if start >= len(runes) {
		return []rune{}
	}
	end := start + length
	if end > len(runes) {
		end = len(runes)
	}
	return runes[start:end]
}

func ReadBlockStr(blockRunes []rune, start int, length int) string {
	if start >= len(blockRunes) {
		return ""
	}
	end := start + length
	if end > len(blockRunes) {
		end = len(blockRunes)
	}
	return strings.TrimSpace(string(blockRunes[start:end]))
}

func  ReadBlockInt(blockRunes []rune, start int, length int) int {
	s := ReadBlockStr(blockRunes, start, length)
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}
		
func ReadBlockFloat100ToDecimal(blockRunes []rune, start int, length int) DecimalString {
    s := ReadBlockStr(blockRunes, start, length)
    i, _ := strconv.ParseInt(s, 10, 64)
    f := float64(i) / 100.0
    return DecimalString(f)
}

func ConvertStringToInt(s string) int {
    i, err := strconv.Atoi(s)
    if err != nil {
        return 0
    }
    return i
}
