package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestPrintJSON(t *testing.T) {
	data := map[string]string{"key": "value"}
	var buf bytes.Buffer

	err := printJSON(&buf, data)
	if err != nil {
		t.Fatalf("printJSON failed: %v", err)
	}

	expected := "{\n  \"key\": \"value\"\n}\n"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestReadContent_File(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "hello world"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	cmd := &cobra.Command{}
	result, err := readContent(tmpFile, cmd)
	if err != nil {
		t.Fatalf("readContent failed: %v", err)
	}

	if result != content {
		t.Errorf("expected %q, got %q", content, result)
	}
}

func TestReadContent_FileError(t *testing.T) {
	cmd := &cobra.Command{}
	_, err := readContent("non-existent-file.txt", cmd)
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestReadContent_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty.txt")
	if err := os.WriteFile(tmpFile, []byte(""), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	cmd := &cobra.Command{}
	_, err := readContent(tmpFile, cmd)
	if err == nil {
		t.Error("expected error for empty file, got nil")
	}
	if err.Error() != "file is empty" {
		t.Errorf("expected 'file is empty' error, got %v", err)
	}
}

func TestReadContent_NonTextFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "binary.bin")
	data := []byte{0x00, 0xFF, 0x10, 0x00, 0x01}
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		t.Fatalf("failed to create binary file: %v", err)
	}

	cmd := &cobra.Command{}
	_, err := readContent(tmpFile, cmd)
	if err == nil {
		t.Fatal("expected error for binary file, got nil")
	}
	if !strings.Contains(err.Error(), "non-text content") {
		t.Fatalf("expected non-text content error, got %v", err)
	}
}

func TestReadContent_Stdin(t *testing.T) {
	content := "hello from stdin"
	buf := bytes.NewBufferString(content)

	cmd := &cobra.Command{}
	cmd.SetIn(buf)

	// Use "-" to force reading from stdin (bypassing os.Stdin check)
	result, err := readContent("-", cmd)
	if err != nil {
		t.Fatalf("readContent failed: %v", err)
	}

	if result != content {
		t.Errorf("expected %q, got %q", content, result)
	}
}

func TestReadContent_EmptyStdin(t *testing.T) {
	buf := bytes.NewBufferString("   ") // Empty after trim

	cmd := &cobra.Command{}
	cmd.SetIn(buf)

	// Use "-" to force reading from stdin
	_, err := readContent("-", cmd)
	if err == nil {
		t.Error("expected error for empty stdin, got nil")
	}
	if err.Error() != "content is empty" {
		t.Errorf("expected 'content is empty' error, got %v", err)
	}
}
