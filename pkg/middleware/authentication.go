package middleware

import (
	"event-management/model"
	"event-management/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *middleware) AuthenticateUser(c *gin.Context) {
	bearer := c.GetHeader("Authorization")
	if bearer == "" {
		response.Error(c, http.StatusUnauthorized, "empty token", nil)
		c.Abort()
		return
	}

	token := strings.Split(bearer, " ")[1]
	userID, err := m.jwtAuth.ValidateToken(token)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "failed to validate token", err)
		c.Abort()
		return
	}

	user, err := m.service.UserService.GetUser(model.GetUserParam{
		UserID: userID,
	})
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "failed to get user", err)
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}
