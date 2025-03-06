-- name: GetEventByID :one
SELECT sqlc.embed(event) FROM event
WHERE event.id = sqlc.arg(id);

-- name: CreateEvent :one
INSERT INTO event (id, event_type, start_time, end_time)
VALUES (sqlc.arg(id), sqlc.arg(event_type), sqlc.arg(start_time), sqlc.arg(end_time)) RETURNING *;

-- name: UpdateEvent :one
UPDATE event
SET event_type = sqlc.arg(event_type), start_time = sqlc.arg(start_time), end_time = sqlc.arg(end_time)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM event
WHERE id = sqlc.arg(id);

-- name: GetAllEvents :many
SELECT sqlc.embed(event) FROM event
ORDER BY event.start_time ASC;

-- name: AddArtistToEvent :exec
INSERT INTO timeslot (id, event_id, artist_id, artist_name_override, sort_order)
VALUES (sqlc.arg(id), sqlc.arg(event_id), sqlc.arg(artist_id), sqlc.narg(artist_name_override), sqlc.arg(sort_order));

-- name: RemoveArtistFromEvent :exec
DELETE FROM timeslot
WHERE event_id = sqlc.arg(event_id) AND artist_id = sqlc.arg(artist_id);

-- name: TimeSlotsByEventID :many
SELECT sqlc.embed(timeslot), sqlc.embed(artist) FROM timeslot
JOIN artist ON timeslot.artist_id = artist.id
WHERE timeslot.event_id = sqlc.arg(event_id)
ORDER BY timeslot.sort_order ASC;
