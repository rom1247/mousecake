package quote

// Quote 模块错误码（50xxx 段）。
const (
	// ErrCodeProviderNotFound 供应商不存在。
	ErrCodeProviderNotFound = 50001
	// ErrCodeInvalidParam 参数校验失败。
	ErrCodeInvalidParam = 50002
	// ErrCodeSwapNotFound swap 记录不存在。
	ErrCodeSwapNotFound = 50003
	// ErrCodeAlreadySubmitted swap 已提交。
	ErrCodeAlreadySubmitted = 50004
	// ErrCodeProviderError 供应商请求失败。
	ErrCodeProviderError = 50005
)
