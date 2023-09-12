package routes

import (
	"github.com/gin-gonic/gin"
	controllers_test "github.com/kishorens18/ecommerce/cmd/client/controllers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AppRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	user := v1.Group("/users")
	{
	user.POST("/signup", controllers_test.HandlerSignUp)
	user.POST("/signin", controllers_test.HandlerSignIn)
	user.POST("/delete", controllers_test.HandlerDeleteCustomer)
	user.POST("/update", controllers_test.HandlerUpdateCustomer)
	user.GET("/getbyid", controllers_test.HandlerGetById)
	user.POST("/reset", controllers_test.HandlerResetPassword)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
