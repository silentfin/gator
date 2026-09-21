-- +goose Up
ALTER TABLE IF EXISTS feeds ADD last_fetched_at timestamp;

-- +goose Down
ALTER TABLE IF EXISTS feeds DROP COLUMN last_fetched_at;
