package handlers

import (
	"net/http"

	"tabler-api/internal/auth"

	"github.com/gin-gonic/gin"
)

type AuthResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    User    *auth.User  `json:"user,omitempty"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

func VerifyAuth(c *gin.Context) {
    user, ok := auth.GetUser(c)
    if !ok {
        c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authentication failed"})
        return
    }

    u := user
    c.JSON(http.StatusOK, AuthResponse{
        Status:  "success",
        Message: "Token is valid",
        User:    &u,
    })
}
