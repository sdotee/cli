package cmd

import "testing"

func TestParseBioLinks(t *testing.T) {
	links, err := parseBioLinks([]string{"Blog=https://example.com", "Profile=https://example.com/me"})
	if err != nil {
		t.Fatalf("parseBioLinks returned an error: %v", err)
	}
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(links))
	}
	if links[0].Title != "Blog" || links[0].URL != "https://example.com" {
		t.Fatalf("unexpected first link: %#v", links[0])
	}
}

func TestParseBioLinksRejectsInvalidValues(t *testing.T) {
	for _, value := range []string{"missing-separator", "=https://example.com", "Blog="} {
		if _, err := parseBioLinks([]string{value}); err == nil {
			t.Errorf("expected %q to be rejected", value)
		}
	}
}

func TestParsePositiveID(t *testing.T) {
	id, err := parsePositiveID("42", "bio page")
	if err != nil || id != 42 {
		t.Fatalf("expected ID 42, got %d, %v", id, err)
	}

	for _, value := range []string{"0", "-1", "invalid"} {
		if _, err := parsePositiveID(value, "bio page"); err == nil {
			t.Errorf("expected %q to be rejected", value)
		}
	}
}
