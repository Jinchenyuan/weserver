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

var ErrEventRecordNotFound = errors.New("event record not found")
var ErrEventNoteNotFound = errors.New("event note not found")

type EventRecord struct {
	bun.BaseModel `bun:"table:event_records"`
	ID            string    `bun:",pk"`
	AccountID     uint32    `bun:"account_id,notnull"`
	Title         string    `bun:",notnull"`
	OccurredAt    time.Time `bun:"occurred_at,notnull"`
	Location      string    `bun:",notnull"`
	Description   string    `bun:",notnull"`
	CoverPhotoURI string    `bun:"cover_photo_uri,notnull"`
	Base
	Notes []*EventNote `bun:"rel:has-many,join:id=event_id"`
}

type EventNote struct {
	bun.BaseModel `bun:"table:event_notes"`
	ID            string         `bun:",pk"`
	EventID       string         `bun:"event_id,notnull"`
	Title         string         `bun:",notnull"`
	Content       string         `bun:",notnull"`
	PhotoURI      sql.NullString `bun:"photo_uri"`
	Base
}

type EventRecordSummaryRow struct {
	ID             string       `bun:"id"`
	Title          string       `bun:"title"`
	Location       string       `bun:"location"`
	OccurredAt     time.Time    `bun:"occurred_at"`
	CoverPhotoURI  string       `bun:"cover_photo_uri"`
	LatestNoteDate sql.NullTime `bun:"latest_note_date"`
	NoteCount      int32        `bun:"note_count"`
	UpdatedAt      time.Time    `bun:"updated_at"`
	SortLatestAt   time.Time    `bun:"sort_latest_at"`
}

func (e *EventRecord) SetDB(db *bun.DB) {
	e.db = db
}

func FindEventRecordByID(ctx context.Context, db *bun.DB, accountID uint32, id string) (*EventRecord, error) {
	event := &EventRecord{}
	err := db.NewSelect().
		Model(event).
		Relation("Notes", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("created_at DESC")
		}).
		Where("event_record.id = ?", id).
		Where("event_record.account_id = ?", accountID).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventRecordNotFound
		}
		return nil, err
	}
	return event, nil
}

func LoadEventRecordForUpdate(ctx context.Context, db *bun.DB, accountID uint32, id string) (*EventRecord, error) {
	event := &EventRecord{}
	err := db.NewSelect().
		Model(event).
		Where("id = ?", id).
		Where("account_id = ?", accountID).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventRecordNotFound
		}
		return nil, err
	}
	return event, nil
}

func ListEventRecordSummaries(ctx context.Context, db *bun.DB, accountID uint32) ([]*EventRecordSummaryRow, error) {
	rows := make([]*EventRecordSummaryRow, 0)
	err := db.NewSelect().
		TableExpr("event_records AS e").
		ColumnExpr("e.id").
		ColumnExpr("e.title").
		ColumnExpr("e.location").
		ColumnExpr("e.occurred_at").
		ColumnExpr("e.cover_photo_uri").
		ColumnExpr("e.updated_at").
		ColumnExpr("COUNT(n.id)::int AS note_count").
		ColumnExpr("MAX(n.created_at) AS latest_note_date").
		ColumnExpr("COALESCE(MAX(n.created_at), e.updated_at) AS sort_latest_at").
		Join("LEFT JOIN event_notes AS n ON n.event_id = e.id").
		Where("e.account_id = ?", accountID).
		GroupExpr("e.id").
		GroupExpr("e.title").
		GroupExpr("e.location").
		GroupExpr("e.occurred_at").
		GroupExpr("e.cover_photo_uri").
		GroupExpr("e.updated_at").
		OrderExpr("sort_latest_at DESC").
		OrderExpr("e.occurred_at DESC").
		Scan(ctx, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func CreateEventRecord(ctx context.Context, db *bun.DB, event *EventRecord) error {
	_, err := db.NewInsert().Model(event).Exec(ctx)
	return err
}

func UpdateEventRecord(ctx context.Context, db *bun.DB, event *EventRecord) error {
	res, err := db.NewUpdate().
		Model(event).
		Column("title", "occurred_at", "location", "description", "cover_photo_uri", "updated_at").
		WherePK().
		Where("account_id = ?", event.AccountID).
		Exec(ctx)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrEventRecordNotFound
	}
	return nil
}

func CreateEventNote(ctx context.Context, db *bun.DB, accountID uint32, eventID string, note *EventNote, updatedAt time.Time) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		ok, err := eventRecordExists(ctx, tx, accountID, eventID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrEventRecordNotFound
		}
		if _, err := tx.NewInsert().Model(note).Exec(ctx); err != nil {
			return err
		}
		_, err = tx.NewUpdate().
			Model((*EventRecord)(nil)).
			Table("event_records").
			Set("updated_at = ?", updatedAt).
			Where("id = ?", eventID).
			Where("account_id = ?", accountID).
			Exec(ctx)
		return err
	})
}

