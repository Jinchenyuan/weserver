package ginhandler

import (
	"errors"
	pb "server/protobuf/gen"
	"strings"
)

func toCreateEventRecordRequest(accountID uint32, req CreateEventRecordRequest) *pb.CreateEventRecordRequest {
	return &pb.CreateEventRecordRequest{
		AccountId:     accountID,
		Title:         req.Title,
		OccurredAt:    req.OccurredAt,
		Location:      req.Location,
		Description:   req.Description,
		CoverPhotoUri: req.CoverPhotoURI,
	}
}

func toUpdateEventRecordRequest(accountID uint32, req UpdateEventRecordRequest) *pb.UpdateEventRecordRequest {
	return &pb.UpdateEventRecordRequest{
		AccountId:     accountID,
		Id:            req.ID,
		Title:         req.Title,
		OccurredAt:    req.OccurredAt,
		Location:      req.Location,
		Description:   req.Description,
		CoverPhotoUri: req.CoverPhotoURI,
	}
}

func toCreateEventNoteRpcRequest(accountID uint32, eventID string, req CreateEventNoteRequest) *pb.CreateEventNoteRpcRequest {
	return &pb.CreateEventNoteRpcRequest{
		AccountId: accountID,
		EventId:   eventID,
		Note: &pb.CreateEventNoteRequest{
			Title:    req.Title,
			Content:  req.Content,
			PhotoUri: normalizeOptionalString(req.PhotoURI),
		},
	}
}

func toUpdateEventNoteRpcRequest(accountID uint32, req UpdateEventNoteRequest) *pb.UpdateEventNoteRpcRequest {
	return &pb.UpdateEventNoteRpcRequest{
		AccountId: accountID,
		Note: &pb.UpdateEventNoteRequest{
			EventId:  req.EventID,
			NoteId:   req.NoteID,
			Title:    req.Title,
			Content:  req.Content,
			PhotoUri: normalizeOptionalString(req.PhotoURI),
		},
	}
}

func toEventListResponse(rsp *pb.ListEventRecordsResponse) ListEventRecordsResponse {
	events := make([]EventRecordSummary, 0, len(rsp.GetEvents()))
	for _, item := range rsp.GetEvents() {
		events = append(events, EventRecordSummary{
			ID:             item.GetId(),
			Title:          item.GetTitle(),
			Location:       item.GetLocation(),
			OccurredAt:     item.GetOccurredAt(),
			CoverPhotoURI:  item.GetCoverPhotoUri(),
			LatestNoteDate: toOptionalString(item.GetLatestNoteDate()),
			NoteCount:      item.GetNoteCount(),
			UpdatedAt:      item.GetUpdatedAt(),
		})
	}
	return ListEventRecordsResponse{Events: events}
}

func toEventMutationResponse(rsp *pb.EventRecordMutationResponse) EventRecordMutationResponse {
	return EventRecordMutationResponse{
		Success: rsp.GetSuccess(),
		Event:   toEventDetail(rsp.GetEvent()),
		Message: toOptionalString(rsp.GetMessage()),
	}
}

func toEventDetail(detail *pb.EventRecordDetail) EventRecordDetail {
	notes := make([]EventNote, 0, len(detail.GetNotes()))
	for _, note := range detail.GetNotes() {
		notes = append(notes, EventNote{
			ID:        note.GetId(),
			CreatedAt: note.GetCreatedAt(),
			Title:     note.GetTitle(),
			Content:   note.GetContent(),
			PhotoURI:  toOptionalString(note.GetPhotoUri()),
		})
	}
	return EventRecordDetail{
		ID:            detail.GetId(),
		Title:         detail.GetTitle(),
		OccurredAt:    detail.GetOccurredAt(),
		Location:      detail.GetLocation(),
		Description:   detail.GetDescription(),
		CoverPhotoURI: detail.GetCoverPhotoUri(),
		CreatedAt:     detail.GetCreatedAt(),
		UpdatedAt:     detail.GetUpdatedAt(),
		Notes:         notes,
	}
}

func validateEventRecordRequest(title, occurredAt, coverPhotoURI string) error {
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

func validateEventNoteRequest(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("event note title is required")
	}
	return nil
}

func normalizeOptionalString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func toOptionalString(value string) *string {
	if value == "" {
		return nil
	}
	ret := value
	return &ret
}
