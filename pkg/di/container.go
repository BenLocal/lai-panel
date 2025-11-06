package di

import (
	"go.uber.org/dig"
)

var Container *dig.Container

func init() {
	Container = dig.New()
}

func Provide(constructor interface{}) error {
	return Container.Provide(constructor)
}

// singleton
func Invoke(function interface{}) error {
	return Container.Invoke(function)
}
