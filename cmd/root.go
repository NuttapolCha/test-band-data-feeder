package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// for custom config file
	configFile string

	rootCmd = &cobra.Command{
		Use:   "data-feeder",
		Short: "this project is only for Band Protocol Interview process.",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ./config/config.yaml)")
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {

		// Search config in home directory with name "config" (without extension).
		viper.AddConfigPath("./config")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
