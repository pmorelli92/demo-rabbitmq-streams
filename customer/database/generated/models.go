// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package gen_sql

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Customer struct {
	ID        string
	Name      string
	Email     string
	Address   string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}
