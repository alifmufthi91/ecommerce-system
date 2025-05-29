BEGIN;

ALTER TABLE orders
    DROP COLUMN IF EXISTS expires_at;

ALTER TABLE stock_locks
    ADD COLUMN expires_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp;
    
COMMIT;