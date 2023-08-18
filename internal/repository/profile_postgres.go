package repository

import (
	"github.com/jmoiron/sqlx"
)

type ProfilePostgres struct {
	db *sqlx.DB
}

func NewProfilePostgres(db *sqlx.DB) *ProfilePostgres {
	return &ProfilePostgres{
		db: db,
	}
}
