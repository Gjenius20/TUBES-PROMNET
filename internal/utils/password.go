package utils

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 12 // minimum yang disyaratkan: 10

// HashPassword menghasilkan hash bcrypt. Panjang password maksimal 72 byte (batas bcrypt).
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword membandingkan password dengan hash bcrypt.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
