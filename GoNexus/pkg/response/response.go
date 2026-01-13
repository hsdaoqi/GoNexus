package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Result 基础返回方法
func Result(c *gin.Context, httpStatus int, code int, data interface{}, msg string) {
	c.JSON(httpStatus, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// Success 成功 (默认 200)
func Success(c *gin.Context, data interface{}) {
	Result(c, http.StatusOK, CodeSuccess, data, GetMsg(CodeSuccess))
}

// Fail 失败 (业务错误，HTTP 状态码依然返回 200，由前端根据 Code 处理)
// 用法：response.Fail(c, response.ErrPasswordError)
// 优点：不需要你手动写 "密码错误" 字符串了！
func Fail(c *gin.Context, code int) {
	Result(c, http.StatusOK, code, nil, GetMsg(code))
}

// FailWithMessage 如果你想自定义错误信息，可以用这个
func FailWithMessage(c *gin.Context, code int, message string) {
	Result(c, http.StatusOK, code, nil, message)
}

// Error HTTP 错误响应 (HTTP 4xx/5xx)
// 用法：response.Error(c, http.StatusInternalServerError, response.ErrSystemError, "数据库炸了")
// 场景：鉴权失败(401)、服务器崩溃(500)、参数绑定失败(400)
func Error(c *gin.Context, httpCode int, code int, message string) {
	Result(c, httpCode, code, nil, message)
}
