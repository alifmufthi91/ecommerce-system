BEGIN;

CREATE UNIQUE INDEX IF NOT EXISTS stock_unique_key ON warehouse_stocks(warehouse_id, product_id);

COMMIT;