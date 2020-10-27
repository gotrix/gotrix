package app

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func readPathsRecursive(dir string, suffix string) (list []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), suffix) {
			list = append(list, path)
		}
		return nil
	})
	return
}

func dirExists(dir string) bool {
	d, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return d.IsDir()
}

type path struct {
	route    *regexp.Regexp
	template *template.Template
	isSlug   bool
}
