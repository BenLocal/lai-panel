package api

import (
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/cloudwego/hertz/pkg/route"
)

type Registry struct {
	routeBindings []func(baseHandler *handler.BaseHandler, router *route.Engine)
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		routeBindings: []func(baseHandler *handler.BaseHandler, router *route.Engine){},
	}
}

func (r *Registry) Add(binding func(baseHandler *handler.BaseHandler, router *route.Engine)) {
	r.routeBindings = append(r.routeBindings, binding)
}

func (r *Registry) Bindings() []func(baseHandler *handler.BaseHandler, router *route.Engine) {
	return r.routeBindings
}
