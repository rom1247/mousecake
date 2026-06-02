-- 回滚 launchpad 模块所有表（按依赖顺序逆序 DROP）

DROP TABLE IF EXISTS launchpad_prepare_txs;
DROP TABLE IF EXISTS launchpad_vesting_releases;
DROP TABLE IF EXISTS launchpad_vesting_schedules;
DROP TABLE IF EXISTS launchpad_harvests;
DROP TABLE IF EXISTS launchpad_user_credit;
DROP TABLE IF EXISTS launchpad_user_pool_state;
DROP TABLE IF EXISTS launchpad_deposits;
DROP TABLE IF EXISTS launchpad_whitelists;
DROP TABLE IF EXISTS launchpad_tier_limits;
DROP TABLE IF EXISTS launchpad_pools;
DROP TABLE IF EXISTS launchpad_sale_meta;
DROP TABLE IF EXISTS launchpad_sales;
DROP TABLE IF EXISTS launchpad_tier_params;
DROP TABLE IF EXISTS launchpad_tokens;
