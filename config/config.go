// Package config 提供基于 Viper 的配置加载功能，支持 YAML 文件和环境变量覆盖。
package config

//go:generate go run ../cmd/schema-gen/main.go

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 是应用全局配置结构体。
type Config struct {
	Server    ServerConfig    `json:"server" mapstructure:"server"`
	Database  DatabaseConfig  `json:"database" mapstructure:"database"`
	JWT       JWTConfig       `json:"jwt" mapstructure:"jwt"`
	RateLimit RateLimitConfig `json:"rate_limit" mapstructure:"rate_limit"`
	Admin     AdminConfig     `json:"admin" mapstructure:"admin"`
	Chains    ChainsConfig    `json:"chains" mapstructure:"chains"`
	Log       LogConfig       `json:"log" mapstructure:"log"`
	Launchpad LaunchpadConfig `json:"launchpad" mapstructure:"launchpad"`
	Quote     QuoteConfig     `json:"quote" mapstructure:"quote"`
	Sync      SyncConfig      `json:"sync" mapstructure:"sync"`
}

// ServerConfig 包含 HTTP 服务器配置。
type ServerConfig struct {
	Port int    `json:"port" mapstructure:"port"`
	Mode string `json:"mode" mapstructure:"mode"`
}

// DatabaseConfig 包含数据库连接和连接池配置。
type DatabaseConfig struct {
	Host            string        `json:"host" mapstructure:"host"`
	Port            int           `json:"port" mapstructure:"port"`
	User            string        `json:"user" mapstructure:"user"`
	Password        string        `json:"password" mapstructure:"password"`
	DBName          string        `json:"dbname" mapstructure:"dbname"`
	SSLMode         string        `json:"sslmode" mapstructure:"sslmode"`
	MaxOpenConns    int           `json:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns" mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time" mapstructure:"conn_max_idle_time"`
	LogLevel        string        `json:"log_level" mapstructure:"log_level"`
	SlowThreshold   time.Duration `json:"slow_threshold" mapstructure:"slow_threshold"`
}

// JWTConfig 包含 JWT 签发配置。
type JWTConfig struct {
	Secret       string        `json:"secret" mapstructure:"secret"`
	WalletExpire time.Duration `json:"wallet_expire" mapstructure:"wallet_expire"`
	AdminExpire  time.Duration `json:"admin_expire" mapstructure:"admin_expire"`
}

// RateLimitConfig 包含限流配置。
type RateLimitConfig struct {
	IP      RateLimitRule `json:"ip" mapstructure:"ip"`
	Account RateLimitRule `json:"account" mapstructure:"account"`
}

// RateLimitRule 定义单个维度的限流规则。
type RateLimitRule struct {
	Rate  float64 `json:"rate" mapstructure:"rate"`
	Burst int     `json:"burst" mapstructure:"burst"`
}

