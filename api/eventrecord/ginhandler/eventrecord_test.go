package ginhandler

import "testing"

func TestValidateEventRecordRequestRejectsInvalidInput(t *testing.T) {
	if err := validateEventRecordRequest("", "2026-05-17T00:00:00", "cover"); err == nil {
		t.Fatal("expected missing title validation error")
	}
	if err := validateEventRecordRequest("title", "", "cover"); err == nil {
		t.Fatal("expected missing occurredAt validation error")
	}
	if err := validateEventRecordRequest("title", "2026-05-17T00:00:00", ""); err == nil {
		t.Fatal("expected missing coverPhotoUri validation error")
	}
}

func TestValidateEventNoteRequestRejectsBlankTitle(t *testing.T) {
	if err := validateEventNoteRequest(""); err == nil {
		t.Fatal("expected missing note title validation error")
	}
}
