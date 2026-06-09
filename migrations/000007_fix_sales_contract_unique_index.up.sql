-- 修复 contract_address 唯一索引：draft sale 的 contract_address 为空字符串，
-- 多条 draft 记录会导致唯一约束冲突。排除空字符串后再建唯一索引。

DROP INDEX IF EXISTS idx_launchpad_sales_contract;
CREATE UNIQUE INDEX idx_launchpad_sales_contract ON launchpad_sales (contract_address) WHERE deleted_at IS NULL AND contract_address != '';
