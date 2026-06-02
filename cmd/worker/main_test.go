// Package main 提供 cmd/worker 的管理 API 测试。
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/mousecake-go/mousecake-go/internal/shared/sync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mousecake-go/mousecake-go/config"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setupWorkerTestDB 创建 SQLite 内存数据库并自动迁移同步相关表。
func setupWorkerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "打开 SQLite 内存数据库失败")
	err = db.AutoMigrate(&sync.ChainEvent{}, &sync.Checkpoint{})
	require.NoError(t, err, "AutoMigrate 失败")
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
	return db
}

// insertTestEvent 插入一条测试事件并返回其 ID。
func insertTestEvent(t *testing.T, db *gorm.DB, status sync.EventStatus) int64 {
	t.Helper()
	event := sync.ChainEvent{
		ChainID:         1,
		BlockNumber:     100,
		TxHash:          "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		TxIndex:         0,
		LogIndex:        0,
		ContractAddress: "0x1234567890abcdef1234567890abcdef12345678",
		EventName:       "TestEvent",
		EventData:       `{"key":"value"}`,
		Status:          status,
		ProcessorID:     "launchpad",
	}
	if status == sync.StatusFailed || status == sync.StatusDeadLetter {
		now := time.Now()
		errMsg := "测试错误"
		event.ErrorMessage = &errMsg
		event.LastFailedAt = &now
	}
	err := db.Create(&event).Error
	require.NoError(t, err, "插入测试事件失败")
	return event.ID
}

// TestHealthzHandler 测试存活探针返回 200 和正确的 JSON 响应。
func TestHealthzHandler(t *testing.T) {
	t.Parallel()

	router := gin.New()
	router.GET("/healthz", healthzHandler)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "ok", body["status"])
}

