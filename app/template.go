package app

import (
	"fmt"

	"github.com/gotrix/gotrix"
)

func (app *App) templateComponent(name string, data *gotrix.PageData, params ...string) string {
	//if comp, ok := app.components[name]; ok {
	//	res := comp.Include(gotrix.NewComponentParams(name, app, params))
	//	if res.Err() != nil {
	//		return app.templateError("failed to include component %s: %v", name, res.Err())
	//	}
	//	data.AddAsyncCSS(res.CSS.Async...)
	//	data.AddCSS(res.CSS.Regular...)
	//	data.AddAsyncJS(res.JS.Async...)
	//	data.AddDeferJS(res.JS.Defer...)
	//	data.AddJS(res.JS.Regular...)
	//	return res.Body
	//}
	return app.templateError("component %s not found", name)
}

func (app *App) templateError(format string, args ...interface{}) string {
	return fmt.Sprintf(
		`<div class="gtx-error">%s</div>`,
		fmt.Sprintf(format, args...))
}
