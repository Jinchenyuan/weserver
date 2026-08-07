package servicehandler

import (
	"context"
	"fmt"
	"server/model"
	pb "server/protobuf/gen"
	"strings"
	"time"

	"github.com/Jinchenyuan/wego"
	"github.com/Jinchenyuan/wego/logger"
)

type EventRecord struct {
	log *logger.Logger
}

func NewEventRecord(log *logger.Logger) *EventRecord {
	return &EventRecord{log: resolveEventRecordLogger(log)}
}

func resolveEventRecordLogger(log *logger.Logger) *logger.Logger {
	if log != nil {
		return log
	}
	if globalLog := wego.GetGlobalLogger(); globalLog != nil {
		return globalLog
	}
	return logger.GetLogger("eventrecord.service")
}

func (e *EventRecord) ListEventRecords(ctx context.Context, req *pb.ListEventRecordsRequest, rsp *pb.ListEventRecordsResponse) error {
	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}

	rows, err := model.ListEventRecordSummaries(ctx, m.DB, req.GetAccountId())
	if err != nil {
		return err
	}

	rsp.Events = make([]*pb.EventRecordSummary, 0, len(rows))
	for _, row := range rows {
		summary := &pb.EventRecordSummary{
			Id:            row.ID,
			Title:         row.Title,
			Location:      row.Location,
			OccurredAt:    row.OccurredAt.Format(time.RFC3339Nano),
			CoverPhotoUri: row.CoverPhotoURI,
			NoteCount:     row.NoteCount,
			UpdatedAt:     row.UpdatedAt.Format(time.RFC3339Nano),
		}
		if row.LatestNoteDate.Valid {
			summary.LatestNoteDate = row.LatestNoteDate.Time.Format(time.RFC3339Nano)
		}
		rsp.Events = append(rsp.Events, summary)
	}
	return nil
}

func (e *EventRecord) GetEventRecord(ctx context.Context, req *pb.GetEventRecordRequest, rsp *pb.EventRecordDetail) error {
	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}
	event, err := model.FindEventRecordByID(ctx, m.DB, req.GetAccountId(), req.GetId())
	if err != nil {
		return err
	}
	fillEventRecordDetail(rsp, event)
	return nil
}

func (e *EventRecord) CreateEventRecord(ctx context.Context, req *pb.CreateEventRecordRequest, rsp *pb.EventRecordMutationResponse) error {
	if err := validateEventRecordMutation(req.GetTitle(), req.GetOccurredAt(), req.GetCoverPhotoUri()); err != nil {
		rsp.Code = 400
		rsp.Message = err.Error()
		return nil
	}

	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}

	occurredAt, err := parseEventRecordTime(req.GetOccurredAt())
	if err != nil {
		rsp.Code = 400
		rsp.Message = err.Error()
		return nil
	}

	now := time.Now().UTC()
	event := model.NewEventRecord(
		req.GetAccountId(),
		strings.TrimSpace(req.GetTitle()),
		occurredAt,
		req.GetLocation(),
		req.GetDescription(),
		req.GetCoverPhotoUri(),
		now,
	)
	event.Notes = []*model.EventNote{}
	if err := model.CreateEventRecord(ctx, m.DB, event); err != nil {
		return err
	}

	rsp.Success = true
	rsp.Code = 201
	rsp.Message = "Event record created"
	rsp.Event = toProtoEventRecordDetail(event)
	return nil
}

func (e *EventRecord) UpdateEventRecord(ctx context.Context, req *pb.UpdateEventRecordRequest, rsp *pb.EventRecordMutationResponse) error {
	if err := validateEventRecordMutation(req.GetTitle(), req.GetOccurredAt(), req.GetCoverPhotoUri()); err != nil {
		rsp.Code = 400
		rsp.Message = err.Error()
		return nil
	}

	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}

	existing, err := model.LoadEventRecordForUpdate(ctx, m.DB, req.GetAccountId(), req.GetId())
	if err != nil {
		return err
	}
	occurredAt, err := parseEventRecordTime(req.GetOccurredAt())
	if err != nil {
		rsp.Code = 400
		rsp.Message = err.Error()
		return nil
	}

	now := time.Now().UTC()
	existing.Title = strings.TrimSpace(req.GetTitle())
	existing.OccurredAt = occurredAt
	existing.Location = req.GetLocation()
	existing.Description = req.GetDescription()
	existing.CoverPhotoURI = req.GetCoverPhotoUri()
	existing.UpdatedAt = now
	if err := model.UpdateEventRecord(ctx, m.DB, existing); err != nil {
		return err
	}

	reloaded, err := model.FindEventRecordByID(ctx, m.DB, req.GetAccountId(), req.GetId())
	if err != nil {
		return err
	}
	rsp.Success = true
	rsp.Code = 200
	rsp.Message = "Event record updated"
	rsp.Event = toProtoEventRecordDetail(reloaded)
	return nil
}

