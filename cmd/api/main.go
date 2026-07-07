package main

import (
	"log"
	"todo_api/internal/config"
	"todo_api/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()
	if err != nil {
		log.Fatal("Failed to load configurration:", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connet(cfg.DatabaseUrl)
	if err != nil {
		log.Println("Could not establish a Pool to the database", err)
	}
	defer pool.Close()
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message":  "TODO API running fine!",
			"status":   "sucess",
			"database": "database Connected",
		})
	})

	router.Run(":" + cfg.Port)

}
