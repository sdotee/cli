//
// Copyright (c) 2025 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: utils.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2025-12-22 22:27:28
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2025-12-23 11:08:17
//

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"
)

// printJSON encodes the given value as JSON and writes it to the writer.
// It uses indentation for better readability.
func printJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// readContent reads content from the specified file path.
// If filePath is "-", it reads from stdin.
// It returns an error if the file or stdin is empty.
func readContent(filePath string, cmd *cobra.Command) (string, error) {
	var (
		b   []byte
		err error
	)

	if filePath != "" && filePath != "-" {
		b, err = os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		if len(b) == 0 {
			return "", errors.New("file is empty")
		}
	} else {
		// Only check if stdin is a terminal if we are not explicitly asked to read from stdin
		if filePath != "-" {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				return "", errors.New("no input: provide --file <path> or pipe content via stdin")
			}
		}

		b, err = io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return "", err
		}
		if len(b) == 0 {
			return "", errors.New("content is empty")
		}
	}

	if err := ensureTextContent(b); err != nil {
		return "", err
	}

	content := string(b)
	if len(strings.TrimSpace(content)) == 0 {
		return "", errors.New("content is empty")
	}
	return content, nil
}

// ensureTextContent rejects non-text input based on MIME sniffing to avoid
// uploading binary data to text endpoints.
func ensureTextContent(b []byte) error {
	mt := mimetype.Detect(b)
	mimeType := mt.String()
	if idx := strings.Index(mimeType, ";"); idx != -1 {
		mimeType = strings.TrimSpace(mimeType[:idx])
	}

	if isAllowedTextMIME(mimeType) {
		return nil
	}

	return fmt.Errorf("non-text content detected (%s); only text input is allowed", mimeType)
}

func isAllowedTextMIME(mimeType string) bool {
	if strings.HasPrefix(mimeType, "text/") {
		return true
	}

	switch mimeType {
	case "application/json",
		"application/xml",
		"application/javascript",
		"application/x-sh",
		"application/x-httpd-php",
		"application/x-yaml":
		return true
	default:
		return false
	}
}
