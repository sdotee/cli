package cmd

import (
	"fmt"
	"strconv"
	"strings"

	seesdk "github.com/sdotee/sdk.go"
	"github.com/spf13/cobra"
)

var (
	bioCreateOpts struct {
		domain      string
		slug        string
		description string
		mastodonURL string
		rssURL      string
		links       []string
	}
	bioUpdateOpts struct {
		title       string
		description string
		mastodonURL string
		rssURL      string
		links       []string
	}
	bioHistoryPage int
)

var bioCmd = &cobra.Command{Use: "bio", Short: "Manage bio pages"}

var bioCreateCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a bio page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		links, err := parseBioLinks(bioCreateOpts.links)
		if err != nil {
			return err
		}
		resp, err := apiClient.CreateBioPage(seesdk.CreateBioPageRequest{
			Title: args[0], Description: bioCreateOpts.description,
			Domain: bioCreateOpts.domain, CustomSlug: bioCreateOpts.slug,
			MastodonURL: bioCreateOpts.mastodonURL, RSSURL: bioCreateOpts.rssURL,
			CustomLinks: links,
		})
		if err != nil {
			return fmt.Errorf("failed to create bio page: %w", err)
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		fmt.Fprintln(cmd.OutOrStdout(), resp.Data.ShortURL)
		return nil
	},
}

var bioUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a bio page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := parsePositiveID(args[0], "bio page")
		if err != nil {
			return err
		}
		if strings.TrimSpace(bioUpdateOpts.title) == "" {
			return fmt.Errorf("--title is required")
		}
		links, err := parseBioLinks(bioUpdateOpts.links)
		if err != nil {
			return err
		}
		resp, err := apiClient.UpdateBioPage(seesdk.UpdateBioPageRequest{
			ID: id, Title: bioUpdateOpts.title, Description: bioUpdateOpts.description,
			MastodonURL: bioUpdateOpts.mastodonURL, RSSURL: bioUpdateOpts.rssURL,
			CustomLinks: links,
		})
		if err != nil {
			return fmt.Errorf("failed to update bio page: %w", err)
		}
		return printMessage(cmd, resp, resp.Message)
	},
}

var bioDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a bio page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := parsePositiveID(args[0], "bio page")
		if err != nil {
			return err
		}
		resp, err := apiClient.DeleteBioPage(id)
		if err != nil {
			return fmt.Errorf("failed to delete bio page: %w", err)
		}
		return printMessage(cmd, resp, resp.Message)
	},
}

var bioHistoryCmd = &cobra.Command{
	Use: "history", Short: "List bio pages", Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.GetBioPageHistory(bioHistoryPage)
		if err != nil {
			return fmt.Errorf("failed to get bio page history: %w", err)
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		for _, bio := range resp.Data.BioPages {
			fmt.Fprintf(cmd.OutOrStdout(), "%d\t%s\t%s\n", bio.ID, bio.Link, bio.Title)
		}
		return nil
	},
}

func parseBioLinks(values []string) ([]seesdk.BioCustomLink, error) {
	links := make([]seesdk.BioCustomLink, 0, len(values))
	for _, value := range values {
		title, url, ok := strings.Cut(value, "=")
		if !ok || strings.TrimSpace(title) == "" || strings.TrimSpace(url) == "" {
			return nil, fmt.Errorf("invalid link %q: expected title=url", value)
		}
		links = append(links, seesdk.BioCustomLink{Title: title, URL: url})
	}
	return links, nil
}

func parsePositiveID(value, resource string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid %s ID %q", resource, value)
	}
	return id, nil
}

func init() {
	bioCmd.AddCommand(bioCreateCmd, bioUpdateCmd, bioDeleteCmd, bioHistoryCmd)
	bioCreateCmd.Flags().StringVar(&bioCreateOpts.domain, "domain", "", "Domain for the bio page URL")
	bioCreateCmd.Flags().StringVar(&bioCreateOpts.slug, "slug", "", "Custom slug")
	bioCreateCmd.Flags().StringVar(&bioCreateOpts.description, "description", "", "Description")
	bioCreateCmd.Flags().StringVar(&bioCreateOpts.mastodonURL, "mastodon-url", "", "Mastodon profile URL")
	bioCreateCmd.Flags().StringVar(&bioCreateOpts.rssURL, "rss-url", "", "RSS feed URL")
	bioCreateCmd.Flags().StringArrayVar(&bioCreateOpts.links, "link", nil, "Custom link as title=url (repeatable)")
	bioUpdateCmd.Flags().StringVar(&bioUpdateOpts.title, "title", "", "Title (required)")
	bioUpdateCmd.Flags().StringVar(&bioUpdateOpts.description, "description", "", "Description")
	bioUpdateCmd.Flags().StringVar(&bioUpdateOpts.mastodonURL, "mastodon-url", "", "Mastodon profile URL")
	bioUpdateCmd.Flags().StringVar(&bioUpdateOpts.rssURL, "rss-url", "", "RSS feed URL")
	bioUpdateCmd.Flags().StringArrayVar(&bioUpdateOpts.links, "link", nil, "Replacement custom link as title=url (repeatable)")
	addPageFlag(bioHistoryCmd, &bioHistoryPage)
}
