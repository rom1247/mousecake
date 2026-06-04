-- 回滚迁移：移除新增的列

ALTER TABLE launchpad_prepare_txs DROP COLUMN IF EXISTS value;
ALTER TABLE launchpad_prepare_txs DROP COLUMN IF EXISTS target_address;
ALTER TABLE launchpad_sales DROP COLUMN IF EXISTS status;
