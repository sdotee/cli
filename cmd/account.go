package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Show account usage",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.GetUsage()
		if err != nil {
			return fmt.Errorf("failed to get account usage: %w", err)
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		data := resp.Data
		fmt.Fprintf(cmd.OutOrStdout(), "API requests (day): %d/%d\n", data.APICountDay, data.APICountDayLimit)
		fmt.Fprintf(cmd.OutOrStdout(), "Links (day): %d/%d\n", data.LinkCountDay, data.LinkCountDayLimit)
		fmt.Fprintf(cmd.OutOrStdout(), "Texts (day): %d/%d\n", data.TextCountDay, data.TextCountDayLimit)
		fmt.Fprintf(cmd.OutOrStdout(), "Uploads (day): %d/%d\n", data.UploadCountDay, data.UploadCountDayLimit)
		fmt.Fprintf(cmd.OutOrStdout(), "Files: %d\n", data.FileCount)
		fmt.Fprintf(cmd.OutOrStdout(), "Storage (MB): %s/%s\n", data.StorageUsageMB, data.StorageUsageLimitMB)
		return nil
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage API tokens",
}

var tokenCheckCmd = &cobra.Command{
	Use:   "check [token]",
	Short: "Check whether an API token is valid",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := rootOpts.apiKey
		if len(args) == 1 {
			token = args[0]
		}
		resp, err := apiClient.CheckToken(token)
		if err != nil {
			return fmt.Errorf("failed to check API token: %w", err)
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Valid: %t\n", resp.Data.Valid)
		fmt.Fprintf(cmd.OutOrStdout(), "Expires At: %d\n", resp.Data.ExpiresAt)
		return nil
	},
}

func init() {
	tokenCmd.AddCommand(tokenCheckCmd)
}
