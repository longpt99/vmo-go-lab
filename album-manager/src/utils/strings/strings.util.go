package strings

import (
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

func RandomNumStr(length uint) string {
	number := strconv.FormatFloat(math.Floor(rand.Float64()*math.Pow(10, float64(length))), 'f', -1, 64)

	for len(number) < int(length) {
		number = "0" + number
	}

	return number
}

func GenerateString(length uint) string {
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	number := "0123456789"
	special := "#$@!%&*?"
	all := uppercase + lowercase + number + special
	strLen := length - 4
	strs := make([]string, strLen)

	for i := 0; i < int(strLen); i++ {
		strs = append(strs, characterRandom(all))
	}

	password := characterRandom(uppercase) +
		characterRandom(lowercase) +
		characterRandom(number) +
		characterRandom(special) +
		strings.Join(strs, "")

	pwSlice := strings.Split(password, "")
	sort.SliceStable(pwSlice, func(i, j int) bool {
		return rand.Float64() > 0.5
	})

	return strings.Join(pwSlice, "")
}

func characterRandom(str string) string {
	num := int(math.Floor(rand.Float64()*100)) % len(str)
	return string(str[num])
}
