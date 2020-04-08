package def

import (
	"strconv"
	"strings"
)

func getBaseValue(val string) (string, int, int) {
	base := 10
	if strings.HasPrefix(val, "0x") {
		base = 16
		val = strings.TrimPrefix(val, "0x")
	} else if strings.HasPrefix(val, "0o") {
		base = 8
		val = strings.TrimPrefix(val, "0o")
	} else if strings.HasPrefix(val, "0b") {
		base = 2
		val = strings.TrimPrefix(val, "0b")
	}

	return val, base, 0
}

func getStringToComplex(r, i string) (complex128, error) {
	nR, err := strconv.ParseFloat(r, 0)
	if err != nil {
		return 0, err
	}
	nI, err := strconv.ParseFloat(i, 0)
	if err != nil {
		return 0, err
	}

	return complex(nR, nI), nil
}