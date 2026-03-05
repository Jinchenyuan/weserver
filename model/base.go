package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Base struct {
	CreatedAt time.Time `bun:",nullzero,notnull"`
	UpdatedAt time.Time `bun:",nullzero,notnull"`
	db        *bun.DB   `bun:"-" json:"-"`
}
