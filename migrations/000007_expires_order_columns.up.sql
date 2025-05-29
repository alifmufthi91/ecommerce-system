BEGIN;

ALTER TABLE stock_locks
    DROP COLUMN IF EXISTS expires_at;

ALTER TABLE orders
    ADD COLUMN expires_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp;

COMMIT;