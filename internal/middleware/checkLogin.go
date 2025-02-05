package middleware

import (
	"QA-System/internal/pkg/code"
	"QA-System/internal/pkg/utils"
	"QA-System/internal/service"
	"github.com/gin-gonic/gin"
)

// CheckLogin 检查登录
func CheckLogin(c *gin.Context) {
	isLogin := service.CheckUserSession(c)
	if !isLogin {
		utils.JsonErrorResponse(c, code.NotLogin.Code, code.NotLogin.Msg)
		c.Abort()
		return
	}
	c.Next()
}
