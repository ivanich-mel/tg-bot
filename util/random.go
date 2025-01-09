package util

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt (min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
func RandomFloat (min, max int64) float64 {
	return float64(min) + rand.Float64()*float64(max-min)
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
    var sb strings.Builder
    k := len(alphabet)

    for i := 0; i < n; i++ {
        c := alphabet[rand.Intn(k)]
        sb.WriteByte(c)
    }

    return sb.String()
}
func RandomName() string {
	return RandomString(10)
}
func RandomBalance()  float64 {
	return RandomFloat(10, 1000)
}