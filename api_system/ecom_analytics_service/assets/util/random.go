package util_assets

// write a function to random string uppercase
import (
	"math/rand"
	"time"
)

func RandomString(length int) string {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rng.Intn(len(charset))]
	}

	return string(result)
}
func RandomNumber(min, max int) int {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)
	return rng.Intn(max-min+1) + min
}
func RandomBool() bool {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)
	if rng.Intn(2)%2 == 0 {
		return true
	} else {
		return false
	}
}
func RandomInArray[T any](array []T) T {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)
	return array[rng.Intn(len(array))]
}
