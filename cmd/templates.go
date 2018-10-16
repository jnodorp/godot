package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
	"strings"
	"text/template"
)

const templateExtension string = ".tmpl"

func processTemplate(dir string, name string, ctx Context) error {
	tmpl, err := template.New(name).ParseFiles(path.Join(dir, name))
	if err != nil {
		return err
	}

	// Determine target name from template name (e.g. 'test.template' becomes '.test').
	targetName := strings.NewReplacer(templateExtension, "").Replace(name)
	if !strings.HasPrefix(targetName, ".") {
		targetName = "." + targetName
	}

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	targetFile := path.Join(home, targetName)

	// Check if file already exists.
	if _, err := os.Stat(targetFile); !os.IsNotExist(err) {
		if !userConfirm(fmt.Sprintf("'%s' already exists. Overwrite?", targetFile), true) {
			return fmt.Errorf("'%s' already exists", targetFile)
		}
	}

	target, err := os.Create(targetFile)
	if err != nil {
		return err
	}

	err = tmpl.Execute(target, ctx)
	if err != nil {
		return err
	}

	return nil
}
