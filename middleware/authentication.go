package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/alyzsa/FinPro3/helper"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		verifyToken, err := helper.VerifyToken(c)
		_ = verifyToken

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthenticated",
				"message": err.Error(),
			})
			return
		}
		c.Set("userData", verifyToken)
		c.Next()
	}
}

type TaskAuthenticationData struct {
	CategoryID float64 `json:"category_id"`
}

func TaskAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawJSON, err := c.GetRawData()
		fmt.Println("Raw JSON request:", string(rawJSON))
		fmt.Println("Request Headers:", c.Request.Header)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Error reading raw JSON data",
			})
			c.Abort()
			return
		}

		var requestBody map[string]interface{}
		if err := json.Unmarshal(rawJSON, &requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid JSON payload",
			})
			c.Abort()
			return
		}

		categoryID, ok := requestBody["category_id"].(float64)
		if !ok || categoryID != float64(uint64(categoryID)) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "'category_id' must be a non-negative integer",
			})
			c.Abort()
			return
		}

		c.Set("categoryData", map[string]interface{}{"id": uint(categoryID)})

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawJSON))

		c.Next()
	}
}
