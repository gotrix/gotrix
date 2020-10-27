package gotrix

// Component interface. Must be implemented by all included components.
type Component interface {
	Include(ComponentParams) (string, error)
}

type ComponentParams interface {
	Name() string
	Params() []string
	App() App
}

func NewComponentParams(name string, app App, params []string) *componentParams {
	return &componentParams{
		name:   name,
		app:    app,
		params: params,
	}
}

type componentParams struct {
	name   string
	params []string
	app    App
}

func (cp *componentParams) Name() string {
	return cp.name
}

func (cp *componentParams) Params() []string {
	return cp.params
}

func (cp *componentParams) App() App {
	return cp.app
}

type ComponentWrapper interface {
	Include(ComponentParams) *ComponentResult
}

type ComponentPathsList interface {
	All() []string
}

type ComponentResult struct {
	Error error
	Hash  string
	Body  string
	CSS   *ComponentCSS
	JS    *ComponentJS
}

func (cr *ComponentResult) Err() error {
	return cr.Error
}

type ComponentCSS struct {
	Regular []string
	Async   []string
}

func (c *ComponentCSS) All() []string {
	list := make([]string, 0, len(c.Regular)+len(c.Async))
	list = append(list, c.Regular...)
	list = append(list, c.Async...)
	return list
}

type ComponentJS struct {
	Regular []string
	Async   []string
	Defer   []string
}

func (c *ComponentJS) All() []string {
	list := make([]string, 0)
	list = append(list, c.Regular...)
	list = append(list, c.Async...)
	list = append(list, c.Defer...)
	return list
}
