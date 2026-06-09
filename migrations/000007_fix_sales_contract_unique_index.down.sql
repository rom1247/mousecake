DROP INDEX IF EXISTS idx_launchpad_sales_contract;
CREATE UNIQUE INDEX idx_launchpad_sales_contract ON launchpad_sales (contract_address) WHERE deleted_at IS NULL;
