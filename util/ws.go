package util

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求，生产环境请限制来源
	},
}

// InitWebsocket 初始化websocket
func InitWebsocket(r *gin.Engine) {
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "WebSocket upgrade failed"})
			return
		}

		log.Println("WebSocket connection established")

		// WebSocket 消息处理循环
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println("WebSocket read error:", err)
				break
			}
			log.Printf("Received message: %s", p)

			// 回写消息给客户端
			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Println("WebSocket write error:", err)
				break
			}
		}
	})
}
