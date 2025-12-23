//
// Copyright (c) 2025 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: text.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2025-12-22 22:27:43
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2025-12-23 11:08:14
//

package cmd

import (
	"fmt"

	seesdk "github.com/sdotee/sdk.go"
	"github.com/spf13/cobra"
)

var (
	// textCreateOpts holds options for creating a text entry
	textCreateOpts struct {
		domain   string
		slug     string
		title    string
		textType string
		password string
		expireAt int64
		tagIDs   []int64
		file     string
	}

	// textUpdateOpts holds options for updating a text entry
	textUpdateOpts struct {
		domain string
		title  string
		file   string
	}

	// textDeleteOpts holds options for deleting a text entry
	textDeleteOpts struct {
		domain string
	}
)

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Manage text/paste entries",
}

func init() {
	textCmd.AddCommand(textCreateCmd)
	textCmd.AddCommand(textUpdateCmd)
	textCmd.AddCommand(textDeleteCmd)

	textCreateCmd.Flags().StringVar(&textCreateOpts.domain, "domain", "s.ee", "Short domain")
	textCreateCmd.Flags().StringVar(&textCreateOpts.slug, "slug", "", "Custom slug")
	textCreateCmd.Flags().StringVar(&textCreateOpts.title, "title", "", "Title")
	textCreateCmd.Flags().StringVar(&textCreateOpts.textType, "type", "", "Syntax highlighting type")
	textCreateCmd.Flags().StringVar(&textCreateOpts.password, "password", "", "Password")
	textCreateCmd.Flags().Int64Var(&textCreateOpts.expireAt, "expire-at", 0, "Expire at (unix seconds)")
	textCreateCmd.Flags().Int64SliceVar(&textCreateOpts.tagIDs, "tag-ids", nil, "Tag IDs")
	textCreateCmd.Flags().StringVar(&textCreateOpts.file, "file", "-", "Input file path, or '-' for stdin")

	textUpdateCmd.Flags().StringVar(&textUpdateOpts.domain, "domain", "s.ee", "Short domain")
	textUpdateCmd.Flags().StringVar(&textUpdateOpts.title, "title", "", "Title")
	textUpdateCmd.Flags().StringVar(&textUpdateOpts.file, "file", "-", "Input file path, or '-' for stdin")

	textDeleteCmd.Flags().StringVar(&textDeleteOpts.domain, "domain", "s.ee", "Short domain")
}

var textCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a text entry (reads from --file or stdin)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := readContent(textCreateOpts.file, cmd)
		if err != nil {
			return err
		}
		resp, err := apiClient.CreateText(seesdk.CreateTextRequest{
			Content:    content,
			Domain:     textCreateOpts.domain,
			CustomSlug: textCreateOpts.slug,
			Title:      textCreateOpts.title,
			TextType:   textCreateOpts.textType,
			Password:   textCreateOpts.password,
			ExpireAt:   textCreateOpts.expireAt,
			TagIDs:     textCreateOpts.tagIDs,
		})
		if err != nil {
			return err
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		fmt.Fprintln(cmd.OutOrStdout(), resp.Data.ShortURL)
		return nil
	},
}

var textUpdateCmd = &cobra.Command{
	Use:   "update <slug>",
	Short: "Update a text entry (reads from --file or stdin)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := readContent(textUpdateOpts.file, cmd)
		if err != nil {
			return err
		}
		resp, err := apiClient.UpdateText(seesdk.UpdateTextRequest{
			Domain:  textUpdateOpts.domain,
			Slug:    args[0],
			Content: content,
			Title:   textUpdateOpts.title,
		})
		if err != nil {
			return err
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp)
		}
		if resp.Message != "" {
			fmt.Fprintln(cmd.OutOrStdout(), resp.Message)
		}
		return nil
	},
}

var textDeleteCmd = &cobra.Command{
	Use:   "delete <slug>",
	Short: "Delete a text entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.DeleteText(seesdk.DeleteTextRequest{
			Domain: textDeleteOpts.domain,
			Slug:   args[0],
		})
		if err != nil {
			return err
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp)
		}
		if resp.Message != "" {
			fmt.Fprintln(cmd.OutOrStdout(), resp.Message)
		}
		return nil
	},
}
