// Package errs 定义业务错误码常量。
package errs

// 认证相关错误码段（40100-40108）
const (
	CodeAuthFailed       = 40100 // 认证失败（通用）
	CodeNonceExpired     = 40101 // nonce 已过期或不存在
	CodeSignInvalid      = 40102 // 签名验证失败
	CodeCredentialBad    = 40103 // 用户名或密码错误
	CodeAccountDisabled  = 40104 // 账号已禁用
	CodeTokenInvalid     = 40105 // Token 无效或已过期
	CodeRateLimited      = 40106 // 请求过于频繁
	CodeSIWEFormat       = 40107 // SIWE 消息格式错误
	CodeChainUnsupported = 40108 // 不支持的链 ID
	CodeInternal         = 50000 // 服务端内部错误
)

// 认证错误码对应的消息映射
var codeMessages = map[int]string{
	CodeAuthFailed:       "认证失败",
	CodeNonceExpired:     "nonce 已过期或不存在",
	CodeSignInvalid:      "签名验证失败",
	CodeCredentialBad:    "用户名或密码错误",
	CodeAccountDisabled:  "账号已禁用",
	CodeTokenInvalid:     "Token 无效或已过期",
	CodeRateLimited:      "请求过于频繁",
	CodeSIWEFormat:       "SIWE 消息格式错误",
	CodeChainUnsupported: "不支持的链 ID",
	CodeInternal:         "服务端内部错误",
}

// GetErrorMessage 根据错误码返回对应的消息文本。
func GetErrorMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
