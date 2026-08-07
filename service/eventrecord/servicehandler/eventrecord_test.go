package servicehandler

import (
	"testing"

	"github.com/Jinchenyuan/wego"
	"github.com/Jinchenyuan/wego/logger"
)

func TestNewEventRecordUsesInjectedLogger(t *testing.T) {
	injected := logger.NewLogger("eventrecord-test")
	handler := NewEventRecord(injected)

	if handler.log != injected {
		t.Fatal("expected handler to use injected logger")
	}
}

func TestNewEventRecordFallsBackToGlobalLogger(t *testing.T) {
	global := logger.NewLogger("eventrecord-global")
	wego.SetGlobalLogger(global)

	handler := NewEventRecord(nil)
	if handler.log != global {
		t.Fatal("expected handler to use global logger")
	}
}

func TestValidateEventRecordMutationRejectsInvalidInput(t *testing.T) {
	if err := validateEventRecordMutation("", "2026-05-17T00:00:00", "cover"); err == nil {
		t.Fatal("expected missing title validation error")
	}

	if err := validateEventRecordMutation("title", "", "cover"); err == nil {
		t.Fatal("expected missing occurredAt validation error")
	}

	if err := validateEventRecordMutation("title", "2026-05-17T00:00:00", ""); err == nil {
		t.Fatal("expected missing coverPhotoUri validation error")
	}
}
