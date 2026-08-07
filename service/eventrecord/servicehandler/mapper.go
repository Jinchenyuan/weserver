package servicehandler

import (
	"database/sql"
	"errors"
	"server/model"
	pb "server/protobuf/gen"
	"strings"
	"time"
)

func validateEventRecordMutation(title, occurredAt, coverPhotoURI string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("event title is required")
	}
	if strings.TrimSpace(occurredAt) == "" {
		return errors.New("event occurredAt is required")
	}
	if strings.TrimSpace(coverPhotoURI) == "" {
		return errors.New("event coverPhotoUri is required")
	}
	return nil
}

func validateEventNoteMutation(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("event note title is required")
	}
	return nil
}

func parseEventRecordTime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05",
	}

	var lastErr error
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed.UTC(), nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

func toProtoEventRecordDetail(event *model.EventRecord) *pb.EventRecordDetail {
	notes := model.CopyEventNotesByCreatedDesc(event.Notes)
	ret := &pb.EventRecordDetail{
		Id:            event.ID,
		Title:         event.Title,
		OccurredAt:    event.OccurredAt.Format(time.RFC3339Nano),
		Location:      event.Location,
		Description:   event.Description,
		CoverPhotoUri: event.CoverPhotoURI,
		CreatedAt:     event.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:     event.UpdatedAt.Format(time.RFC3339Nano),
		Notes:         make([]*pb.EventNote, 0, len(notes)),
	}
	for _, note := range notes {
		ret.Notes = append(ret.Notes, &pb.EventNote{
			Id:        note.ID,
			CreatedAt: note.CreatedAt.Format(time.RFC3339Nano),
			Title:     note.Title,
			Content:   note.Content,
			PhotoUri:  nullStringValue(note.PhotoURI),
		})
	}
	return ret
}

func fillEventRecordDetail(rsp *pb.EventRecordDetail, event *model.EventRecord) {
	detail := toProtoEventRecordDetail(event)
	*rsp = *detail
}

func toNullableString(value string) *string {
	if value == "" {
		return nil
	}
	ret := value
	return &ret
}

func toNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func nullStringValue(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
