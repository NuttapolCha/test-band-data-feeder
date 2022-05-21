package cmd

import (
	"context"

	"github.com/NuttapolCha/test-band-data-feeder/app"
	"github.com/NuttapolCha/test-band-data-feeder/log"
	"github.com/spf13/cobra"
)

var autoFeederCmd = &cobra.Command{
	Use:   "auto-feeder",
	Short: "feeds coins pricing data from data source to destination service",
	RunE: func(cmd *cobra.Command, args []string) error {

		application := app.New(log.NewLogger(), context.Background())

		return application.DataAutomaticFeeder()
	},
}

func init() {
	rootCmd.AddCommand(autoFeederCmd)
}
