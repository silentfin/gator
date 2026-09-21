-- +goose Up
CREATE TABLE feed_follows (
    id uuid primary key,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    user_id uuid NOT NULL references users(id) ON DELETE CASCADE,
    feed_id uuid NOT NULL references feeds(id) ON DELETE CASCADE,
    UNIQUE(user_id, feed_id)
);


-- +goose Down
DROP TABLE IF EXISTS feed_follows;
