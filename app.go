package gotrix

import (
	"gopkg.in/reform.v1"
)

// App interface.
type App interface {
	DB() *reform.DB
	Components() map[string]ComponentWrapper
}

// AppConfig struct.
type AppConfig struct {
	ComponentPaths []string
	TemplatePaths  []string
	StaticPaths    []string
}
