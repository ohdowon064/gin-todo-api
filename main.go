package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strconv"
	"time"
)

type Todo struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	category    string `gorm:"null;default:''"`
	description string `gorm:"null;default:''"`
	completed   bool   `gorm:"default:false"`
	createdAt   time.Time
	updatedAt   time.Time
}

type TodoRes struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Category    *string   `json:"category"`
	Description *string   `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type GetTitle struct {
	Title string `form:"title"`
}

type CreateTodoRequestBody struct {
	Title       string `json:"title" binding:"required"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

func main() {
	dsn := "host=localhost dbname=gin_todo_app port=5432 sslmode=disable TimeZone=Asia/Seoul"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	r := gin.Default()
	r.GET("/todos", func(c *gin.Context) {
		var title GetTitle
		var todos []TodoRes
		c.Bind(&title)

		db.Raw("select * from todos where title ilike ?", "%"+title.Title+"%").Scan(&todos)
		c.JSON(200, todos)
	})

	r.POST("/todos", func(c *gin.Context) {
		var body CreateTodoRequestBody
		if err := c.Bind(&body); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{
				"message": "invalid request body",
			})
			return
		}

		var todo TodoRes
		db.Raw("insert into todos (title, category, description) values (?, ?, ?) returning id, title, category, description, completed, created_at, updated_at", body.Title, body.Category, body.Description).Scan(&todo)

		c.JSON(200, todo)
	})

	r.DELETE("/todos/:id", func(c *gin.Context) {
		ID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "id must be number",
			})
			return
		}

		var deletedID uint
		db.Raw("delete from todos where id = ? returning id", ID).Scan(&deletedID)
		if deletedID == 0 {
			c.JSON(404, gin.H{
				"message": "not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "deleted",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
