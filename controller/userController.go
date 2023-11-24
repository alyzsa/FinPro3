package controller

import "github.com/gin-gonic/gin"

type UserHandler interface {
	UserRegister(*gin.Context)
	UserLogin(*gin.Context)
	UserUpdate(*gin.Context)
	UserDelete(*gin.Context)
}
