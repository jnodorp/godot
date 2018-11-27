package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "godot",
	Short: "A dotfile manager written in Go",
	Long:  `Godot is a dotfile manager written in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		templates()
	},
}

// Execute the root command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Initialize configuration.
	cobra.OnInitialize(initConfig)

	// Setup configuration flags.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.godot/godot.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config directory in home directory with name ".godot" (without extension).
		viper.AddConfigPath(path.Join(home, ".godot"))
		viper.SetConfigName("godot")
	}

	// Automatically read configuration from environment variables named 'GODOT_*'.
	viper.SetEnvPrefix("godot")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using configuration file:", viper.ConfigFileUsed())
	}
}
