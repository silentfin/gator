-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at=NOW(), updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * from feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
