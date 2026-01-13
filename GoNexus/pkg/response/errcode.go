package response

// 1. 定义错误码 (const)
const (
	CodeSuccess = 200

	// 1000~1999: 通用错误
	ErrSystemError  = 1001 // 系统内部错误 (Database error, etc.)
	ErrParamInvalid = 1002 // 参数校验失败

	// 2000~2999: 用户模块错误
	ErrUserExist     = 2001 // 用户已存在
	ErrUserNotExist  = 2002 // 用户不存在
	ErrPasswordError = 2003 // 密码错误
	ErrAuthFailed    = 2004 // 鉴权失败 (Token 无效)
)

// 2. 定义错误码对应的消息 (Map)
var msgMap = map[int]string{
	CodeSuccess:      "success",
	ErrSystemError:   "系统内部错误",
	ErrParamInvalid:  "参数非法",
	ErrUserExist:     "该用户名已被注册",
	ErrUserNotExist:  "用户不存在",
	ErrPasswordError: "密码错误",
	ErrAuthFailed:    "认证失败，请重新登录",
}

// GetMsg 辅助函数：给个 code，还你个 string
func GetMsg(code int) string {
	msg, ok := msgMap[code]
	if ok {
		return msg
	}
	return msgMap[ErrSystemError] // 找不到就返回系统错误
}
