package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomInt32 generates a random integer between min and max
func RandomInt32(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return strings.ToLower(RandomString(6))
}

func RandomUserId() int32 {
	return RandomInt32(1, 1024)
}

func RandomImgId() int32 {
	return RandomInt32(1, 1024)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
