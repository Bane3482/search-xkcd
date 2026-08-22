CREATE SCHEMA IF NOT EXISTS search_schema;

CREATE TABLE IF NOT EXISTS search_schema.keywords (
    keyword TEXT NOT NULL PRIMARY KEY,
    ids INTEGER[] NOT NULL
);

CREATE INDEX idx_keyword_hash ON search_schema.keywords USING hash (keyword); 