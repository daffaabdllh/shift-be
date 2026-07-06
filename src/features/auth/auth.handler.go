package auth

import (
	"net/http"
	"net/mail"
	"shift-be/src/lib"

	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	psql sq.StatementBuilderType
}

func NewAuthHandler(psql sq.StatementBuilderType) *AuthHandler {
	return &AuthHandler{psql: psql}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginSchema

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		status, formattedErr := lib.FormatValidationError(err)
		c.JSON(status, formattedErr)
		return
	}

	_, err := mail.ParseAddress(req.Identifier)
	isEmail := err == nil

	searchColumn := "username"
	if isEmail {
		searchColumn = "email"
	}

	var user User
	err = h.
		psql.
		Select("id", "fullname", "username", "email", "password").
		From("users").
		Where(sq.Eq{searchColumn: req.Identifier}).
		Limit(1).
		QueryRow().
		Scan(&user.ID, &user.Fullname, &user.Username, &user.Email, &user.Password)

	if err != nil || !lib.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	accessToken, err := lib.GenerateToken(lib.TokenUser{
		ID:       user.ID,
		Fullname: user.Fullname,
		Username: user.Username,
		Email:    user.Email,
	}, "access")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := lib.GenerateToken(lib.TokenUser{
		ID:       user.ID,
		Fullname: user.Fullname,
		Username: user.Username,
		Email:    user.Email,
	}, "refresh")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	SetAuthCookies(c, accessToken, refreshToken)

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"data":    user,
	})
}

func (h *AuthHandler) Userinfo(c *gin.Context) {
	val, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	claims := val.(*lib.JWTClaims)

	c.JSON(http.StatusOK, gin.H{
		"message": "User info retrieved successfully",
		"data":    claims,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Panggil helper untuk menghapus cookie
	ClearAuthCookies(c)

	c.JSON(http.StatusNoContent, nil)
}
