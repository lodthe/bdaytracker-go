BEGIN;

CREATE TABLE IF NOT EXISTS states (
    telegram_id bigint primary key not null,
    created_at timestamp without time zone default now() not null,
    updated_at timestamp without time zone default now() not null,
    state jsonb
);

CREATE INDEX IF NOT EXISTS state_created_at_idx ON states USING btree(created_at);
CREATE INDEX IF NOT EXISTS state_updated_at_idx ON states USING btree(updated_at);

COMMIT;