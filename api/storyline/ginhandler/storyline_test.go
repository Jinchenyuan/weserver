package ginhandler

import "testing"

func TestValidateNodesRejectsEmptyInput(t *testing.T) {
	if err := validateNodes(nil); err == nil {
		t.Fatal("expected empty nodes validation error")
	}
}

func TestValidateNodesRejectsBlankTitle(t *testing.T) {
	err := validateNodes([]StorylineNodeInput{{
		Title: "",
		Date:  "2026-05-05T00:00:00Z",
	}})
	if err == nil {
		t.Fatal("expected blank title validation error")
	}
}
