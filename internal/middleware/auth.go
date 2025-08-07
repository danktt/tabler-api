package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"tabler-api/config"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

// BetterAuthClaims representa as claims do Better Auth
type BetterAuthClaims struct {
	Sub       string `json:"sub"`
	Iss       string `json:"iss"`
	Aud       string `json:"aud"`
	Exp       int64  `json:"exp"`
	Iat       int64  `json:"iat"`
	Email     string `json:"email,omitempty"`
	Name      string `json:"name,omitempty"`
	EmailVerified bool `json:"emailVerified,omitempty"`
	Image     string `json:"image,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// BetterAuthMiddleware é o middleware principal para autenticação JWT com Better Auth
type BetterAuthMiddleware struct {
	config     *config.BetterAuthConfig
	keySet     jwk.Set
	lastUpdate time.Time
}

// NewBetterAuthMiddleware cria uma nova instância do middleware Better Auth
func NewBetterAuthMiddleware(cfg *config.BetterAuthConfig) (*BetterAuthMiddleware, error) {
	middleware := &BetterAuthMiddleware{
		config:     cfg,
		lastUpdate: time.Time{}, // Zero time to force initial load
	}

	// Carrega as chaves JWK iniciais
	if err := middleware.refreshJWKS(); err != nil {
		return nil, fmt.Errorf("failed to load initial JWKS: %w", err)
	}

	return middleware, nil
}

// refreshJWKS atualiza o conjunto de chaves JWK do Better Auth
func (m *BetterAuthMiddleware) refreshJWKS() error {
	// Verifica se precisa atualizar (a cada 24 horas)
	if time.Since(m.lastUpdate) < 24*time.Hour && m.keySet != nil {
		return nil
	}

	// Carrega as chaves JWK do endpoint do Better Auth
	fmt.Printf("Fetching JWKS from: %s\n", m.config.JWKSEndpoint)
	keySet, err := jwk.Fetch(context.Background(), m.config.JWKSEndpoint)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	m.keySet = keySet
	m.lastUpdate = time.Now()
	fmt.Printf("JWKS loaded successfully with %d keys\n", keySet.Len())
	return nil
}

// Authenticate é o middleware Gin para autenticação
func (m *BetterAuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai o token do header Authorization
		tokenString, err := extractTokenFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Atualiza as chaves JWK se necessário
		if err := m.refreshJWKS(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": "Failed to refresh authentication keys",
			})
			c.Abort()
			return
		}

		// Valida e decodifica o token
		claims, err := m.validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": err.Error(),
				"debug": gin.H{
					"issuer":   m.config.Issuer,
					"audience": m.config.Audience,
					"jwks_url": m.config.JWKSEndpoint,
				},
			})
			c.Abort()
			return
		}

		// Adiciona as claims ao contexto
		c.Set("user", claims)
		c.Next()
	}
}

// extractTokenFromHeader extrai o token JWT do header Authorization
func extractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	// Verifica se o header tem o formato "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("authorization header must be in format 'Bearer <token>'")
	}

	token := parts[1]
	if token == "" {
		return "", fmt.Errorf("token is required")
	}

	return token, nil
}

// validateToken valida o token JWT usando as chaves JWK do Better Auth
func (m *BetterAuthMiddleware) validateToken(tokenString string) (*BetterAuthClaims, error) {
	// Parse e valida o token
	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKeySet(m.keySet),
		jwt.WithValidate(true),
		jwt.WithIssuer(m.config.Issuer),
		jwt.WithAudience(m.config.Audience),
		jwt.WithAcceptableSkew(30*time.Second), // Tolerância de 30 segundos para diferença de relógio
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extrai as claims
	var claims BetterAuthClaims
	
	var sub string
	if err := token.Get("sub", &sub); err != nil {
		return nil, fmt.Errorf("invalid token: missing sub claim")
	}
	claims.Sub = sub
	
	var iss string
	if err := token.Get("iss", &iss); err != nil {
		return nil, fmt.Errorf("invalid token: missing iss claim")
	}
	claims.Iss = iss
	
	// Handle audience claim - it can be either a string or string array
	var audInterface interface{}
	if err := token.Get("aud", &audInterface); err != nil {
		return nil, fmt.Errorf("invalid token: missing aud claim")
	}
	
	// Convert audience to string
	switch v := audInterface.(type) {
	case string:
		claims.Aud = v
	case []string:
		if len(v) > 0 {
			claims.Aud = v[0] // Use the first audience
		} else {
			return nil, fmt.Errorf("invalid token: empty aud claim")
		}
	default:
		return nil, fmt.Errorf("invalid token: unexpected aud claim type")
	}
	
	// Handle expiration claim
	var expInterface interface{}
	if err := token.Get("exp", &expInterface); err != nil {
		return nil, fmt.Errorf("invalid token: missing exp claim")
	}
	
	// Convert expiration to int64
	switch v := expInterface.(type) {
	case float64:
		claims.Exp = int64(v)
	case int64:
		claims.Exp = v
	case int:
		claims.Exp = int64(v)
	default:
		return nil, fmt.Errorf("invalid token: unexpected exp claim type")
	}
	
	// Handle issued at claim
	var iatInterface interface{}
	if err := token.Get("iat", &iatInterface); err != nil {
		return nil, fmt.Errorf("invalid token: missing iat claim")
	}
	
	// Convert issued at to int64
	switch v := iatInterface.(type) {
	case float64:
		claims.Iat = int64(v)
	case int64:
		claims.Iat = v
	case int:
		claims.Iat = int64(v)
	default:
		return nil, fmt.Errorf("invalid token: unexpected iat claim type")
	}

	// Claims opcionais do Better Auth
	var email string
	if err := token.Get("email", &email); err == nil {
		claims.Email = email
	}
	
	var name string
	if err := token.Get("name", &name); err == nil {
		claims.Name = name
	}
	
	var emailVerified bool
	if err := token.Get("emailVerified", &emailVerified); err == nil {
		claims.EmailVerified = emailVerified
	}
	
	var image string
	if err := token.Get("image", &image); err == nil {
		claims.Image = image
	}
	
	var createdAt string
	if err := token.Get("createdAt", &createdAt); err == nil {
		claims.CreatedAt = createdAt
	}
	
	var updatedAt string
	if err := token.Get("updatedAt", &updatedAt); err == nil {
		claims.UpdatedAt = updatedAt
	}

	return &claims, nil
}

// GetUserFromContext extrai o usuário do contexto Gin
func GetUserFromContext(c *gin.Context) (*BetterAuthClaims, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	claims, ok := user.(*BetterAuthClaims)
	return claims, ok
} 