// TestHealthzHandler_MethodNotAllowed 测试非 GET 方法访问 healthz 返回 404。
func TestHealthzHandler_MethodNotAllowed(t *testing.T) {
	t.Parallel()

	router := gin.New()
	router.GET("/healthz", healthzHandler)

	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Gin 对未注册的方法返回 404
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestMetricsHandler 测试 Prometheus 指标端点返回 200 和 text/plain 内容类型。
func TestMetricsHandler(t *testing.T) {
	t.Parallel()

	router := gin.New()
	router.GET("/metrics", metricsHandler)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Prometheus handler 默认返回 text/plain 内容
	contentType := w.Header().Get("Content-Type")
	assert.True(t, strings.Contains(contentType, "text/plain"),
		"期望 Content-Type 包含 text/plain，实际为 %s", contentType)
}

// TestSyncRetryHandler_MissingID 测试重试请求缺少 ID 参数时返回 400。
// 传 nil EventStore 是安全的：参数绑定失败时 handler 在访问 store 之前就返回了。
func TestSyncRetryHandler_MissingID(t *testing.T) {
	t.Parallel()

	var store *sync.EventStore // nil，参数校验路径不会访问 store

	tests := []struct {
		name           string
		body           string
		wantStatus     int
		wantErrContain string
	}{
		{
			name:           "空 JSON",
			body:           `{}`,
			wantStatus:     http.StatusBadRequest,
			wantErrContain: "缺少 id 参数",
		},
		{
			name:           "请求体为空",
			body:           ``,
			wantStatus:     http.StatusBadRequest,
			wantErrContain: "缺少 id 参数",
		},
		{
			name:           "id 为零值",
			body:           `{"id":0}`,
			wantStatus:     http.StatusBadRequest,
			wantErrContain: "缺少 id 参数",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := gin.New()
			r.POST("/admin/sync/events/retry", syncRetryHandler(store))

			req := httptest.NewRequest(http.MethodPost, "/admin/sync/events/retry",
				strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var resp map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Contains(t, resp["error"], tt.wantErrContain)
		})
	}
}

// TestSyncRetryHandler_EventNotFound 测试重试不存在的事件返回 404。
func TestSyncRetryHandler_EventNotFound(t *testing.T) {
	t.Parallel()

	db := setupWorkerTestDB(t)
	store := sync.NewEventStore(db)

	r := gin.New()
	r.POST("/admin/sync/events/retry", syncRetryHandler(store))

	body := `{"id":999}`
	req := httptest.NewRequest(http.MethodPost, "/admin/sync/events/retry",
		strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "事件不存在")
}

// TestSyncRetryHandler_AlreadyPending 测试重试已经是 pending 状态的事件返回 500。
// ResetToPending 要求状态为 failed 或 dead_letter，pending 不匹配会返回错误。
func TestSyncRetryHandler_AlreadyPending(t *testing.T) {
	t.Parallel()

	db := setupWorkerTestDB(t)
	store := sync.NewEventStore(db)

	// 插入一个 pending 状态的事件
	eventID := insertTestEvent(t, db, sync.StatusPending)

	r := gin.New()
	r.POST("/admin/sync/events/retry", syncRetryHandler(store))

	body := fmt.Sprintf(`{"id":%d}`, eventID)
	req := httptest.NewRequest(http.MethodPost, "/admin/sync/events/retry",
		strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// ResetToPending 对 pending 状态事件返回错误（状态不允许重置）
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "不存在或状态不允许重置")
}

// TestSyncRetryHandler_Success 测试成功重试 failed 状态的事件。
func TestSyncRetryHandler_Success(t *testing.T) {
	t.Parallel()

	db := setupWorkerTestDB(t)
	store := sync.NewEventStore(db)

	// 插入一个 failed 状态的事件
	eventID := insertTestEvent(t, db, sync.StatusFailed)

	r := gin.New()
	r.POST("/admin/sync/events/retry", syncRetryHandler(store))

	body := fmt.Sprintf(`{"id":%d}`, eventID)
	req := httptest.NewRequest(http.MethodPost, "/admin/sync/events/retry",
		strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["message"], "事件已重置为 pending")
}

// TestSyncEventsHandler_DeadLetter 测试查看 dead_letter 事件列表，验证分页和总数。
func TestSyncEventsHandler_DeadLetter(t *testing.T) {
	t.Parallel()

	db := setupWorkerTestDB(t)
	store := sync.NewEventStore(db)

	// 插入 15 个 dead_letter 事件
	for i := 0; i < 15; i++ {
		event := sync.ChainEvent{
			ChainID:         1,
			BlockNumber:     int64(100 + i),
			TxHash:          fmt.Sprintf("0x%064d", i),
			TxIndex:         0,
			LogIndex:        0,
			ContractAddress: "0x1234567890abcdef1234567890abcdef12345678",
			EventName:       "TestEvent",
			EventData:       `{"key":"value"}`,
			Status:          sync.StatusDeadLetter,
			ProcessorID:     "launchpad",
		}
		errMsg := "测试死信"
		now := time.Now()
		event.ErrorMessage = &errMsg
		event.LastFailedAt = &now
		err := db.Create(&event).Error
		require.NoError(t, err, "插入 dead_letter 事件失败")
	}

	r := gin.New()
	r.GET("/admin/sync/events", syncEventsHandler(store))

	req := httptest.NewRequest(http.MethodGet,
		"/admin/sync/events?status=dead_letter&chain_id=1&page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// 验证 events 数组
	events, ok := resp["events"].([]any)
	assert.True(t, ok, "响应应包含 events 数组")
	assert.Len(t, events, 10, "第一页应返回 10 条记录")

	// 验证 total
	total, ok := resp["total"].(float64)
	assert.True(t, ok, "响应应包含 total 字段")
	assert.Equal(t, float64(15), total, "总数应为 15")
}

// TestSyncEventsHandler_DefaultPending 测试默认查询 pending 事件列表。
// ListPending 使用 processorID 过滤，handler 传空字符串，因此匹配 processorID 为空的事件。
func TestSyncEventsHandler_DefaultPending(t *testing.T) {
	t.Parallel()

	db := setupWorkerTestDB(t)
	store := sync.NewEventStore(db)

	// 插入 3 个 pending 事件，processorID 为空（与 handler 中的查询一致）
	for i := 0; i < 3; i++ {
		event := sync.ChainEvent{
			ChainID:         1,
			BlockNumber:     int64(200 + i),
			TxHash:          fmt.Sprintf("0x%064d", 100+i),
			TxIndex:         0,
			LogIndex:        0,
			ContractAddress: "0x1234567890abcdef1234567890abcdef12345678",
			EventName:       "TestEvent",
			EventData:       `{"key":"value"}`,
			Status:          sync.StatusPending,
			ProcessorID:     "", // handler 中 ListPending 传空字符串
		}
		err := db.Create(&event).Error
		require.NoError(t, err, "插入 pending 事件失败")
	}

	r := gin.New()
	r.GET("/admin/sync/events", syncEventsHandler(store))

	// 不指定 status，默认走 ListPending 逻辑
	req := httptest.NewRequest(http.MethodGet, "/admin/sync/events?chain_id=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	events, ok := resp["events"].([]any)
	assert.True(t, ok, "响应应包含 events 数组")
	assert.Len(t, events, 3, "应返回 3 条 pending 记录")
}

// TestSyncStatusHandler 测试同步状态查询，验证返回包含 chains 字段。
// SyncManager.GetStatus 依赖 EventStore.CountByStatus 和 CheckpointRepository，
// 空 SyncConfig.Chains 时 GetStatus 返回 nil（JSON 序列化为 null）。
func TestSyncStatusHandler(t *testing.T) {
	t.Parallel()

	db := setupWorkerTestDB(t)

	// 使用空的 SyncConfig 创建 SyncManager
	cfg := config.SyncConfig{
		Chains: []config.SyncChainConfig{},
	}
	syncMgr, err := sync.NewSyncManager(cfg, db, nil)
	require.NoError(t, err, "创建 SyncManager 失败")

	r := gin.New()
	r.GET("/admin/sync/status", syncStatusHandler(syncMgr))

	req := httptest.NewRequest(http.MethodGet, "/admin/sync/status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// 无链配置时 GetStatus 返回 nil slice，JSON 序列化为 null
	_, exists := resp["chains"]
	assert.True(t, exists, "响应应包含 chains 字段")
}

// TestWorker_GracefulShutdownComponents 测试关停流程中各组件的行为。
func TestWorker_GracefulShutdownComponents(t *testing.T) {
	t.Parallel()

	t.Run("closeDB_normal", func(t *testing.T) {
		t.Parallel()
		// 测试正常 db 的关闭流程不 panic
		db := setupWorkerTestDB(t)
		// 不应 panic
		closeDB(db)
		// 验证数据库已关闭：再次 Ping 应失败
		sqlDB, err := db.DB()
		if err == nil && sqlDB != nil {
			assert.Error(t, sqlDB.Ping(), "关闭后 Ping 应失败")
		}
	})

	t.Run("admin_server_shutdown", func(t *testing.T) {
		t.Parallel()
		router := gin.New()
		router.GET("/healthz", healthzHandler)
		srv := &http.Server{
			Addr:    "127.0.0.1:0",
			Handler: router,
		}

		// 启动服务器
		go func() {
			srv.ListenAndServe()
		}()

		// 给服务器一点启动时间
		time.Sleep(50 * time.Millisecond)

		// 优雅关停
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		assert.NoError(t, err, "管理端口关停不应返回错误")
	})
}

// TestWorker_StartupValidation 测试配置加载和验证逻辑。
func TestWorker_StartupValidation(t *testing.T) {
	t.Parallel()

	t.Run("empty_sync_chains", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			JWT:   config.JWTConfig{Secret: "test-secret"},
			Admin: config.AdminConfig{Password: "test-password"},
			Chains: config.ChainsConfig{
				AllowedChainIDs: []int{1},
			},
		}
		err := cfg.ValidateSync()
		assert.Error(t, err, "空 chains 配置应返回错误")
		assert.Contains(t, err.Error(), "sync.chains")
	})

	t.Run("missing_chain_id", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			JWT:   config.JWTConfig{Secret: "test-secret"},
			Admin: config.AdminConfig{Password: "test-password"},
			Chains: config.ChainsConfig{
				AllowedChainIDs: []int{1},
			},
			Sync: config.SyncConfig{
				Chains: []config.SyncChainConfig{
					{
						ChainID:     0, // 缺少 chain_id
						StartBlock:  0,
						ProcessorID: "launchpad",
						Contracts: config.SyncContractsConfig{
							MouseTier: "0x1234",
						},
						Nodes: []config.ChainNodeConfig{
							{HTTPURL: "http://localhost:8545"},
						},
					},
				},
			},
		}
		err := cfg.ValidateSync()
		assert.Error(t, err, "缺少 chain_id 应返回错误")
		assert.Contains(t, err.Error(), "chain_id")
	})

	t.Run("negative_start_block", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			JWT:   config.JWTConfig{Secret: "test-secret"},
			Admin: config.AdminConfig{Password: "test-password"},
			Chains: config.ChainsConfig{
				AllowedChainIDs: []int{1},
			},
			Sync: config.SyncConfig{
				Chains: []config.SyncChainConfig{
					{
						ChainID:     1,
						StartBlock:  -1,
						ProcessorID: "launchpad",
						Contracts: config.SyncContractsConfig{
							MouseTier: "0x1234",
						},
						Nodes: []config.ChainNodeConfig{
							{HTTPURL: "http://localhost:8545"},
						},
					},
				},
			},
		}
		err := cfg.ValidateSync()
		assert.Error(t, err, "负数 start_block 应返回错误")
		assert.Contains(t, err.Error(), "start_block")
	})

	t.Run("missing_nodes", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			JWT:   config.JWTConfig{Secret: "test-secret"},
			Admin: config.AdminConfig{Password: "test-password"},
			Chains: config.ChainsConfig{
				AllowedChainIDs: []int{1},
			},
			Sync: config.SyncConfig{
				Chains: []config.SyncChainConfig{
					{
						ChainID:     1,
						StartBlock:  0,
						ProcessorID: "launchpad",
						Contracts: config.SyncContractsConfig{
							MouseTier: "0x1234",
						},
						Nodes: []config.ChainNodeConfig{},
					},
				},
			},
		}
		err := cfg.ValidateSync()
		assert.Error(t, err, "缺少 nodes 应返回错误")
		assert.Contains(t, err.Error(), "nodes")
	})

	t.Run("missing_contracts", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			JWT:   config.JWTConfig{Secret: "test-secret"},
			Admin: config.AdminConfig{Password: "test-password"},
			Chains: config.ChainsConfig{
				AllowedChainIDs: []int{1},
			},
			Sync: config.SyncConfig{
				Chains: []config.SyncChainConfig{
					{
						ChainID:     1,
						StartBlock:  0,
						ProcessorID: "launchpad",
						Contracts:   config.SyncContractsConfig{},
						Nodes: []config.ChainNodeConfig{
							{HTTPURL: "http://localhost:8545"},
						},
					},
				},
			},
		}
		err := cfg.ValidateSync()
		assert.Error(t, err, "缺少合约地址应返回错误")
		assert.Contains(t, err.Error(), "contracts")
	})

	t.Run("valid_config", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			JWT:   config.JWTConfig{Secret: "test-secret"},
			Admin: config.AdminConfig{Password: "test-password"},
			Chains: config.ChainsConfig{
				AllowedChainIDs: []int{1},
			},
			Sync: config.SyncConfig{
				Chains: []config.SyncChainConfig{
					{
						ChainID:     1,
						StartBlock:  0,
						ProcessorID: "launchpad",
						Contracts: config.SyncContractsConfig{
							MouseTier: "0x1234",
						},
						Nodes: []config.ChainNodeConfig{
							{HTTPURL: "http://localhost:8545"},
						},
					},
				},
			},
		}
		err := cfg.ValidateSync()
		assert.NoError(t, err, "有效配置不应返回错误")
	})
}

