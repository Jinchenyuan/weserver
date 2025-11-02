package model

import (
	"context"
	"database/sql"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func TestAddS3KV(t *testing.T) {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN("postgres://user:password@localhost:5432/land_contract?sslmode=disable")))
	db := bun.NewDB(sqlDb, pgdialect.New())
	defer db.Close()

	s3kv := &S3KV{
		Key:   "example_key",
		Value: "example_value",
	}
	s3kv.SetDB(db)

	ctx := context.Background()
	err := s3kv.Create(ctx)
	if err != nil {
		t.Errorf("failed to create s3kv: %v\n", err)
	}
}
