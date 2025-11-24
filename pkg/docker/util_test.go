package docker

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func NewTestListener(tb testing.TB) net.Listener {
	tb.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		tb.Fatalf("failed to create test listener: %s", err)
	}
	return ln
}

type RouteEngine interface {
	IsRunning() bool
}

func WaitEngineRunning(e RouteEngine) {
	for i := 0; i < 100; i++ {
		if e.IsRunning() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	panic("not running")
}

func runEngine(onCreate func(*route.Engine)) (string, int, *route.Engine) {
	ln := NewTestListener(&testing.T{})
	opt := config.NewOptions(nil)
	opt.Listener = ln
	engine := route.NewEngine(opt)
	onCreate(engine)
	go engine.Run()
	WaitEngineRunning(engine)
	host := ln.Addr().(*net.TCPAddr).IP.String()
	port := ln.Addr().(*net.TCPAddr).Port
	return host, port, engine
}

func TestLocalDockerClient(t *testing.T) {
	dockerClient, err := LocalDockerClient()
	if err != nil {
		t.Fatalf("failed to create local docker client: %v", err)
	}
	defer dockerClient.Close()

	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		t.Fatalf("failed to list containers: %v", err)
	}
	t.Logf("containers: %v", containers)

	images, err := dockerClient.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		t.Fatalf("failed to list images: %v", err)
	}
	t.Logf("images: %v", images)
}

func TestAgentDockerClient(t *testing.T) {
	// start agent server
	dh := client.DefaultDockerHost
	dd, _ := NewDockerProxy(dh, "/docker.proxy")
	host, port, engine := runEngine(func(engine *route.Engine) {
		engine.Any("/docker.proxy/*path", func(ctx context.Context, c *app.RequestContext) {
			dd.HandleProxy(ctx, c)
		})
	})
	defer engine.Close()

	dockerClient, err := AgentDockerClient(host, port)
	if err != nil {
		t.Fatalf("failed to create agent docker client: %v", err)
	}
	defer dockerClient.Close()

	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		t.Fatalf("failed to list containers: %v", err)
	}
	t.Logf("containers: %v", containers)

	images, err := dockerClient.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		t.Fatalf("failed to list images: %v", err)
	}
	t.Logf("images: %v", images)
}
