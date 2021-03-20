package utils

import "math/rand"

func GetRandomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bts = make([]byte, n)
	rand.Read(bts)
	for i, b := range bts {
		bts[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bts)
}

func GenerateApiKey() string {
	return GetRandomString(24)
}
