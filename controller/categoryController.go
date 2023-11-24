package controller

import "github.com/gin-gonic/gin"

type CategoryHandler interface {
	CategoryCreate(*gin.Context)
	CategoryGet(*gin.Context)
	CategoryUpdate(*gin.Context)
	CategoryDelete(*gin.Context)
}
