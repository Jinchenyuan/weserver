package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type S3KV struct {
	bun.BaseModel `bun:"table:s3kv,alias:kv"`
	Key           string    `bun:",pk"`
	Value         string    `bun:"column:value,notnull"`
	CreatedAt     time.Time `bun:",nullzero,notnull"`
	UpdatedAt     time.Time `bun:",nullzero,notnull"`
	db            *bun.DB   `bun:"-" json:"-"`
}

type S3KVRepository interface {
	Create(ctx context.Context) error
	Update(ctx context.Context, column string) error
	Delete(ctx context.Context, key string) error
	FindByKey(ctx context.Context, key string) (*S3KV, error)
}

func (s *S3KV) SetDB(db *bun.DB) {
	s.db = db
}

func (s *S3KV) Create(ctx context.Context) error {
	_, err := s.db.NewInsert().Model(s).Exec(ctx)
	return err
}

func (s *S3KV) Update(ctx context.Context, column string) error {
	_, err := s.db.NewUpdate().Model(s).Column(column).WherePK().Exec(ctx)
	return err
}

func (s *S3KV) Delete(ctx context.Context, key string) error {
	_, err := s.db.NewDelete().Model(s).Where("key = ?", key).Exec(ctx)
	return err
}

func (s *S3KV) FindByKey(ctx context.Context, key string) (*S3KV, error) {
	s3kv := &S3KV{}
	err := s.db.NewSelect().Model(s3kv).Where("key = ?", key).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return s3kv, nil
}
