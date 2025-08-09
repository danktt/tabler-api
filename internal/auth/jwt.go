package auth

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

// User represents the authenticated user extracted from the JWT.
type User struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}

// UserContextKey is the key used to store the user in the Gin context.
const UserContextKey = "user"

var (
    // ErrMissingUserID is returned when the token has no subject claim.
    ErrMissingUserID = errors.New("missing user id")
)

// NewAuthMiddleware creates a Gin middleware that validates JWTs using the provided JWKS URL
// and stores the authenticated user in the request context.
func NewAuthMiddleware(jwksURL string, logger *slog.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        keyset, err := jwk.Fetch(c.Request.Context(), jwksURL)
        if err != nil {
            logger.Error("fetch jwks", slog.Any("error", err))
            c.AbortWithStatusJSON(401, gin.H{"error": "Authentication failed"})
            return
        }

        token, err := jwt.ParseRequest(c.Request, jwt.WithKeySet(keyset))
        if err != nil {
            logger.Error("parse jwt", slog.Any("error", err))
            c.AbortWithStatusJSON(401, gin.H{"error": "Authentication failed"})
            return
        }

        sub, ok := token.Subject()
        if !ok {
            c.AbortWithStatusJSON(401, gin.H{"error": ErrMissingUserID.Error()})
            return
        }

        var email string
        var name string
        _ = token.Get("email", &email)
        _ = token.Get("name", &name)

        user := User{ID: sub, Email: email, Name: name}
        c.Set(UserContextKey, user)
        c.Next()
    }
}

// GetUser returns the authenticated user from the Gin context, if present.
func GetUser(c *gin.Context) (User, bool) {
    v, exists := c.Get(UserContextKey)
    if !exists {
        return User{}, false
    }
    user, ok := v.(User)
    return user, ok
}
