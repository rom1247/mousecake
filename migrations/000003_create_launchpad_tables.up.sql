-- Launchpad 模块数据库迁移：创建 14 张表
-- 包含 IDO 销售、池子、申购、结算、vesting、代币元信息和 Prepare 交易记录

-- 1. 代币元信息表
CREATE TABLE IF NOT EXISTS launchpad_tokens (
    id BIGSERIAL PRIMARY KEY,
    address CHAR(42) NOT NULL,
    chain_id INT NOT NULL,
    name VARCHAR(64) NOT NULL,
    symbol VARCHAR(32) NOT NULL,
    decimals INT NOT NULL DEFAULT 18,
    logo_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_tokens_address_chain ON launchpad_tokens (address, chain_id);

COMMENT ON TABLE launchpad_tokens IS '代币元信息：存储链上代币的展示信息';
COMMENT ON COLUMN launchpad_tokens.address IS '代币合约地址';
COMMENT ON COLUMN launchpad_tokens.chain_id IS '链 ID';
COMMENT ON COLUMN launchpad_tokens.name IS '代币名称';
COMMENT ON COLUMN launchpad_tokens.symbol IS '代币符号';
COMMENT ON COLUMN launchpad_tokens.decimals IS '代币精度';
COMMENT ON COLUMN launchpad_tokens.logo_url IS '代币 Logo URL';

-- 2. MouseTier 参数表
CREATE TABLE IF NOT EXISTS launchpad_tier_params (
    id BIGSERIAL PRIMARY KEY,
    chain_id INT NOT NULL,
    ceiling NUMERIC(78,0) NOT NULL DEFAULT 0,
    multiplier NUMERIC(78,0) NOT NULL DEFAULT 0,
    tier_base_amount NUMERIC(78,0) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE launchpad_tier_params IS 'MouseTier 合约参数：存储等级计算的配置参数';
COMMENT ON COLUMN launchpad_tier_params.ceiling IS '满乘数阈值（锁仓时长达到此值可获得完整乘数）';
COMMENT ON COLUMN launchpad_tier_params.multiplier IS '积分乘数参数';
COMMENT ON COLUMN launchpad_tier_params.tier_base_amount IS '等级基准门槛值';

-- 3. IDO 销售合约表（链上状态）
CREATE TABLE IF NOT EXISTS launchpad_sales (
    id BIGSERIAL PRIMARY KEY,
    contract_address CHAR(42) NOT NULL,
    chain_id INT NOT NULL,
    deployer_address CHAR(42) NOT NULL,
    owner_address CHAR(42) NOT NULL,
    raise_token_address CHAR(42) NOT NULL,
    offering_token_address CHAR(42) NOT NULL,
    mouse_tier_address CHAR(42) NOT NULL,
    start_block BIGINT NOT NULL DEFAULT 0,
    end_block BIGINT NOT NULL DEFAULT 0,
    vesting_start_time BIGINT NOT NULL DEFAULT 0,
    vesting_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    max_buffer_blocks BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_sales_contract ON launchpad_sales (contract_address) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_launchpad_sales_chain ON launchpad_sales (chain_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE launchpad_sales IS 'IDO 销售合约：记录每个 MousePadByTier 实例的链上状态';
COMMENT ON COLUMN launchpad_sales.contract_address IS 'MousePadByTier 合约地址';
COMMENT ON COLUMN launchpad_sales.deployer_address IS 'Deployer 工厂合约地址';
COMMENT ON COLUMN launchpad_sales.owner_address IS 'IDO 管理员地址';
COMMENT ON COLUMN launchpad_sales.raise_token_address IS '募资币合约地址';
COMMENT ON COLUMN launchpad_sales.offering_token_address IS '发售币合约地址';
COMMENT ON COLUMN launchpad_sales.mouse_tier_address IS 'MouseTier 合约地址';
COMMENT ON COLUMN launchpad_sales.start_block IS '销售开始区块';
COMMENT ON COLUMN launchpad_sales.end_block IS '销售结束区块';
COMMENT ON COLUMN launchpad_sales.vesting_start_time IS '全局 vesting 起始时间戳（首次 harvest 时设置）';
COMMENT ON COLUMN launchpad_sales.vesting_revoked IS '管理员是否已撤销全部 vesting';
COMMENT ON COLUMN launchpad_sales.max_buffer_blocks IS '最大缓冲区块范围';

-- 4. 销售元信息表（后台管理信息）
CREATE TABLE IF NOT EXISTS launchpad_sale_meta (
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL REFERENCES launchpad_sales(id),
    title VARCHAR(256),
    description TEXT,
    banner_url TEXT,
    logo_url TEXT,
    website_url TEXT,
    social_links JSONB,
    visibility VARCHAR(16) NOT NULL DEFAULT 'hidden',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_sale_meta_sale ON launchpad_sale_meta (sale_id);

COMMENT ON TABLE launchpad_sale_meta IS '销售元信息：管理后台维护的展示信息，与链上状态分表存储';
COMMENT ON COLUMN launchpad_sale_meta.title IS '销售标题';
COMMENT ON COLUMN launchpad_sale_meta.description IS '项目描述';
COMMENT ON COLUMN launchpad_sale_meta.banner_url IS 'Banner 图片 URL';
COMMENT ON COLUMN launchpad_sale_meta.logo_url IS '项目 Logo URL';
COMMENT ON COLUMN launchpad_sale_meta.website_url IS '项目官网 URL';
COMMENT ON COLUMN launchpad_sale_meta.social_links IS '社交链接（JSON 格式）';
COMMENT ON COLUMN launchpad_sale_meta.visibility IS '可见性：public=公开, hidden=隐藏';
COMMENT ON COLUMN launchpad_sale_meta.sort_order IS '排序权重（数值越大越靠前）';

-- 5. 池子配置表
CREATE TABLE IF NOT EXISTS launchpad_pools (
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL REFERENCES launchpad_sales(id),
    pool_index INT NOT NULL,
    raising_amount NUMERIC(78,0) NOT NULL DEFAULT 0,
    offering_amount NUMERIC(78,0) NOT NULL DEFAULT 0,
    limit_per_user NUMERIC(78,0) NOT NULL DEFAULT 0,
    is_special_sale BOOLEAN NOT NULL DEFAULT FALSE,
    has_tax BOOLEAN NOT NULL DEFAULT FALSE,
    tax_rate NUMERIC(78,0) NOT NULL DEFAULT 0,
    vesting_percentage INT NOT NULL DEFAULT 0,
    vesting_cliff BIGINT NOT NULL DEFAULT 0,
    vesting_duration BIGINT NOT NULL DEFAULT 0,
    vesting_slice_period BIGINT NOT NULL DEFAULT 0,
    total_amount NUMERIC(78,0) NOT NULL DEFAULT 0,
    total_tax NUMERIC(78,0) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_pools_sale_pool ON launchpad_pools (sale_id, pool_index);

COMMENT ON TABLE launchpad_pools IS '销售池子：每个 Sale 固定 2 个池，存储池子配置和累计统计';
COMMENT ON COLUMN launchpad_pools.pool_index IS '池子索引（0=普通池, 1=特殊池）';
COMMENT ON COLUMN launchpad_pools.raising_amount IS '募资目标金额（uint256）';
COMMENT ON COLUMN launchpad_pools.offering_amount IS '发售代币总量（uint256）';
COMMENT ON COLUMN launchpad_pools.limit_per_user IS '单用户最大申购金额（0 表示不限）';
COMMENT ON COLUMN launchpad_pools.is_special_sale IS '是否为白名单池（true=特殊池, false=普通 Tier 池）';
COMMENT ON COLUMN launchpad_pools.has_tax IS '超募时是否对退款部分征税';
COMMENT ON COLUMN launchpad_pools.tax_rate IS '超募税率';
COMMENT ON COLUMN launchpad_pools.vesting_percentage IS '锁仓比例（0-100，TGE 释放 100-vesting_percentage）';
COMMENT ON COLUMN launchpad_pools.vesting_cliff IS 'Vesting cliff 时长（秒）';
COMMENT ON COLUMN launchpad_pools.vesting_duration IS 'Vesting 总时长（秒）';
COMMENT ON COLUMN launchpad_pools.vesting_slice_period IS 'Vesting 切片周期（秒）';
COMMENT ON COLUMN launchpad_pools.total_amount IS '反规范化：池子累计申购总额';
COMMENT ON COLUMN launchpad_pools.total_tax IS '反规范化：池子累计超募税总额';

-- 6. Tier 额度表
CREATE TABLE IF NOT EXISTS launchpad_tier_limits (
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL REFERENCES launchpad_sales(id),
    tier INT NOT NULL,
    credit_limit NUMERIC(78,0) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_tier_limits_sale_tier ON launchpad_tier_limits (sale_id, tier);

COMMENT ON TABLE launchpad_tier_limits IS 'Tier 额度：每个 Sale 各档 Tier（0-5）的最大申购额度';
COMMENT ON COLUMN launchpad_tier_limits.tier IS 'Tier 档位（0-5）';
COMMENT ON COLUMN launchpad_tier_limits.credit_limit IS '该档 Tier 允许的最大申购额度';

-- 7. 白名单表
CREATE TABLE IF NOT EXISTS launchpad_whitelists (
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL REFERENCES launchpad_sales(id),
    address CHAR(42) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    block_number BIGINT NOT NULL DEFAULT 0,
    tx_index INT NOT NULL DEFAULT 0,
    log_index INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_whitelists_event ON launchpad_whitelists (block_number, tx_index, log_index);
CREATE INDEX IF NOT EXISTS idx_launchpad_whitelists_sale_addr ON launchpad_whitelists (sale_id, address, is_active);

COMMENT ON TABLE launchpad_whitelists IS '白名单：特殊池的准入地址列表，由同步程序写入';
COMMENT ON COLUMN launchpad_whitelists.address IS '白名单地址';
COMMENT ON COLUMN launchpad_whitelists.is_active IS '是否在白名单中（addWhitelist=true, removeWhitelist=false）';

-- 8. 申购记录表（事件明细）
CREATE TABLE IF NOT EXISTS launchpad_deposits (
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL REFERENCES launchpad_sales(id),
    pool_index INT NOT NULL,
    user_address CHAR(42) NOT NULL,
    amount NUMERIC(78,0) NOT NULL,
    tx_hash CHAR(66) NOT NULL,
    block_number BIGINT NOT NULL,
    tx_index INT NOT NULL,
    log_index INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_deposits_event ON launchpad_deposits (block_number, tx_index, log_index);
CREATE INDEX IF NOT EXISTS idx_launchpad_deposits_user_sale ON launchpad_deposits (sale_id, pool_index, user_address);

COMMENT ON TABLE launchpad_deposits IS '申购记录：每笔链上 deposit 事件明细，由同步程序写入';
COMMENT ON COLUMN launchpad_deposits.amount IS '申购金额（uint256）';
COMMENT ON COLUMN launchpad_deposits.tx_hash IS '链上交易哈希';

-- 9. 用户池内累计状态表
CREATE TABLE IF NOT EXISTS launchpad_user_pool_state (
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL REFERENCES launchpad_sales(id),
    pool_index INT NOT NULL,
    user_address CHAR(42) NOT NULL,
    total_deposited NUMERIC(78,0) NOT NULL DEFAULT 0,
    claimed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_user_pool_state_unique ON launchpad_user_pool_state (sale_id, pool_index, user_address);

COMMENT ON TABLE launchpad_user_pool_state IS '用户池内累计状态：按用户+池维度的汇总数据，查询用';
COMMENT ON COLUMN launchpad_user_pool_state.total_deposited IS '用户在该池的累计申购金额';
COMMENT ON COLUMN launchpad_user_pool_state.claimed IS '用户在该池是否已结算（harvest）';

-- 10. 用户信用使用表
CREATE TABLE IF NOT EXISTS launchpad_user_credit (
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL REFERENCES launchpad_sales(id),
    user_address CHAR(42) NOT NULL,
    credit_used NUMERIC(78,0) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_user_credit_unique ON launchpad_user_credit (sale_id, user_address);

COMMENT ON TABLE launchpad_user_credit IS '用户信用使用：记录用户在某 sale 的所有普通池累计信用消耗';
COMMENT ON COLUMN launchpad_user_credit.credit_used IS '用户在该 sale 所有普通池的累计信用使用量';

-- 11. 结算记录表
CREATE TABLE IF NOT EXISTS launchpad_harvests (
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL REFERENCES launchpad_sales(id),
    pool_index INT NOT NULL,
    user_address CHAR(42) NOT NULL,
    is_overflow BOOLEAN NOT NULL DEFAULT FALSE,
    offering_amount NUMERIC(78,0) NOT NULL DEFAULT 0,
    pay_amount NUMERIC(78,0) NOT NULL DEFAULT 0,
    raise_refund NUMERIC(78,0) NOT NULL DEFAULT 0,
    tax_amount NUMERIC(78,0) NOT NULL DEFAULT 0,
    tge_amount NUMERIC(78,0) NOT NULL DEFAULT 0,
    vesting_amount NUMERIC(78,0) NOT NULL DEFAULT 0,
    tx_hash CHAR(66) NOT NULL,
    block_number BIGINT NOT NULL,
    tx_index INT NOT NULL,
    log_index INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_harvests_event ON launchpad_harvests (block_number, tx_index, log_index);
CREATE INDEX IF NOT EXISTS idx_launchpad_harvests_user_sale ON launchpad_harvests (sale_id, user_address);

COMMENT ON TABLE launchpad_harvests IS '结算记录：用户 harvest 操作的结果，由同步程序写入';
COMMENT ON COLUMN launchpad_harvests.is_overflow IS '是否超募配售';
COMMENT ON COLUMN launchpad_harvests.offering_amount IS '配售发售币数量';
COMMENT ON COLUMN launchpad_harvests.pay_amount IS '有效支付金额';
COMMENT ON COLUMN launchpad_harvests.raise_refund IS '退款金额';
COMMENT ON COLUMN launchpad_harvests.tax_amount IS '超募税额';
COMMENT ON COLUMN launchpad_harvests.tge_amount IS 'TGE 立即到账量';
COMMENT ON COLUMN launchpad_harvests.vesting_amount IS '锁仓量';

-- 12. Vesting 计划表
CREATE TABLE IF NOT EXISTS launchpad_vesting_schedules (
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL REFERENCES launchpad_sales(id),
    pool_index INT NOT NULL,
    schedule_id BIGINT NOT NULL,
    beneficiary CHAR(42) NOT NULL,
    amount_total NUMERIC(78,0) NOT NULL DEFAULT 0,
    released NUMERIC(78,0) NOT NULL DEFAULT 0,
    tx_hash CHAR(66) NOT NULL,
    block_number BIGINT NOT NULL,
    tx_index INT NOT NULL,
    log_index INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_vesting_schedules_event ON launchpad_vesting_schedules (block_number, tx_index, log_index);
CREATE INDEX IF NOT EXISTS idx_launchpad_vesting_schedules_beneficiary ON launchpad_vesting_schedules (beneficiary);
CREATE INDEX IF NOT EXISTS idx_launchpad_vesting_schedules_chain_id ON launchpad_vesting_schedules (schedule_id);

COMMENT ON TABLE launchpad_vesting_schedules IS 'Vesting 锁仓计划：每条记录对应一个用户在一个池的锁仓释放计划';
COMMENT ON COLUMN launchpad_vesting_schedules.schedule_id IS '链上 vesting schedule ID';
COMMENT ON COLUMN launchpad_vesting_schedules.beneficiary IS '受益人地址';
COMMENT ON COLUMN launchpad_vesting_schedules.amount_total IS '锁仓总量';
COMMENT ON COLUMN launchpad_vesting_schedules.released IS '已释放量';

-- 13. Vesting 释放记录表
CREATE TABLE IF NOT EXISTS launchpad_vesting_releases (
    id BIGSERIAL PRIMARY KEY,
    schedule_id BIGINT NOT NULL REFERENCES launchpad_vesting_schedules(id),
    amount NUMERIC(78,0) NOT NULL,
    tx_hash CHAR(66) NOT NULL,
    block_number BIGINT NOT NULL,
    tx_index INT NOT NULL,
    log_index INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_vesting_releases_event ON launchpad_vesting_releases (block_number, tx_index, log_index);
CREATE INDEX IF NOT EXISTS idx_launchpad_vesting_releases_schedule ON launchpad_vesting_releases (schedule_id);

COMMENT ON TABLE launchpad_vesting_releases IS 'Vesting 释放记录：每次 release 操作的事件明细';
COMMENT ON COLUMN launchpad_vesting_releases.amount IS '本次释放数量';

-- 14. Prepare 交易记录表
CREATE TABLE IF NOT EXISTS launchpad_prepare_txs (
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT,
    pool_index INT,
    operation_type VARCHAR(32) NOT NULL,
    caller_address CHAR(42) NOT NULL,
    calldata TEXT NOT NULL,
    calldata_hash CHAR(66) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    tx_hash CHAR(66),
    block_number BIGINT,
    error_message TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    confirmed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_launchpad_prepare_txs_hash ON launchpad_prepare_txs (calldata_hash, status);
CREATE INDEX IF NOT EXISTS idx_launchpad_prepare_txs_caller ON launchpad_prepare_txs (caller_address, status);
CREATE INDEX IF NOT EXISTS idx_launchpad_prepare_txs_poll ON launchpad_prepare_txs (status, updated_at);
CREATE INDEX IF NOT EXISTS idx_launchpad_prepare_txs_sale ON launchpad_prepare_txs (sale_id);

-- 去重：活跃状态的 calldata_hash 唯一约束，防止并发创建重复 PrepareTx
CREATE UNIQUE INDEX IF NOT EXISTS idx_launchpad_prepare_txs_active_hash ON launchpad_prepare_txs (calldata_hash) WHERE status IN ('pending', 'signed', 'broadcast');

COMMENT ON TABLE launchpad_prepare_txs IS 'Prepare 交易记录：Go 后端生成的链上交易 calldata 和状态追踪';
COMMENT ON COLUMN launchpad_prepare_txs.operation_type IS '操作类型：create_sale/set_pool/set_tier_limits/add_whitelist/remove_whitelist/set_start_end_block/revoke/final_withdraw/recover_token/deposit/harvest/release';
COMMENT ON COLUMN launchpad_prepare_txs.caller_address IS '发起者地址';
COMMENT ON COLUMN launchpad_prepare_txs.calldata IS 'ABI 编码的交易数据';
COMMENT ON COLUMN launchpad_prepare_txs.calldata_hash IS 'calldata 的 keccak256 哈希（去重用）';
COMMENT ON COLUMN launchpad_prepare_txs.status IS '状态：pending/signed/broadcast/confirmed/reverted/expired/failed';
COMMENT ON COLUMN launchpad_prepare_txs.expires_at IS '过期时间（创建后 30 分钟）';
COMMENT ON COLUMN launchpad_prepare_txs.confirmed_at IS '链上确认时间';
