BEGIN;

CREATE TABLE shops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO shops (id, name, address) VALUES
('122a579e-e7b6-4f78-8979-7556fd66b59e','Shop A', '123 Main St'),
('f47ac10b-58cc-4372-a567-0e02b2c3d479','Shop B', '456 Elm St'),
('c56a4180-65aa-42ec-a945-5fd21dec0538','Shop C', '789 Oak St');


COMMIT;