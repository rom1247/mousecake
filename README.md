# mousecake-go

区块链后端服务，提供链上链下状态索引同步、钱包登录（SIWE）、IDO 生命周期管理、质押池、代币管理等功能。

## 技术栈

- **语言**: Go 1.22+
- **HTTP 框架**: Gin
- **ORM**: Gorm（PostgreSQL）
- **数据库迁移**: golang-migrate
- **链交互**: ethereum/go-ethereum（ethclient、事件订阅、ABI 绑定）
- **认证**: SIWE（Sign-In with Ethereum）+ JWT
- **API 文档**: Swagger（swaggo/gin-swagger）
- **配置**: Viper
- **日志**: slog（标准库）
- **可观测性**: Prometheus + OpenTelemetry
- **架构**: Package-by-Feature + DDD（领域驱动设计）

## 架构概览

### 双入口模式

区块链后端有两个入口，通过 `service.go` 编排业务逻辑：

1. **HTTP Handler**：处理前端/管理后台 API 请求
2. **Event Worker**：订阅链上合约事件，处理链上状态变更

### 依赖方向

```
handler.go / worker.go（入口层）
       ↓
service.go（用例编排）
       ↓           ↓           ↓
domain/（领域模型）  repository  其他模块 service
       ↑           ↑
chain.go       数据库
```

## 项目结构

```text
cmd/
    cli/                # CLI 工具入口
    schema-gen/         # 数据库 Schema 生成工具
    server/             # HTTP 服务器入口
    worker/             # 链上事件消费者入口
config/                 # 配置文件
docs/                   # 项目文档
internal/
    chain/              # 链交互基础设施（跨模块共享）
        contract/       # 合约 ABI 与绑定
        node.go         # ethclient 封装 + 连接管理
        node_pool.go    # 节点池管理
        subscriber.go   # 统一链上事件订阅
    launchpad/          # Launchpad IDO 模块
        domain/         # 领域模型（实体、值对象、状态机）
        service.go      # 用例编排
        repository.go   # 数据访问
        handler.go      # HTTP handler
        worker.go       # 链上事件消费者
    quote/              # 代币报价聚合模块（OKX、0x）
        domain/
        service.go
        handler.go
    user/               # 用户认证与钱包登录
        domain/
        service.go
        repository.go
        handler.go
    shared/             # 跨模块共享基础设施
        auth/           # JWT 服务
        database/       # 数据库连接
        errs/           # 通用错误
        logger/         # slog 配置
        middleware/     # Gin 中间件
        response/       # 统一响应
        sync/           # 同步原语
migrations/             # 数据库迁移文件

```

## 快速开始

### 前置条件

- Go 1.22+
- PostgreSQL
-（可选）区块链 RPC 节点访问权限

### 配置

1. 复制环境变量模板：

```bash
cp .env.example .env
```

2. 根据环境修改 `.env` 文件中的数据库连接、RPC 端点、JWT Secret 等配置。

### 数据库迁移

```bash
# 应用所有迁移
migrate -path migrations -database "postgres://user:pass@localhost/dbname?sslmode=disable" up

# 回滚最近一次迁移
migrate -path migrations -database "postgres://user:pass@localhost/dbname?sslmode=disable" down 1
```

### 运行

```bash
# 构建
go build -o bin/server ./cmd/server

# 运行 HTTP 服务器
./bin/server

# 运行事件消费者
go run ./cmd/worker
```

或直接运行：

```bash
go run ./cmd/server
```

服务默认监听 `:8080`，Swagger 文档访问 `http://localhost:8080/swagger/index.html`。

### 测试

```bash
# 运行所有测试
go test ./...

# 带竞态检测的测试
go test -race ./...
```

## API 概览

| 模块 | 路径前缀 | 说明 |
|------|---------|------|
| 用户认证 | `/api/v1/auth` | SIWE 钱包登录、Nonce 获取、JWT 签发 |
| 报价聚合 | `/api/v1` | 多源代币报价（OKX、0x）查询与比较 |
| Launchpad | - | IDO 项目创建、认购、结算 |

## 开发

### 构建 & 验证

```bash
go build ./...
go vet ./...
go test -race ./...
```

### 生成 Swagger 文档

```bash
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

生成的文档输出到 `docs/` 目录（`docs.go`、`swagger.json`、`swagger.yaml`）。开发模式下访问 `http://localhost:8080/swagger/index.html` 查看 Swagger UI。

### Lint

```bash
golangci-lint run
```

### 安全检查

```bash
# Go 漏洞检测
govulncheck ./...

# 静态安全分析
gosec ./...
```

## License

[License](LICENSE)
