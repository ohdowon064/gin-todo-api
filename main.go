package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Todo struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	category    *string   `gorm:"null;default:''"`
	description *string   `gorm:"null;default:''"`
	completed   bool      `gorm:"default:false"`
	createdAt   time.Time `gorm:"autoCreateTime"`
	updatedAt   time.Time
}

type GetTitle struct {
	Title string `form:"title"`
}

func main() {
	dsn := "host=localhost dbname=gin_todo_app port=5432 sslmode=disable TimeZone=Asia/Seoul"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	r := gin.Default()
	r.GET("/todos", func(c *gin.Context) {
		var title GetTitle
		var todos []Todo
		c.Bind(&title)

		db.Find(&todos)
		fmt.Println(todos)
		c.JSON(200, gin.H{
			"todos": todos,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
