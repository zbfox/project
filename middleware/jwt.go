package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// CreateToken  创建token

const SIGNINGKEY = "your-very-long-and-complex-secret-key-at-least-32-characters-long"

var (
	initRedis *redis.Client
	ctx       = context.Background()
	once      sync.Once
)

const (
	AccessTokenExpiration  = 15 * time.Minute
	RefreshTokenExpiration = 7 * 24 * time.Hour
)

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Claims 自定义Claims结构
type Claims struct {
	UserID   string `json:"sub"`
	Username string `json:"name,omitempty"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

// InitJWTMiddleware 获取Redis客户端
func InitJWTMiddleware(redis *redis.Client) {
	initRedis = redis
}

// 生成唯一ID
func generateUniqueID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// CreateAccessToken 创建访问令牌 (短期 - 15分钟)
// CreateAccessToken 创建访问令牌 (短期 - 15分钟)
func CreateAccessToken(userID, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SIGNINGKEY))
}

// CreateRefreshToken 创建刷新令牌 (长期 - 7天)
func CreateRefreshToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SIGNINGKEY))
}

func Login(UUID, username string) (*TokenPair, error) {
	accessToken, err := CreateAccessToken(UUID, username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := CreateRefreshToken(UUID)
	if err != nil {
		return nil, err
	}
	userKey := fmt.Sprintf("jwt:%s", UUID)
	err = initRedis.HMSet(ctx, userKey, map[string]interface{}{
		"userUUID":    UUID,
		"username":    username,
		"login_time":  time.Now().Unix(),
		"last_active": time.Now().Unix(),
	}).Err()
	if err != nil {
		log.Printf("缓存用户信息失败: %v", err)
	}
	initRedis.Expire(ctx, userKey, AccessTokenExpiration)
	// 缓存Refresh Token (7天过期)
	refreshKey := fmt.Sprintf("refresh:%s", UUID)
	err = initRedis.Set(ctx, refreshKey, refreshToken, RefreshTokenExpiration).Err()
	if err != nil {
		log.Printf("缓存刷新令牌失败: %v", err)
	}

	log.Printf("✅ 用户 %s 登录成功，信息已缓存\n", username)
	//添加到响应头中
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// 需要排除的路由
var excludePaths = []string{
	"/api/user/login",
	"/api/user/refresh",
}

func isExcludedPath(path string) bool {
	for _, p := range excludePaths {
		if path == p {
			return true
		}
	}
	return false
}

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isExcludedPath(c.Request.URL.Path) {
			return
		}
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证令牌"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌格式错误"})
			c.Abort()
			return
		}

		// 解析令牌
		claims, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌无效或者过期", "code": http.StatusUnauthorized})
			c.Abort()
			return
		}

		// 验证令牌类型
		if claims.Type != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌类型错误"})
			c.Abort()
			return
		}

		// 检查用户是否仍然存在于Redis中
		userKey := fmt.Sprintf("jwt:%s", claims.UserID)
		exists, err := initRedis.Exists(ctx, userKey).Result()
		if err != nil || exists == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户会话已过期"})
			c.Abort()
			return
		}

		// 更新用户最后活跃时间
		initRedis.HSet(ctx, userKey, "last_active", time.Now().Unix())
		initRedis.Expire(ctx, userKey, AccessTokenExpiration)

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// RefreshAccessToken 使用刷新令牌刷新访问令牌
func RefreshAccessToken(refreshTokenString string) (*TokenPair, error) {
	//1.解析令牌
	claims, err := ParseToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("刷新令牌解析失败: %v", err)
	}
	//1.2验证令牌类型
	if claims.Type != "refresh" {
		return nil, errors.New("刷新令牌无效")
	}
	//2.从缓存中获取对应的刷新令牌
	refreshKey := fmt.Sprintf("refresh:%s", claims.UserID)
	storedRefreshToken, err := initRedis.Get(ctx, refreshKey).Result()
	if err != nil {
		log.Printf("无法获取刷新令牌: %v", err)
		return nil, errors.New("刷新令牌已过期")
	}
	//3.对比刷新令牌
	if storedRefreshToken != refreshTokenString {
		log.Printf("刷新令牌不匹配")
		// 立即撤销所有该用户的令牌
		err = RevokeAllUserTokens(claims.UserID)
		if err != nil {
			log.Printf("撤销用户令牌失败: %v", err)
		}
		return nil, errors.New("检测到安全威胁，所有令牌已被撤销")
	}
	//获取到用户名
	tokenPair, err1 := Login(claims.UserID, claims.Username)
	if err != nil {
		return nil, err1
	}
	return tokenPair, nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法%v\n", token.Header["alg"])
		}
		return []byte(SIGNINGKEY), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// ValidateToken 验证令牌有效性
func ValidateToken(tokenString string) (*Claims, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	//// 检查令牌是否过期,同时检查缓存中的两个令牌是否被删除
	//userKey := fmt.Sprintf("user:%s", claims.UserID)
	//_, err = initRedis.Get(ctx, userKey).Result()
	//if err != nil {
	//	return nil, errors.New("访问令牌已过期")
	//}
	//refreshKey := fmt.Sprintf("refresh:%s", claims.UserID)
	//_, err = initRedis.Get(ctx, refreshKey).Result()
	//if err != nil {
	//	return nil, errors.New("刷新令牌已过期")
	//}
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, errors.New("令牌已过期")
	}

	return claims, nil
}

// RevokeAllUserTokens 撤销用户所有令牌的辅助函数
func RevokeAllUserTokens(UUID string) error {
	userKey := fmt.Sprintf("jwt:%s", UUID)
	refreshKey := fmt.Sprintf("refresh:%s", UUID)

	// 删除用户会话信息
	err := initRedis.Del(ctx, userKey).Err()
	if err != nil {
		return fmt.Errorf("删除用户会话失败: %v", err)
	}

	// 删除刷新令牌
	err = initRedis.Del(ctx, refreshKey).Err()
	if err != nil {
		return fmt.Errorf("删除刷新令牌失败: %v", err)
	}

	log.Printf("🔒 用户 %s 的所有令牌已被撤销", UUID)
	return nil
}

// ParseUUID 解析用户uuid
func ParseUUID(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	claims, err := ParseToken(parts[1])
	if err != nil {
		return ""
	}
	return claims.UserID
}
