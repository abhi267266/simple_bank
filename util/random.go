package util

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInd(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

func RandomString(n int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder

	for i:=0 ; i < n ; i++{
		b := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(b)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInd(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	return currencies[rand.Intn(len(currencies))]
}
