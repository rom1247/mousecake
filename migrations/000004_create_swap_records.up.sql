-- Quote 模块数据库迁移：创建 swap_records 表

CREATE TABLE IF NOT EXISTS swap_records (
    id BIGINT PRIMARY KEY,
    provider VARCHAR(32) NOT NULL,
    chain_id INT NOT NULL,
    from_token CHAR(42) NOT NULL,
    to_token CHAR(42) NOT NULL,
    from_amount VARCHAR(78) NOT NULL,
    to_amount VARCHAR(78) NOT NULL,
    slippage_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
    swap_mode VARCHAR(16) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    tx_hash CHAR(66),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_swap_records_provider ON swap_records (provider);
CREATE INDEX IF NOT EXISTS idx_swap_records_chain_id ON swap_records (chain_id);
CREATE INDEX IF NOT EXISTS idx_swap_records_status ON swap_records (status);

COMMENT ON TABLE swap_records IS '代币兑换记录：存储报价请求和前端回传的交易哈希';
COMMENT ON COLUMN swap_records.id IS 'Snowflake 风格主键（17-19 位数字）';
COMMENT ON COLUMN swap_records.provider IS '供应商名称（如 okx、zerox）';
COMMENT ON COLUMN swap_records.chain_id IS '链 ID';
COMMENT ON COLUMN swap_records.from_token IS '源代币合约地址';
COMMENT ON COLUMN swap_records.to_token IS '目标代币合约地址';
COMMENT ON COLUMN swap_records.from_amount IS '源代币数量（wei）';
COMMENT ON COLUMN swap_records.to_amount IS '目标代币数量（wei）';
COMMENT ON COLUMN swap_records.slippage_percent IS '滑点百分比';
COMMENT ON COLUMN swap_records.swap_mode IS '兑换模式：exactIn / exactOut';
COMMENT ON COLUMN swap_records.status IS '状态：pending / submitted';
COMMENT ON COLUMN swap_records.tx_hash IS '链上交易哈希（前端回传）';
