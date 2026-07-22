package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"todo_api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken := ctx.GetHeader("Authorization")
		if authToken == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authToken, "Bearer ")
		if tokenString == "" || tokenString == authToken {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v ", t.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil

		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or exopred token"})
			ctx.Abort()
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			ctx.Abort()
			return

		}

		UserID, ok := claim["user_id"].(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			ctx.Abort()
			return

		}
		if exp, ok := claim["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token Expired"})
				ctx.Abort()
				return

			}

		}
		ctx.Set("user_id", UserID)
		ctx.Next()

	}
}
