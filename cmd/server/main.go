// Package main 提供 HTTP 服务启动入口。
package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/mousecake-go/mousecake-go/config"
	_ "github.com/mousecake-go/mousecake-go/docs"
	"github.com/mousecake-go/mousecake-go/internal/chain"
	"github.com/mousecake-go/mousecake-go/internal/launchpad"
	"github.com/mousecake-go/mousecake-go/internal/quote"
	"github.com/mousecake-go/mousecake-go/internal/quote/okx"
	"github.com/mousecake-go/mousecake-go/internal/quote/zerox"
	"github.com/mousecake-go/mousecake-go/internal/shared/auth"
	"github.com/mousecake-go/mousecake-go/internal/shared/database"
	"github.com/mousecake-go/mousecake-go/internal/shared/logger"
	"github.com/mousecake-go/mousecake-go/internal/shared/middleware"
	"github.com/mousecake-go/mousecake-go/internal/user"
)

// @title           mousecake-go API
// @version         1.0
// @description     mousecake-go 区块链后端服务 API 文档，包含用户认证、报价聚合等模块。
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http https
// @securityDefinitions.apikey  BearerAuth
// @in                            header
// @name                          Authorization
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

	logger.Init(logger.LogConfig{
		Level:     cfg.Log.Level,
		Format:    cfg.Log.Format,
		AddSource: cfg.Log.AddSource,
	})

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

	jwtSvc := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.WalletExpire, cfg.JWT.AdminExpire)
	repo := user.NewUserRepository(db)
	svc := user.NewService(repo, jwtSvc, cfg.Chains.AllowedChainIDs, cfg.Admin.Username, cfg.Admin.Password, "mousecake-go")
	handler := user.NewHandler(svc, jwtSvc)

	if err := svc.SeedAdmin(context.Background()); err != nil {
		slog.Error("初始化管理员失败", "error", err)
		os.Exit(1)
	}

	switch cfg.Server.Mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(
		gin.Recovery(),
		middleware.RequestID(),
		middleware.AccessLog(slog.Default()),
		middleware.NewCORS(cfg.CORS),
	)
	handler.RegisterRoutes(r.Group("/api/v1/auth"))

	// Quote 模块初始化
	quoteRegistry := quote.NewProviderRegistry()
	for name, pc := range cfg.Quote.Providers {
		if !pc.Enabled {
			continue
		}
		switch name {
		case "okx":
			quoteRegistry.Register(okx.NewProvider(okx.Config{
				APIKey:     pc.APIKey,
				SecretKey:  pc.SecretKey,
				Passphrase: pc.Passphrase,
				BaseURL:    pc.BaseURL,
			}))
		case "zerox":
			quoteRegistry.Register(zerox.NewProvider(zerox.Config{
				APIKey:  pc.APIKey,
				BaseURL: pc.BaseURL,
			}))
		}
	}

	cacheTTL := time.Duration(cfg.Quote.CacheTTL) * time.Second
	if cacheTTL <= 0 {
		cacheTTL = 10 * time.Second
	}
	quoteCache := quote.NewMemoryCache(cacheTTL)
	quoteRepo := quote.NewSwapRecordRepository(db)
	quoteSvc := quote.NewQuoteService(quoteRegistry, quoteCache, quoteRepo, 1)
	quoteHandler := quote.NewHandler(quoteSvc)
	quoteHandler.RegisterRoutes(r.Group("/api/v1"))

	// Launchpad 模块初始化
	launchpadChain, err := findLaunchpadChain(cfg.Sync.Chains)
	if err != nil {
		slog.Error("查找 launchpad 链配置失败", "error", err)
		os.Exit(1)
	}

	nodePool, err := chain.NewNodePool(launchpadChain.Nodes, launchpadChain.ChainID)
	if err != nil {
		slog.Error("创建 NodePool 失败", "error", err)
		os.Exit(1)
	}
	defer nodePool.Close()

	chainReader := launchpad.NewChainReader(nodePool)
	encoder := launchpad.NewContractCodeEncoder()

	// 创建所有 Repository
	saleRepo := launchpad.NewSaleRepository(db)
	saleMetaRepo := launchpad.NewSaleMetaRepository(db)
	poolRepo := launchpad.NewPoolRepository(db)
	tierLimitRepo := launchpad.NewTierLimitRepository(db)
	whitelistRepo := launchpad.NewWhitelistRepository(db)
	depositRepo := launchpad.NewDepositRepository(db)
	userPoolRepo := launchpad.NewUserPoolStateRepository(db)
	creditRepo := launchpad.NewUserCreditRepository(db)
	harvestRepo := launchpad.NewHarvestRepository(db)
	vestingRepo := launchpad.NewVestingScheduleRepository(db)
	releaseRepo := launchpad.NewVestingReleaseRepository(db)
	tokenRepo := launchpad.NewTokenRepository(db)
	prepareRepo := launchpad.NewPrepareTxRepository(db)

	// 创建 Prepare Service
	prepareExpiresIn := cfg.Launchpad.PrepareExpiresIn
	if prepareExpiresIn == 0 {
		prepareExpiresIn = 30 * time.Minute
	}
	prepareSvc := launchpad.NewPrepareService(prepareRepo, chainReader, encoder, prepareExpiresIn)

	// 创建其他 Service
	adminSvc := launchpad.NewAdminService(prepareSvc, encoder, "MousePadByTierDeployer", "MousePadByTier", saleRepo, cfg.Launchpad.DeployerAddress)
	querySvc := launchpad.NewUserQueryService(
		saleRepo, poolRepo, saleMetaRepo, tierLimitRepo, whitelistRepo,
		depositRepo, userPoolRepo, creditRepo, harvestRepo, vestingRepo, releaseRepo,
		chainReader,
	)
	userSvc := launchpad.NewUserService(
		prepareSvc, querySvc, encoder,
		launchpadChain.Contracts.MousePadByTier,
		chainReader, saleRepo, poolRepo, tierLimitRepo, whitelistRepo,
		userPoolRepo, creditRepo, vestingRepo,
	)
	metaSvc := launchpad.NewSaleMetaService(saleMetaRepo, saleRepo)
	tokenSvc := launchpad.NewTokenService(tokenRepo)

	// 创建 ChainRefreshService
	chainRefreshRepo := launchpad.NewChainRefreshRepository(db)
	chainRefreshSvc := launchpad.NewChainRefreshService(chainReader, chainRefreshRepo, saleRepo)

	// 创建 Signer（非 release 模式且有私钥时）
	var signer *chain.Signer
	if gin.Mode() != gin.ReleaseMode && cfg.Launchpad.AdminPrivateKey != "" {
		signerChainID := big.NewInt(int64(launchpadChain.ChainID))
		signer, err = chain.NewSigner(cfg.Launchpad.AdminPrivateKey, nodePool, signerChainID)
		if err != nil {
			slog.Error("创建 Signer 失败", "error", err)
			os.Exit(1)
		}
		slog.Info("Signer 初始化完成", "address", signer.Address())
	}

	// 创建 DevExecuteService
	var devExecuteSvc launchpad.DevExecutor
	if signer != nil {
		devExecuteSvc = launchpad.NewDevExecuteService(prepareRepo, signer)
	}

	// 创建通用合约查询服务（仅非 release 模式）
	if gin.Mode() != gin.ReleaseMode {
		contractRegistry := chain.NewABIRegistry()
		addresses := map[string]string{
			"MousePadByTier":         cfg.Launchpad.MousePadByTierAddress,
			"MousePadByTierDeployer": cfg.Launchpad.DeployerAddress,
			"MouseTier":              cfg.Launchpad.MouseTierAddress,
		}
		chain.RegisterAllContracts(contractRegistry, addresses)
		eventParser := chain.NewEventParser()
		devContractSvc := chain.NewDevContractService(contractRegistry, nodePool, signer, eventParser)
		devContractHandler := chain.NewDevContractHandler(devContractSvc)
		devContractHandler.RegisterRoutes(r.Group("/api/v1"))
		slog.Info("通用合约查询服务已注册", "contracts", contractRegistry.ListContracts())
	}

	// 创建 Handler 并注册路由
	launchpadHandler := launchpad.NewHandler(adminSvc, userSvc, prepareSvc, querySvc, metaSvc, tokenSvc, chainRefreshSvc, devExecuteSvc)
	launchpadHandler.RegisterRoutes(r.Group("/api/v1/launchpad"))

	// 注册兜底轮询定时任务
	pollInterval := cfg.Launchpad.PollInterval
	if pollInterval == 0 {
		pollInterval = 1 * time.Minute
	}
	//go runPreparePollTask(context.Background(), prepareSvc, pollInterval)

	slog.Info("Launchpad 模块初始化完成")

	// Swagger UI 仅在非 release 模式下注册
	if gin.Mode() != gin.ReleaseMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("启动 HTTP 服务", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP 服务异常退出", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("收到终止信号，开始优雅关停")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("服务关停超时", "error", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("获取数据库连接失败", "error", err)
	} else if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			slog.Error("关闭数据库连接失败", "error", err)
		}
	}

	slog.Info("服务已优雅关停")
}

// findLaunchpadChain 查找 launchpad 对应的链配置。
func findLaunchpadChain(chains []config.SyncChainConfig) (config.SyncChainConfig, error) {
	for _, c := range chains {
		if c.ProcessorID == "launchpad" {
			return c, nil
		}
	}
	return config.SyncChainConfig{}, fmt.Errorf("未找到 processor_id 为 launchpad 的链配置")
}

// runPreparePollTask 兜底轮询定时任务，定期扫描超时的 Prepare 交易。
func runPreparePollTask(ctx context.Context, svc *launchpad.PrepareService, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	slog.Info("启动 Prepare 兜底轮询任务", "interval", interval)

	for {
		select {
		case <-ctx.Done():
			slog.Info("停止 Prepare 兜底轮询任务")
			return
		case <-ticker.C:
			if err := svc.PollTimeout(ctx); err != nil {
				slog.Error("Prepare 轮询失败", "error", err)
			}
		}
	}
}