// AdminConfig 包含初始管理员 seed 配置。
type AdminConfig struct {
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

// ChainsConfig 包含区块链相关配置。
type ChainsConfig struct {
	AllowedChainIDs []int `json:"allowed_chain_ids" mapstructure:"allowed_chain_ids"`
}

// LogConfig 包含日志配置。
type LogConfig struct {
	// Level 日志级别：debug/info/warn/error。
	Level string `json:"level" mapstructure:"level"`
	// Format 日志格式：json/text。
	Format string `json:"format" mapstructure:"format"`
	// AddSource 是否记录源码调用位置（文件名:行号）。
	AddSource bool `json:"add_source" mapstructure:"add_source"`
}

// QuoteConfig 包含 Quote 报价聚合模块配置。
type QuoteConfig struct {
	// CacheTTL 缓存 TTL（秒）。
	CacheTTL int `json:"cache_ttl" mapstructure:"cache_ttl"`
	// Providers 各供应商配置。
	Providers map[string]ProviderConfig `json:"providers" mapstructure:"providers"`
}

// ProviderConfig 单个报价供应商的配置。
type ProviderConfig struct {
	// Enabled 是否启用。
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	// APIKey API 密钥。
	APIKey string `json:"api_key" mapstructure:"api_key"`
	// SecretKey API Secret（OKX 专用）。
	SecretKey string `json:"secret_key" mapstructure:"secret_key"`
	// Passphrase API 密钥密码短语（OKX 专用）。
	Passphrase string `json:"passphrase" mapstructure:"passphrase"`
	// BaseURL 供应商 API 基础 URL。
	BaseURL string `json:"base_url" mapstructure:"base_url"`
}

// LaunchpadConfig 包含 Launchpad IDO 模块配置。
type LaunchpadConfig struct {
	// RPCURL 以太坊 RPC 节点地址。
	RPCURL string `json:"rpc_url" mapstructure:"rpc_url"`
	// ChainID 链 ID。
	ChainID int `json:"chain_id" mapstructure:"chain_id"`
	// MousePadByTierABI MousePadByTier 合约 ABI 文件路径。
	MousePadByTierABI string `json:"mouse_pad_by_tier_abi" mapstructure:"mouse_pad_by_tier_abi"`
	// MouseTierABI MouseTier 合约 ABI 文件路径。
	MouseTierABI string `json:"mouse_tier_abi" mapstructure:"mouse_tier_abi"`
	// MouseTierAddress MouseTier 合约地址。
	MouseTierAddress string `json:"mouse_tier_address" mapstructure:"mouse_tier_address"`
	// MousePadByTierAddress MousePadByTier 合约地址。
	MousePadByTierAddress string `json:"mouse_pad_by_tier_address" mapstructure:"mouse_pad_by_tier_address"`
	// DeployerAddress Deployer 工厂合约地址，用于 createSale 操作的 to 地址。
	DeployerAddress string `json:"deployer_address" mapstructure:"deployer_address"`
	// AdminPrivateKey 管理员钱包私钥，仅开发环境用于签名广播交易。
	AdminPrivateKey string `json:"admin_private_key" mapstructure:"admin_private_key"`
	// PrepareExpiresIn Prepare 交易过期时间。
	PrepareExpiresIn time.Duration `json:"prepare_expires_in" mapstructure:"prepare_expires_in"`
	// PollInterval 兜底轮询间隔。
	PollInterval time.Duration `json:"poll_interval" mapstructure:"poll_interval"`
	// RPCTimeout RPC 调用超时。
	RPCTimeout time.Duration `json:"rpc_timeout" mapstructure:"rpc_timeout"`
}

// SyncConfig 包含链上同步框架配置。
type SyncConfig struct {
	// Chains 多链同步配置，key 为 chain_id。
	Chains []SyncChainConfig `json:"chains" mapstructure:"chains"`
	// Backfill 历史回填配置。
	Backfill BackfillConfig `json:"backfill" mapstructure:"backfill"`
	// Projector 异步投影配置。
	Projector ProjectorConfig `json:"projector" mapstructure:"projector"`
	// Worker Worker 进程配置。
	Worker WorkerConfig `json:"worker" mapstructure:"worker"`
}

// SyncChainConfig 单条链的同步配置。
type SyncChainConfig struct {
	// ChainID 链 ID（1=ETH 主网, 5=Goerli, 11155111=Sepolia）。
	ChainID int `json:"chain_id" mapstructure:"chain_id"`
	// StartBlock 冷启动回填起始区块号。
	StartBlock int64 `json:"start_block" mapstructure:"start_block"`
	// ConfirmationBlocks 确认区块数（防止链重组）。
	ConfirmationBlocks int64 `json:"confirmation_blocks" mapstructure:"confirmation_blocks"`
	// BlockInterval 出块间隔（用于 WS 僵死检测）。
	BlockInterval time.Duration `json:"block_interval" mapstructure:"block_interval"`
	// ProcessorID 处理器标识（如 launchpad）。
	ProcessorID string `json:"processor_id" mapstructure:"processor_id"`
	// Contracts 合约地址配置。
	Contracts SyncContractsConfig `json:"contracts" mapstructure:"contracts"`
	// Nodes RPC 节点列表（按 priority 排序）。
	Nodes []ChainNodeConfig `json:"nodes" mapstructure:"nodes"`
}

// SyncContractsConfig 合约地址配置。
type SyncContractsConfig struct {
	// MouseTier MouseTier 合约地址。
	MouseTier string `json:"mouse_tier" mapstructure:"mouse_tier"`
	// MousePadByTier MousePadByTier 合约地址。
	MousePadByTier string `json:"mouse_pad_by_tier" mapstructure:"mouse_pad_by_tier"`
}

// ChainNodeConfig 单个 RPC 节点配置。
type ChainNodeConfig struct {
	// Name 节点名称（用于日志和指标标签）。
	Name string `json:"name" mapstructure:"name"`
	// WSURL WebSocket URL。
	WSURL string `json:"ws_url" mapstructure:"ws_url"`
	// HTTPURL HTTP RPC URL。
	HTTPURL string `json:"http_url" mapstructure:"http_url"`
	// Priority 优先级（数值越小优先级越高）。
	Priority int `json:"priority" mapstructure:"priority"`
	// Timeout RPC 调用超时。
	Timeout time.Duration `json:"timeout" mapstructure:"timeout"`
	// RateLimit 令牌桶限流（RPS），0 表示不限流。
	RateLimit float64 `json:"rate_limit" mapstructure:"rate_limit"`
	// CircuitBreaker 熔断器配置。
	CircuitBreaker CircuitBreakerConfig `json:"circuit_breaker" mapstructure:"circuit_breaker"`
}

// CircuitBreakerConfig 熔断器配置。
type CircuitBreakerConfig struct {
	// FailureThreshold 连续失败次数阈值，触发熔断。
	FailureThreshold uint32 `json:"failure_threshold" mapstructure:"failure_threshold"`
	// Timeout OPEN 状态持续时间。
	Timeout time.Duration `json:"timeout" mapstructure:"timeout"`
	// MaxRequests HALF-OPEN 状态最大探测请求数。
	MaxRequests uint32 `json:"max_requests" mapstructure:"max_requests"`
}

// BackfillConfig 历史回填配置。
type BackfillConfig struct {
	// InitialBatchSize 初始批次大小。
	InitialBatchSize int `json:"initial_batch_size" mapstructure:"initial_batch_size"`
	// MinBatchSize 最小批次大小。
	MinBatchSize int `json:"min_batch_size" mapstructure:"min_batch_size"`
	// MaxBatchSize 最大批次大小。
	MaxBatchSize int `json:"max_batch_size" mapstructure:"max_batch_size"`
	// GrowthFactor 成功时批次增长因子。
	GrowthFactor float64 `json:"growth_factor" mapstructure:"growth_factor"`
}

// ProjectorConfig 异步投影配置。
type ProjectorConfig struct {
	// MaxWorkers 最大并发 worker 数。
	MaxWorkers int `json:"max_workers" mapstructure:"max_workers"`
	// MaxRetries 最大重试次数。
	MaxRetries int `json:"max_retries" mapstructure:"max_retries"`
	// ProcessingTimeout processing 状态超时时间。
	ProcessingTimeout time.Duration `json:"processing_timeout" mapstructure:"processing_timeout"`
	// PollInterval 轮询间隔。
	PollInterval time.Duration `json:"poll_interval" mapstructure:"poll_interval"`
}

// WorkerConfig Worker 进程配置。
type WorkerConfig struct {
	// ShutdownTimeout 优雅关停超时时间。
	ShutdownTimeout time.Duration `json:"shutdown_timeout" mapstructure:"shutdown_timeout"`
	// HealthPort 管理端口（健康检查 + Prometheus 指标）。
	HealthPort int `json:"health_port" mapstructure:"health_port"`
}

// Load 从指定路径加载配置文件，环境变量使用 MOUSECAKE_ 前缀覆盖。
// 配置文件不存在时不报错，使用默认值。
// 支持 ${ENV_VAR} 语法在 YAML 中引用环境变量，启动时自动展开。
// 启动前自动加载项目根目录的 .env 文件，已存在的环境变量不会被覆盖。
func Load(path string) (*Config, error) {
	// 加载 .env 文件（不覆盖已存在的环境变量，文件不存在也不报错）
	_ = godotenv.Load()

	v := viper.New()

	setDefaults(v)

	v.SetEnvPrefix("MOUSECAKE")
	v.SetEnvKeyReplacer(NewEnvReplacer())
	v.AutomaticEnv()

	// Viper 的 AutomaticEnv 对嵌套 key 需要显式绑定
	bindEnvVars(v)

	raw, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("读取配置文件 %s 失败: %w", path, err)
		}
		// 文件不存在，使用默认值
	} else {
		// 展开 YAML 中的 ${ENV_VAR} 引用
		expanded := os.ExpandEnv(string(raw))
		v.SetConfigType("yaml")
		if err := v.ReadConfig(strings.NewReader(expanded)); err != nil {
			return nil, fmt.Errorf("解析配置文件 %s 失败: %w", path, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate 验证配置中的必要字段。
func (c *Config) Validate() error {
	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret 不能为空，请设置环境变量 MOUSECAKE_JWT_SECRET")
	}
	if c.Admin.Password == "" {
		return fmt.Errorf("admin.password 不能为空，请设置环境变量 MOUSECAKE_ADMIN_PASSWORD")
	}
	if len(c.Chains.AllowedChainIDs) == 0 {
		return fmt.Errorf("chains.allowed_chain_ids 不能为空")
	}
	return nil
}

// ValidateSync 验证同步配置的完整性（Worker 进程启动时调用）。
func (c *Config) ValidateSync() error {
	if len(c.Sync.Chains) == 0 {
		return fmt.Errorf("sync.chains 不能为空")
	}
	for i, chain := range c.Sync.Chains {
		if chain.ChainID == 0 {
			return fmt.Errorf("sync.chains[%d]: chain_id 不能为空", i)
		}
		if chain.StartBlock < 0 {
			return fmt.Errorf("sync.chains[%d]: start_block 不能为负数", i)
		}
		if len(chain.Nodes) == 0 {
			return fmt.Errorf("sync.chains[%d]: nodes 不能为空", i)
		}
		if chain.Contracts.MouseTier == "" && chain.Contracts.MousePadByTier == "" {
			return fmt.Errorf("sync.chains[%d]: contracts 至少需要一个合约地址", i)
		}
		for j, node := range chain.Nodes {
			if node.WSURL == "" && node.HTTPURL == "" {
				return fmt.Errorf("sync.chains[%d].nodes[%d]: 至少需要 ws_url 或 http_url", i, j)
			}
		}
	}
	return nil
}

// NewEnvReplacer 创建环境变量 key 中 . 到 _ 的替换器。
func NewEnvReplacer() *strings.Replacer {
	return strings.NewReplacer(".", "_")
}

func bindEnvVars(v *viper.Viper) {
	envKeys := []string{
		"server.port", "server.mode",
		"database.host", "database.port", "database.user", "database.password",
		"database.dbname", "database.sslmode", "database.max_open_conns",
		"database.max_idle_conns", "database.conn_max_lifetime", "database.conn_max_idle_time",
		"database.log_level", "database.slow_threshold",
		"jwt.secret", "jwt.wallet_expire", "jwt.admin_expire",
		"rate_limit.ip.rate", "rate_limit.ip.burst",
		"rate_limit.account.rate", "rate_limit.account.burst",
		"admin.username", "admin.password",
		"log.level", "log.format", "log.add_source",
		"launchpad.rpc_url", "launchpad.chain_id", "launchpad.mouse_pad_by_tier_abi",
		"launchpad.mouse_tier_abi", "launchpad.mouse_tier_address", "launchpad.mouse_pad_by_tier_address",
		"launchpad.deployer_address", "launchpad.admin_private_key",
		"launchpad.prepare_expires_in", "launchpad.poll_interval", "launchpad.rpc_timeout",
		"quote.cache_ttl",
	}
	for _, key := range envKeys {
		v.MustBindEnv(key)
	}
}

func setDefaults(v *viper.Viper) {
	// 服务器
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")

	// 数据库
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", 30*time.Minute)
	v.SetDefault("database.conn_max_idle_time", 5*time.Minute)
	v.SetDefault("database.log_level", "info")
	v.SetDefault("database.slow_threshold", 200*time.Millisecond)

	// JWT
	v.SetDefault("jwt.wallet_expire", 15*time.Minute)
	v.SetDefault("jwt.admin_expire", 15*time.Minute)

	// 限流
	v.SetDefault("rate_limit.ip.rate", 20)
	v.SetDefault("rate_limit.ip.burst", 30)
	v.SetDefault("rate_limit.account.rate", 10)
	v.SetDefault("rate_limit.account.burst", 15)

	// 链
	v.SetDefault("chains.allowed_chain_ids", []int{1, 5, 11155111})

	// 日志
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.add_source", false)

	// Launchpad
	v.SetDefault("launchpad.prepare_expires_in", 30*time.Minute)
	v.SetDefault("launchpad.poll_interval", 1*time.Minute)
	v.SetDefault("launchpad.rpc_timeout", 10*time.Second)

	// Quote
	v.SetDefault("quote.cache_ttl", 10)

	// Sync
	v.SetDefault("sync.backfill.initial_batch_size", 5000)
	v.SetDefault("sync.backfill.min_batch_size", 500)
	v.SetDefault("sync.backfill.max_batch_size", 10000)
	v.SetDefault("sync.backfill.growth_factor", 1.2)
	v.SetDefault("sync.projector.max_workers", 4)
	v.SetDefault("sync.projector.max_retries", 3)
	v.SetDefault("sync.projector.processing_timeout", 5*time.Minute)
	v.SetDefault("sync.projector.poll_interval", 1*time.Second)
	v.SetDefault("sync.worker.shutdown_timeout", 30*time.Second)
	v.SetDefault("sync.worker.health_port", 9090)
}
