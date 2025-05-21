package main

import (
	"database/sql"
	"log"
)

// 全局变量
var db *sql.DB

// 初始化数据库
func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		return err
	}

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		return err
	}

	// 创建表
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            role TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            customer_id INTEGER NOT NULL,
            staff_id INTEGER,
            status TEXT NOT NULL DEFAULT 'open',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (customer_id) REFERENCES users (id),
            FOREIGN KEY (staff_id) REFERENCES users (id)
        )`,
		`CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            session_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            content TEXT NOT NULL,
            is_staff BOOLEAN NOT NULL DEFAULT 0,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (session_id) REFERENCES sessions (id),
            FOREIGN KEY (user_id) REFERENCES users (id)
        )`,
		`CREATE TABLE IF NOT EXISTS faqs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            question TEXT NOT NULL,
            answer TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS system_configs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            key TEXT NOT NULL UNIQUE,
            value TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS api_params (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            value TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return err
		}
	}

	return nil
}

// 创建默认配置
func createDefaultConfig() {
	// 检查是否已有配置
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM system_configs").Scan(&count)
	if err != nil {
		log.Printf("查询配置数量失败: %v", err)
		return
	}

	// 如果没有配置，创建默认配置
	if count == 0 {
		configs := []struct {
			Key   string
			Value string
		}{
			{"system_name", "智能客服系统"},
			{"logo_url", "/static/logo.png"},
			{"welcome_message", "您好！有什么可以帮助您的？"},
			{"working_hours", "周一至周五 9:00-18:00"},
			{"max_response_time", "30分钟"},
		}

		for _, config := range configs {
			_, err := db.Exec(`
                INSERT INTO system_configs (key, value)
                VALUES (?, ?)
            `, config.Key, config.Value)
			if err != nil {
				log.Printf("创建默认配置失败: %v", err)
			}
		}
	}
}
