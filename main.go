package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kalong1014/kefusystem/database"
	"github.com/kalong1014/kefusystem/routes"
)

func main() {
	// 初始化数据库
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer db.Close()

	// 创建默认配置
	if err := database.CreateDefaultConfig(db); err != nil {
		log.Printf("创建默认配置失败: %v", err)
	}

	// 设置Gin模式
	if os.Getenv("GIN_MODE") != "release" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 静态文件服务
	r.Static("/static", "./static")
	r.Static("/uploads", "./uploads")

	// 注册路由
	routes.RegisterAPIRoutes(r, db)
	routes.RegisterPageRoutes(r)
	routes.RegisterWebSocketRoute(r)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("服务器启动在 :%s", port)
	log.Fatal(r.Run(":" + port))
}
