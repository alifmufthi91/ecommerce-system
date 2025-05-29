BEGIN;

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    shop_id UUID NOT NULL,
    shop_name TEXT NOT NULL DEFAULT '',
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_shop_id ON products (shop_id);

INSERT INTO products (id, name, description, shop_id, shop_name, price) VALUES
('9a2b7c93-7c27-4e20-842f-24bf4df95bf0', 'Product A', 'Description for Product A', '122a579e-e7b6-4f78-8979-7556fd66b59e', 'Shop A', 19.99),
('14c0374f-0fa3-4a02-baff-04e226910d3b', 'Product B', 'Description for Product B', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Shop B', 29.99),
('2ae686f2-fd3e-4672-aef2-7cc1e4b5f3b0', 'Product C', 'Description for Product C', 'c56a4180-65aa-42ec-a945-5fd21dec0538', 'Shop C', 39.99);

COMMIT;