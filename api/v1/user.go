package v1

import (
	"action-camera/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	var userRegister service.UserService
	if err := c.ShouldBind(&userRegister); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		resp := userRegister.Register(c.Request.Context())
		c.JSON(http.StatusOK, resp)
	}
}
func UserVaildEmail(c *gin.Context) {
	var sendService service.SendEmailService
	if err := c.ShouldBind(&sendService); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		resp := sendService.Send(c.Request.Context())
		c.JSON(http.StatusOK, resp)
	}
}
func UserLogin(c *gin.Context) {
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		resp := userLogin.Login(c.Request.Context())
		c.JSON(http.StatusOK, resp)
	}
}
