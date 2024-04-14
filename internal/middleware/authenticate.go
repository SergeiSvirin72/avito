package middleware

import (
	"net/http"
	"strings"

	"avito/internal/service"

	"github.com/gin-gonic/gin"
)

const KeyUser = "user"

type Authenticate struct {
	authService service.Auth
}

func NewAuthenticate(authService service.Auth) *Authenticate {
	return &Authenticate{
		authService: authService,
	}
}

func (m *Authenticate) Handle(c *gin.Context) {
	headers, ok := c.Request.Header["Authorization"]
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header not found"})
		return
	}

	if len(headers) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is empty"})
		return
	}

	tokenString := strings.TrimPrefix(headers[0], service.TokenTypeBearer+" ")

	user, err := m.authService.ParseToken(c.Request.Context(), tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Set(KeyUser, user)
	c.Next()
}
