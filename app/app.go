package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/gotrix/gotrix"
	"gopkg.in/reform.v1"
)

func New(cnf *gotrix.AppConfig) (*App, error) {
	a := &App{
		cnf:         cnf,
		components:  make(map[string]gotrix.ComponentWrapper, 0),
		pages:       make(map[string]*template.Template, 0),
		paths:       make([]*path, 0),
		staticPaths: make(map[string]http.Handler, len(cnf.StaticPaths)),
	}
	if err := a.init(); err != nil {
		return nil, err
	}
	return a, nil
}

type App struct {
	cnf *gotrix.AppConfig

	// components
	componentPaths []string
	components     map[string]gotrix.ComponentWrapper

	// pages
	pagesPaths []string
	pages      map[string]*template.Template
	layout     *template.Template

	// paths
	paths []*path

	// static
	staticPaths map[string]http.Handler
}

func (app *App) init() error {
	if err := app.loadComponents(); err != nil {
		return err
	}
	if err := app.loadTemplates(); err != nil {
		return err
	}
	if err := app.buildStatic(); err != nil {
		return err
	}
	if err := app.buildPaths(); err != nil {
		return err
	}
	return nil
}

func (app *App) buildStatic() error {
	list := make([]string, 0)
	copy(list, app.cnf.StaticPaths)
	if len(list) == 0 {
		app.cnf.StaticPaths = []string{
			"static:./static",
		}
	}
	list = append(list, "gotrix:./.gotrix/static")
	for _, p := range list {
		parts := strings.Split(p, ":")
		if len(parts) != 2 {
			log.Printf("warn: static path '%s' is invalid, "+
				"please use 'url_path:localpath' notation, ignoring\n", p)
			continue
		}
		path, err := filepath.Abs(parts[1])
		if err != nil {
			return err
		}
		if dirExists(path) {
			sp := parts[0]
			app.staticPaths["/"+parts[0]] = http.StripPrefix("/"+sp+"/", http.FileServer(http.Dir(path)))
			log.Printf("registered static path /%s to serve from %s",
				parts[0], path)
		} else {
			log.Printf("filepath %s does not exist\n", path)
		}
	}
	return nil
}

func (app *App) loadComponents() error {
	log.Println("loading components")
	app.componentPaths = make([]string, 0, len(app.cnf.ComponentPaths)+2)
	if app.cnf.ComponentPaths != nil {
		app.componentPaths = append(app.componentPaths, app.cnf.ComponentPaths...)
	}
	localComponents, err := filepath.Abs("./components")
	if err != nil {
		return err
	}
	app.componentPaths = append(app.componentPaths,
		"/usr/local/gotrix/components",
		localComponents,
	)
	getName := func(d, f string) string {
		n := strings.TrimPrefix(f, d)
		n = strings.TrimLeft(n, "/")
		l := strings.Split(n, "/")
		return strings.Join(l[0:len(l)-1], "/")
	}
	for _, d := range app.componentPaths {
		files, err := readPathsRecursive(d, ".so")
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("failed to list %s: %v", d, err)
		}
		for _, f := range files {
			name := getName(d, f)
			if _, ok := app.components[name]; ok {
				log.Printf("ignoring %s, component %s already loaded", f, name)
			}
			log.Println(f)
			plug, err := plugin.Open(f)
			if err != nil {
				return err
			}
			com, err := plug.Lookup("Component")
			if err != nil {
				return err
			}
			if c, ok := com.(gotrix.ComponentWrapper); ok {
				app.components[name] = c
				log.Printf("loaded components %s from %s\n", name, f)
			} else {
				return fmt.Errorf("%s is not a component at %s", name, f)
			}
		}
	}
	return nil
}

