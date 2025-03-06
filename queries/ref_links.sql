-- name: GetReferenceLinkByID :one
SELECT sqlc.embed(reference_link) FROM reference_link
WHERE reference_link.id = $1;

-- name: GetReferenceLinkByToken :one
SELECT sqlc.embed(reference_link) FROM reference_link
WHERE reference_link.token = $1;

-- name: CreateReferenceLink :one
INSERT INTO reference_link (id, link_id, link_type, token, expires_at) 
VALUES (sqlc.arg(id), sqlc.arg(link_id), sqlc.arg(link_type), sqlc.arg(token), sqlc.arg(expires_at)) RETURNING *;

-- name: DeleteReferenceLink :one
DELETE FROM reference_link WHERE reference_link.id = $1 RETURNING *;