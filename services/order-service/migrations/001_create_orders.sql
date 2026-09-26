CREATE TABLE IF NOT EXISTS orders (
    id          VARCHAR(64) PRIMARY KEY,
    customer_id VARCHAR(64) NOT NULL,
    address     TEXT NOT NULL,
    lat         DOUBLE PRECISION,
    lng         DOUBLE PRECISION,
    status      VARCHAR(32) NOT NULL,
    description TEXT,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_status      ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_created_at  ON orders(created_at DESC);