package route

import (
	"github.com/fasthttp/router"
)

type Registry struct {
	routeBindings []func(router *router.Router)
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		routeBindings: []func(router *router.Router){},
	}
}

func (r *Registry) Add(binding func(router *router.Router)) {
	r.routeBindings = append(r.routeBindings, binding)
}

func (r *Registry) Bindings() []func(router *router.Router) {
	return r.routeBindings
}
