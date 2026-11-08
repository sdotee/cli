package cmd

import (
	"strings"
	"testing"
)

func TestTextTitleOrDefault(t *testing.T) {
	tests := map[string]string{
		"":        "Untitled",
		"  ":      "Untitled",
		"My note": "My note",
	}
	for in, want := range tests {
		if got := textTitleOrDefault(in); got != want {
			t.Errorf("textTitleOrDefault(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestTextUpdateCmdRequiresTitle(t *testing.T) {
	textUpdateOpts.title = " "
	t.Cleanup(func() { textUpdateOpts.title = "" })

	err := textUpdateCmd.RunE(textUpdateCmd, []string{"slug"})
	if err == nil || !strings.Contains(err.Error(), "--title is required") {
		t.Errorf("expected --title error, got %v", err)
	}
}
