package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB 初始化数据库连接
func InitDB() (*sql.DB, error) {
    // 创建数据库文件目录（如果不存在）
    if _, err := os.Stat("./data"); os.IsNotExist(err) {
        os.Mkdir("./data", 0755)
    }

    // 打开数据库连接
    db, err := sql.Open("sqlite3", "./data/kefu.db")
    if err != nil {
        return nil, fmt.Errorf("打开数据库失败: %v", err)
    }

    // 测试连接
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("连接数据库失败: %v", err)
    }

    // 创建表（如果不存在）
    err = createTables(db)
    if err != nil {
        return nil, fmt.Errorf("创建表失败: %v", err)
    }

    return db, nil
}

// createTables 创建数据库表
func createTables(db *sql.DB) error {
    // 创建用户表
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            role TEXT NOT NULL,
            created_at TEXT NOT NULL,
            updated_at TEXT NOT NULL
        )
    `)
    if err != nil {
        return err
    }

    // 创建会话表
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            user_name TEXT NOT NULL,
            status TEXT NOT NULL,
            created_at TEXT NOT NULL,
            updated_at TEXT NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users (id)
        )
    `)
    if err != nil {
        return err
    }

    // 创建消息表
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            session_id INTEGER NOT NULL,
            sender_id INTEGER NOT NULL,
            sender_name TEXT NOT NULL,
            content TEXT NOT NULL,
            message_type TEXT NOT NULL,
            created_at TEXT NOT NULL,
            FOREIGN KEY (session_id) REFERENCES sessions (id)
        )
    `)
    if err != nil {
        return err
    }

    // 创建FAQ表
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS faqs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            question TEXT NOT NULL,
            answer TEXT NOT NULL,
            created_at TEXT NOT NULL,
            updated_at TEXT NOT NULL
        )
    `)
    if err != nil {
        return err
    }

    // 创建系统配置表
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS config (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            welcome_text TEXT,
            welcome_image TEXT,
            banner_image TEXT,
            default_language TEXT,
            created_at TEXT NOT NULL,
            updated_at TEXT NOT NULL
        )
    `)
    if err != nil {
        return err
    }

    // 创建API参数表
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS api_params (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            param_key TEXT UNIQUE NOT NULL,
            param_value TEXT NOT NULL,
            description TEXT,
            created_at TEXT NOT NULL,
            updated_at TEXT NOT NULL
        )
    `)
    if err != nil {
        return err
    }

    return nil
}

// CreateDefaultConfig 创建默认配置
func CreateDefaultConfig(db *sql.DB) error {
    // 检查是否已有配置
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM config").Scan(&count)
    if err != nil && err != sql.ErrNoRows {
        return err
    }

    // 如果没有配置，则创建默认配置
    if count == 0 {
        now := time.Now().Format("2006-01-02 15:04:05")
        _, err := db.Exec(`
            INSERT INTO config (welcome_text, welcome_image, banner_image, default_language, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?)
        `, "欢迎使用在线客服系统", "", "", "zh-CN", now, now)
        if err != nil {
            return err
        }
        log.Println("创建默认配置成功")
    }

    // 检查是否有默认客服账号
    var userCount int
    err = db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&userCount)
    if err != nil && err != sql.ErrNoRows {
        return err
    }

    // 如果没有管理员账号，则创建默认管理员
    if userCount == 0 {
        now := time.Now().Format("2006-01-02 15:04:05")
        _, err := db.Exec(`
            INSERT INTO users (name, email, password, role, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?)
        `, "管理员", "admin@example.com", "admin123", "admin", now, now)
        if err != nil {
            return err
        }
        log.Println("创建默认管理员账号成功，默认邮箱: admin@example.com，密码: admin123")
    }

    return nil
}