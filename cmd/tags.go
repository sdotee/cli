//
// Copyright (c) 2025 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: tags.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2025-12-22 22:29:25
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2025-12-23 11:08:11
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List available tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.GetTags()
		if err != nil {
			return err
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data.Tags)
		}
		for _, t := range resp.Data.Tags {
			fmt.Fprintf(cmd.OutOrStdout(), "%d\t%s\n", t.ID, t.Name)
		}
		return nil
	},
}
