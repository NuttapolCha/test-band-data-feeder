package cmd

import (
	"github.com/NuttapolCha/test-band-data-feeder/app"
	"github.com/NuttapolCha/test-band-data-feeder/connector"
	"github.com/NuttapolCha/test-band-data-feeder/log"
	"github.com/spf13/cobra"
)

var autoFeederCmd = &cobra.Command{
	Use:   "auto-feeder",
	Short: "feeds coins pricing data from data source to destination service",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := log.NewLogger()
		if err != nil {
			panic(err)
		}
		httpClient := connector.NewCustomHttpClient(logger)
		application := app.New(logger, httpClient)
		return application.StartDataAutomaticFeeder()
	},
}

func init() {
	rootCmd.AddCommand(autoFeederCmd)
}
