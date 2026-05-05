package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/uptrace/bun"
)

var ErrStorylineNotFound = errors.New("storyline not found")

type Storyline struct {
	bun.BaseModel `bun:"table:storylines"`
	ID            string         `bun:",pk"`
	AccountID     uint32         `bun:"account_id,notnull"`
	Title         string         `bun:",notnull"`
	Description   string         `bun:",notnull"`
	CoverPhotoURI sql.NullString `bun:"cover_photo_uri"`
	Base
	Nodes []*StorylineNode `bun:"rel:has-many,join:id=storyline_id"`
}

type StorylineNode struct {
	bun.BaseModel `bun:"table:storyline_nodes"`
	ID            string         `bun:",pk"`
	StorylineID   string         `bun:"storyline_id,notnull"`
	Title         string         `bun:",notnull"`
	Date          time.Time      `bun:",notnull"`
	Note          string         `bun:",notnull"`
	Location      string         `bun:",notnull"`
	PhotoURI      sql.NullString `bun:"photo_uri"`
	SortOrder     int32          `bun:"sort_order,notnull"`
	Base
}

type StorylineSummaryRow struct {
	ID              string         `bun:"id"`
	Title           string         `bun:"title"`
	Description     string         `bun:"description"`
	CoverPhotoURI   sql.NullString `bun:"cover_photo_uri"`
	NodeCount       int32          `bun:"node_count"`
	LatestNodeTitle sql.NullString `bun:"latest_node_title"`
	LatestNodeDate  sql.NullTime   `bun:"latest_node_date"`
	UpdatedAt       time.Time      `bun:"updated_at"`
	SortLatestAt    time.Time      `bun:"sort_latest_at"`
}

type StorylineInput struct {
	ID        string
	Title     string
	Date      time.Time
	Note      string
	Location  string
	PhotoURI  *string
	SortOrder int32
}

func (s *Storyline) SetDB(db *bun.DB) {
	s.db = db
}

func FindStorylineByID(ctx context.Context, db *bun.DB, accountID uint32, id string) (*Storyline, error) {
	storyline := &Storyline{}
	err := db.NewSelect().
		Model(storyline).
		Relation("Nodes", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("sort_order ASC")
		}).
		Where("storyline.id = ?", id).
		Where("storyline.account_id = ?", accountID).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStorylineNotFound
		}
		return nil, err
	}
	return storyline, nil
}

func ListStorylineSummaries(ctx context.Context, db *bun.DB, accountID uint32) ([]*StorylineSummaryRow, error) {
	rows := make([]*StorylineSummaryRow, 0)
	err := db.NewSelect().
		TableExpr("storylines AS s").
		ColumnExpr("s.id").
		ColumnExpr("s.title").
		ColumnExpr("s.description").
		ColumnExpr("s.cover_photo_uri").
		ColumnExpr("s.updated_at").
		ColumnExpr("COUNT(n.id)::int AS node_count").
		ColumnExpr("MAX(n.date) AS latest_node_date").
		ColumnExpr(`(
			SELECT n2.title
			FROM storyline_nodes AS n2
			WHERE n2.storyline_id = s.id
			ORDER BY n2.date DESC, n2.sort_order DESC
			LIMIT 1
		) AS latest_node_title`).
		ColumnExpr("COALESCE(MAX(n.date), s.updated_at) AS sort_latest_at").
		Join("LEFT JOIN storyline_nodes AS n ON n.storyline_id = s.id").
		Where("s.account_id = ?", accountID).
		GroupExpr("s.id").
		GroupExpr("s.title").
		GroupExpr("s.description").
		GroupExpr("s.cover_photo_uri").
		GroupExpr("s.updated_at").
		OrderExpr("sort_latest_at DESC").
		OrderExpr("s.updated_at DESC").
		Scan(ctx, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func CreateStoryline(ctx context.Context, db *bun.DB, storyline *Storyline, nodes []*StorylineNode) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(storyline).Exec(ctx); err != nil {
			return err
		}
		if len(nodes) == 0 {
			return nil
		}
		if _, err := tx.NewInsert().Model(&nodes).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
}

func ReplaceStorylineNodes(ctx context.Context, db *bun.DB, storyline *Storyline, nodes []*StorylineNode) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(storyline).
			Column("title", "description", "cover_photo_uri", "updated_at").
			WherePK().
			Where("account_id = ?", storyline.AccountID).
			Exec(ctx)
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return ErrStorylineNotFound
		}

		if _, err := tx.NewDelete().Model((*StorylineNode)(nil)).Where("storyline_id = ?", storyline.ID).Exec(ctx); err != nil {
			return err
		}
		if len(nodes) == 0 {
			return nil
		}
		if _, err := tx.NewInsert().Model(&nodes).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
}

func LoadStorylineForUpdate(ctx context.Context, db *bun.DB, accountID uint32, id string) (*Storyline, error) {
	storyline := &Storyline{}
	err := db.NewSelect().
		Model(storyline).
		Where("id = ?", id).
		Where("account_id = ?", accountID).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStorylineNotFound
		}
		return nil, err
	}
	return storyline, nil
}

func NormalizeStorylineNodes(storylineID string, now time.Time, inputs []StorylineInput) []*StorylineNode {
	sortedInputs := append([]StorylineInput(nil), inputs...)
	sort.SliceStable(sortedInputs, func(i, j int) bool {
		if sortedInputs[i].SortOrder == sortedInputs[j].SortOrder {
			return sortedInputs[i].Date.Before(sortedInputs[j].Date)
		}
		return sortedInputs[i].SortOrder < sortedInputs[j].SortOrder
	})

	nodes := make([]*StorylineNode, 0, len(sortedInputs))
	for idx, input := range sortedInputs {
		nodeID := input.ID
		if nodeID == "" {
			nodeID = newStorylineID()
		}
		node := &StorylineNode{
			ID:          nodeID,
			StorylineID: storylineID,
			Title:       input.Title,
			Date:        input.Date,
			Note:        input.Note,
			Location:    input.Location,
			SortOrder:   int32(idx),
			Base: Base{
				CreatedAt: now,
				UpdatedAt: now,
			},
		}
		if input.PhotoURI != nil {
			node.PhotoURI = sql.NullString{String: *input.PhotoURI, Valid: true}
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func NewStorylineRecord(accountID uint32, title, description string, coverPhotoURI *string, now time.Time) *Storyline {
	storyline := &Storyline{
		ID:          newStorylineID(),
		AccountID:   accountID,
		Title:       title,
		Description: description,
		Base: Base{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	if coverPhotoURI != nil {
		storyline.CoverPhotoURI = sql.NullString{String: *coverPhotoURI, Valid: true}
	}
	return storyline
}

func newStorylineID() string {
	return fmt.Sprintf("%s", NewStringID())
}
