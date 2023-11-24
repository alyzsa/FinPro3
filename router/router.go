package router

import (
	"github.com/alyzsa/FinPro3/controller"
	"github.com/alyzsa/FinPro3/middleware"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/user")
	{
		userRouter.POST("/register", controller.NewUserHandlerImpl().UserRegister)

		userRouter.POST("/login", controller.NewUserHandlerImpl().UserLogin)

		userRouter.Use(middleware.Authentication())
		userRouter.PUT("/update-account", controller.NewUserHandlerImpl().UserUpdate)

		userRouter.DELETE("/:userID", middleware.UserAuthorization(), controller.NewUserHandlerImpl().UserDelete)
	}

	categoryRouter := r.Group("/category")
	{
		categoryRouter.Use(middleware.Authentication())
		categoryRouter.Use(middleware.AdminAuthMiddleware())
		categoryRouter.POST("/create", controller.NewCategoryHandlerImpl().CategoryCreate)
		categoryRouter.GET("/get", controller.NewCategoryHandlerImpl().CategoryGet)
		categoryRouter.PATCH("/update/:categoryID", controller.NewCategoryHandlerImpl().CategoryUpdate)
		categoryRouter.DELETE("/delete/:categoryID", controller.NewCategoryHandlerImpl().CategoryDelete)

	}

	taskRouter := r.Group("/task")
	{
		taskRouter.Use(middleware.Authentication())
		taskRouter.POST("/create", middleware.TaskAuthentication(), controller.NewTaskHandlerImpl().TaskCreate)
		taskRouter.GET("/get", controller.NewTaskHandlerImpl().TaskGet)
		taskRouter.PUT("/update/:taskID", controller.NewTaskHandlerImpl().TaskUpdate)
		taskRouter.DELETE("/delete/:taskID", controller.NewTaskHandlerImpl().TaskDelete)
		taskRouter.PATCH("/update-status/:taskID", controller.NewTaskHandlerImpl().TaskStatusUpdate)
		taskRouter.PATCH("/update-category/:taskID", middleware.TaskAuthentication(), controller.NewTaskHandlerImpl().TaskCategoryUpdate)
	}

	return r
}
