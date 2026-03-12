package servicehandler

import (
	"testing"

	"github.com/Jinchenyuan/wego"
	"github.com/Jinchenyuan/wego/logger"
)

func TestNewAccountUsesInjectedLogger(t *testing.T) {
	injected := logger.NewLogger("account-test")
	handler := NewAccount(injected)

	if handler.log != injected {
		t.Fatal("expected handler to use injected logger")
	}
}

func TestNewAccountFallsBackToGlobalLogger(t *testing.T) {
	global := logger.NewLogger("global-test")
	wego.SetGlobalLogger(global)

	handler := NewAccount(nil)
	if handler.log != global {
		t.Fatal("expected handler to use global logger")
	}
}

func TestResolveLoggerFallsBackToPackageLogger(t *testing.T) {
	wego.SetGlobalLogger(nil)

	resolved := resolveLogger(nil)
	if resolved == nil {
		t.Fatal("expected fallback logger")
	}

	if resolved != logger.GetLogger("account.service") {
		t.Fatal("expected account service fallback logger")
	}
}
