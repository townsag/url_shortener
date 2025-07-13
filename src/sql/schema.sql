CREATE TABLE url_mapping (
    id VARCHAR(8) PRIMARY KEY,
    long_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    visits INTEGER DEFAULT 0
);