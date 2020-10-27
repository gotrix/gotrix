package gotrix

// ComponentWrapper interface.
type ComponentWrapper interface {
	Include(Context, ComponentParams) error
}

// Component interface. Must be implemented by all included components.
type Component interface {
	Component(Context, ComponentParams) (map[string]interface{}, error)
}

type ComponentParams interface {
	Name() string
	Params() []interface{}
}

func NewComponentParams(name string, params []interface{}) *componentParams {
	return &componentParams{
		name:   name,
		params: params,
	}
}

type componentParams struct {
	name   string
	params []interface{}
}

func (cp *componentParams) Name() string {
	return cp.name
}

func (cp *componentParams) Params() []interface{} {
	return cp.params
}
