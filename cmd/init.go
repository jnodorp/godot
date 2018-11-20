package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [git clone url]",
	Short: "Initialize godot from a dotfile repository",
	Long: `The init command initializes godot from a given
dotfile repository. The repository should contain a godot.yaml
file with configuration and multiple tmpl files.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create a temporary directory to clone the repository.
		tmpDir, err := ioutil.TempDir("", "godot")
		if err != nil {
			log.Fatal("Could not create temporary directory.")
		}

		// Clone the git repository provided in the argument to a temporary directory.
		log.Printf("Cloning %s to %s...", args[0], tmpDir)
		git := exec.Command("git", "clone", args[0], tmpDir)
		git.Run()

		// Build context.
		ctx := NewContext()

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Create $HOME/.godot directory.
		os.Mkdir(path.Join(home, ".godot"), 0755)

		// Process the godot configuration template.
		log.Printf("Processing the godot configuration template.")
		err = processTemplateTo(tmpDir, "godot.tmpl", path.Join(home, ".godot", "godot.yaml"), *ctx)
		if err != nil {
			log.Printf("Error processing the godot configuration template: %s", err)
		} else {
			log.Printf("Successfully processed the godot configuration template.")
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func copy(src string, dst string) {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	checkErr(err)
	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0644)
	checkErr(err)
}
