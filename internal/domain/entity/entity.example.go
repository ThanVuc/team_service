package entity

import "github.com/jackc/pgx/v5/pgtype"

type Group struct {
	ID          pgtype.UUID
	Name        string
	Description pgtype.Text
	OwnerID     pgtype.UUID
	CreatedAt   pgtype.Timestamptz
}
