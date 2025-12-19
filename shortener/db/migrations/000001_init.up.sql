CREATE TABLE IF NOT EXISTS links
(
    id         SERIAL PRIMARY KEY,
    source_url TEXT      NOT NULL,
    short_url  TEXT      NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_short_url ON links (short_url);