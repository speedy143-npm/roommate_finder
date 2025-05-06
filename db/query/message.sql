-- name: CreateUser :one
INSERT INTO users (fname,lname,phoneno,email,password,bio,preferences,profile_picture)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUserById :many
SELECT * FROM users
WHERE id = $1;

-- name: CreateMatch :one
INSERT INTO matches (user1_id, user2_id, match_score, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserMatchById :many
SELECT * FROM users
WHERE id1 = $1, id2 = $2;