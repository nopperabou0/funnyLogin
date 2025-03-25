package middleware

import (
	"net/http"
	"strings"

	"enigmacamp.com/unit-test-starter-pack/model"
	"enigmacamp.com/unit-test-starter-pack/utils/service"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	RequireToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
}

type authHeader struct {
	AuthorizationHeader string `header:"Authorization" binding:"required"`
}

func (a *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var aH authHeader
		if err := ctx.ShouldBindHeader(&aH); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		token := strings.Replace(aH.AuthorizationHeader, "Bearer ", "", 1)
		tokenClaim, err := a.jwtService.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		ctx.Set("user", model.UserCredential{Id: tokenClaim.ID, Role: tokenClaim.Role})
		validRole := false
		for _, role := range roles {
			if role == tokenClaim.Role {
				validRole = true
				break
			}
		}

		if !validRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden Resource"})
			return
		}
		ctx.Next()
	}
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
