package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// TokenAuthMiddleware verifies the JWT token using the public key
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Check if the header is in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := parts[1]

		// Retrieve the public key from the environment variables
		publicKeyStr := os.Getenv("JWT_PUBLIC_KEY")
		if publicKeyStr == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT public key not configured"})
			c.Abort()
			return
		}

		// Parse the public key
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyStr))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid public key format"})
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token method is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Optionally pass claims to the next handler
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["regnumber"]) // Change "regnumber" to match your token claim key
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Token is valid; proceed with the request
		c.Next()
	}
}
