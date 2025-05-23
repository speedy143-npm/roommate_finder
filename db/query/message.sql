-- name: CreateUser :one
INSERT INTO users (fname,lname,phoneno,email,password,bio,preferences,profile_picture)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUserById :many
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: CreateMatch :one
INSERT INTO matches (user1_id, user2_id, match_score, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUserProfile :one
UPDATE users
SET
    fname = COALESCE($2, fname),
    lname = COALESCE($3, lname),
    phoneno = COALESCE($4, phoneno),
    email = COALESCE($5, email),
    password = COALESCE($6, password),
    bio = COALESCE($7, bio),
    preferences = COALESCE($8, preferences),
    profile_picture = COALESCE($9, profile_picture)
WHERE id = $1
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: ForgotPassword :one
INSERT INTO password_resets (user_id, token)
VALUES ($1, $2)
RETURNING *;

-- name: GetResetToken :many
SELECT * FROM password_resets
WHERE token = $1
LIMIT 1;

-- name: UpdateUserPassword :one
UPDATE users
SET password = $2
WHERE id = $1
RETURNING *;

-- name: DeleteResetToken :exec
DELETE FROM password_resets WHERE token = $1;

-- name: DeleteExpiredTokens :exec
DELETE FROM password_resets WHERE expiry < NOW();




