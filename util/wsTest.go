package util

/**

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求，生产环境请限制来源
	},
}

type Client struct {
	ID   string
	Conn *websocket.Conn
}

var (
	clients     = make(map[string]*Client) // 使用ID作为键
	clientCount = 0                        // 客户端计数器
	mu          sync.Mutex                 // 保护对clients的并发访问

)

// InitWebsocket 初始化WebSocket模块（预留，如需）
// 不再在此注册路由，由 api 层注册路由到 HandleWebsocket
func InitWebsocket(r *gin.Engine) {
	// 预留：如需在此进行心跳、Hub初始化等
}

// HandleWebsocket  ws
// @Summary WebSocket连接
// @Description 升级并处理WebSocket连接
// @Tags WebSocket
// @Accept application/json
// @Produce application/json
// @Param userId query string false "用户ID"
// @Success 200 {object} string "成功"
// @Router /ws [GET]
func HandleWebsocket(c *gin.Context) {
	userID := c.Query("id")
	if userID == "" {
		userID = c.GetHeader("X-Chat-Id")
	}
	if userID == "" {
		userID = uuid.NewString()
	}

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("升级WebSocket连接失败:", err)
		return
	}
	// 生成唯一ID
	mu.Lock()
	clientCount++
	clientID := strconv.Itoa(clientCount) // 将计数器转换为字符串作为ID
	mu.Unlock()

	client := &Client{ID: clientID, Conn: conn}

	defer func(conn *websocket.Conn) {
		mu.Lock()                 // 锁定
		delete(clients, clientID) // 从客户端列表中删除
		mu.Unlock()               // 解锁
		//关闭WebSocket连接
		err := conn.Close()
		if err != nil {
			log.Println("关闭WebSocket连接失败:", err)
			return
		} else {
			log.Println("关闭WebSocket连接成功")
		}
	}(conn)
	mu.Lock()
	clients[clientID] = client
	mu.Unlock()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取WebSocket消息失败:", err)
			break
		}
		var data map[string]interface{}
		err = json.Unmarshal([]byte(message), &data)
		if err != nil {
			return
		}
		if name, ok := data["name"].(string); ok {
			fmt.Println("Name:", name)
		}
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("发送WebSocket消息失败:", err)
			break
		}
	}
}
**/
