-- name: GetImageByID :one
SELECT sqlc.embed(images) FROM images
WHERE images.id = sqlc.arg(id);

-- name: CreateImage :one
INSERT INTO images (id, bucket_name, object_id, width, height, file_size) 
VALUES (sqlc.arg(id), sqlc.arg(bucket_name), sqlc.arg(object_id), sqlc.arg(width), sqlc.arg(height), sqlc.arg(file_size)) RETURNING *;
