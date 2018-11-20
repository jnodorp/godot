package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	homedir "github.com/mitchellh/go-homedir"
)

const templateExtension string = ".tmpl"

func processTemplate(dir, name string, ctx Context) error {
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

	return processTemplateTo(dir, name, path.Join(home, targetName), ctx)
}

func processTemplateTo(dir, name, targetFile string, ctx Context) error {
	tmpl, err := template.New(name).ParseFiles(path.Join(dir, name))
	if err != nil {
		return err
	}

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