// TestWorker_ShutdownTimeout 测试 context 超时机制。
func TestWorker_ShutdownTimeout(t *testing.T) {
	t.Parallel()

	t.Run("context_deadline_exceeded", func(t *testing.T) {
		t.Parallel()
		// 创建极短超时的 context，模拟关停超时
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		// 等待超时
		<-ctx.Done()

		assert.Equal(t, context.DeadlineExceeded, ctx.Err(),
			"超时后 context 错误应为 DeadlineExceeded")
	})

	t.Run("shutdown_with_timeout", func(t *testing.T) {
		t.Parallel()
		shutdownTimeout := 50 * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// 模拟关停流程：启动 HTTP server 后立即关闭
		router := gin.New()
		router.GET("/healthz", healthzHandler)
		srv := &http.Server{
			Addr:    "127.0.0.1:0",
			Handler: router,
		}

		go func() {
			srv.ListenAndServe()
		}()
		time.Sleep(20 * time.Millisecond)

		// 正常关停应在超时前完成
		err := srv.Shutdown(ctx)
		assert.NoError(t, err, "在超时前关停不应返回错误")
	})
}

// TestReadyzHandler 测试就绪探针检查数据库连通性。
func TestReadyzHandler(t *testing.T) {
	t.Parallel()

	t.Run("db_available", func(t *testing.T) {
		t.Parallel()
		db := setupWorkerTestDB(t)

		r := gin.New()
		r.GET("/readyz", readyzHandler(db))

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "ok", resp["status"])
	})

	t.Run("db_closed", func(t *testing.T) {
		t.Parallel()
		db := setupWorkerTestDB(t)
		// 关闭底层连接模拟数据库不可用
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close()

		r := gin.New()
		r.GET("/readyz", readyzHandler(db))

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var resp map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp["error"], "Ping")
	})
}

