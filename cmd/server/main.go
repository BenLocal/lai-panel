package main

import (
	_ "github.com/benlocal/lai-panel/cmd/server/route"
	"github.com/benlocal/lai-panel/pkg/server"
)

func main() {
	runtime := server.NewRuntime()

	if err := runtime.Start(); err != nil {
		panic(err)
	}
}
