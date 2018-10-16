package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "godot",
	Short: "A dotfile manager written in Go",
	Long:  `Godot is a dotfile manager written in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		location := viper.GetString("location")
		if location == "" {
			log.Fatal("Property 'location' not set.")
		}

		// Expand tilde in the configured location.
		dir, err := homedir.Expand(viper.GetString("location"))
		if err != nil {
			log.Printf("Error expanding location: %s", err)
			dir = viper.GetString("location")
		}

		// Find all files in the configured location.
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Printf("Error reading files from location: %s", err)
		}

		// Build context.
		ctx := NewContext()

		// Process all known files in location.
		log.Printf("Processing files in '%s'.", dir)
		for _, f := range files {
			if strings.HasSuffix(f.Name(), templateExtension) {
				log.Printf("Processing template '%s'.", f.Name())
				err := processTemplate(dir, f.Name(), *ctx)
				if err != nil {
					log.Printf("Error processing template '%s': %s", f.Name(), err)
				} else {
					log.Printf("Successfully processed template '%s'.", f.Name())
				}
			} else {
				continue
			}
		}
	},
}

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
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.godot.yaml)")
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

		// Search config in home directory with name ".godot" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".godot")
	}

	// Automatically read configuration from environment variables named 'GODOT_*'.
	viper.SetEnvPrefix("godot")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using configuration file:", viper.ConfigFileUsed())
	}
}
