package middleware

import (
	"strings"

	"tabler-api/config"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura os headers CORS para permitir requisições cross-origin
func CORSMiddleware(cfg *config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// Configurar origem permitida
		if len(cfg.AllowedOrigins) > 0 {
			// Se a primeira origem for "*", permitir todas as origens
			if cfg.AllowedOrigins[0] == "*" {
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				// Verificar se a origem da requisição está na lista de origens permitidas
				originAllowed := false
				for _, allowedOrigin := range cfg.AllowedOrigins {
					if allowedOrigin == origin {
						c.Header("Access-Control-Allow-Origin", origin)
						originAllowed = true
						break
					}
				}
				
				// Se a origem não estiver na lista, usar a primeira origem permitida
				if !originAllowed {
					c.Header("Access-Control-Allow-Origin", cfg.AllowedOrigins[0])
				}
			}
		}
		
		// Configurar headers permitidos
		if len(cfg.AllowedHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
		}
		
		// Configurar métodos permitidos
		if len(cfg.AllowedMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
		}
		
		// Permitir credenciais (cookies, headers de autorização)
		c.Header("Access-Control-Allow-Credentials", "true")
		
		// Expor headers customizados
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization, X-Total-Count")
		
		// Adicionar headers de debug para desenvolvimento
		if gin.Mode() == gin.DebugMode {
			c.Header("X-CORS-Debug", "enabled")
			c.Header("X-Request-Origin", origin)
		}
		
		// Se for uma requisição OPTIONS (preflight), responder imediatamente
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
} 