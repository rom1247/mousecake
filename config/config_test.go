package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_FromYAML(t *testing.T) {
	// 显式设置环境变量，确保测试不依赖 .env 文件
	t.Setenv("MOUSECAKE_DATABASE_HOST", "localhost")
	t.Setenv("MOUSECAKE_DATABASE_PASSWORD", "test_password")
	t.Setenv("MOUSECAKE_JWT_SECRET", "test-jwt-secret")
	t.Setenv("MOUSECAKE_ADMIN_PASSWORD", "test_admin_pw")

	cfg, err := Load("config/app.yaml")
	require.NoError(t, err)

	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "test_password", cfg.Database.Password)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Server.Mode)
	assert.Equal(t, "test-jwt-secret", cfg.JWT.Secret)
	assert.Equal(t, 16*time.Minute, cfg.JWT.WalletExpire)
	assert.Equal(t, 15*time.Minute, cfg.JWT.AdminExpire)
	assert.Equal(t, "admin", cfg.Admin.Username)
	assert.Equal(t, "test_admin_pw", cfg.Admin.Password)
	assert.Equal(t, []int{1, 5, 11155111}, cfg.Chains.AllowedChainIDs)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, false, cfg.Log.AddSource)
}

func TestLoad_EnvOverride(t *testing.T) {
	t.Setenv("MOUSECAKE_DATABASE_HOST", "db.production")
	t.Setenv("MOUSECAKE_SERVER_PORT", "9090")
	t.Setenv("MOUSECAKE_LOG_LEVEL", "debug")
	t.Setenv("MOUSECAKE_LOG_FORMAT", "text")

	cfg, err := Load("config/app.yaml")
	require.NoError(t, err)

	assert.Equal(t, "db.production", cfg.Database.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "text", cfg.Log.Format)
}

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load("nonexistent.yaml")
	require.NoError(t, err)

	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Server.Mode)
	assert.Equal(t, 25, cfg.Database.MaxOpenConns)
	assert.Equal(t, 10, cfg.Database.MaxIdleConns)
	assert.Equal(t, 30*time.Minute, cfg.Database.ConnMaxLifetime)
	assert.Equal(t, 5*time.Minute, cfg.Database.ConnMaxIdleTime)
	assert.Equal(t, 15*time.Minute, cfg.JWT.WalletExpire)
	assert.Equal(t, 15*time.Minute, cfg.JWT.AdminExpire)
	assert.Equal(t, float64(20), cfg.RateLimit.IP.Rate)
	assert.Equal(t, 30, cfg.RateLimit.IP.Burst)
	assert.Equal(t, float64(10), cfg.RateLimit.Account.Rate)
	assert.Equal(t, 15, cfg.RateLimit.Account.Burst)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, false, cfg.Log.AddSource)
}

func TestLoad_MissingFile(t *testing.T) {
	cfg, err := Load("nonexistent.yaml")
	require.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestMain(m *testing.M) {
	os.Chdir("..")
	os.Exit(m.Run())
}

func TestValidateSync(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{
			name: "同步配置为空时报错",
			config: Config{
				Sync: SyncConfig{},
			},
			wantErr: "sync.chains 不能为空",
		},
		{
			name: "缺少 chain_id 报错",
			config: Config{
				Sync: SyncConfig{
					Chains: []SyncChainConfig{
						{StartBlock: 100, Nodes: []ChainNodeConfig{{HTTPURL: "http://rpc"}},
							Contracts: SyncContractsConfig{MouseTier: "0x1"}},
					},
				},
			},
			wantErr: "chain_id 不能为空",
		},
		{
			name: "缺少 start_block 报错",
			config: Config{
				Sync: SyncConfig{
					Chains: []SyncChainConfig{
						{ChainID: 1, StartBlock: -1, Nodes: []ChainNodeConfig{{HTTPURL: "http://rpc"}},
							Contracts: SyncContractsConfig{MouseTier: "0x1"}},
					},
				},
			},
			wantErr: "start_block 不能为负数",
		},
		{
			name: "缺少节点报错",
			config: Config{
				Sync: SyncConfig{
					Chains: []SyncChainConfig{
						{ChainID: 1, StartBlock: 100,
							Contracts: SyncContractsConfig{MouseTier: "0x1"}},
					},
				},
			},
			wantErr: "nodes 不能为空",
		},
		{
			name: "缺少合约地址报错",
			config: Config{
				Sync: SyncConfig{
					Chains: []SyncChainConfig{
						{ChainID: 1, StartBlock: 100, Nodes: []ChainNodeConfig{{HTTPURL: "http://rpc"}}},
					},
				},
			},
			wantErr: "contracts 至少需要一个合约地址",
		},
		{
			name: "缺少节点 URL 报错",
			config: Config{
				Sync: SyncConfig{
					Chains: []SyncChainConfig{
						{ChainID: 1, StartBlock: 100,
							Contracts: SyncContractsConfig{MouseTier: "0x1"},
							Nodes:     []ChainNodeConfig{{Name: "test"}}},
					},
				},
			},
			wantErr: "至少需要 ws_url 或 http_url",
		},
		{
			name: "有效配置通过验证",
			config: Config{
				Sync: SyncConfig{
					Chains: []SyncChainConfig{
						{
							ChainID:    1,
							StartBlock: 100,
							Contracts:  SyncContractsConfig{MouseTier: "0x1"},
							Nodes:      []ChainNodeConfig{{Name: "test", HTTPURL: "http://rpc"}},
						},
					},
				},
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.config.ValidateSync()
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}
