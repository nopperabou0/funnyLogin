package middleware

import (
	"funny-login/model"
	jwtservice "funny-login/utils/jwt_service"
	role "funny-login/utils/role"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type authHeader struct {
	AuthorizationHeader string `header:"Authorization" binding:"required"`
}

func RequireToken(roles ...role.Roles) gin.HandlerFunc {
	return func(c *gin.Context) {
		var aH authHeader
		if err := c.ShouldBindHeader(&aH); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized : " + err.Error()})
			return
		}
		token := strings.Replace(aH.AuthorizationHeader, "Bearer ", "", 1)
		tokenClaim, err := jwtservice.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized : " + err.Error()})
			return
		}
		c.Set("user", model.User{
			Id:   tokenClaim.ID,
			Role: tokenClaim.Role,
		})

		validRole := false

		for _, role := range roles {
			if string(role) == tokenClaim.Role {
				validRole = true
				break
			}
		}

		if !validRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Forbidden Resource"})
		}

		c.Next()

	}
}
