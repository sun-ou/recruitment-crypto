package wallet

import (
	"fmt"
	"strconv"
)

func String2Cent(param string) uint {
	f, _ := strconv.ParseFloat(param, 64)
	return uint(f * 100)
}

func Cent2String(param uint) string {
	return fmt.Sprintf("%.2f", float64(param)/100)
}
