package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomFloat(min, max float64) float64 {
	return float64(RandomInt(int64(min), int64(max))) + rand.Float64()
}

func RandomString(stringLength int) string {
	var sb strings.Builder
	alphabetsLength := len(alphabets)
	for i := 0; i < stringLength; i++ {
		char := alphabets[rand.Intn(alphabetsLength)]
		sb.WriteByte(char)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(8)
}

func RandomAmount() string {
	return decimal.NewFromFloat(randomFloat(0, 1000)).String()
}

func RandomCurrency() string {
	currencies := []string{
		INR,
		USD,
		CAD,
		EUR,
	}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomOwner())
}
