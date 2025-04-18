// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: image.sql

package models

import (
	"context"

	"github.com/google/uuid"
)

const createImage = `-- name: CreateImage :one
INSERT INTO images (id, bucket_name, object_id, width, height, file_size) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, bucket_name, object_id, height, width, file_size, created_at, updated_at, version
`

type CreateImageParams struct {
	ID         uuid.UUID `json:"id"`
	BucketName string    `json:"bucket_name"`
	ObjectID   string    `json:"object_id"`
	Width      int32     `json:"width"`
	Height     int32     `json:"height"`
	FileSize   int32     `json:"file_size"`
}

func (q *Queries) CreateImage(ctx context.Context, arg CreateImageParams) (Image, error) {
	row := q.db.QueryRow(ctx, createImage,
		arg.ID,
		arg.BucketName,
		arg.ObjectID,
		arg.Width,
		arg.Height,
		arg.FileSize,
	)
	var i Image
	err := row.Scan(
		&i.ID,
		&i.BucketName,
		&i.ObjectID,
		&i.Height,
		&i.Width,
		&i.FileSize,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Version,
	)
	return i, err
}

const getImageByID = `-- name: GetImageByID :one
SELECT images.id, images.bucket_name, images.object_id, images.height, images.width, images.file_size, images.created_at, images.updated_at, images.version FROM images
WHERE images.id = $1
`

type GetImageByIDRow struct {
	Image Image `json:"image"`
}

func (q *Queries) GetImageByID(ctx context.Context, id uuid.UUID) (GetImageByIDRow, error) {
	row := q.db.QueryRow(ctx, getImageByID, id)
	var i GetImageByIDRow
	err := row.Scan(
		&i.Image.ID,
		&i.Image.BucketName,
		&i.Image.ObjectID,
		&i.Image.Height,
		&i.Image.Width,
		&i.Image.FileSize,
		&i.Image.CreatedAt,
		&i.Image.UpdatedAt,
		&i.Image.Version,
	)
	return i, err
}
