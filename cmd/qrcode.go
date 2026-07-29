package cmd

import (
	"fmt"
	"strings"

	seesdk "github.com/sdotee/sdk.go"
	"github.com/spf13/cobra"
)

var (
	qrcodeCreateOpts   struct{ title, domain, slug string }
	qrcodeDeleteDomain string
	qrcodeHistoryPage  int
)

var qrcodeCmd = &cobra.Command{Use: "qrcode", Short: "Manage dynamic QR codes"}

var qrcodeCreateCmd = &cobra.Command{
	Use: "create <target-url>", Short: "Create a dynamic QR code", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(qrcodeCreateOpts.title) == "" {
			return fmt.Errorf("--title is required")
		}
		resp, err := apiClient.CreateQRCode(seesdk.CreateQRCodeRequest{
			TargetURL: args[0], Title: qrcodeCreateOpts.title,
			Domain: qrcodeCreateOpts.domain, CustomSlug: qrcodeCreateOpts.slug,
		})
		if err != nil {
			return fmt.Errorf("failed to create QR code: %w", err)
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Short URL: %s\nPNG: %s\nSVG: %s\nPDF: %s\n", resp.Data.ShortURL, resp.Data.PNGURL, resp.Data.SVGURL, resp.Data.PDFURL)
		return nil
	},
}

var qrcodeDeleteCmd = &cobra.Command{
	Use: "delete <slug>", Short: "Delete a QR code", Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.DeleteQRCode(seesdk.DeleteQRCodeRequest{Domain: qrcodeDeleteDomain, Slug: args[0]})
		if err != nil {
			return fmt.Errorf("failed to delete QR code: %w", err)
		}
		return printMessage(cmd, resp, resp.Message)
	},
}

var qrcodeHistoryCmd = &cobra.Command{
	Use: "history", Short: "List QR codes", Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := apiClient.GetQRCodeHistory(qrcodeHistoryPage)
		if err != nil {
			return fmt.Errorf("failed to get QR code history: %w", err)
		}
		if rootOpts.jsonOutput {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		for _, qr := range resp.Data.QRCodes {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%d\t%s\n", qr.ShortURL, qr.ScanCount, qr.Title)
		}
		return nil
	},
}

func init() {
	qrcodeCmd.AddCommand(qrcodeCreateCmd, qrcodeDeleteCmd, qrcodeHistoryCmd)
	qrcodeCreateCmd.Flags().StringVar(&qrcodeCreateOpts.title, "title", "", "QR code title (required)")
	qrcodeCreateCmd.Flags().StringVar(&qrcodeCreateOpts.domain, "domain", "", "Domain for the short URL")
	qrcodeCreateCmd.Flags().StringVar(&qrcodeCreateOpts.slug, "slug", "", "Custom slug")
	qrcodeDeleteCmd.Flags().StringVar(&qrcodeDeleteDomain, "domain", "s.ee", "Short domain")
	addPageFlag(qrcodeHistoryCmd, &qrcodeHistoryPage)
}
