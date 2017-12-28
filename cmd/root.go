package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ovh/fossil/listener"
	"github.com/ovh/fossil/writer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool
var listen string
var flushDir string

// Fossil init - define command line arguments.
func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file to use")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().StringVarP(&listen, "listen", "l", ":2003", "listen address")
	RootCmd.PersistentFlags().StringVarP(&flushDir, "directory", "d", "./sources", "directory to write metrics file")

	viper.BindPFlags(RootCmd.Flags())
}

// Load config - initialize defaults and read config.
func initConfig() {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	/*if thotToken := os.Getenv("THOT_TOKEN"); thotToken != "" {
		config := thot.NewConfig()
		config.Extra[thot.OVHToken] = thotToken

		hook, err := thot.NewHook(config)
		if err != nil {
			log.Errorf("Failed to setup THOT: %s", err)
		} else {
			log.AddHook(hook)
			log.Info("THOT OK")
		}
	}*/

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
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
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

		//writer := writer.New

		graphite := listener.NewGraphite(listen)
		err := graphite.OpenTCPServer()
		if err != nil {
			panic(err)
		}

		wr := writer.NewWriter(viper.GetString("directory"))
		wr.Write(graphite.Output)

		log.Info("Fossil started")
		select {}
	},
}
