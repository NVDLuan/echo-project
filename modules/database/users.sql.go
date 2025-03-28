// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
)

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, password FROM users WHERE email= $1 LIMIT 1
`

type GetUserByEmailRow struct {
	ID       int32
	Password string
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(&i.ID, &i.Password)
	return i, err
}
