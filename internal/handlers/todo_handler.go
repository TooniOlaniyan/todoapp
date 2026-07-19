package handlers

import (
	"net/http"
	"strconv"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTodoInput struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateTodoInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}

		todo, err := repository.CreateTodo(pool, input.Title, input.Completed)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return

		}
		c.JSON(http.StatusCreated, todo)

	}

}

func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := repository.GetAllTodos(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todos)
	}

}

func GetTodoByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		todo, err := repository.GetTodoById(pool, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return

		}
		ctx.JSON(http.StatusOK, todo)
	}

}

func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idstr := ctx.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todo Id"})
			return
		}

		var input UpdateTodoInput

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}
		if input.Title == nil && input.Completed == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one filed (title , completed) has to be filled"})
			return

		}
		existing, err := repository.GetTodoById(pool, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo was not found"})
				return
			}
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		}
		title := existing.Title
		if input.Title != nil {
			title = *input.Title

		}
		completed := existing.Completed
		if input.Completed != nil {
			completed = *input.Completed

		}

		todo, err := repository.UpdateTodo(pool, id, title, completed)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		}

		ctx.JSON(http.StatusOK, todo)

	}

}
func DeleteTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idstr := ctx.Param("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		err = repository.DeleteTodo(pool, id)
		if err != nil {
			if err.Error() == "todo with id"+idstr+"not found" {
				ctx.JSON(http.StatusNotFound, gin.H{"error": " Todo not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Todo deleted Successfuly"})
	}
}
