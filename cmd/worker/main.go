// Package main 提供链上同步 Worker 进程入口。
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mousecake-go/mousecake-go/internal/shared/sync"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"

	"github.com/mousecake-go/mousecake-go/config"
	"github.com/mousecake-go/mousecake-go/internal/chain"
	"github.com/mousecake-go/mousecake-go/internal/launchpad"
	"github.com/mousecake-go/mousecake-go/internal/shared/database"
	"github.com/mousecake-go/mousecake-go/internal/shared/logger"
)

func main() {
	cfg, err := config.Load("config/app.yaml")
	if err != nil {
		slog.Error("加载配置失败", "error", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		slog.Error("配置验证失败", "error", err)
		os.Exit(1)
	}

	if err := cfg.ValidateSync(); err != nil {
		slog.Error("同步配置验证失败", "error", err)
		os.Exit(1)
	}

	logger.Init(logger.LogConfig{
		Level:     cfg.Log.Level,
		Format:    cfg.Log.Format,
		AddSource: cfg.Log.AddSource,
	})

	switch cfg.Server.Mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		slog.Error("连接数据库失败", "error", err)
		os.Exit(1)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.DBName, cfg.Database.SSLMode,
	)
	if err := database.RunMigrations(dsn, "migrations"); err != nil {
		slog.Error("执行 migration 失败", "error", err)
		os.Exit(1)
	}

	prepareTxRepo := launchpad.NewPrepareTxRepository(db)
	eventSvc := launchpad.NewEventService(db, prepareTxRepo)

	syncMgr, err := sync.NewSyncManager(cfg.Sync, db, eventSvc)
	if err != nil {
		slog.Error("创建 SyncManager 失败", "error", err)
		os.Exit(1)
	}

	// 管理端口路由
	store := sync.NewEventStore(db)

	adminRouter := gin.New()
	adminRouter.Use(gin.Recovery())
	adminRouter.GET("/healthz", healthzHandler)
	adminRouter.GET("/readyz", readyzHandler(db))
	adminRouter.GET("/metrics", metricsHandler)

	syncGroup := adminRouter.Group("/admin/sync")
	{
		syncGroup.GET("/status", syncStatusHandler(syncMgr))
		syncGroup.GET("/events", syncEventsHandler(store))
		syncGroup.POST("/events/retry", syncRetryHandler(store))
		syncGroup.POST("/events/replay", syncReplayHandler(store))
	}

	adminSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Sync.Worker.HealthPort),
		Handler:      adminRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("管理端口启动", "addr", adminSrv.Addr)
		if err := adminSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("管理端口异常退出", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	slog.Info("Worker 进程启动", "chains", len(cfg.Sync.Chains))
	if err := syncMgr.Start(ctx); err != nil {
		slog.Error("SyncManager 启动失败", "error", err)
		os.Exit(1)
	}

	<-ctx.Done()
	slog.Info("收到终止信号，开始优雅关停")

	shutdownTimeout := cfg.Sync.Worker.ShutdownTimeout
	if shutdownTimeout == 0 {
		shutdownTimeout = 15 * time.Second
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// 关停超时后强制退出
	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			slog.Warn("关停超时，强制退出")
			os.Exit(1)
		}
	}()

	syncMgr.Stop()

	if err := adminSrv.Shutdown(shutdownCtx); err != nil {
		slog.Error("管理端口关停超时", "error", err)
	}

	closeDB(db)
	slog.Info("Worker 进程已优雅关停")
}

// closeDB 关闭数据库连接。
func closeDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil && sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			slog.Error("关闭数据库连接失败", "error", err)
		}
	}
}

// healthzHandler 存活探针。
func healthzHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// readyzHandler 就绪探针，检查数据库连通性。
func readyzHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sqlDB *sql.DB
		var err error
		if sqlDB, err = db.DB(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "error": "数据库连接失败"})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "error": "数据库 Ping 失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// metricsHandler Prometheus 指标暴露端点。
func metricsHandler(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

// syncStatusHandler 同步状态查询 API。
func syncStatusHandler(mgr *sync.SyncManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		statuses, err := mgr.GetStatus(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"chains": statuses})
	}
}

// syncEventsHandler 查看事件列表 API。
func syncEventsHandler(store *sync.EventStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Query("status")
		chainID := 0
		if v := c.Query("chain_id"); v != "" {
			id, err := strconv.Atoi(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "chain_id 参数无效"})
				return
			}
			chainID = id
		}
		page := 1
		if v := c.Query("page"); v != "" {
			p, err := strconv.Atoi(v)
			if err != nil || p < 1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "page 参数无效"})
				return
			}
			page = p
		}
		pageSize := 20
		if v := c.Query("page_size"); v != "" {
			ps, err := strconv.Atoi(v)
			if err != nil || ps < 1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "page_size 参数无效"})
				return
			}
			pageSize = ps
		}

		if status == "dead_letter" {
			events, total, err := store.ListFailed(c.Request.Context(), chainID, sync.StatusDeadLetter, page, pageSize)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"events": events, "total": total, "status": status})
			return
		}

		events, err := store.ListPending(c.Request.Context(), chainID, "", pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"events": events})
	}
}

// syncRetryHandler 重试单个失败事件 API。
func syncRetryHandler(store *sync.EventStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID int64 `json:"id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 id 参数"})
			return
		}

		event, err := store.GetByID(c.Request.Context(), req.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "事件不存在"})
			return
		}

		if err := store.ResetToPending(c.Request.Context(), event.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "事件已重置为 pending", "id": event.ID})
	}
}

// syncReplayHandler 重播区块范围 API。
func syncReplayHandler(store *sync.EventStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ChainID   int   `json:"chain_id" binding:"required"`
			FromBlock int64 `json:"from_block" binding:"required"`
			ToBlock   int64 `json:"to_block" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
			return
		}

		if req.FromBlock > req.ToBlock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "from_block 不能大于 to_block"})
			return
		}

		affected, err := store.ResetBlockRange(c.Request.Context(), req.ChainID, req.FromBlock, req.ToBlock)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":    "区块范围事件已重置",
			"chain_id":   req.ChainID,
			"from_block": req.FromBlock,
			"to_block":   req.ToBlock,
			"affected":   affected,
		})
	}
}

// 确保 chain.NodePool 被引用（Worker 间接通过 SyncManager 使用）
var _ chain.NodePool

// 确保 launchpad.EventService 满足 syncpkg.EventService 接口
var _ sync.EventService = (*launchpad.EventService)(nil)
