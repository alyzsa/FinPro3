package controller

import (
	"encoding/json"
	"github.com/alyzsa/FinPro3/database"
	"github.com/alyzsa/FinPro3/entity"
	"github.com/alyzsa/FinPro3/helper"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandlerImpl struct{}

func NewTaskHandlerImpl() TaskHandler {
	return &TaskHandlerImpl{}
}

func (s *TaskHandlerImpl) TaskCreate(c *gin.Context) {
	var db = database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	categoryData := c.MustGet("categoryData").(map[string]interface{})
	contentType := helper.GetContentType(c)

	userID := uint(userData["id"].(float64))
	categoryID := uint(categoryData["id"].(uint))
	Task := entity.Task{}
	Task.Status = false

	rawJSON, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Error reading raw JSON data",
		})
		return
	}

	fmt.Println("Raw JSON request:", string(rawJSON))
	fmt.Println("Request Headers:", c.Request.Header)

	if contentType == appJSON {
		if err := json.Unmarshal(rawJSON, &Task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid JSON payload for Comment",
			})
			return
		}
	} else {
		c.ShouldBind(&Task)
	}

	Task.UserID = userID
	Task.CategoryID = categoryID

	if err := Task.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	err = db.Debug().Create(&Task).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          Task.ID,
		"title":       Task.Title,
		"status":      Task.Status,
		"description": Task.Description,
		"user_id":     userID,
		"category_id": categoryID,
		"created_at":  Task.CreatedAt,
	})
}

func (s *TaskHandlerImpl) TaskGet(c *gin.Context) {
	var db = database.GetDB()
	contentType := helper.GetContentType(c)

	Task := entity.Task{}

	if contentType == appJSON {
		c.ShouldBindJSON(&Task)
	} else {
		c.ShouldBind(&Task)
	}

	err := db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, full_name, email")
	}).Find(&Task).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          Task.ID,
		"title":       Task.Title,
		"status":      Task.Status,
		"description": Task.Description,
		"user_id":     Task.UserID,
		"category_id": Task.CategoryID,
		"created_at":  Task.CreatedAt,
		"User": gin.H{
			"id":        Task.UserID,
			"email":     Task.User.Email,
			"full_name": Task.User.Full_Name,
		},
	})
}

func (s *TaskHandlerImpl) TaskUpdate(c *gin.Context) {
	var db = database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helper.GetContentType(c)
	_, _ = db, contentType

	userID := uint(userData["id"].(float64))
	Task := entity.Task{}

	taskID, _ := strconv.Atoi(c.Param("taskID"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Task)
	} else {
		c.ShouldBind(&Task)
	}

	Task.ID = uint(taskID)
	Task.UserID = userID

	err := db.Model(&Task).Where("id = ?", taskID).Updates(
		entity.Task{
			Title:       Task.Title,
			Description: Task.Description}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	updatedTask := entity.Task{}
	if err := db.Preload("Category").First(&updatedTask, taskID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          Task.ID,
		"title":       Task.Title,
		"description": Task.Description,
		"status":      Task.Status,
		"user_id":     Task.UserID,
		"category_id": updatedTask.Category.ID,
		"updated_at":  Task.UpdatedAt,
	})
}

func (s *TaskHandlerImpl) TaskStatusUpdate(c *gin.Context) {
	var db = database.GetDB()
	contentType := helper.GetContentType(c)
	_, _ = db, contentType

	Task := entity.Task{}

	taskID, _ := strconv.Atoi(c.Param("taskID"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Task)
	} else {
		c.ShouldBind(&Task)
	}

	Task.ID = uint(taskID)

	err := db.Model(&Task).Where("id = ?", taskID).Update("status", Task.Status).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	updatedTask := entity.Task{}
	if err := db.Preload("Category").First(&updatedTask, taskID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          updatedTask.ID,
		"title":       updatedTask.Title,
		"description": updatedTask.Description,
		"status":      updatedTask.Status,
		"user_id":     updatedTask.UserID,
		"category_id": updatedTask.CategoryID,
		"updated_at":  updatedTask.UpdatedAt,
	})
}

func (s *TaskHandlerImpl) TaskCategoryUpdate(c *gin.Context) {
	var db = database.GetDB()
	categoryData := c.MustGet("categoryData").(map[string]interface{})
	contentType := helper.GetContentType(c)
	_, _ = db, contentType

	Task := entity.Task{}
	categoryID := uint(categoryData["id"].(uint))

	taskID, _ := strconv.Atoi(c.Param("taskID"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Task)
	} else {
		c.ShouldBind(&Task)
	}

	Task.ID = uint(taskID)
	Task.CategoryID = categoryID

	err := db.Model(&Task).Where("id = ?", taskID).Update("category_id", Task.CategoryID).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	updatedTask := entity.Task{}
	if err := db.Preload("Category").First(&updatedTask, taskID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          updatedTask.ID,
		"title":       updatedTask.Title,
		"description": updatedTask.Description,
		"status":      updatedTask.Status,
		"user_id":     updatedTask.UserID,
		"category_id": updatedTask.CategoryID,
		"updated_at":  updatedTask.UpdatedAt,
	})
}

func (s *TaskHandlerImpl) TaskDelete(c *gin.Context) {
	var db = database.GetDB()
	contentType := helper.GetContentType(c)

	Task := entity.Task{}

	taskID, _ := strconv.Atoi(c.Param("taskID"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Task)
	} else {
		c.ShouldBind(&Task)
	}

	Task.ID = uint(taskID)

	err := db.Model(&Task).Where("id = ?", taskID).Delete(&Task).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task has been successfully deleted",
	})
}
