-- +goose Up
CREATE TABLE posts (
    id uuid primary key,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    title text NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at timestamp,
    feed_id uuid NOT NULL references feeds(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE IF EXISTS posts;
