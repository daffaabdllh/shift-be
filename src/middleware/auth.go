package middleware

import (
	"net/http"
	"shift-be/src/lib"

	"github.com/gin-gonic/gin"
)

// IsLoggedIn adalah middleware untuk memastikan pengguna sudah login
func IsLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil access_token dari Cookie
		cookie, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Access token missing"})
			c.Abort() // Menghentikan request agar tidak lanjut ke handler
			return
		}

		// 2. Validasi access_token menggunakan helper ParseToken
		claims, err := lib.ParseToken(cookie, "access")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid or expired token"})
			c.Abort()
			return
		}

		// 3. Simpan data user claims ke Context Gin agar bisa diakses oleh handler selanjutnya
		c.Set("user", claims)

		c.Next() // Lanjut ke handler berikutnya
	}
}

func IsNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil access_token dari Cookie
		cookie, _ := c.Cookie("access_token")

		// 2. Validasi access_token menggunakan helper ParseToken
		_, err := lib.ParseToken(cookie, "access")
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request: Access token still exists"})
			c.Abort()
			return
		}

		c.Next() // Lanjut ke handler berikutnya
	}
}
