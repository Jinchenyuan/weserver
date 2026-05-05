package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"server/utils"

	"github.com/Jinchenyuan/wego"
	"github.com/gin-gonic/gin"
)

const accountIDContextKey = "account_id"

func AuthMiddleware(excludePaths ...string) gin.HandlerFunc {
	excluded := make(map[string]struct{}, len(excludePaths))
	for _, path := range excludePaths {
		excluded[path] = struct{}{}
	}

	return func(c *gin.Context) {
		if _, ok := excluded[c.FullPath()]; ok {
			c.Next()
			return
		}

		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		m := wego.GetGlobalMesa()
		if m == nil || m.Redis == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "runtime redis unavailable"})
			c.Abort()
			return
		}

		cacheKey := fmt.Sprintf("token:%d", claims.AccountID)
		cachedToken, err := m.Redis.Get(context.Background(), cacheKey).Result()
		if err != nil || cachedToken != token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set(accountIDContextKey, claims.AccountID)
		c.Next()
	}
}

func GetAccountID(c *gin.Context) (uint32, bool) {
	value, ok := c.Get(accountIDContextKey)
	if !ok {
		return 0, false
	}
	accountID, ok := value.(uint32)
	return accountID, ok
}
