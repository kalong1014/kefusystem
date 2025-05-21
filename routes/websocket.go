package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
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

// 客户端结构
type Client struct {
	ID        string          `json:"id"`
	Conn      *websocket.Conn `json:"-"`
	Send      chan []byte     `json:"-"`
	IsAgent   bool            `json:"is_agent"`
	SessionID string          `json:"session_id,omitempty"`
}

// 消息结构
type Message struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	MsgType   string `json:"msg_type"` // text, image, system
}

// 会话结构
type Session struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	AgentID   string    `json:"agent_id,omitempty"`
	Status    string    `json:"status"` // open, closed, pending
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WebSocket管理器
type WebSocketManager struct {
	Clients    map[string]*Client
	Sessions   map[string]*Session
	Register   chan *Client
	Unregister chan *Client
	Message    chan *Message
	DB         *sql.DB
	Mutex      sync.Mutex
}

// 创建新的WebSocket管理器
func NewWebSocketManager(db *sql.DB) *WebSocketManager {
	return &WebSocketManager{
		Clients:    make(map[string]*Client),
		Sessions:   make(map[string]*Session),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Message:    make(chan *Message),
		DB:         db,
	}
}

// 启动WebSocket管理器
func (m *WebSocketManager) Start() {
	for {
		select {
		case client := <-m.Register:
			m.registerClient(client)

		case client := <-m.Unregister:
			m.unregisterClient(client)

		case message := <-m.Message:
			m.handleMessage(message)
		}
	}
}

// 注册客户端
func (m *WebSocketManager) registerClient(client *Client) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	m.Clients[client.ID] = client
	log.Printf("新客户端已连接: %s (agent: %v)", client.ID, client.IsAgent)

	// 创建新会话或分配现有会话
	if !client.IsAgent {
		sessionID := fmt.Sprintf("session-%d", time.Now().UnixNano())
		session := &Session{
			ID:        sessionID,
			ClientID:  client.ID,
			Status:    "pending",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		client.SessionID = sessionID
		m.Sessions[sessionID] = session

		// 保存会话到数据库
		_, err := m.DB.Exec(`
            INSERT INTO sessions (id, user_id, user_name, status, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?)
        `, sessionID, client.ID, client.ID, "pending",
			session.CreatedAt.Format("2006-01-02 15:04:05"),
			session.UpdatedAt.Format("2006-01-02 15:04:05"))

		if err != nil {
			log.Printf("保存会话失败: %v", err)
		}

		// 向客户端发送欢迎消息
		welcomeMsg := &Message{
			Sender:    "system",
			Content:   "欢迎使用在线客服系统，我们将尽快为您分配客服人员",
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			MsgType:   "system",
		}
		client.Send <- m.serializeMessage(welcomeMsg)
	} else {
		// 向客服发送在线客户列表
		m.sendAgentClientList(client)
	}
}

// 序列化消息为JSON
func (m *WebSocketManager) serializeMessage(msg *Message) []byte {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("序列化消息失败: %v", err)
		return []byte(`{"sender":"system","content":"消息格式错误","timestamp":"` +
			time.Now().Format("2006-01-02 15:04:05") + `","msg_type":"system"}`)
	}
	return data
}

// 发送客服客户端列表
func (m *WebSocketManager) sendAgentClientList(agent *Client) {
	// 构建客户端列表消息
	clientsData := struct {
		Type    string    `json:"type"`
		Clients []*Client `json:"clients"`
	}{
		Type:    "client_list",
		Clients: make([]*Client, 0),
	}

	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	// 修复：使用正确的变量名
	for clientID, client := range m.Clients {
		if !client.IsAgent {
			clientsData.Clients = append(clientsData.Clients, &Client{
				ID:        clientID, // 使用正确的clientID变量
				IsAgent:   client.IsAgent,
				SessionID: client.SessionID,
			})
		}
	}

	data, err := json.Marshal(clientsData)
	if err != nil {
		log.Printf("序列化客户端列表失败: %v", err)
		return
	}

	agent.Send <- data
}

// 注销客户端
func (m *WebSocketManager) unregisterClient(client *Client) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if _, ok := m.Clients[client.ID]; ok {
		// 关闭发送通道
		close(client.Send)

		// 关闭WebSocket连接
		if client.Conn != nil {
			client.Conn.Close()
		}

		// 删除客户端
		delete(m.Clients, client.ID)
		log.Printf("客户端已断开连接: %s", client.ID)

		// 如果是客户断开连接，更新会话状态
		if !client.IsAgent && client.SessionID != "" {
			if session, ok := m.Sessions[client.SessionID]; ok {
				session.Status = "closed"
				session.UpdatedAt = time.Now()

				// 更新数据库中的会话状态
				_, err := m.DB.Exec(`
                    UPDATE sessions 
                    SET status = ?, updated_at = ? 
                    WHERE id = ?
                `, "closed",
					session.UpdatedAt.Format("2006-01-02 15:04:05"),
					session.ID)

				if err != nil {
					log.Printf("更新会话状态失败: %v", err)
				}

				// 通知所有客服会话已关闭
				closeMsg := &Message{
					Sender:    "system",
					Content:   fmt.Sprintf("会话 %s 已关闭", client.SessionID),
					Timestamp: time.Now().Format("2006-01-02 15:04:05"),
					MsgType:   "system",
				}

				for _, agent := range m.Clients {
					if agent.IsAgent {
						agent.Send <- m.serializeMessage(closeMsg)
					}
				}
			}
		}
	}
}

