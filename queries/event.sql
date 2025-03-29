-- name: GetEventByID :one
SELECT sqlc.embed(event), COALESCE(json_agg(timeslot_marker.*) FILTER (WHERE timeslot_marker.event_id IS NOT NULL), '[]')::json as markers FROM event
LEFT JOIN timeslot_marker ON event.id = timeslot_marker.event_id
WHERE event.id = sqlc.arg(id)
GROUP BY event.id;

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
SELECT sqlc.embed(event), COALESCE(json_agg(timeslot_marker.*) FILTER (WHERE timeslot_marker.event_id IS NOT NULL), '[]')::json as markers FROM event
LEFT JOIN timeslot_marker ON event.id = timeslot_marker.event_id
GROUP BY event.id
ORDER BY event.start_time ASC;

-- name: AddArtistToEvent :exec
INSERT INTO timeslot (id, event_id, artist_id, artist_name_override, sort_key)
VALUES (sqlc.arg(id), sqlc.arg(event_id), sqlc.arg(artist_id), sqlc.narg(artist_name_override), sqlc.arg(sort_key));

-- name: RemoveArtistFromEvent :exec
DELETE FROM timeslot
WHERE event_id = sqlc.arg(event_id) AND artist_id = sqlc.arg(artist_id);

-- name: TimeSlotsByEventID :many
SELECT sqlc.embed(timeslot), sqlc.embed(artist) FROM timeslot
JOIN artist ON timeslot.artist_id = artist.id
WHERE timeslot.event_id = sqlc.arg(event_id)
ORDER BY timeslot.sort_key ASC;

-- name: UpdateTimeSlot :many
UPDATE timeslot
SET artist_name_override = sqlc.arg(artist_name_override), sort_key = sqlc.arg(sort_key)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: CreateTimeslotMarker :one
INSERT INTO timeslot_marker (id, event_id, marker_type, marker_value, timeslot_index)
VALUES (sqlc.arg(id), sqlc.arg(event_id), sqlc.arg(marker_type), sqlc.arg(marker_value), sqlc.arg(timeslot_index)) RETURNING *;

-- name: DeleteTimeslotMarker :exec
DELETE FROM timeslot_marker
WHERE id = sqlc.arg(id);

-- name: UpdateTimeslotMarker :one
UPDATE timeslot_marker
SET marker_type = sqlc.arg(marker_type), marker_value = sqlc.arg(marker_value), timeslot_index = sqlc.arg(timeslot_index)
WHERE id = sqlc.arg(id) RETURNING *;