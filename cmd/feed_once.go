package cmd

import (
	"github.com/NuttapolCha/test-band-data-feeder/app"
	"github.com/NuttapolCha/test-band-data-feeder/connector"
	"github.com/NuttapolCha/test-band-data-feeder/log"
	"github.com/spf13/cobra"
)

var feedOne = &cobra.Command{
	Use:   "feed-once",
	Short: "feeds coins pricing data from data source to destination service only once",
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := log.NewLogger()
		if err != nil {
			panic(err)
		}
		httpClient := connector.NewCustomHttpClient(logger)
		application := app.New(logger, httpClient)
		application.Feed()
	},
}

func init() {
	rootCmd.AddCommand(feedOne)
}
