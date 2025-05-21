package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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

type Config struct {
	ID        int    `json:"id,omitempty"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// 注册API路由
func registerAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/login", loginHandler)
		api.POST("/upload", uploadHandler)

		// FAQ管理
		api.GET("/faqs", getFAQsHandler)
		api.POST("/faq", createFAQHandler)
		api.PUT("/faq/:id", updateFAQHandler)
		api.DELETE("/faq/:id", deleteFAQHandler)

		// 会话管理
		api.GET("/sessions", getSessionsHandler)
		api.GET("/session/:id/messages", getSessionMessagesHandler)
		api.POST("/session", createSessionHandler)
		api.POST("/message", sendMessageHandler)

		// 用户管理
		api.GET("/user/:id", getUserHandler)

		// 系统配置
		api.POST("/config", updateConfigHandler)
		api.GET("/config", getConfigHandler)

		// API参数
		api.POST("/param", saveAPIParamHandler)
	}
}

// API处理函数实现
func loginHandler(c *gin.Context) {
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

	// 检查密码（实际项目中应该使用哈希比较）
	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 生成令牌（实际项目中应该使用JWT）
	token := generateToken(user.ID, user.Role)

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// 生成简单令牌（实际项目中应使用JWT）
func generateToken(userID int, role string) string {
	return fmt.Sprintf("%d_%s_%s", userID, role, time.Now().Unix())
}

func uploadHandler(c *gin.Context) {
	// 文件上传逻辑
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
}

func getFAQsHandler(c *gin.Context) {
	// 获取常见问题列表
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
}

func createFAQHandler(c *gin.Context) {
	// 创建常见问题
	var faq FAQ
	if err := c.ShouldBindJSON(&faq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        INSERT INTO faqs (question, answer, created_at, updated_at)
        VALUES (?, ?, ?, ?)
    `, faq.Question, faq.Answer, time.Now(), time.Now())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取插入ID失败"})
		return
	}

	faq.ID = int(id)
	faq.CreatedAt = time.Now().Format(time.RFC3339)
	faq.UpdatedAt = faq.CreatedAt

	c.JSON(http.StatusCreated, faq)
}

func updateFAQHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的FAQ ID"})
		return
	}

	var faq FAQ
	if err := c.ShouldBindJSON(&faq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查FAQ是否存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM faqs WHERE id = ?", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "FAQ不存在"})
		return
	}

	// 更新FAQ
	_, err = db.Exec(`
        UPDATE faqs 
        SET question = ?, answer = ?, updated_at = ? 
        WHERE id = ?
    `, faq.Question, faq.Answer, time.Now(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	// 返回更新后的FAQ
	err = db.QueryRow(`
        SELECT id, question, answer, created_at, updated_at 
        FROM faqs 
        WHERE id = ?
    `, id).Scan(
		&faq.ID,
		&faq.Question,
		&faq.Answer,
		&faq.CreatedAt,
		&faq.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	c.JSON(http.StatusOK, faq)
}

func deleteFAQHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的FAQ ID"})
		return
	}

	// 检查FAQ是否存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM faqs WHERE id = ?", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "FAQ不存在"})
		return
	}

	// 删除FAQ
	_, err = db.Exec("DELETE FROM faqs WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "FAQ已删除"})
}

func getSessionsHandler(c *gin.Context) {
	userIDStr := c.Query("user_id")
	status := c.Query("status")

	query := "SELECT id, user_id, staff_id, status, created_at, updated_at FROM sessions"
	var args []interface{}
	whereClauses := []string{}

	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}
		whereClauses = append(whereClauses, "user_id = ?")
		args = append(args, userID)
	}

	if status != "" {
		whereClauses = append(whereClauses, "status = ?")
		args = append(args, status)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}
	defer rows.Close()

	var sessions []struct {
		ID        int    `json:"id"`
		UserID    int    `json:"user_id"`
		StaffID   int    `json:"staff_id"`
		Status    string `json:"status"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	for rows.Next() {
		var session struct {
			ID        int
			UserID    int
			StaffID   sql.NullInt64
			Status    string
			CreatedAt string
			UpdatedAt string
		}

		if err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.StaffID,
			&session.Status,
			&session.CreatedAt,
			&session.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
			return
		}

		sessions = append(sessions, struct {
			ID        int    `json:"id"`
			UserID    int    `json:"user_id"`
			StaffID   int    `json:"staff_id"`
			Status    string `json:"status"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		}{
			ID:        session.ID,
			UserID:    session.UserID,
			StaffID:   int(session.StaffID.Int64),
			Status:    session.Status,
			CreatedAt: session.CreatedAt,
			UpdatedAt: session.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, sessions)
}

func getSessionMessagesHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的会话ID"})
		return
	}

	// 检查会话是否存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sessions WHERE id = ?", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "会话不存在"})
		return
	}

	// 获取消息
	rows, err := db.Query(`
        SELECT id, session_id, sender_id, content, is_staff, created_at 
        FROM messages 
        WHERE session_id = ? 
        ORDER BY created_at ASC
    `, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}
	defer rows.Close()

	var messages []struct {
		ID        int    `json:"id"`
		SessionID int    `json:"session_id"`
		SenderID  int    `json:"sender_id"`
		Content   string `json:"content"`
		IsStaff   bool   `json:"is_staff"`
		CreatedAt string `json:"created_at"`
	}

	for rows.Next() {
		var message struct {
			ID        int
			SessionID int
			SenderID  int
			Content   string
			IsStaff   int
			CreatedAt string
		}

		if err := rows.Scan(
			&message.ID,
			&message.SessionID,
			&message.SenderID,
			&message.Content,
			&message.IsStaff,
			&message.CreatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
			return
		}

		messages = append(messages, struct {
			ID        int    `json:"id"`
			SessionID int    `json:"session_id"`
			SenderID  int    `json:"sender_id"`
			Content   string `json:"content"`
			IsStaff   bool   `json:"is_staff"`
			CreatedAt string `json:"created_at"`
		}{
			ID:        message.ID,
			SessionID: message.SessionID,
			SenderID:  message.SenderID,
			Content:   message.Content,
			IsStaff:   message.IsStaff == 1,
			CreatedAt: message.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, messages)
}

func createSessionHandler(c *gin.Context) {
	var session struct {
		UserID int    `json:"user_id" binding:"required"`
		Status string `json:"status" default:"open"`
	}

	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否存在
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", session.UserID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 创建会话
	result, err := db.Exec(`
        INSERT INTO sessions (user_id, status, created_at, updated_at)
        VALUES (?, ?, ?, ?)
    `, session.UserID, session.Status, time.Now(), time.Now())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取插入ID失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         id,
		"user_id":    session.UserID,
		"status":     session.Status,
		"created_at": time.Now().Format(time.RFC3339),
		"updated_at": time.Now().Format(time.RFC3339),
	})
}

func sendMessageHandler(c *gin.Context) {
	var message struct {
		SessionID int    `json:"session_id" binding:"required"`
		SenderID  int    `json:"sender_id" binding:"required"`
		Content   string `json:"content" binding:"required"`
		IsStaff   bool   `json:"is_staff"`
	}

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查会话是否存在
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sessions WHERE id = ?", message.SessionID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "会话不存在"})
		return
	}

	// 检查用户是否存在
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", message.SenderID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "发送者不存在"})
		return
	}

	// 发送消息
	result, err := db.Exec(`
        INSERT INTO messages (session_id, sender_id, content, is_staff, created_at)
        VALUES (?, ?, ?, ?, ?)
    `, message.SessionID, message.SenderID, message.Content, message.IsStaff, time.Now())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取插入ID失败"})
		return
	}

	// 更新会话状态为"active"
	_, err = db.Exec(`
        UPDATE sessions 
        SET status = 'active', updated_at = ? 
        WHERE id = ?
    `, time.Now(), message.SessionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新会话状态失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         id,
		"session_id": message.SessionID,
		"sender_id":  message.SenderID,
		"content":    message.Content,
		"is_staff":   message.IsStaff,
		"created_at": time.Now().Format(time.RFC3339),
	})
}

func getUserHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var user User
	err = db.QueryRow(`
        SELECT id, name, email, role, created_at, updated_at 
        FROM users 
        WHERE id = ?
    `, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func updateConfigHandler(c *gin.Context) {
	var config Config
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查配置是否存在
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM system_configs WHERE key = ?", config.Key).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	if count > 0 {
		// 更新配置
		_, err := db.Exec(`
            UPDATE system_configs 
            SET value = ?, updated_at = ? 
            WHERE key = ?
        `, config.Value, time.Now(), config.Key)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
			return
		}
	} else {
		// 创建新配置
		_, err := db.Exec(`
            INSERT INTO system_configs (key, value, created_at, updated_at)
            VALUES (?, ?, ?, ?)
        `, config.Key, config.Value, time.Now(), time.Now())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
			return
		}
	}

	// 返回配置
	err = db.QueryRow(`
        SELECT id, key, value, created_at, updated_at 
        FROM system_configs 
        WHERE key = ?
    `, config.Key).Scan(
		&config.ID,
		&config.Key,
		&config.Value,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	c.JSON(http.StatusOK, config)
}

func getConfigHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		// 获取所有配置
		rows, err := db.Query("SELECT id, key, value, created_at, updated_at FROM system_configs")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
			return
		}
		defer rows.Close()

		var configs []Config
		for rows.Next() {
			var config Config
			if err := rows.Scan(
				&config.ID,
				&config.Key,
				&config.Value,
				&config.CreatedAt,
				&config.UpdatedAt,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
				return
			}
			configs = append(configs, config)
		}

		c.JSON(http.StatusOK, configs)
	} else {
		// 获取单个配置
		var config Config
		err := db.QueryRow(`
            SELECT id, key, value, created_at, updated_at 
            FROM system_configs 
            WHERE key = ?
        `, key).Scan(
			&config.ID,
			&config.Key,
			&config.Value,
			&config.CreatedAt,
			&config.UpdatedAt,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
			}
			return
		}

		c.JSON(http.StatusOK, config)
	}
}

func saveAPIParamHandler(c *gin.Context) {
	var param struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 这里应该实现保存API参数的逻辑
	// 示例中简单返回成功
	c.JSON(http.StatusOK, gin.H{
		"key":   param.Key,
		"value": param.Value,
		"saved": true,
	})
}
