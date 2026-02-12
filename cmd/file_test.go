//
// Copyright (c) 2026 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: file_test.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2026-01-19 18:36:48
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2026-01-19 19:25:51
//

package cmd

import (
	"testing"
)

func TestFileCmd_Structure(t *testing.T) {
	if fileCmd.Use != "file" {
		t.Errorf("expected use 'file', got '%s'", fileCmd.Use)
	}

	commands := fileCmd.Commands()
	hasUpload := false
	hasDelete := false
	hasDomains := false
	hasHistory := false

	for _, c := range commands {
		switch c.Name() {
		case "upload":
			hasUpload = true
		case "delete":
			hasDelete = true
		case "domains":
			hasDomains = true
		case "history":
			hasHistory = true
		}
	}

	if !hasUpload {
		t.Error("file command missing 'upload' subcommand")
	}
	if !hasDelete {
		t.Error("file command missing 'delete' subcommand")
	}
	if !hasDomains {
		t.Error("file command missing 'domains' subcommand")
	}
	if !hasHistory {
		t.Error("file command missing 'history' subcommand")
	}
}

func TestFileUploadCmd_Flags(t *testing.T) {
	// Reset flags
	fileUploadOpts.file = ""
	fileUploadOpts.name = ""
	fileUploadOpts.isPrivate = 0

	f := fileUploadCmd.Flags()
	if f.Lookup("file") == nil {
		t.Error("upload command missing 'file' flag")
	}
	if f.Lookup("name") == nil {
		t.Error("upload command missing 'name' flag")
	}
	if f.Lookup("is-private") == nil {
		t.Error("upload command missing 'is-private' flag")
	}
}

func TestFileHistoryCmd_Flags(t *testing.T) {
	fileHistoryOpts.page = 1

	f := fileHistoryCmd.Flags()
	if f.Lookup("page") == nil {
		t.Error("history command missing 'page' flag")
	}
}

func TestFileDeleteCmd_Args(t *testing.T) {
	// Test that delete command requires at least 1 argument
	err := fileDeleteCmd.Args(fileDeleteCmd, []string{})
	if err == nil {
		t.Error("expected error when no args provided to delete command")
	}

	err = fileDeleteCmd.Args(fileDeleteCmd, []string{"key123"})
	if err != nil {
		t.Errorf("unexpected error with 1 arg: %v", err)
	}

	err = fileDeleteCmd.Args(fileDeleteCmd, []string{"key1", "key2"})
	if err != nil {
		t.Errorf("unexpected error with 2 args: %v", err)
	}
}
