CREATE TABLE IF NOT EXISTS links (
    code text PRIMARY KEY,
    long_url text NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_links_code
ON links (code);
CREATE UNIQUE INDEX IF NOT EXISTS idx_links_long_url
ON links (long_url);