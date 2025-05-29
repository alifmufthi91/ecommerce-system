BEGIN;

CREATE TABLE warehouses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_warehouses_status ON warehouses (status);

CREATE TABLE warehouse_stocks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    warehouse_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    reserved INTEGER NOT NULL CHECK (reserved >= 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_warehouse_stocks_warehouse_id ON warehouse_stocks (warehouse_id);
CREATE INDEX idx_warehouse_stocks_product_id ON warehouse_stocks (product_id);

CREATE TABLE stock_transfers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_warehouse_id UUID NOT NULL,
    to_warehouse_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_stock_transfers_from_warehouse_id ON stock_transfers (from_warehouse_id);
CREATE INDEX idx_stock_transfers_to_warehouse_id ON stock_transfers (to_warehouse_id);
CREATE INDEX idx_stock_transfers_product_id ON stock_transfers (product_id);

CREATE TABLE shop_warehouses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shop_id UUID NOT NULL,
    warehouse_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_shop_warehouses_shop_id ON shop_warehouses (shop_id);
CREATE INDEX idx_shop_warehouses_warehouse_id ON shop_warehouses (warehouse_id);

INSERT INTO warehouses (id, name, address, status) VALUES
('8f1cc115-4434-4829-81c4-23fb01aa0dc0', 'Warehouse A', '123 Warehouse St', 'active'),
('a0ebb46d-6482-405c-a340-c4a144591fce', 'Warehouse B', '456 Warehouse Ave', 'active'),
('c3c2c921-8abd-4d1c-ba56-4a9e7d8c9df6', 'Warehouse C', '789 Warehouse Blvd', 'inactive');

INSERT INTO warehouse_stocks (id, warehouse_id, product_id, quantity, reserved) VALUES
(uuid_generate_v4(), '8f1cc115-4434-4829-81c4-23fb01aa0dc0', '9a2b7c93-7c27-4e20-842f-24bf4df95bf0', 100, 0),
(uuid_generate_v4(), 'a0ebb46d-6482-405c-a340-c4a144591fce', '14c0374f-0fa3-4a02-baff-04e226910d3b', 200, 0),
(uuid_generate_v4(), 'c3c2c921-8abd-4d1c-ba56-4a9e7d8c9df6', '2ae686f2-fd3e-4672-aef2-7cc1e4b5f3b0', 150, 0);

INSERT INTO shop_warehouses (id, shop_id, warehouse_id) VALUES
(uuid_generate_v4(), '122a579e-e7b6-4f78-8979-7556fd66b59e', '8f1cc115-4434-4829-81c4-23fb01aa0dc0'),
(uuid_generate_v4(), 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'a0ebb46d-6482-405c-a340-c4a144591fce'),
(uuid_generate_v4(), 'c56a4180-65aa-42ec-a945-5fd21dec0538', 'c3c2c921-8abd-4d1c-ba56-4a9e7d8c9df6');

COMMIT;