func (app *App) loadTemplates() error {
	log.Println("loading pages")
	app.pagesPaths = make([]string, 0, len(app.cnf.TemplatePaths)+2)
	if app.cnf.TemplatePaths != nil {
		app.pagesPaths = append(app.pagesPaths, app.cnf.TemplatePaths...)
	}
	localTemplates, err := filepath.Abs("./pages")
	if err != nil {
		return err
	}
	app.pagesPaths = append(app.pagesPaths,
		localTemplates,
	)
	getName := func(d, f string) string {
		n := strings.TrimPrefix(f, d)
		n = strings.TrimLeft(n, "/")
		n = strings.TrimSuffix(n, ".html")
		return n
	}
	for _, d := range app.pagesPaths {
		files, err := readPathsRecursive(d, ".html")
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("failed to list %s: %v", d, err)
		}
		for _, f := range files {
			name := getName(d, f)
			if _, ok := app.pages[name]; ok {
				log.Printf("ignoring %s, template %s already loaded\n", f, name)
				continue
			}
			contents, err := ioutil.ReadFile(f)
			if err != nil {
				return err
			}
			tpl, err := template.
				New(name).
				Parse(string(contents))
			if err != nil {
				return err
			}
			if name == "layout" {
				app.layout = tpl
				log.Printf("loaded layout from %s\n", f)
			} else {
				app.pages[name] = tpl
				log.Printf("loaded template %s from %s\n", name, f)
			}
		}
	}
	return nil
}

func (app *App) buildPaths() error {
	prepare := func(v string) (string, bool) {
		isSlug := false
		if strings.HasSuffix(v, "[slug]") {
			v = strings.Replace(v, "[slug]", "([a-zA-Z0-9-]+)", 1)
			isSlug = true
		}
		if strings.HasSuffix(v, "index") {
			v = strings.Replace(v, "index", "", 1)
		}
		return "^/" + strings.TrimRight(v, "/") + "$", isSlug
	}
	col := func(v *regexp.Regexp) error {
		s := strings.TrimLeft(v.String(), "^")
		s = strings.TrimRight(s, "$")
		for _, p := range app.paths {
			if p.route.MatchString(s) {
				return fmt.Errorf("route %s collides with %s",
					v.String(), p.route.String())
			}
		}
		return nil
	}
	for n, t := range app.pages {
		r, isSlug := prepare(n)
		route, err := regexp.Compile(r)
		if err != nil {
			return err
		}
		if err := col(route); err != nil {
			return fmt.Errorf("route collition detected: %s", err)
		}
		app.paths = append(app.paths, &path{
			route:    route,
			template: t,
			isSlug:   isSlug,
		})
	}
	for _, p := range app.paths {
		log.Println(p.route.String())
	}
	sort.Slice(app.paths, func(i, j int) bool {
		return app.paths[i].route.String() < app.paths[j].route.String()
	})
	return nil
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v\n", r)
		}
	}()
	// redirect tailing slash
	if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
		w.Header().Set("location", r.URL.Path[0:len(r.URL.Path)-1])
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}
	// serve static
	for p, f := range app.staticPaths {
		if strings.HasPrefix(r.URL.Path, p) {
			log.Println("static", r.URL.Path)
			f.ServeHTTP(w, r)
			return
		}
	}
	// walk through paths
	var (
		body    = bytes.NewBuffer([]byte{})
		data    = gotrix.NewPageData(body, r, app)
		matched = false
	)
	for _, p := range app.paths {
		if p.route.MatchString(r.URL.Path) {
			//rd := &renderData{
			//	Path: r.URL.Path,
			//	Data: data,
			//}
			if p.isSlug {
				m := p.route.FindStringSubmatch(r.URL.Path)
				if len(m) < 2 {
					continue
				}
				data.SetSlug(m[1])
			}
			if err := p.template.Execute(body, data); err != nil {
				log.Println(err)
				return
			}
			matched = true
			break
		}
	}
	if !matched {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	data.SetBody(body.String())
	//rd := &renderData{
	//	Body: body.String(),
	//	Data: data,
	//	//ComponentWrapper: func(s ...string) string {
	//	//	return `<div>` + strings.Join(s, ";") + `</div>`
	//	//},
	//}
	if err := app.layout.Execute(w, data); err != nil {
		log.Printf("failed to render layout: %s\n", err)
	}
}

func (app *App) Run() error {
	port := app.cnf.Port
	if port < 1 {
		port = 8080
	}
	return http.ListenAndServe(
		fmt.Sprintf(":%d", port), app)
}

func (app *App) DB() *reform.DB {
	return nil
}

func (app *App) Components() map[string]gotrix.ComponentWrapper {
	return app.components
}
