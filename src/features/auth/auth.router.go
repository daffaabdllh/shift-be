package auth

import (
	"shift-be/src/middleware"

	sq "github.com/Masterminds/squirrel" // Tambahkan import squirrel
	"github.com/gin-gonic/gin"
)

// Tambahkan parameter psql ke fungsi RegisterRouter
func RegisterRouter(r *gin.RouterGroup, psql sq.StatementBuilderType) {
	// 1. Buat instance AuthHandler dengan menyuntikkan psql
	h := NewAuthHandler(psql)

	// 2. Hubungkan ke method h.Login (menggunakan pointer h)
	r.POST("/auth/login", middleware.IsNotLoggedIn(), h.Login)
	r.GET("/auth/userinfo", middleware.IsLoggedIn(), h.Userinfo)
	r.DELETE("/auth/logout", middleware.IsLoggedIn(), h.Logout)
}
