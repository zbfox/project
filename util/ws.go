package util

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求，生产环境请限制来源
	},
}

var (
	wsConns = make(map[string]*websocket.Conn)
	wsMu    sync.RWMutex
)

// InitWebsocket 初始化WebSocket模块（预留，如需）
// 不再在此注册路由，由 api 层注册路由到 HandleWebsocket
func InitWebsocket(r *gin.Engine) {
	// 预留：如需在此进行心跳、Hub初始化等
}

// HandleWebsocket 升级并处理WebSocket连接
// 支持通过 query/header 传入 userId：优先 query: userId，其次 header: X-User-Id；都无则分配随机ID
func HandleWebsocket(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		userID = c.GetHeader("X-User-Id")
	}
	if userID == "" {
		userID = uuid.NewString()
	}

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "WebSocket upgrade failed"})
		return
	}

	wsMu.Lock()
	wsConns[userID] = conn
	wsMu.Unlock()
	log.Printf("WebSocket connection established, user=%s", userID)

	defer func() {
		wsMu.Lock()
		delete(wsConns, userID)
		wsMu.Unlock()
		_ = conn.Close()
		log.Printf("WebSocket connection closed, user=%s", userID)
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		log.Printf("Received message from %s: %s", userID, p)

		// 简单路由：当文本消息内容为 "targetId:message" 时发送给目标用户，否则回写给自己
		if messageType == websocket.TextMessage {
			parts := strings.SplitN(string(p), ":", 2)
			if len(parts) == 2 {
				targetID := strings.TrimSpace(parts[0])
				msg := strings.TrimSpace(parts[1])
				if err := SendMessageToUser(targetID, []byte(msg)); err != nil {
					log.Println("Failed to send to target user:", err)
				}
				continue
			}
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}

// SendMessageToUser 发送消息给指定用户
func SendMessageToUser(userID string, data []byte) error {
	wsMu.RLock()
	conn, ok := wsConns[userID]
	wsMu.RUnlock()
	if !ok {
		return http.ErrNoCookie // 复用一个错误，表示未找到连接
	}
	return conn.WriteMessage(websocket.TextMessage, data)
}
