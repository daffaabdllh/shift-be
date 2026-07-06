package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// SetAuthCookies menyetel cookie access_token dan refresh_token berdasarkan file .env
func SetAuthCookies(c *gin.Context, accessToken string, refreshToken string) {
	// 1. Baca konfigurasi dinamis dari .env
	domain := os.Getenv("AUTH_COOKIE_DOMAIN")
	secure := os.Getenv("AUTH_COOKIE_SECURE") == "true"
	sameSiteStr := os.Getenv("AUTH_COOKIE_SAMESITE")

	// 2. Map string SameSite dari env ke tipe bawaan Go (http.SameSite)
	var sameSite http.SameSite
	switch sameSiteStr {
	case "Lax":
		sameSite = http.SameSiteLaxMode
	case "Strict":
		sameSite = http.SameSiteStrictMode
	case "None":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteDefaultMode
	}

	// 3. Set SameSite pada context Gin
	c.SetSameSite(sameSite)

	// 4. Set Cookie untuk Access Token (aktif 15 menit = 900 detik)
	c.SetCookie("access_token", accessToken, 900, "/", domain, secure, true)

	// 5. Set Cookie untuk Refresh Token (aktif 7 hari = 604800 detik)
	c.SetCookie("refresh_token", refreshToken, 604800, "/", domain, secure, true)
}

// ClearAuthCookies menghapus cookie access_token dan refresh_token dengan menyetel MaxAge ke -1
func ClearAuthCookies(c *gin.Context) {
	domain := os.Getenv("AUTH_COOKIE_DOMAIN")
	secure := os.Getenv("AUTH_COOKIE_SECURE") == "true"

	// Menyetel MaxAge ke -1 akan menghapus cookie di browser secara instan
	c.SetCookie("access_token", "", -1, "/", domain, secure, true)
	c.SetCookie("refresh_token", "", -1, "/", domain, secure, true)
}
