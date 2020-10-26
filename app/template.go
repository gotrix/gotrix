package app

import (
	"fmt"

	"github.com/gotrix/gotrix"
)

func (app *App) templateComponent(name string, params ...string) string {
	if comp, ok := app.components[name]; ok {
		res, err := comp.Include(gotrix.NewComponentParams(name, app, params))
		if err != nil {
			return app.templateError("failed to include component %s: %v", name, err)
		}
		return res
	}
	return app.templateError("component %s not found", name)
}

func (app *App) templateError(format string, args ...interface{}) string {
	return fmt.Sprintf(
		`<div class="gtx-error">%s</div>`,
		fmt.Sprintf(format, args...))
}
