-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS teams (
    team_name VARCHAR(255) NOT NULL UNIQUE,
    user_id VARCHAR(255) REFERENCES users(user_id)

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
