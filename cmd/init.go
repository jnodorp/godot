package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
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

		// Clone the git repo (provided as arg) to a tmp dir.
		log.Printf("cloning %s to %s", args[0], tmpDir)
		git, err := exec.Command("git", "clone", args[0], tmpDir).CombinedOutput()
		if err != nil {
			log.Fatalf("cloning repository failed: %s", git)
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
		cp, err := exec.Command("cp", "-R", tmpDir, location).CombinedOutput()
		if err != nil {
			log.Fatalf("failed to copy repository to location: %s", cp)
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
		err = copyFile(cfgSrc, cfgFile)
		if err != nil {
			log.Fatalf("copying the configuration file to %s failed: %s", cfgDir, err)
		}

		fmt.Println("successfully initialized godot")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
