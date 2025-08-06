package middleware

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	redisChea "github.com/redis/go-redis/v9"
	"log"
	"time"
)

type CacheOptions struct {
	RedisClient *redisChea.Client
	TTL         time.Duration
	KeyFunc     func(c *gin.Context) string // 自定义 key（可选）
}

// RedisCacheMiddleware 返回一个 Gin 中间件用于自动缓存
func RedisCacheMiddleware(opts CacheOptions, handler gin.HandlerFunc) gin.HandlerFunc {

	if opts.KeyFunc == nil {
		opts.KeyFunc = defaultKeyFunc
	}

	return func(c *gin.Context) {
		ctx := context.Background()
		cacheKey := opts.KeyFunc(c)

		// 查询缓存
		cached, err := opts.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			// 命中缓存
			c.Data(200, "application/json", []byte(cached))
			c.Abort()
			return
		} else {
			// 响应拦截器
			writer := &bodyWriter{ResponseWriter: c.Writer, body: bytes.NewBuffer(nil)}
			c.Writer = writer
			handler(c)
			// 如果状态码是 200，写入缓存
			if c.Writer.Status() == 200 {
				go func() {
					err := opts.RedisClient.Set(ctx, cacheKey, writer.body.String(), opts.TTL).Err()
					if err == nil {
						log.Printf("缓存成功: %s", cacheKey)
					} else {
						log.Printf("缓存失败: %s", cacheKey)
					}
				}()
			}
		}

	}
}

// 默认的缓存 key（基于完整 URL）
func defaultKeyFunc(c *gin.Context) string {
	raw := c.Request.Method + ":" + c.Request.URL.RequestURI() + c.Request.Header.Get("token")
	sum := sha1.Sum([]byte(raw))
	return "cache:" + hex.EncodeToString(sum[:])
}

// bodyWriter 用于捕获响应内容
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写 Write 方法在调用 ResponseWriter.Write 时将内容写入 body
func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
