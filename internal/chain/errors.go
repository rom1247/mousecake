// Package chain 提供 RPC 节点管理基础设施，包括多节点池、熔断器、限流、健康检查。
package chain

import "errors"

// ErrAllNodesUnavailable 表示所有节点均不可用。
var ErrAllNodesUnavailable = errors.New("chain: 所有节点不可用")

// ErrRateLimited 表示节点返回 429 限流响应。
var ErrRateLimited = errors.New("chain: 节点限流")
