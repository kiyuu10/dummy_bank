package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "qwertyuiopasdfghjklzxcvbnm"

// func init() *rand.Rand {
// 	source := rand.NewSource(time.Now().UnixNano())
// 	rng := rand.New(source)
// 	return rng
// }

var (
	source = rand.NewSource(time.Now().UnixNano())
	rnd    = rand.New(source)
)

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rnd.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rnd.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Random Owner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// Random Money generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 10000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencis := []string{
		EUR, USD, CAD,
	}
	n := len(currencis)
	return currencis[rnd.Intn(n)]
}
