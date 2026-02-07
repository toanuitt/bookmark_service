package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toanuitt/bookmark_service/pkg/jwtutils"
)

type JWTAuth interface {
	JWTAuth() gin.HandlerFunc
}

type jwtAuth struct {
	jwtValidator jwtutils.JWTValidator
}

func NewJWTAuth(jwtValidator jwtutils.JWTValidator) JWTAuth {
	return &jwtAuth{jwtValidator: jwtValidator}
}

func (j *jwtAuth) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get auth header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		//get token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format is wrong"})
			c.Abort()
			return
		}
		tokenString := parts[1]
		// validate token
		tokenContent, err := j.jwtValidator.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}
		//verify token
		user_id, ok := tokenContent["sub"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}
		//store user_id to context
		c.Set("userID", user_id)
		c.Next()
	}
}
