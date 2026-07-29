//
// Copyright (c) 2025 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: shorturl.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2025-12-22 22:25:46
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2025-12-23 11:08:08
//

package cmd

import (
	"fmt"
	"strings"

	seesdk "github.com/sdotee/sdk.go"
	"github.com/spf13/cobra"
)

var (
	// shortCreateOpts holds options for creating a short URL
	shortCreateOpts struct {
		domain                string
		slug                  string
		title                 string
		password              string
		expireAt              int64
		tagIDs                []int64
		expirationRedirectURL string
	}

	// shortUpdateOpts holds options for updating a short URL
	shortUpdateOpts struct {
		domain    string
		targetURL string
		title     string
	}

	// shortDeleteOpts holds options for deleting a short URL
	shortDeleteOpts struct {
		domain string
	}

	shortStatOpts struct {
		domain string
		period string
	}

	shortHistoryPage int
)

var shorturlCmd = &cobra.Command{
	Use:   "shorturl",
	Short: "Manage short URLs",
}

func init() {
	shorturlCmd.AddCommand(shorturlCreateCmd)
	shorturlCmd.AddCommand(shorturlUpdateCmd)
	shorturlCmd.AddCommand(shorturlDeleteCmd)
	shorturlCmd.AddCommand(shorturlStatCmd)
	shorturlCmd.AddCommand(shorturlHistoryCmd)

	shorturlCreateCmd.Flags().StringVar(&shortCreateOpts.domain, "domain", "s.ee", "Short domain")
	shorturlCreateCmd.Flags().StringVar(&shortCreateOpts.slug, "slug", "", "Custom slug")
	shorturlCreateCmd.Flags().StringVar(&shortCreateOpts.title, "title", "", "Title")
	shorturlCreateCmd.Flags().StringVar(&shortCreateOpts.password, "password", "", "Password")
	shorturlCreateCmd.Flags().Int64Var(&shortCreateOpts.expireAt, "expire-at", 0, "Expire at (unix seconds)")
	shorturlCreateCmd.Flags().Int64SliceVar(&shortCreateOpts.tagIDs, "tag-ids", nil, "Tag IDs")
	shorturlCreateCmd.Flags().StringVar(&shortCreateOpts.expirationRedirectURL, "expiration-redirect-url", "", "Redirect URL after expiration")

	shorturlUpdateCmd.Flags().StringVar(&shortUpdateOpts.domain, "domain", "s.ee", "Short domain")
	shorturlUpdateCmd.Flags().StringVar(&shortUpdateOpts.targetURL, "target-url", "", "New target URL")
	shorturlUpdateCmd.Flags().StringVar(&shortUpdateOpts.title, "title", "", "Title")

	shorturlDeleteCmd.Flags().StringVar(&shortDeleteOpts.domain, "domain", "s.ee", "Short domain")

	shorturlStatCmd.Flags().StringVar(&shortStatOpts.domain, "domain", "s.ee", "Short domain")
	shorturlStatCmd.Flags().StringVar(&shortStatOpts.period, "period", seesdk.VisitStatPeriodTotally, "Statistics period: daily, monthly, or totally")
	addPageFlag(shorturlHistoryCmd, &shortHistoryPage)
}

var shorturlCreateCmd = &cobra.Command{
	Use:   "create <target-url>",
	Short: "Create a short URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		req := seesdk.CreateShortURLRequest{
			TargetURL:             args[0],
			Domain:                shortCreateOpts.domain,
			CustomSlug:            shortCreateOpts.slug,
			Title:                 shortCreateOpts.title,
			Password:              shortCreateOpts.password,
			ExpireAt:              shortCreateOpts.expireAt,
			TagIDs:                shortCreateOpts.tagIDs,
			ExpirationRedirectURL: shortCreateOpts.expirationRedirectURL,
		}

		resp, err := apiClient.CreateShortURL(req)
		if err != nil {
			return fmt.Errorf("failed to create short URL: %w", err)
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		fmt.Fprintln(cmd.OutOrStdout(), resp.Data.ShortURL)
		return nil
	},
}

var shorturlUpdateCmd = &cobra.Command{
	Use:   "update <slug>",
	Short: "Update an existing short URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(shortUpdateOpts.targetURL) == "" {
			return fmt.Errorf("--target-url is required")
		}
		resp, err := apiClient.UpdateShortURL(seesdk.UpdateShortURLRequest{
			Domain:    shortUpdateOpts.domain,
			Slug:      args[0],
			TargetURL: shortUpdateOpts.targetURL,
			Title:     shortUpdateOpts.title,
		})
		if err != nil {
			return fmt.Errorf("failed to update short URL: %w", err)
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

var shorturlDeleteCmd = &cobra.Command{
	Use:   "delete <slug>",
	Short: "Delete a short URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.DeleteShortURL(seesdk.DeleteURLRequest{
			Domain: shortDeleteOpts.domain,
			Slug:   args[0],
		})
		if err != nil {
			return fmt.Errorf("failed to delete short URL: %w", err)
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

var shorturlStatCmd = &cobra.Command{
	Use:   "stats <slug>",
	Short: "Get visit statistics for a short URL",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		switch shortStatOpts.period {
		case seesdk.VisitStatPeriodDaily, seesdk.VisitStatPeriodMonthly, seesdk.VisitStatPeriodTotally:
			return nil
		default:
			return fmt.Errorf("invalid period %q: use daily, monthly, or totally", shortStatOpts.period)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.GetLinkVisitStat(shortStatOpts.domain, args[0], shortStatOpts.period)
		if err != nil {
			return fmt.Errorf("failed to get short URL statistics: %w", err)
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		fmt.Fprintln(cmd.OutOrStdout(), resp.Data.VisitCount)
		return nil
	},
}

var shorturlHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "List short URL history",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.GetLinkHistory(shortHistoryPage)
		if err != nil {
			return fmt.Errorf("failed to get short URL history: %w", err)
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		for _, link := range resp.Data {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%d\t%s\n", link.ShortURL, link.TargetURL, link.VisitCount, link.Title)
		}
		return nil
	},
}
