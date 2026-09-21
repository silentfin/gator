-- +goose Up
CREATE TABLE feeds (
    id uuid primary key,
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    user_id uuid NOT NULL references users(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE IF EXISTS feeds;
