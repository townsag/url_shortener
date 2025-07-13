-- name: InsertMapping :one
INSERT INTO url_mapping (id, long_url, created_at)
VALUES ($1, $2, $3)
ON CONFLICT (id) DO NOTHING
RETURNING id;

-- name: SelectMapping :one
SELECT * FROM url_mapping
WHERE id = $1 LIMIT 1;