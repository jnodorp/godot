package cmd

import (
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
			// TODO: This could be handled gracefully by ask the user for permission to overwrite.
			log.Fatalf("godot is already initailized (%s exists)", cfgFile)
		}

		// Create a tmp dir to clone the repo.
		tmpDir, err := ioutil.TempDir("", "godot")
		if err != nil {
			log.Fatalf("creating temporary directory failed: %s", err)
		}
		defer func() {
			log.Printf("removing %s", tmpDir)
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
			// TODO: This could be handled gracefully by ask the user for permission to overwrite.
			log.Fatalf("location %s already exists", location)
		}

		// Copy the tmp dir to the location.
		log.Printf("copying %s to %s", tmpDir, location)
		cp, err := exec.Command("cp", "-R", tmpDir, location).CombinedOutput()
		if err != nil {
			log.Fatalf("failed to copy repository to location: %s", cp)
		}

		// Create config dir.
		log.Printf("creating configuration directory %s", cfgDir)
		if err := os.Mkdir(cfgDir, 0755); err != nil {
			log.Fatalf("creating confiuration directory failed: %s", err)
		}

		// Copy the config file to the config dir.
		cfgSrc := path.Join(tmpDir, cfgFileName)
		log.Printf("copying %s to %s", cfgSrc, cfgFile)
		cpCfg, err := exec.Command("cp", cfgSrc, cfgFile).CombinedOutput()
		if err != nil {
			log.Fatalf("copying the configuration file to %s failed: %s", cfgDir, cpCfg)
		}

		log.Print("successfully initialized godot")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
