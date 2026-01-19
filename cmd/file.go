//
// Copyright (c) 2026 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: file.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2026-01-19 18:36:26
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2026-01-19 19:25:41
//

package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	fileUploadOpts struct {
		file string
		name string
	}
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Manage file uploads",
}

var fileUploadCmd = &cobra.Command{
	Use:   "upload [file...]",
	Short: "Upload one or more files",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Collect all files to upload
		var filesToUpload []string
		filesToUpload = append(filesToUpload, args...)
		if fileUploadOpts.file != "" && fileUploadOpts.file != "-" {
			filesToUpload = append(filesToUpload, fileUploadOpts.file)
		}

		// Case 1: Stdin
		if len(filesToUpload) == 0 || (len(filesToUpload) == 1 && filesToUpload[0] == "-") {
			if fileUploadOpts.name == "" {
				return fmt.Errorf("filename must be provided via --name when reading from stdin")
			}
			return uploadReader(cmd, fileUploadOpts.name, cmd.InOrStdin())
		}

		// Case 2: Multiple files
		if len(filesToUpload) > 1 && fileUploadOpts.name != "" {
			return fmt.Errorf("cannot use --name with multiple files")
		}

		for _, filePath := range filesToUpload {
			if filePath == "-" {
				continue // Skip explicit stdin marker in multi-file mode or handle it? simplified to skip/error
			}

			f, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("failed to open file %q: %w", filePath, err)
			}
			// We defer close inside a loop, which is not ideal for many files, but fine for CLI typical usage.
			// Better to wrap in func.
			err = func() error {
				defer f.Close()
				filename := filepath.Base(filePath)
				// If strictly single file and name is provided (already checked above for >1), use it.
				// But we checked >1. If len==1, we can use name.
				if len(filesToUpload) == 1 && fileUploadOpts.name != "" {
					filename = fileUploadOpts.name
				}
				return uploadReader(cmd, filename, f)
			}()
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func uploadReader(cmd *cobra.Command, filename string, reader io.Reader) error {
	resp, err := apiClient.UploadFile(filename, reader)
	if err != nil {
		return err
	}

	if rootOpts.jsonOutput {
		// If multiple files are uploaded in JSON mode, this will produce multiple JSON objects
		// concatenated. This is "JSON Lines" format usually.
		return printJSON(cmd.OutOrStdout(), resp.Data)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "File uploaded successfully: %s\n", filename)
	fmt.Fprintf(cmd.OutOrStdout(), "URL: %s\n", resp.Data.URL)
	fmt.Fprintf(cmd.OutOrStdout(), "Delete Key: %s\n", resp.Data.Delete)
	fmt.Fprintf(cmd.OutOrStdout(), "Page: %s\n", resp.Data.Page)
	fmt.Fprintln(cmd.OutOrStdout(), "---")
	return nil
}

var fileDeleteCmd = &cobra.Command{
	Use:   "delete <key...>",
	Short: "Delete one or more files",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, deleteKey := range args {
			resp, err := apiClient.DeleteFile(deleteKey)
			if err != nil {
				return fmt.Errorf("failed to delete file with key %q: %w", deleteKey, err)
			}

			if rootOpts.jsonOutput {
				printJSON(cmd.OutOrStdout(), resp)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "File with key arg %q deleted successfully\n", deleteKey)
			}
		}
		return nil
	},
}

var fileDomainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "List available file domains",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.GetFileDomains()
		if err != nil {
			return err
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data.Domains)
		}
		for _, d := range resp.Data.Domains {
			fmt.Fprintln(cmd.OutOrStdout(), d)
		}
		return nil
	},
}

func init() {
	fileCmd.AddCommand(fileUploadCmd)
	fileCmd.AddCommand(fileDeleteCmd)
	fileCmd.AddCommand(fileDomainsCmd)

	fileUploadCmd.Flags().StringVarP(&fileUploadOpts.file, "file", "f", "", "Path to file to upload (default stdin if not provided or -)")
	fileUploadCmd.Flags().StringVarP(&fileUploadOpts.name, "name", "n", "", "Filename to use (required for stdin, optional override for file)")
}
