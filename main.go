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
	// 设置生产模式
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// 加载HTML模板
	router.LoadHTMLGlob("templates/*.html")

	// 配置信任的代理
	router.SetTrustedProxies([]string{"127.0.0.1"}) // 根据实际情况配置

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

	// 静态文件服务 - 分别为不同类型的静态资源设置路径
	r.Static("/static", filepath.Join(cwd, "/templates/static")) // 第三方库
	r.Static("/css", filepath.Join(cwd, "templates/css"))        // 自定义CSS
	r.Static("/js", filepath.Join(cwd, "templates/js"))          // 自定义JS
	r.Static("/locales", filepath.Join(cwd, "templates/locales"))
	r.StaticFile("/favicon.ico", filepath.Join(cwd, "templates/favicon.ico"))

	// 上传文件服务
	uploadsPath := filepath.Join(cwd, "uploads")
	r.Static("/uploads", uploadsPath)

	// 注册路由
	routes.RegisterPageRoutes(r, db)
	routes.RegisterWebSocketRoute(r, db)
	routes.RegisterAPIRoutes(r, db)
	RegisterRoutes(router)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("服务器启动在 :%s", port)
	log.Fatal(r.Run(":" + port))
}
