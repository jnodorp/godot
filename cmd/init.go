package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [git clone url]",
	Short: "Initialize godot from a dotfile repository",
	Long: `The init command initializes godot from a given
dotfile repository. The repository should contain a godot.yaml
file together with all other required files.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfgDir := expand("~/.godot")
		cfgFileName := "godot.yaml"
		cfgFile := path.Join(cfgDir, cfgFileName)

		// Check if godot is already initialized.
		if _, err := os.Stat(cfgFile); !os.IsNotExist(err) {
			if !userConfirm(fmt.Sprintf("godot is already initailized (%s exists). Continue?", cfgFile), false) {
				os.Exit(0)
			}
		}

		// Create the OS tmp dir if it does not exist.
		if _, err := os.Stat(os.TempDir()); os.IsNotExist(err) {
			err := os.MkdirAll(os.TempDir(), os.ModeDir)
			if err != nil {
				log.Fatalf("failed to create missing OS temporary directory %s", os.TempDir())
			}
		}

		// Create a tmp dir to clone the repo.
		tmpDir, err := ioutil.TempDir("", "godot")
		if err != nil {
			log.Fatalf("creating temporary directory failed: %s", err)
		}
		defer func() {
			if err = os.RemoveAll(tmpDir); err != nil {
				log.Printf("removing temporary directory failed: %s", err)
			}
		}()

		// TODO: Support private key auth.
		var gitAuth transport.AuthMethod
		if viper.GetString("username") != "" || viper.GetString("password") != "" {
			gitAuth = &http.BasicAuth{
				Username: viper.GetString("username"),
				Password: viper.GetString("password"),
			}
		} else {
			gitAuth = nil
		}

		// Clone the git repo (provided as arg) to a tmp dir.
		log.Printf("cloning %s to %s", args[0], tmpDir)
		_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
			URL:      args[0],
			Auth:     gitAuth,
			Progress: os.Stdout,
		})
		if err != nil {
			log.Fatalf("cloning repository failed: %s", err)
		}

		// Utilize the config from the repo cloned to the tmp dir.
		viper.AddConfigPath(tmpDir)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("reading 'godot.yaml' from the repository failed: %s", err)
		}

		// Read and expand the location property.
		location := viper.GetString("location")
		if location == "" {
			log.Fatal("property 'location' not set")
		} else {
			location = expand(location)
		}

		// Check if location path exists.
		if _, err := os.Stat(location); !os.IsNotExist(err) {
			if !userConfirm(fmt.Sprintf("location %s already exists. Remove?", location), false) {
				os.Exit(0)
			}

			// Remove existing directory.
			err := os.RemoveAll(location)
			if err != nil {
				log.Fatalf("failed to remove %s", location)
			}
		}

		// Copy the tmp dir to the location.
		log.Printf("copying %s to %s", tmpDir, location)
		err = copy.Copy(tmpDir, location)
		if err != nil {
			log.Fatalf("failed to copy repository to location: %s", err)
		}

		// Create config dir if it does not exist yet.
		if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
			log.Printf("creating configuration directory %s", cfgDir)
			if err := os.Mkdir(cfgDir, os.ModeDir); err != nil {
				log.Fatalf("creating confiuration directory failed: %s", err)
			}
		}

		// Copy the config file to the config dir.
		cfgSrc := path.Join(tmpDir, cfgFileName)
		log.Printf("copying %s to %s", cfgSrc, cfgFile)
		err = copy.Copy(cfgSrc, cfgFile)
		if err != nil {
			log.Fatalf("copying the configuration file to %s failed: %s", cfgDir, err)
		}

		fmt.Println("successfully initialized godot")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Setup configuration flags.
	initCmd.Flags().StringP("username", "u", "", "the git username to use")
	viper.BindPFlag("username", initCmd.Flags().Lookup("username"))

	initCmd.Flags().StringP("password", "p", "", "the git password to use")
	viper.BindPFlag("password", initCmd.Flags().Lookup("password"))
}
