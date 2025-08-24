package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Account struct {
	ID        int64     `bun:",pk,autoincrement"`
	OwnerID   int64     `bun:",notnull"`
	Balance   float64   `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,notnull"`
	UpdatedAt time.Time `bun:",nullzero,notnull"`
	db        *bun.DB   `bun:"-" json:"-"`
}

type AccountRepository interface {
	Create(ctx context.Context) error
	Update(ctx context.Context, column string) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*Account, error)
	FindAll(ctx context.Context) ([]*Account, error)
}

func (a *Account) SetDB(db *bun.DB) {
	a.db = db
}

func (a *Account) Create(ctx context.Context) error {
	_, err := a.db.NewInsert().Model(a).Exec(ctx)
	return err
}

func (a *Account) Update(ctx context.Context, column string) error {
	_, err := a.db.NewUpdate().Model(a).Column(column).WherePK().Exec(ctx)
	return err
}

func (a *Account) Delete(ctx context.Context, id int64) error {
	_, err := a.db.NewDelete().Model(a).Where("id = ?", id).Exec(ctx)
	return err
}

func (a *Account) FindByID(ctx context.Context, id int64) (*Account, error) {
	account := &Account{}
	err := a.db.NewSelect().Model(account).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a *Account) FindAll(ctx context.Context) ([]*Account, error) {
	var accounts []*Account
	err := a.db.NewSelect().Model(&accounts).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
