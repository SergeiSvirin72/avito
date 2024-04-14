package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"avito/internal/dto"
	"avito/internal/service"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	authService service.Auth
}

func NewAuth(authService service.Auth) *Auth {
	return &Auth{
		authService: authService,
	}
}

func (h *Auth) Auth(c *gin.Context) {
	var input dto.AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.authService.IssueToken(c.Request.Context(), &input)
	if err != nil && (errors.Is(err, sql.ErrNoRows) || errors.Is(err, service.ErrAuth)) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
