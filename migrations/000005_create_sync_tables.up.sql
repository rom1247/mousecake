-- 同步框架数据库迁移：创建 chain_events 和 sync_checkpoints 表
-- chain_events 存储链上原始事件，Projector 异步消费并写入业务投影表
-- sync_checkpoints 记录每条链的同步进度

-- 1. 链上事件存储表
CREATE TABLE IF NOT EXISTS chain_events (
    id BIGSERIAL PRIMARY KEY,
    chain_id INT NOT NULL,
    block_number BIGINT NOT NULL,
    tx_hash CHAR(66) NOT NULL,
    tx_index INT NOT NULL,
    log_index INT NOT NULL,
    contract_address CHAR(42) NOT NULL,
    event_name VARCHAR(64) NOT NULL,
    event_data JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    retry_count INT NOT NULL DEFAULT 0,
    error_message TEXT,
    processor_id VARCHAR(64) NOT NULL DEFAULT '',
    last_failed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 链上事件唯一约束：同一链同一区块同一交易同一日志只存一条
CREATE UNIQUE INDEX IF NOT EXISTS idx_chain_events_unique ON chain_events (chain_id, block_number, tx_index, log_index);

-- Projector 消费用：按状态查询待处理事件
CREATE INDEX IF NOT EXISTS idx_chain_events_status ON chain_events (status, id) WHERE status IN ('pending', 'failed');

-- 处理器维度查询
CREATE INDEX IF NOT EXISTS idx_chain_events_processor ON chain_events (processor_id, status);

-- 死信队列查询
CREATE INDEX IF NOT EXISTS idx_chain_events_dead_letter ON chain_events (chain_id, status) WHERE status = 'dead_letter';

-- processing 超时扫描
CREATE INDEX IF NOT EXISTS idx_chain_events_processing ON chain_events (status, updated_at) WHERE status = 'processing';

COMMENT ON TABLE chain_events IS '链上事件存储：所有链上事件统一写入此表，Projector 异步消费并写入业务投影表';
COMMENT ON COLUMN chain_events.chain_id IS '链 ID（1=ETH 主网, 5=Goerli, 11155111=Sepolia）';
COMMENT ON COLUMN chain_events.block_number IS '区块号';
COMMENT ON COLUMN chain_events.tx_hash IS '交易哈希';
COMMENT ON COLUMN chain_events.tx_index IS '交易在区块中的索引';
COMMENT ON COLUMN chain_events.log_index IS '日志在交易中的索引';
COMMENT ON COLUMN chain_events.contract_address IS '发出事件的合约地址';
COMMENT ON COLUMN chain_events.event_name IS '事件名称（如 Deposited, SaleCreated）';
COMMENT ON COLUMN chain_events.event_data IS '事件参数数据（JSON 格式）';
COMMENT ON COLUMN chain_events.status IS '事件状态：pending/processing/processed/failed/dead_letter';
COMMENT ON COLUMN chain_events.retry_count IS '已重试次数';
COMMENT ON COLUMN chain_events.error_message IS '最后失败原因';
COMMENT ON COLUMN chain_events.processor_id IS '处理器标识（如 launchpad）';
COMMENT ON COLUMN chain_events.last_failed_at IS '最后失败时间';

-- 2. 同步进度检查点表
CREATE TABLE IF NOT EXISTS sync_checkpoints (
    id BIGSERIAL PRIMARY KEY,
    chain_id INT NOT NULL,
    processor_id VARCHAR(64) NOT NULL,
    last_synced_block BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 每条链每个处理器只有一个 checkpoint
CREATE UNIQUE INDEX IF NOT EXISTS idx_sync_checkpoints_unique ON sync_checkpoints (chain_id, processor_id);

COMMENT ON TABLE sync_checkpoints IS '同步进度检查点：记录每条链每个处理器的已同步区块号';
COMMENT ON COLUMN sync_checkpoints.chain_id IS '链 ID';
COMMENT ON COLUMN sync_checkpoints.processor_id IS '处理器标识（如 launchpad）';
COMMENT ON COLUMN sync_checkpoints.last_synced_block IS '已同步到的最大区块号';
