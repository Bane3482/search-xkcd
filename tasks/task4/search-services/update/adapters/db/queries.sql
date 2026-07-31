-- name: AddComics :exec
INSERT INTO comics (
    comics_id,
    comics_url,
    words
)
VALUES (
    sqlc.arg(comics_id),
    sqlc.arg(comics_url),
    sqlc.arg(words)
)
ON CONFLICT (comics_id) DO NOTHING;

-- name: GetWords :many
SELECT words
FROM comics;

-- name: GetIDs :many
SELECT comics_id
FROM comics;

-- name: DropComics :exec
DELETE FROM comics;