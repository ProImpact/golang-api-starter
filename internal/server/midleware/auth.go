package midleware

import (
	auth "apistarter/internal/security"
	"apistarter/internal/server/response"
	"apistarter/pkg/model"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(
				c,
				http.StatusUnauthorized,
				model.UNAUTHENTICATED,
				"Authorization header is missing",
				nil,
			)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(
				c,
				http.StatusUnauthorized,
				model.UNAUTHENTICATED,
				"Invalid Authorization header format. Expected 'Bearer <token>'",
				nil,
			)
			c.Abort()
			return
		}

		tokenStr := parts[1]
		if tokenStr == "" {
			response.Error(
				c,
				http.StatusUnauthorized,
				model.UNAUTHENTICATED,
				"Token is empty",
				nil,
			)
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			var code model.ErrorCode
			var message string

			switch {
			case errors.Is(err, auth.ErrTokenExpired):
				code = model.UNAUTHENTICATED
				message = "Token has expired"
			case errors.Is(err, auth.ErrInvalidToken):
				code = model.UNAUTHENTICATED
				message = "Invalid token signature or malformed token"
			case errors.Is(err, auth.ErrInvalidIssuer):
				code = model.UNAUTHENTICATED
				message = "Token issued by an untrusted source"
			case errors.Is(err, auth.ErrInvalidSubject):
				code = model.UNAUTHENTICATED
				message = "Token subject is invalid"
			case errors.Is(err, auth.ErrInvalidAlgorithm):
				code = model.UNAUTHENTICATED
				message = "Token signed with an unsupported algorithm"
			default:
				code = model.UNAUTHENTICATED
				message = "Authentication failed"
			}

			response.Error(c, http.StatusUnauthorized, code, message, nil)
			c.Abort()
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("claims", claims)

		c.Next()
	}
}
