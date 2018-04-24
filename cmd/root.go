package cmd

import (
	"github.com/ovh/fossil/listener"
	"github.com/ovh/fossil/writer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Fossil init - define command line arguments.
func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().String("config", "", "config file to use")
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	RootCmd.Flags().StringP("listen", "l", ":2003", "listen address")
	RootCmd.Flags().IntP("batch", "b", 10000, "batch count per file")
	RootCmd.Flags().IntP("timeout", "t", 5, "batch timeout for flushing datapoints")
	RootCmd.Flags().StringP("directory", "d", "./sources", "directory to write metrics file")
	RootCmd.Flags().BoolP("parse", "p", true, "parse metric name to auto fill tags")

	viper.BindPFlags(RootCmd.Flags())
	viper.BindPFlags(RootCmd.PersistentFlags())
}

// Load config - initialize defaults and read config.
func initConfig() {
	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}

	// Bind environment variables
	viper.SetEnvPrefix("fossil")
	viper.AutomaticEnv()

	// Set config search path
	viper.AddConfigPath("/etc/fossil/")
	viper.AddConfigPath("$HOME/.fossil")
	viper.AddConfigPath(".")

	// Load config
	viper.SetConfigName("config")
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("No config file found")
		} else {
			log.Panicf("Fatal error in config file: %v \n", err)
		}
	}

	// Load user defined config
	cfg := viper.GetString("config")
	if cfg != "" {
		viper.SetConfigFile(cfg)
		viper.SetConfigType("json")
		err := viper.ReadInConfig()
		if err != nil {
			log.Panicf("Fatal error in config file: %v \n", err)
		}
	}
}

// RootCmd launch the aggregator agent.
var RootCmd = &cobra.Command{
	Use:   "fossil",
	Short: "Fossil fossil Graphite to beamium forwarder",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Fossil starting")

		wr := writer.NewWriter(viper.GetString("directory"))

		graphite := listener.NewGraphite(viper.GetString("listen"), wr, viper.GetBool("parse"))
		err := graphite.OpenTCPServer()
		if err != nil {
			panic(err)
		}
	},
}
