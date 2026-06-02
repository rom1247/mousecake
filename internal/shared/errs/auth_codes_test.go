package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetErrorMessage_KnownCodes(t *testing.T) {
	cases := []struct {
		code    int
		message string
	}{
		{CodeAuthFailed, "认证失败"},
		{CodeNonceExpired, "nonce 已过期或不存在"},
		{CodeSignInvalid, "签名验证失败"},
		{CodeCredentialBad, "用户名或密码错误"},
		{CodeAccountDisabled, "账号已禁用"},
		{CodeTokenInvalid, "Token 无效或已过期"},
		{CodeRateLimited, "请求过于频繁"},
		{CodeSIWEFormat, "SIWE 消息格式错误"},
		{CodeChainUnsupported, "不支持的链 ID"},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.message, GetErrorMessage(tc.code))
	}
}

func TestGetErrorMessage_UnknownCode(t *testing.T) {
	assert.Equal(t, "未知错误", GetErrorMessage(99999))
}