func UpdateEventNote(ctx context.Context, db *bun.DB, accountID uint32, eventID string, note *EventNote, updatedAt time.Time) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		ok, err := eventRecordExists(ctx, tx, accountID, eventID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrEventRecordNotFound
		}

		res, err := tx.NewUpdate().
			Model(note).
			Column("title", "content", "photo_uri", "updated_at").
			WherePK().
			Where("event_id = ?", eventID).
			Exec(ctx)
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return ErrEventNoteNotFound
		}

		_, err = tx.NewUpdate().
			Model((*EventRecord)(nil)).
			Table("event_records").
			Set("updated_at = ?", updatedAt).
			Where("id = ?", eventID).
			Where("account_id = ?", accountID).
			Exec(ctx)
		return err
	})
}

func DeleteEventNote(ctx context.Context, db *bun.DB, accountID uint32, eventID, noteID string, updatedAt time.Time) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		ok, err := eventRecordExists(ctx, tx, accountID, eventID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrEventRecordNotFound
		}

		res, err := tx.NewDelete().
			Model((*EventNote)(nil)).
			Where("id = ?", noteID).
			Where("event_id = ?", eventID).
			Exec(ctx)
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return ErrEventNoteNotFound
		}

		_, err = tx.NewUpdate().
			Model((*EventRecord)(nil)).
			Table("event_records").
			Set("updated_at = ?", updatedAt).
			Where("id = ?", eventID).
			Where("account_id = ?", accountID).
			Exec(ctx)
		return err
	})
}

func eventRecordExists(ctx context.Context, db bun.IDB, accountID uint32, eventID string) (bool, error) {
	count, err := db.NewSelect().
		Model((*EventRecord)(nil)).
		Where("id = ?", eventID).
		Where("account_id = ?", accountID).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func NewEventRecord(accountID uint32, title string, occurredAt time.Time, location, description, coverPhotoURI string, now time.Time) *EventRecord {
	return &EventRecord{
		ID:            newEventRecordID(),
		AccountID:     accountID,
		Title:         title,
		OccurredAt:    occurredAt,
		Location:      location,
		Description:   description,
		CoverPhotoURI: coverPhotoURI,
		Base: Base{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func NewEventNote(eventID, title, content string, photoURI *string, now time.Time) *EventNote {
	note := &EventNote{
		ID:      newEventNoteID(),
		EventID: eventID,
		Title:   title,
		Content: content,
		Base: Base{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	if photoURI != nil {
		note.PhotoURI = sql.NullString{String: *photoURI, Valid: true}
	}
	return note
}

func CopyEventNotesByCreatedDesc(notes []*EventNote) []*EventNote {
	copied := append([]*EventNote(nil), notes...)
	sort.SliceStable(copied, func(i, j int) bool {
		return copied[i].CreatedAt.After(copied[j].CreatedAt)
	})
	return copied
}

func newEventRecordID() string {
	return fmt.Sprintf("%s", NewStringID())
}

func newEventNoteID() string {
	return fmt.Sprintf("%s", NewStringID())
}
