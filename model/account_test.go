package model

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func TestAddAccount(t *testing.T) {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN("postgres://user:password@localhost:5432/land_contract?sslmode=disable")))
	db := bun.NewDB(sqlDb, pgdialect.New())
	defer db.Close()

	account := &Account{
		ID:        1,
		Account:   "testaccount",
		Name:      "testuser",
		Email:     "testuser@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	account.SetDB(db)
	err := account.Create(context.Background())
	if err != nil {
		t.Errorf("failed to create account: %v\n", err)
	}
}

func TestUpdateAccount(t *testing.T) {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN("postgres://user:password@10.4.11.140:5432/land_contract?sslmode=disable")))
	db := bun.NewDB(sqlDb, pgdialect.New())
	defer db.Close()

	account := &Account{
		ID:        1,
		Account:   "testaccount",
		Name:      "testuser",
		Email:     "testuser@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	account.SetDB(db)
	err := account.Create(context.Background())
	if err != nil {
		t.Errorf("failed to create account: %v\n", err)
	}

	account.Email = "newemail@example.com"
	err = account.Update(context.Background(), "email")
	if err != nil {
		t.Errorf("failed to update account: %v\n", err)
	}
}

func TestDeleteAccount(t *testing.T) {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN("postgres://user:password@10.4.11.140:5432/land_contract?sslmode=disable")))
	db := bun.NewDB(sqlDb, pgdialect.New())
	defer db.Close()

	account := &Account{
		ID: 1,
	}
	account.SetDB(db)
	err := account.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("failed to delete account: %v\n", err)
	}
}

func TestFindAccountByID(t *testing.T) {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN("postgres://user:password@10.4.11.140:5432/land_contract?sslmode=disable")))
	db := bun.NewDB(sqlDb, pgdialect.New())
	defer db.Close()

	account := &Account{
		ID: 0,
	}
	account.SetDB(db)
	result, err := FindAccountByID(context.Background(), db, 2)
	if err != nil {
		t.Errorf("failed to find account by ID: %v\n", err)
	}
	t.Logf("found account: %+v\n", result)
}

func TestFindAllAccounts(t *testing.T) {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN("postgres://user:password@10.4.11.140:5432/land_contract?sslmode=disable")))
	db := bun.NewDB(sqlDb, pgdialect.New())
	defer db.Close()

	account := &Account{
		ID:      0,
		Account: "testaccount",
	}
	account.SetDB(db)
	result, err := FindAllAccount(context.Background(), db)
	if err != nil {
		t.Errorf("failed to find any account: %v\n", err)
	}
	for _, a := range result {
		t.Logf("account: %+v\n", a)
	}
}
