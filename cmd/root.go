//
// Copyright (c) 2025 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: root.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2025-12-22 22:23:57
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2025-12-23 11:08:05
//

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	seesdk "github.com/sdotee/sdk.go"
	"github.com/spf13/cobra"
)

var (
	// apiClient is the global SDK client instance used by commands
	apiClient *seesdk.Client

	// rootOpts holds the global command-line options
	rootOpts struct {
		baseURL    string
		apiKey     string
		timeout    time.Duration
		jsonOutput bool
	}
)

var rootCmd = &cobra.Command{
	Use:           "see",
	Short:         "CLI for S.EE Content Share Platform",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if rootOpts.apiKey == "" {
			return errors.New("missing API key: use --api-key or set SEE_API_KEY")
		}
		apiClient = seesdk.NewClient(seesdk.Config{
			BaseURL: rootOpts.baseURL,
			APIKey:  rootOpts.apiKey,
			Timeout: rootOpts.timeout,
		})
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	defaultBaseURL := os.Getenv("SEE_BASE_URL")
	if defaultBaseURL == "" {
		defaultBaseURL = seesdk.DefaultBaseURL
	}
	rootCmd.PersistentFlags().StringVar(&rootOpts.baseURL, "base-url", defaultBaseURL, "API base URL")

	defaultAPIKey := os.Getenv("SEE_API_KEY")
	rootCmd.PersistentFlags().StringVar(&rootOpts.apiKey, "api-key", defaultAPIKey, "API key (or set SEE_API_KEY env)")

	defaultTimeout := seesdk.DefaultTimeout
	if envTimeout := os.Getenv("SEE_TIMEOUT"); envTimeout != "" {
		if d, err := strconv.ParseUint(envTimeout, 10, 64); err == nil {
			defaultTimeout = time.Duration(d) * time.Second
		}
	}
	rootCmd.PersistentFlags().BoolVar(&rootOpts.jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().DurationVar(&rootOpts.timeout, "timeout", defaultTimeout, "HTTP timeout")

	rootCmd.AddCommand(domainsCmd)
	rootCmd.AddCommand(tagsCmd)
	rootCmd.AddCommand(shorturlCmd)
	rootCmd.AddCommand(textCmd)
}
