package cmd

import "testing"

func TestShortURLCmdURLAlias(t *testing.T) {
	found, _, err := rootCmd.Find([]string{"url", "create"})
	if err != nil {
		t.Fatalf("unexpected error resolving 'url create': %v", err)
	}
	if found != shorturlCreateCmd {
		t.Errorf("expected 'url create' to resolve to shorturl create, got %q", found.CommandPath())
	}
}
