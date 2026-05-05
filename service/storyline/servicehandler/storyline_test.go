package servicehandler

import (
	"testing"

	pb "server/protobuf/gen"

	"github.com/Jinchenyuan/wego"
	"github.com/Jinchenyuan/wego/logger"
)

func TestNewStorylineUsesInjectedLogger(t *testing.T) {
	injected := logger.NewLogger("storyline-test")
	handler := NewStoryline(injected)

	if handler.log != injected {
		t.Fatal("expected handler to use injected logger")
	}
}

func TestNewStorylineFallsBackToGlobalLogger(t *testing.T) {
	global := logger.NewLogger("storyline-global")
	wego.SetGlobalLogger(global)

	handler := NewStoryline(nil)
	if handler.log != global {
		t.Fatal("expected handler to use global logger")
	}
}

func TestValidateStorylineMutationRejectsInvalidNodes(t *testing.T) {
	err := validateStorylineMutation("title", []*pb.StorylineNodeInput{})
	if err == nil {
		t.Fatal("expected empty nodes validation error")
	}

	err = validateStorylineMutation("title", []*pb.StorylineNodeInput{{Title: "", Date: "2026-05-05T00:00:00Z"}})
	if err == nil {
		t.Fatal("expected empty node title validation error")
	}
}
