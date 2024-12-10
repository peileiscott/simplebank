CREATE TABLE IF NOT EXISTS accounts (
    id         BIGSERIAL   PRIMARY KEY,
    owner      TEXT        NOT NULL,
    balance    BIGINT      NOT NULL,
    currency   VARCHAR(3)  NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS entries (
    id         BIGSERIAL   PRIMARY KEY,
    account_id BIGINT      NOT NULL REFERENCES accounts (id),
    amount     BIGINT      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transfers (
    id              BIGSERIAL   PRIMARY KEY,
    from_account_id BIGINT      NOT NULL REFERENCES accounts (id),
    to_account_id   BIGINT      NOT NULL REFERENCES accounts (id),
    amount          BIGINT      NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT check_amount_positive CHECK (amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_entries_owner ON accounts (owner);

CREATE INDEX IF NOT EXISTS idx_entries_account_id ON entries (account_id);

CREATE INDEX IF NOT EXISTS idx_transfers_from_account_id ON transfers (from_account_id);

CREATE INDEX IF NOT EXISTS idx_transfers_to_account_id ON transfers (to_account_id);

CREATE INDEX IF NOT EXISTS idx_transfers_from_to_account_id ON transfers (from_account_id, to_account_id);

COMMENT ON COLUMN entries.amount IS 'can be negative or positive';

COMMENT ON COLUMN transfers.amount IS 'must be positive';