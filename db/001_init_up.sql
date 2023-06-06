CREATE TABLE IF NOT EXISTS link (
    original_url    VARCHAR(512),
    abbreviated_url VARCHAR(10) UNIQUE
);
