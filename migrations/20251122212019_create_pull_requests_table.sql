-- +goose Up
-- +goose StatementBegin
CREATE TYPE pull_request_status AS ENUM ('OPEN', 'MERGED');
CREATE TABLE IF NOT EXISTS pull_requests (
    pr_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) REFERENCES users(user_id),
    status pull_request_status NOT NULL DEFAULT 'OPEN'
    merged_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pull_requests;
-- +goose StatementEnd
