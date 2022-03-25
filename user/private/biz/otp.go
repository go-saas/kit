package biz

import (
	"crypto/rand"
)

const otpChars = "1234567890"
const otpLength = 6

func GenerateOtp() (string, error) {
	buffer := make([]byte, otpLength)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	otpCharsLength := len(otpChars)
	for i := 0; i < otpLength; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}
	return string(buffer), nil
}
