package handlers

import (
	"net/http"
	"strings"
	"time"
	"todo_api/internal/config"
	"todo_api/internal/models"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var registerRequest RegisterRequest
		if err := ctx.BindJSON(&registerRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}
		if len(registerRequest.Password) < 6 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Passowrd must be at least 6 charcaters long"})
			return

		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user := &models.USER{
			Email:    registerRequest.Email,
			Password: string(hashedPassword),
		}

		createdUser, err := repository.CreateUser(pool, user)
		if err != nil {
			if strings.Contains(err.Error(), "duplicaate") || strings.Contains(err.Error(), "unique") {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
				return

			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, createdUser)

	}

}

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginRequest LoginRequest
		if err := ctx.BindJSON(&loginRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}
		user, err := repository.GetUserByEmail(pool, loginRequest.Email)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
			return
		}
		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token" + err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, LoginResponse{Token: tokenString})

	}
}