func (e *EventRecord) CreateEventNote(ctx context.Context, req *pb.CreateEventNoteRpcRequest, rsp *pb.EventRecordMutationResponse) error {
	if req.GetNote() == nil {
		rsp.Code = 400
		rsp.Message = "event note is required"
		return nil
	}
	if err := validateEventNoteMutation(req.GetNote().GetTitle()); err != nil {
		rsp.Code = 400
		rsp.Message = err.Error()
		return nil
	}

	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}
	now := time.Now().UTC()
	note := model.NewEventNote(
		req.GetEventId(),
		strings.TrimSpace(req.GetNote().GetTitle()),
		req.GetNote().GetContent(),
		toNullableString(req.GetNote().GetPhotoUri()),
		now,
	)
	if err := model.CreateEventNote(ctx, m.DB, req.GetAccountId(), req.GetEventId(), note, now); err != nil {
		return err
	}

	reloaded, err := model.FindEventRecordByID(ctx, m.DB, req.GetAccountId(), req.GetEventId())
	if err != nil {
		return err
	}
	rsp.Success = true
	rsp.Code = 201
	rsp.Message = "Event note created"
	rsp.Event = toProtoEventRecordDetail(reloaded)
	return nil
}

func (e *EventRecord) UpdateEventNote(ctx context.Context, req *pb.UpdateEventNoteRpcRequest, rsp *pb.EventRecordMutationResponse) error {
	if req.GetNote() == nil {
		rsp.Code = 400
		rsp.Message = "event note is required"
		return nil
	}
	if err := validateEventNoteMutation(req.GetNote().GetTitle()); err != nil {
		rsp.Code = 400
		rsp.Message = err.Error()
		return nil
	}

	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}
	existing, err := model.FindEventRecordByID(ctx, m.DB, req.GetAccountId(), req.GetNote().GetEventId())
	if err != nil {
		return err
	}

	var createdAt time.Time
	found := false
	for _, note := range existing.Notes {
		if note.ID == req.GetNote().GetNoteId() {
			createdAt = note.CreatedAt
			found = true
			break
		}
	}
	if !found {
		return model.ErrEventNoteNotFound
	}

	now := time.Now().UTC()
	note := &model.EventNote{
		ID:       req.GetNote().GetNoteId(),
		EventID:  req.GetNote().GetEventId(),
		Title:    strings.TrimSpace(req.GetNote().GetTitle()),
		Content:  req.GetNote().GetContent(),
		PhotoURI: toNullString(req.GetNote().GetPhotoUri()),
		Base: model.Base{
			CreatedAt: createdAt,
			UpdatedAt: now,
		},
	}
	if err := model.UpdateEventNote(ctx, m.DB, req.GetAccountId(), req.GetNote().GetEventId(), note, now); err != nil {
		return err
	}

	reloaded, err := model.FindEventRecordByID(ctx, m.DB, req.GetAccountId(), req.GetNote().GetEventId())
	if err != nil {
		return err
	}
	rsp.Success = true
	rsp.Code = 200
	rsp.Message = "Event note updated"
	rsp.Event = toProtoEventRecordDetail(reloaded)
	return nil
}

func (e *EventRecord) DeleteEventNote(ctx context.Context, req *pb.DeleteEventNoteRequest, rsp *pb.EventRecordMutationResponse) error {
	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}
	now := time.Now().UTC()
	if err := model.DeleteEventNote(ctx, m.DB, req.GetAccountId(), req.GetEventId(), req.GetNoteId(), now); err != nil {
		return err
	}

	reloaded, err := model.FindEventRecordByID(ctx, m.DB, req.GetAccountId(), req.GetEventId())
	if err != nil {
		return err
	}
	rsp.Success = true
	rsp.Code = 200
	rsp.Message = "Event note deleted"
	rsp.Event = toProtoEventRecordDetail(reloaded)
	return nil
}
