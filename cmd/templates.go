package cmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/otiai10/copy"
	"github.com/spf13/viper"
)

func templates(ctx Context) {
	// Read and expand the location property.
	location := viper.GetString("location")
	if location == "" {
		log.Fatal("property 'location' not set")
	} else {
		location = expand(location)
	}

	// Get template targets.
	targets := viper.GetStringMapString("templates.targets")

	// Create a temporary directory to write the templates.
	tmpDir, err := ioutil.TempDir("", "godot")
	if err != nil {
		log.Fatal("Could not create temporary directory.")
	}

	// Process all templates.
	for src, target := range targets {
		log.Printf("Processing template '%s'.", src)

		// Fallback to default target (~/.{{ template }}).
		if target == "" {
			target = path.Join(homeDir(), ".", src)
		} else {
			target = expand(target)
		}

		// Process the template.
		err := processTemplate(location, src, target, tmpDir, ctx)
		if err != nil {
			log.Printf("Error processing template '%s': %s", src, err)
		} else {
			log.Printf("Successfully processed template '%s'.", src)
		}
	}
}

func processTemplate(dir, src, target, tmpDir string, ctx Context) error {
	// Parse the template file.
	tmpl, err := template.New(src).ParseFiles(path.Join(dir, src))
	if err != nil {
		return err
	}

	// Create the temporary target file.
	tmp := path.Join(tmpDir, src)
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}

	// Render the template to the temporary target file.
	err = tmpl.Execute(out, ctx)
	if err != nil {
		return err
	}

	// Check if target file already exists.
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		// Show diff if diff is installed. If diff is empty: skip confirmation.
		diff, err := exec.Command("diff", tmp, target).CombinedOutput()
		if err != nil {
			log.Printf("%s", diff)
		}

		if len(diff) != 0 && !userConfirm(fmt.Sprintf("'%s' already exists. Overwrite?", target), true) {
			return fmt.Errorf("'%s' changed but already exists", target)
		}
	}

	copy.Copy(tmp, target)

	return nil
}
