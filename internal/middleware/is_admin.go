package middleware

import (
	"net/http"

	"avito/internal/model"

	"github.com/gin-gonic/gin"
)

type IsAdmin struct {
}

func NewIsAdmin() *IsAdmin {
	return &IsAdmin{}
}

func (*IsAdmin) Handle(c *gin.Context) {
	user, ok := c.Get(KeyUser)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	userModel, ok := user.(*model.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "can't type assert user"})
		return
	}

	if userModel.Role != model.UserRoleAdmin {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Next()
}
