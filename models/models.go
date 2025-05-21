package main

import "time"

// 定义数据模型
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Role     string `json:"role"` // admin, staff, customer
    Email    string `json:"email"`
    Password string `json:"password,omitempty"`
}

type Session struct {
    ID          int       `json:"id"`
    CustomerID  int       `json:"customer_id"`
    StaffID     int       `json:"staff_id,omitempty"`
    Status      string    `json:"status"` // open, closed, pending
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Message struct {
    ID         int       `json:"id"`
    SessionID  int       `json:"session_id"`
    UserID     int       `json:"user_id"`
    Content    string    `json:"content"`
    IsStaff    bool      `json:"is_staff"`
    CreatedAt  time.Time `json:"created_at"`
}

type FAQ struct {
    ID         int       `json:"id"`
    Question   string    `json:"question"`
    Answer     string    `json:"answer"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type Config struct {
    ID         int       `json:"id"`
    Key        string    `json:"key"`
    Value      string    `json:"value"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type APIParam struct {
    ID         int       `json:"id"`
    Name       string    `json:"name"`
    Value      string    `json:"value"`
    CreatedAt  time.Time `json:"created_at"`
}    