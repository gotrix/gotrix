package gotrix

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

func NewPageData() *PageData {
	return &PageData{}
}

type PageData struct {
	req        *http.Request
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

func (cd *PageData) Request() *http.Request {
	return cd.req
}

func (cd *PageData) AddAsyncCSS(css ...string) {
	cd.asyncCssMu.Lock()
	cd.asyncCss = append(cd.asyncCss, css...)
	cd.asyncCssMu.Unlock()
}

func (cd *PageData) AddCSS(css ...string) {
	cd.cssMu.Lock()
	cd.css = append(cd.css, css...)
	cd.cssMu.Unlock()
}

func (cd *PageData) AddAsyncJS(js ...string) {
	cd.asyncJsMu.Lock()
	cd.asyncJs = append(cd.asyncJs, js...)
	cd.asyncJsMu.Unlock()
}

func (cd *PageData) AddDeferJS(js ...string) {
	cd.deferJsMu.Lock()
	cd.deferJs = append(cd.deferJs, js...)
	cd.deferJsMu.Unlock()
}

func (cd *PageData) AddJS(js ...string) {
	cd.jsMu.Lock()
	cd.js = append(cd.js, js...)
	cd.jsMu.Unlock()
}

func (cd *PageData) JS() string {
	l := make([]string, 0)
	l = append(l, unique(cd.asyncJs, func(s string) string {
		return toJsTag(s, "async")
	})...)
	l = append(l, unique(cd.deferJs, func(s string) string {
		return toJsTag(s, "defer")
	})...)
	l = append(l, unique(cd.js, func(s string) string {
		return toJsTag(s, "")
	})...)
	return strings.Join(l, "")
}

func (cd *PageData) CSS() string {
	return strings.Join(unique(cd.css, func(s string) string {
		return toCssTag(s)
	}), "")
}

func (cd *PageData) AsyncCSS() string {
	return `["` + strings.Join(unique(cd.asyncCss, nil), `","`) + `"]`
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
