-- name: CreateFeedFollow :one
WITH feed_follow_insert AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT 
    feed_follow_insert.*,
    users.name AS user , 
    feeds.name AS feed
FROM  feed_follow_insert
INNER JOIN users
ON feed_follow_insert.user_id = users.id
INNER JOIN feeds
ON feed_follow_insert.feed_id = feeds.id;

-- name: GetFeedFollowByUser :many
SELECT 
    feed_follows.*,
    users.name AS user, 
    feeds.name AS feed
FROM  feed_follows
INNER JOIN users
ON feed_follows.user_id = users.id
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
WHERE users.id = $1;

-- name: DeleteFeedFollowById :exec
DELETE FROM  feed_follows
WHERE user_id = $1 AND feed_id = $2;