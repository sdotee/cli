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
	Use:   "upload",
	Short: "Upload a file",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			filename string
			reader   io.Reader
		)

		// Determine file source
		filePath := fileUploadOpts.file
		// If argument provided, use it as file path
		if len(args) > 0 {
			filePath = args[0]
		}

		if filePath == "" || filePath == "-" {
			// Read from stdin
			if fileUploadOpts.name == "" {
				return fmt.Errorf("filename must be provided via --name when reading from stdin")
			}
			filename = fileUploadOpts.name
			reader = cmd.InOrStdin()
		} else {
			// Read from file
			f, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer f.Close()
			reader = f

			if fileUploadOpts.name != "" {
				filename = fileUploadOpts.name
			} else {
				filename = filepath.Base(filePath)
			}
		}

		resp, err := apiClient.UploadFile(filename, reader)
		if err != nil {
			return err
		}

		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "File uploaded successfully:\n")
		fmt.Fprintf(cmd.OutOrStdout(), "URL: %s\n", resp.Data.URL)
		fmt.Fprintf(cmd.OutOrStdout(), "Delete Key: %s\n", resp.Data.Delete)
		fmt.Fprintf(cmd.OutOrStdout(), "Page: %s\n", resp.Data.Page)
		return nil
	},
}

var fileDeleteCmd = &cobra.Command{
	Use:   "delete <key>",
	Short: "Delete a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deleteKey := args[0]
		resp, err := apiClient.DeleteFile(deleteKey)
		if err != nil {
			return err
		}

		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "File deleted successfully")
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
