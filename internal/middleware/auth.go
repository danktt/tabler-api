package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"tabler-api/config"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Auth0Claims representa as claims padrão do Auth0
type Auth0Claims struct {
	Sub       string `json:"sub"`
	Iss       string `json:"iss"`
	Aud       string `json:"aud"`
	Exp       int64  `json:"exp"`
	Iat       int64  `json:"iat"`
	Email     string `json:"email,omitempty"`
	Name      string `json:"name,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	Picture   string `json:"picture,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// Auth0Middleware é o middleware principal para autenticação JWT
type Auth0Middleware struct {
	config     *config.Auth0Config
	keySet     jwk.Set
	lastUpdate time.Time
}

// NewAuth0Middleware cria uma nova instância do middleware Auth0
func NewAuth0Middleware(cfg *config.Auth0Config) (*Auth0Middleware, error) {
	middleware := &Auth0Middleware{
		config:     cfg,
		lastUpdate: time.Time{}, // Zero time to force initial load
	}

	// Carrega as chaves JWK iniciais
	if err := middleware.refreshJWKS(); err != nil {
		return nil, fmt.Errorf("failed to load initial JWKS: %w", err)
	}

	return middleware, nil
}

// refreshJWKS atualiza o conjunto de chaves JWK do Auth0
func (m *Auth0Middleware) refreshJWKS() error {
	// Verifica se precisa atualizar (a cada 24 horas)
	if time.Since(m.lastUpdate) < 24*time.Hour && m.keySet != nil {
		return nil
	}

	// Carrega as chaves JWK do endpoint do Auth0
	keySet, err := jwk.Fetch(context.Background(), m.config.JWKSEndpoint)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	m.keySet = keySet
	m.lastUpdate = time.Now()
	return nil
}

// Authenticate é o middleware Gin para autenticação
func (m *Auth0Middleware) Authenticate() gin.HandlerFunc {
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

// validateToken valida o token JWT usando as chaves JWK do Auth0
func (m *Auth0Middleware) validateToken(tokenString string) (*Auth0Claims, error) {
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
	var claims Auth0Claims
	
	if sub, ok := token.Get("sub"); ok {
		if subStr, ok := sub.(string); ok {
			claims.Sub = subStr
		} else {
			return nil, fmt.Errorf("invalid token: sub claim is not a string")
		}
	} else {
		return nil, fmt.Errorf("invalid token: missing sub claim")
	}
	
	if iss, ok := token.Get("iss"); ok {
		if issStr, ok := iss.(string); ok {
			claims.Iss = issStr
		} else {
			return nil, fmt.Errorf("invalid token: iss claim is not a string")
		}
	} else {
		return nil, fmt.Errorf("invalid token: missing iss claim")
	}
	
	if aud, ok := token.Get("aud"); ok {
		if audStr, ok := aud.(string); ok {
			claims.Aud = audStr
		} else {
			return nil, fmt.Errorf("invalid token: aud claim is not a string")
		}
	} else {
		return nil, fmt.Errorf("invalid token: missing aud claim")
	}
	
	if exp, ok := token.Get("exp"); ok {
		if expNum, ok := exp.(float64); ok {
			claims.Exp = int64(expNum)
		} else {
			return nil, fmt.Errorf("invalid token: exp claim is not a number")
		}
	} else {
		return nil, fmt.Errorf("invalid token: missing exp claim")
	}
	
	if iat, ok := token.Get("iat"); ok {
		if iatNum, ok := iat.(float64); ok {
			claims.Iat = int64(iatNum)
		} else {
			return nil, fmt.Errorf("invalid token: iat claim is not a number")
		}
	} else {
		return nil, fmt.Errorf("invalid token: missing iat claim")
	}

	// Claims opcionais
	if email, ok := token.Get("email"); ok {
		if emailStr, ok := email.(string); ok {
			claims.Email = emailStr
		}
	}
	if name, ok := token.Get("name"); ok {
		if nameStr, ok := name.(string); ok {
			claims.Name = nameStr
		}
	}
	if nickname, ok := token.Get("nickname"); ok {
		if nicknameStr, ok := nickname.(string); ok {
			claims.Nickname = nicknameStr
		}
	}
	if picture, ok := token.Get("picture"); ok {
		if pictureStr, ok := picture.(string); ok {
			claims.Picture = pictureStr
		}
	}
	if updatedAt, ok := token.Get("updated_at"); ok {
		if updatedAtStr, ok := updatedAt.(string); ok {
			claims.UpdatedAt = updatedAtStr
		}
	}

	return &claims, nil
}

// GetUserFromContext extrai o usuário do contexto Gin
func GetUserFromContext(c *gin.Context) (*Auth0Claims, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	claims, ok := user.(*Auth0Claims)
	return claims, ok
} 