package formatutil

import (
	"fmt"
	"strconv"
)

func MustFloat2int(f float64) int {
	i, _ := strconv.Atoi(fmt.Sprintf("%.f", f))
	return i
}

func MustFloat(f float64) float64 {
	f, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", f), 64)
	return f
}
