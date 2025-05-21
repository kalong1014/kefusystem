package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// API模型定义
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type FAQ struct {
	ID        int    `json:"id"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// RegisterAPIRoutes 注册API路由
func RegisterAPIRoutes(r *gin.Engine, db *sql.DB) {
	api := r.Group("/api")
	{
		// 登录接口
		api.POST("/login", func(c *gin.Context) {
			var req struct {
				Email    string `json:"email" binding:"required"`
				Password string `json:"password" binding:"required"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// 查询用户
			var user User
			err := db.QueryRow(`
                SELECT id, name, email, password, role, created_at, updated_at 
                FROM users 
                WHERE email = ?
            `, req.Email).Scan(
				&user.ID,
				&user.Name,
				&user.Email,
				&user.Password,
				&user.Role,
				&user.CreatedAt,
				&user.UpdatedAt,
			)

			if err != nil {
				if err == sql.ErrNoRows {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
				}
				return
			}

			// 检查密码
			if user.Password != req.Password {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
				return
			}

			// 生成令牌
			token := generateToken(user.ID, user.Role)

			c.JSON(http.StatusOK, gin.H{
				"user":  user,
				"token": token,
			})
		})

		// 文件上传接口
		api.POST("/upload", func(c *gin.Context) {
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// 创建uploads目录（如果不存在）
			if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
				os.Mkdir("./uploads", 0755)
			}

			// 保存文件
			fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
			filePath := fmt.Sprintf("./uploads/%s", fileName)
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"url":  fmt.Sprintf("/uploads/%s", fileName),
				"name": file.Filename,
			})
		})

		// FAQ管理接口
		api.GET("/faqs", func(c *gin.Context) {
			rows, err := db.Query("SELECT id, question, answer, created_at, updated_at FROM faqs ORDER BY created_at DESC")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
				return
			}
			defer rows.Close()

			var faqs []FAQ
			for rows.Next() {
				var faq FAQ
				if err := rows.Scan(&faq.ID, &faq.Question, &faq.Answer, &faq.CreatedAt, &faq.UpdatedAt); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
					return
				}
				faqs = append(faqs, faq)
			}

			c.JSON(http.StatusOK, faqs)
		})

		api.POST("/faq", func(c *gin.Context) {
			var faq FAQ
			if err := c.ShouldBindJSON(&faq); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			now := time.Now().Format("2006-01-02 15:04:05")
			result, err := db.Exec(`
                INSERT INTO faqs (question, answer, created_at, updated_at)
                VALUES (?, ?, ?, ?)
            `, faq.Question, faq.Answer, now, now)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
				return
			}

			id, _ := result.LastInsertId()
			faq.ID = int(id)
			faq.CreatedAt = now
			faq.UpdatedAt = now

			c.JSON(http.StatusOK, faq)
		})

		// 其他API接口保持不变...

	}
}

// 生成简单令牌
func generateToken(userID int, role string) string {
	return fmt.Sprintf("%d_%s_%s", userID, role, time.Now().Unix())
}
