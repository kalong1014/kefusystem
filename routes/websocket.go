package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
}

// 客户端管理器
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// 客户端结构
type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

// 消息结构
type Message struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

// 注册WebSocket路由
func registerWebSocketRoute(r *gin.Engine) {
	r.GET("/ws", handleWebSocket)

	// 启动客户端管理器
	go manager.start()
}

// 处理WebSocket连接
func handleWebSocket(c *gin.Context) {
	// 升级HTTP连接为WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.Error(c.Writer, "无法升级连接", http.StatusInternalServerError)
		return
	}

	// 生成唯一客户端ID
	clientID := generateClientID()

	// 创建客户端
	client := &Client{
		id:     clientID,
		socket: conn,
		send:   make(chan []byte, 256),
	}

	// 注册客户端
	manager.register <- client

	// 启动读取和写入goroutine
	go client.read()
	go client.write()
}

// 生成客户端ID
func generateClientID() string {
	return fmt.Sprintf("client-%d", time.Now().UnixNano())
}

// 客户端读取消息
func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("错误: %v", err)
			}
			break
		}

		// 广播消息
		manager.broadcast <- message
	}
}

// 客户端写入消息
func (c *Client) write() {
	defer c.socket.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// 发送通道关闭
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// 刷新缓冲区
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// 启动客户端管理器
func (m *ClientManager) start() {
	for {
		select {
		case client := <-m.register:
			m.clients[client] = true
			log.Printf("新客户端已连接: %s", client.id)

			// 向客户端发送欢迎消息
			welcomeMsg := []byte(fmt.Sprintf(`{"sender":"system","content":"欢迎，您的ID是 %s"}`, client.id))
			client.send <- welcomeMsg

		case client := <-m.unregister:
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
				log.Printf("客户端已断开连接: %s", client.id)
			}

		case message := <-m.broadcast:
			// 向所有客户端广播消息
			for client := range m.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(m.clients, client)
				}
			}
		}
	}
}
