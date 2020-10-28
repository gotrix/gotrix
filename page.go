package gotrix

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"text/template"
)

type Page interface {
	Page(Context) error
	Template() *template.Template
}

func NewPageData(w io.Writer, r *http.Request, app App) *PageData {
	return &PageData{
		w:        w,
		r:        r,
		app:      app,
		data:     make(map[string]interface{}),
		css:      make([]string, 0),
		asyncCss: make([]string, 0),
		js:       make([]string, 0),
		asyncJs:  make([]string, 0),
		deferJs:  make([]string, 0),
	}
}

type PageData struct {
	w          io.Writer
	r          *http.Request
	data       map[string]interface{}
	app        App
	slug       string
	body       string
	cssMu      sync.Mutex
	css        []string
	asyncCssMu sync.Mutex
	asyncCss   []string
	jsMu       sync.Mutex
	js         []string
	asyncJsMu  sync.Mutex
	asyncJs    []string
	deferJsMu  sync.Mutex
	deferJs    []string
}

func (pd *PageData) Extend() Context {
	return &PageData{
		w:        pd.Writer(),
		r:        pd.Request(),
		app:      pd.app,
		css:      make([]string, 0),
		asyncCss: make([]string, 0),
		js:       make([]string, 0),
		asyncJs:  make([]string, 0),
		deferJs:  make([]string, 0),
	}
}

func (pd *PageData) Data(k string) interface{} {
	if v, ok := pd.data[k]; ok {
		return v
	}
	return nil
}

func (pd *PageData) SetData(k string, v interface{}) {
	pd.data[k] = v
}

func (pd *PageData) Body() string {
	return pd.body
}

func (pd *PageData) SetBody(body string) {
	pd.body = body
}

func (pd *PageData) Slug() string {
	return pd.slug
}

func (pd *PageData) SetSlug(slug string) {
	pd.slug = slug
}

func (pd *PageData) Request() *http.Request {
	return pd.r
}

func (pd *PageData) Writer() io.Writer {
	return pd.w
}

func (pd *PageData) AddAsyncCSS(css ...string) {
	pd.asyncCssMu.Lock()
	pd.asyncCss = append(pd.asyncCss, css...)
	pd.asyncCssMu.Unlock()
}

func (pd *PageData) AddCSS(css ...string) {
	pd.cssMu.Lock()
	pd.css = append(pd.css, css...)
	pd.cssMu.Unlock()
}

func (pd *PageData) AddAsyncJS(js ...string) {
	pd.asyncJsMu.Lock()
	pd.asyncJs = append(pd.asyncJs, js...)
	pd.asyncJsMu.Unlock()
}

func (pd *PageData) AddDeferJS(js ...string) {
	pd.deferJsMu.Lock()
	pd.deferJs = append(pd.deferJs, js...)
	pd.deferJsMu.Unlock()
}

func (pd *PageData) AddJS(js ...string) {
	pd.jsMu.Lock()
	pd.js = append(pd.js, js...)
	pd.jsMu.Unlock()
}

func (pd *PageData) JS() string {
	l := make([]string, 0)
	l = append(l, unique(pd.asyncJs, func(s string) string {
		return toJsTag(s, "async")
	})...)
	l = append(l, unique(pd.deferJs, func(s string) string {
		return toJsTag(s, "defer")
	})...)
	l = append(l, unique(pd.js, func(s string) string {
		return toJsTag(s, "")
	})...)
	return strings.Join(l, "")
}

func (pd *PageData) CSS() string {
	return strings.Join(unique(pd.css, func(s string) string {
		return toCssTag(s)
	}), "")
}

func (pd *PageData) AsyncCSS() string {
	return `["` + strings.Join(unique(pd.asyncCss, nil), `","`) + `"]`
}

func (pd *PageData) EndHead() string {
	eb := "<!-- endhead -->"
	if ac := pd.AsyncCSS(); ac != `[""]` {
		eb += "<script>((g,o,t,r)=>{" +
			"t=g.getElementsByTagName('head')[0];" +
			"o.map(s=>{" +
			"r=g.createElement('link');" +
			"r.href=s;" +
			"r.type='text/css';" +
			"r.rel='stylesheet';" +
			"r.onload=_=>console.log(`gotrix: 🎨 loaded async stylesheet ${s}`);" +
			"r.onerror=e=>e.preventDefault()||console.error(`gotrix: 🛑 failed to load async stylesheet ${s}`);" +
			"t.appendChild(r)" +
			"})" +
			"})(document," + ac + ");</script>"
	}
	eb += pd.CSS()
	eb += "<!-- /endhead -->"
	return eb
}

func (pd *PageData) EndBody() string {
	eb := "<!-- endbody -->"
	eb += pd.JS()
	eb += "<!-- /endbody -->"
	return eb
}

func (pd *PageData) Component(name string, params ...interface{}) string {
	if comp, ok := pd.app.Components()[name]; ok {
		if err := comp.Include(pd, NewComponentParams(name, params)); err != nil {
			return pd.Error("failed to include component %s: %v", name, err)
		}
		return ""
	}
	return pd.Error("component %s not found", name)
}

func (pd *PageData) Error(format string, args ...interface{}) string {
	return fmt.Sprintf(
		`<div class="error">%s</div>`,
		fmt.Sprintf(format, args...))
}

func unique(v []string, w func(string) string) []string {
	if w == nil {
		w = func(s string) string { return s }
	}
	m := make(map[string]bool, 0)
	for _, s := range v {
		m[s] = true
	}
	l := make([]string, 0, len(m))
	for s := range m {
		l = append(l, w(s))
	}
	return l
}

func toJsTag(s, t string) string {
	if t != "" {
		t = " " + t
	}
	return fmt.Sprintf(`<script%s src="%s"></script>`, t, s)
}

func toCssTag(s string) string {
	return fmt.Sprintf(`<link rel="stylesheet" type="text/css" href="%s" />`, s)
}
