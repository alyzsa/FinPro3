package middleware

import (
	"github.com/alyzsa/FinPro3/database"
	"github.com/alyzsa/FinPro3/entity"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func UserAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()
		userID, err := strconv.Atoi(c.Param("userID"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "invalid parameter",
			})
			return
		}

		user := entity.User{}

		err = db.Select("id").First(&user, uint(userID)).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Data Not Found",
				"message": "data doesnt exist",
			})
			return
		}

	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var db = database.GetDB()
		userData := c.MustGet("userData").(jwt.MapClaims)

		userID := uint(userData["id"].(float64))

		User := entity.User{}

		err := db.Where("id = ?", userID).Find(&User).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}

		if User.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		c.Next()
	}
}
