package handlers

import (
	"net/http"

	"tabler-api/internal/auth"

	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
    if _, ok := auth.GetUser(c); !ok {
        c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authentication failed"})
        return
    }

    users := []auth.User{
        {ID: "1", Email: "alice@example.com", Name: "Alice"},
        {ID: "2", Email: "bob@example.com", Name: "Bob"},
    }

    c.JSON(http.StatusOK, users)
}
