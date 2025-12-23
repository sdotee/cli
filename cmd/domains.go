//
// Copyright (c) 2025 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: domains.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2025-12-22 22:29:22
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2025-12-23 11:07:58
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "List available domains",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := client.GetDomains()
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