// 处理消息
func (m *WebSocketManager) handleMessage(message *Message) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	// 保存消息到数据库
	if message.MsgType != "system" {
		_, err := m.DB.Exec(`
            INSERT INTO messages (session_id, sender_id, sender_name, content, message_type, created_at)
            VALUES (?, ?, ?, ?, ?, ?)
        `, message.Receiver, message.Sender, message.Sender,
			message.Content, message.MsgType,
			message.Timestamp)

		if err != nil {
			log.Printf("保存消息失败: %v", err)
		}
	}

	// 查找接收者
	if receiver, ok := m.Clients[message.Receiver]; ok {
		receiver.Send <- m.serializeMessage(message)
	} else {
		// 接收者不在线，发送离线消息
		offlineMsg := &Message{
			Sender:    "system",
			Content:   "对方当前不在线，您的消息可能无法及时送达",
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			MsgType:   "system",
		}

		if sender, ok := m.Clients[message.Sender]; ok {
			sender.Send <- m.serializeMessage(offlineMsg)
		}
	}
}

// 注册WebSocket路由
func RegisterWebSocketRoute(r *gin.Engine, db *sql.DB) {
	manager := NewWebSocketManager(db)
	go manager.Start()

	r.GET("/ws", func(c *gin.Context) {
		// 升级HTTP连接为WebSocket连接
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket升级失败: %v", err)
			return
		}

		// 获取客户端ID和是否为客服标识
		clientID := c.Query("client_id")
		if clientID == "" {
			clientID = fmt.Sprintf("guest_%d", time.Now().UnixNano())
		}

		isAgent := c.Query("is_agent") == "true"

		// 创建客户端
		client := &Client{
			ID:      clientID,
			Conn:    conn,
			Send:    make(chan []byte, 256),
			IsAgent: isAgent,
		}

		// 注册客户端
		manager.Register <- client

		// 启动读取和写入goroutine
		go client.read(manager)
		go client.write()
	})
}

// 客户端读取消息
func (c *Client) read(manager *WebSocketManager) {
	defer func() {
		manager.Unregister <- c
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("客户端读取错误: %v", err)
			}
			break
		}

		// 解析消息
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("解析消息失败: %v", err)
			continue
		}

		// 设置发送者和时间戳
		msg.Sender = c.ID
		msg.Timestamp = time.Now().Format("2006-01-02 15:04:05")

		// 如果是客户发送的消息且没有接收者，设置接收者为所有客服
		if !c.IsAgent && msg.Receiver == "" {
			// 查找可用客服
			var agentID string
			for id, client := range manager.Clients {
				if client.IsAgent {
					agentID = id
					break
				}
			}

			if agentID != "" {
				msg.Receiver = agentID

				// 如果会话还没有分配客服，分配当前客服
				if c.SessionID != "" {
					if session, ok := manager.Sessions[c.SessionID]; ok {
						session.AgentID = agentID
						session.Status = "open"
						session.UpdatedAt = time.Now()

						// 更新数据库中的会话
						_, err := manager.DB.Exec(`
                            UPDATE sessions 
                            SET agent_id = ?, status = ?, updated_at = ? 
                            WHERE id = ?
                        `, agentID, "open",
							session.UpdatedAt.Format("2006-01-02 15:04:05"),
							session.ID)

						if err != nil {
							log.Printf("更新会话失败: %v", err)
						}
					}
				}
			} else {
				// 没有可用客服，发送系统消息
				systemMsg := Message{
					Sender:    "system",
					Content:   "当前没有可用客服，请稍后再试",
					Timestamp: time.Now().Format("2006-01-02 15:04:05"),
					MsgType:   "system",
				}
				c.Send <- manager.serializeMessage(&systemMsg)
				continue
			}
		}

		// 发送消息到管理器
		manager.Message <- &msg
	}
}

// 客户端写入消息
func (c *Client) write() {
	defer func() {
		if c.Conn != nil {
			c.Conn.Close()
		}
	}()

	for message := range c.Send {
		// 写入消息
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("写入消息失败: %v", err)
			return
		}
	}

	// 发送通道关闭，关闭WebSocket连接
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}
