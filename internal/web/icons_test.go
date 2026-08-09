package web

import (
	"strings"
	"testing"
)

func TestTypeIcon(t *testing.T) {
	// Every type the classifier can produce should have an icon.
	for _, dt := range []string{"Phone", "Computer", "Printer", "TV", "Media Player", "Router", "Server", "IoT", "Game Console"} {
		html := string(typeIcon(dt))
		if !strings.Contains(html, "<svg") {
			t.Errorf("type %q has no icon", dt)
		}
	}
	// An unknown type yields no icon, not broken markup.
	if got := typeIcon("Nonsense"); got != "" {
		t.Errorf("unknown type should yield empty, got %q", got)
	}
	if got := typeIcon(""); got != "" {
		t.Errorf("empty type should yield empty, got %q", got)
	}
}
