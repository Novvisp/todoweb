package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Task struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var db *gorm.DB

func main() {
	dsn := "host=localhost user=todo_user password=todo_password dbname=todo_db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к БД")
	}

	db.AutoMigrate(&Task{})

	r := gin.Default()

	r.POST("/tasks", createTask)
	r.GET("/tasks", getTasks)
	r.GET("/tasks/:id", getTaskByID)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)

	r.Run(":8080")
}

func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&task)
	c.JSON(http.StatusCreated, task)
}

func getTasks(c *gin.Context) {
	var tasks []Task
	db.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func getTaskByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var task Task
	if err := db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func updateTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var task Task
	if err := db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&task)
	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	db.Delete(&Task{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Задача удалена"})
}
