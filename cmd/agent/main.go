package main

import (
	_ "github.com/benlocal/lai-panel/cmd/agent/route"
	"github.com/benlocal/lai-panel/pkg/agent"
)

func main() {
	runtime := agent.NewRuntime()

	if err := runtime.Start(); err != nil {
		panic(err)
	}
}
