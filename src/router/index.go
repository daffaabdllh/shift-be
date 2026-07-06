package router

import (
	"shift-be/src/features/auth"

	sq "github.com/Masterminds/squirrel" // Tambahkan import squirrel
	"github.com/gin-gonic/gin"
)

// Tambahkan parameter psql sq.StatementBuilderType di sini
func SetupRouter(r *gin.Engine, appPath string, psql sq.StatementBuilderType) {
	api := r.Group(appPath)

	// Kirim parameter psql ke RegisterRouter
	auth.RegisterRouter(api, psql)
}
