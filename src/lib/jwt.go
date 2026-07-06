package lib

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenUser struct {
	ID       string
	Fullname string
	Username string
	Email    string
}

type JWTClaims struct {
	ID       string `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken membuat 1 token JWT berdasarkan tipe: "access" atau "refresh"
func GenerateToken(user TokenUser, tokenType string) (string, error) {
	var secret string
	var duration time.Duration

	switch tokenType {
	case "access":
		secret = os.Getenv("JWT_ACCESS_SECRET")
		if secret == "" {
			secret = "supersecret_access_key"
		}
		duration = 15 * time.Minute

	case "refresh":
		secret = os.Getenv("JWT_REFRESH_SECRET")
		if secret == "" {
			secret = "supersecret_refresh_key"
		}
		duration = 7 * 24 * time.Hour

	default:
		return "", errors.New("invalid token type: must be 'access' or 'refresh'")
	}

	claims := JWTClaims{
		ID:       user.ID,
		Fullname: user.Fullname,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenObj.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseToken(tokenStr string, tokenType string) (*JWTClaims, error) {
	var secret string
	if tokenType == "access" {
		secret = os.Getenv("JWT_ACCESS_SECRET")
		if secret == "" {
			secret = "supersecret_access_key"
		}
	} else {
		secret = os.Getenv("JWT_REFRESH_SECRET")
		if secret == "" {
			secret = "supersecret_refresh_key"
		}
	}
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi metode enkripsi token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
