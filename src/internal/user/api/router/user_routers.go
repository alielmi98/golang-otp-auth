package router

import (
	"github.com/alielmi98/golang-otp-auth/internal/user/api/handler"
	"github.com/alielmi98/golang-otp-auth/pkg/config"
	"github.com/gin-gonic/gin"
)

func Users(router *gin.RouterGroup, cfg *config.Config, handler *handler.UsersHandler) {

	router.POST("/send-otp", handler.SendOtp)
	router.POST("/login-by-mobile", handler.RegisterLoginByMobileNumber)
	router.GET("/:mobile_number", handler.GetUserByMobileNumber)
	router.GET("/", handler.GetUsers)

}