// TestSyncReplayHandler_InvalidBlockRange 测试 from_block 大于 to_block 时返回 400。
// 传 nil EventStore 是安全的：区块范围校验在访问 store 之前完成。
func TestSyncReplayHandler_InvalidBlockRange(t *testing.T) {
	t.Parallel()

	var store *sync.EventStore // nil，区块范围校验不访问 store

	r := gin.New()
	r.POST("/admin/sync/events/replay", syncReplayHandler(store))

	body := `{"chain_id":1, "from_block":100, "to_block":50}`
	req := httptest.NewRequest(http.MethodPost, "/admin/sync/events/replay",
		strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "from_block")
}

// TestSyncReplayHandler_MissingParams 测试重播请求缺少必要参数时返回 400。
// 传 nil EventStore 是安全的：参数绑定失败时 handler 在访问 store 之前就返回了。
func TestSyncReplayHandler_MissingParams(t *testing.T) {
	t.Parallel()

	var store *sync.EventStore // nil，参数校验路径不访问 store

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "空 JSON",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "缺少 chain_id",
			body:       `{"from_block":1, "to_block":100}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "缺少 from_block",
			body:       `{"chain_id":1, "to_block":100}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "缺少 to_block",
			body:       `{"chain_id":1, "from_block":1}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "所有字段为零值",
			body:       `{"chain_id":0, "from_block":0, "to_block":0}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := gin.New()
			r.POST("/admin/sync/events/replay", syncReplayHandler(store))

			req := httptest.NewRequest(http.MethodPost, "/admin/sync/events/replay",
				strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var resp map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Contains(t, resp["error"], "参数错误")
		})
	}
}

