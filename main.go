package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"kefusystem/database"
	"kefusystem/routes"
)

func main() {
	// 获取当前工作目录（项目根目录）
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前工作目录失败: %v", err)
	}
	log.Printf("项目根目录: %s", cwd)

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

	// 加载templates目录下的所有HTML文件
	templatePath := filepath.Join(cwd, "templates/*.html")
	log.Printf("尝试加载模板: %s", templatePath)
	r.LoadHTMLGlob(templatePath)

	// 静态文件服务 - 注意这里指向templates目录下的静态资源
	staticPath := filepath.Join(cwd, "templates")
	uploadsPath := filepath.Join(cwd, "uploads")
	r.Static("/static", staticPath)
	r.Static("/uploads", uploadsPath)

	// 注册路由
	routes.RegisterPageRoutes(r, db)
	routes.RegisterWebSocketRoute(r, db)
	routes.RegisterAPIRoutes(r, db)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("服务器启动在 :%s", port)
	log.Fatal(r.Run(":" + port))
}
