-- Launchpad 模块数据库迁移：添加部署状态和 PrepareTx 目标字段
-- 为 launchpad_sales 表添加 status 列，支持 Draft 模式
-- 为 launchpad_prepare_txs 表添加 target_address 和 value 列，支持后端签名广播

-- 1. 为 launchpad_sales 表添加 status 列
ALTER TABLE launchpad_sales ADD COLUMN IF NOT EXISTS status VARCHAR(16) NOT NULL DEFAULT 'deploying';

-- 2. 为 launchpad_prepare_txs 表添加 target_address 列
ALTER TABLE launchpad_prepare_txs ADD COLUMN IF NOT EXISTS target_address TEXT NOT NULL DEFAULT '';

-- 3. 为 launchpad_prepare_txs 表添加 value 列
ALTER TABLE launchpad_prepare_txs ADD COLUMN IF NOT EXISTS value TEXT NOT NULL DEFAULT '0';

-- 添加 COMMENT
COMMENT ON COLUMN launchpad_sales.status IS '销售合约部署状态：deploying（部署中）或 deployed（已部署）';
COMMENT ON COLUMN launchpad_prepare_txs.target_address IS '目标合约地址（to），签名广播时使用';
COMMENT ON COLUMN launchpad_prepare_txs.value IS '原生代币数量（wei），签名广播时使用';
