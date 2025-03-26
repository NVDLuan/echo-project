-- name: GetUserByEmail :one
SELECT id, password FROM users WHERE email= $1 LIMIT 1;

