package routes

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine) {
	// 注册页面路由
	RegisterPageRoutes(router)

	// 注册API路由
	RegisterAPIRoutes(router)

	// 注册WebSocket路由
	RegisterWebSocketRoute(router)
}

// 其他路由注册函数实现...
