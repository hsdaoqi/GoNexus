package middleware

import (
	"github.com/gin-gonic/gin"
	"go-nexus/pkg/response"
	"go-nexus/pkg/utils"
	"strings"
)

// Auth JWT鉴权中间件
func Auth() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 1. 获取 Authorization Header
		// 标准格式：Authorization: Bearer <token>
		tokenString := c.GetHeader("Authorization")

		// 2. 校验格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			// 401 Unauthorized
			response.Error(c, 401, response.ErrAuthFailed, "未登录或非法访问")
			c.Abort() // ⛔ 阻止后续处理，直接退回
			return
		}

		// 3. 截取 Token 部分 (去掉 "Bearer " 前缀)
		// tokenString[7:] 意思是取第7个字符及其后面的所有内容
		tokenString = tokenString[7:]

		// 4. 解析 Token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			response.Error(c, 401, response.ErrAuthFailed, "Token 已失效")
			c.Abort()
			return
		}

		// 5. 【关键】将当前请求的 UserID 信息保存到 Context 中
		// 后续的 Controller 可以通过 c.Get("userID") 拿到这个值
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		// 6. 放行
		c.Next()
	}

}