// TestSyncReplayHandler_ValidBlockRange 测试有效区块范围通过参数校验。
// 由于需要真实 DB 才能完成实际操作，此测试使用 Recovery 中间件捕获 nil store 引发的 panic，
// 验证参数校验通过（不会返回 400）。
func TestSyncReplayHandler_ValidBlockRange(t *testing.T) {
	t.Parallel()

	var store *sync.EventStore // nil，会导致后续 store.ResetBlockRange panic

	r := gin.New()
	r.Use(gin.Recovery()) // 捕获 panic，返回 500
	r.POST("/admin/sync/events/replay", syncReplayHandler(store))

	body := `{"chain_id":1, "from_block":1, "to_block":100}`
	req := httptest.NewRequest(http.MethodPost, "/admin/sync/events/replay",
		strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 参数校验通过后，nil store 会导致 panic，被 Recovery 捕获返回 500。
	// 如果返回 400 说明参数校验拦截了有效参数，是 bug。
	assert.Equal(t, http.StatusInternalServerError, w.Code,
		"有效区块范围应通过参数校验（返回 500 而非 400）")
}

// TestSyncReplayHandler_WithRealDB 测试重播区块范围使用真实数据库。
func TestSyncReplayHandler_WithRealDB(t *testing.T) {
	t.Parallel()

	db := setupWorkerTestDB(t)
	store := sync.NewEventStore(db)

	// 插入 2 个 failed 事件在区块范围 [100, 110] 内
	for i := 0; i < 2; i++ {
		event := sync.ChainEvent{
			ChainID:         1,
			BlockNumber:     int64(100 + i*5),
			TxHash:          fmt.Sprintf("0x%064d", 200+i),
			TxIndex:         0,
			LogIndex:        0,
			ContractAddress: "0x1234567890abcdef1234567890abcdef12345678",
			EventName:       "TestEvent",
			EventData:       `{"key":"value"}`,
			Status:          sync.StatusFailed,
			ProcessorID:     "launchpad",
		}
		errMsg := "测试错误"
		now := time.Now()
		event.ErrorMessage = &errMsg
		event.LastFailedAt = &now
		err := db.Create(&event).Error
		require.NoError(t, err)
	}

	// 插入 1 个 dead_letter 事件在范围外
	outEvent := sync.ChainEvent{
		ChainID:         1,
		BlockNumber:     200,
		TxHash:          "0x" + strings.Repeat("0", 63) + "1",
		TxIndex:         0,
		LogIndex:        0,
		ContractAddress: "0x1234567890abcdef1234567890abcdef12345678",
		EventName:       "TestEvent",
		EventData:       `{"key":"value"}`,
		Status:          sync.StatusDeadLetter,
		ProcessorID:     "launchpad",
	}
	outErrMsg := "范围外错误"
	outNow := time.Now()
	outEvent.ErrorMessage = &outErrMsg
	outEvent.LastFailedAt = &outNow
	err := db.Create(&outEvent).Error
	require.NoError(t, err)

	r := gin.New()
	r.POST("/admin/sync/events/replay", syncReplayHandler(store))

	body := `{"chain_id":1, "from_block":100, "to_block":110}`
	req := httptest.NewRequest(http.MethodPost, "/admin/sync/events/replay",
		strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(2), resp["affected"], "应影响区块范围内的 2 个事件")
}
