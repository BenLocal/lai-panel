package route

import (
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/fasthttp/router"
)

type Registry struct {
	routeBindings []func(baseHandler *handler.BaseHandler, router *router.Router)
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		routeBindings: []func(baseHandler *handler.BaseHandler, router *router.Router){},
	}
}

func (r *Registry) Add(binding func(baseHandler *handler.BaseHandler, router *router.Router)) {
	r.routeBindings = append(r.routeBindings, binding)
}

func (r *Registry) Bindings() []func(baseHandler *handler.BaseHandler, router *router.Router) {
	return r.routeBindings
}
