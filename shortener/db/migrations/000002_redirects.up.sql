CREATE TABLE IF NOT EXISTS redirects
(
    short_url varchar(6), -- no PK! super nice
    click_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_agent TEXT NOT NULL
);