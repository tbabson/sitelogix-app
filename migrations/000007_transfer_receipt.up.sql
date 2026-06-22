ALTER TABLE inventory_transactions
    ADD COLUMN received    BOOLEAN     NOT NULL DEFAULT false,
    ADD COLUMN received_by UUID        REFERENCES users(id),
    ADD COLUMN received_at TIMESTAMPTZ;

-- Fast lookup for the "pending deliveries" query
CREATE INDEX idx_inv_tx_pending ON inventory_transactions(project_id, received)
    WHERE transaction_type = 'transfer_in' AND received = false;
