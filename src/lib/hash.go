package lib

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword menerima password mentah dan mengembalikan string hash
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost saat ini bernilai 10 (rekomendasi standar)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash membandingkan password mentah dengan password hash di DB
// Mengembalikan nilai true jika cocok, dan false jika salah
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
