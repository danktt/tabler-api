package handler

import (
	"net/http"

	"tabler-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// ProfileHandler gerencia as rotas relacionadas ao perfil do usuário
type ProfileHandler struct{}

// NewProfileHandler cria uma nova instância do ProfileHandler
func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

// GetProfile retorna os dados do usuário autenticado extraídos do token JWT
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Extrai o usuário do contexto (definido pelo middleware de autenticação)
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not found in context",
		})
		return
	}

	// Retorna os dados do usuário
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data": gin.H{
			"sub":            user.Sub,
			"email":          user.Email,
			"name":           user.Name,
			"emailVerified":  user.EmailVerified,
			"image":          user.Image,
			"createdAt":      user.CreatedAt,
			"updatedAt":      user.UpdatedAt,
			"issuer":         user.Iss,
			"audience":       user.Aud,
			"expires_at":     user.Exp,
			"issued_at":      user.Iat,
		},
	})
} 