package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由，接收数据库连接
func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	// 注册页面路由，传递数据库连接
	RegisterPageRoutes(router, db)

	// 注册API路由，传递数据库连接
	RegisterAPIRoutes(router, db)

	// 注册WebSocket路由，传递数据库连接
	RegisterWebSocketRoute(router, db)
}
