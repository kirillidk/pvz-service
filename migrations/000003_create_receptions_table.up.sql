CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS receptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date_time TIMESTAMP NOT NULL,
    pvz_id UUID NOT NULL REFERENCES pvz(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('in_progress', 'close'))
);

CREATE INDEX IF NOT EXISTS idx_receptions_pvz_id ON receptions (pvz_id);