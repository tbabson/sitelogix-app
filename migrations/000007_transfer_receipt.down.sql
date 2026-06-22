DROP INDEX IF EXISTS idx_inv_tx_pending;

ALTER TABLE inventory_transactions
    DROP COLUMN IF EXISTS received,
    DROP COLUMN IF EXISTS received_by,
    DROP COLUMN IF EXISTS received_at;
