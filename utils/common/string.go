package common

import "math/rand"

func GenRandomString(length int) string {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := ""
	for i := 0; i < length; i++ {
		randomPos := rand.Intn(length)
		result += string(characters[randomPos])
	}

	return result
}
