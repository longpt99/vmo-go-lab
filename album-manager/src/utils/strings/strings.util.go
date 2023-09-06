package strings

import (
	"math"
	"math/rand"
	"strconv"
)

func RandomNumStr(length uint) string {
	number := strconv.FormatFloat(math.Floor(rand.Float64()*math.Pow(10, float64(length))), 'f', -1, 64)

	for len(number) < int(length) {
		number = "0" + number
	}

	return number
}
