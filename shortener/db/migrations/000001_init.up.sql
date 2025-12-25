CREATE TABLE IF NOT EXISTS links
(
    source_url TEXT                     NOT NULL,  -- yeah I do edit migration files before I ever deployed my DB
    short_url  VARCHAR(6)               NOT NULL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_short_url ON links (short_url);