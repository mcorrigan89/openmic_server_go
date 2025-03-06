-- name: GetArtistByID :one
SELECT sqlc.embed(artist) FROM artist
WHERE artist.id = $1;

-- name: GetAllArtists :many
SELECT sqlc.embed(artist) FROM artist
ORDER BY artist.artist_title ASC;

-- name: GetArtistsByTitle :many
SELECT sqlc.embed(artist) FROM artist
WHERE similarity(artist.artist_title, sqlc.arg(title)) > sqlc.arg(min_similarity)
ORDER BY similarity(artist.artist_title, sqlc.arg(title)) DESC;

-- name: CreateArtist :one
INSERT INTO artist (id, artist_title, artist_subtitle, bio, avatar_id)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateArtist :one
UPDATE artist
SET artist_title = $2, artist_subtitle = $3, bio = $4, avatar_id = $5
WHERE id = $1 RETURNING *;