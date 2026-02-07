package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toanuitt/bookmark_service/pkg/jwtutils"
)

// JWTAuth defines an interface for JWT authentication middleware.
// Implementations should return a gin.HandlerFunc that validates JWT tokens
// and injects authenticated user information into the request context.
type JWTAuth interface {
	// JWTAuth returns a Gin middleware handler that:
	// - Reads the Authorization header
	// - Validates the JWT token
	// - Extracts the user identifier from the token
	// - Stores the user ID in the Gin context for downstream handlers
	JWTAuth() gin.HandlerFunc
}

// jwtAuth is the concrete implementation of JWTAuth.
// It uses a jwtutils.JWTValidator to validate and parse JWT tokens.
type jwtAuth struct {
	jwtValidator jwtutils.JWTValidator
}

// NewJWTAuth creates a new JWTAuth middleware instance using the provided
// JWTValidator. The validator is responsible for verifying token signatures
// and returning token claims.
func NewJWTAuth(jwtValidator jwtutils.JWTValidator) JWTAuth {
	return &jwtAuth{jwtValidator: jwtValidator}
}

// JWTAuth returns a Gin middleware handler that enforces JWT-based authentication.
//
// Behavior:
//   - Expects the Authorization header in the form: "Bearer <token>"
//   - Returns 401 Unauthorized if the header is missing or malformed
//   - Validates the token using the configured JWTValidator
//   - Returns 401 Unauthorized if the token is invalid
//   - Extracts the "sub" (subject) claim as the user ID
//   - Stores the user ID in the Gin context under the key "userID"
//   - Calls the next handler if authentication succeeds
func (j *jwtAuth) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Expect format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format is wrong"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		tokenContent, err := j.jwtValidator.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		// Extract user ID from token claims (subject)
		userID, ok := tokenContent["sub"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		// Store user ID in Gin context for downstream handlers
		c.Set("userID", userID)
		c.Next()
	}
}
