package routes

import (
	"action-camera/middleware"

	api "action-camera/api/v1"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	v1 := r.Group("/api/v1")
	{
		//用户操作
		v1.POST("user/login", api.UserLogin)
		v1.POST("user/vaild-email", api.UserVaildEmail)
		v1.POST("user/register", api.UserRegister)
	}
	return r
}
