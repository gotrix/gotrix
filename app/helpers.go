package app

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/gotrix/gotrix"
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

type pageData struct {
	Path string
	Body string
	Slug string
	Data *gotrix.PageData
}

func (pd *pageData) EndHead() string {
	eb := "<!-- gotrix end head -->\n"
	if ac := pd.Data.AsyncCSS(); ac != `[""]` {
		eb += `<script>((d,l,h,e)=>{` +
			`h=d.getElementsByTagName("head")[0];` +
			`l.forEach(s=>{` +
			`e=d.createElement("link");` +
			`e.href=s;` +
			`e.type="text/css";` +
			`e.rel="stylesheet";` +
			`h.appendChild(e);` +
			`})` +
			`})(document,` + ac + `);</script>`
	}
	eb += pd.Data.CSS()
	return eb
}

func (pd *pageData) EndBody() string {
	eb := "<!-- gotrix end body -->\n"

	eb += pd.Data.JS()
	return eb
}

type path struct {
	route    *regexp.Regexp
	template *template.Template
	isSlug   bool
}
