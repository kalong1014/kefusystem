package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterPageRoutes 注册页面路由（修正函数名导出）
func RegisterPageRoutes(r *gin.Engine, db *sql.DB) {
	// 静态页面路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 客服工作台
	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", nil)
	})
	// 会话管理
	r.GET("/sessions", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sessions.html", nil)
	})

	// FAQ管理
	r.GET("/faqs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "faqs.html", nil)
	})

	// 用户管理
	r.GET("/users", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users.html", nil)
	})

	// 系统配置
	r.GET("/settings", func(c *gin.Context) {
		c.HTML(http.StatusOK, "settings.html", nil)
	})
}
