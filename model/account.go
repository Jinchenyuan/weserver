package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Account struct {
	ID        uint32    `bun:",pk"`
	Account   string    `bun:",notnull,unique"`
	Name      string    `bun:",notnull"`
	Email     string    `bun:",notnull,unique"`
	Password  string    `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,notnull"`
	UpdatedAt time.Time `bun:",nullzero,notnull"`
	db        *bun.DB   `bun:"-" json:"-"`
}

type AccountRepository interface {
	Create(ctx context.Context) error
	Update(ctx context.Context, column string) error
	Delete(ctx context.Context, id int64) error
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

func FindAccountByID(ctx context.Context, db *bun.DB, id int64) (*Account, error) {
	account := &Account{}
	err := db.NewSelect().Model(account).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func FindAllAccount(ctx context.Context, db *bun.DB) ([]*Account, error) {
	var accounts []*Account
	err := db.NewSelect().Model(&accounts).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func FindAccountByAccount(ctx context.Context, db *bun.DB, ac string) (*Account, error) {
	account := &Account{}
	err := db.NewSelect().Model(account).Where("account = ?", ac).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return account, nil
}
