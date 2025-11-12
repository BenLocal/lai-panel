package docker

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
)

type DockerProxy struct {
	prefix     string
	socketPath string

	handler app.HandlerFunc
}

func NewDockerProxy(socketPath string, prefix string) (*DockerProxy, error) {
	// 检查 Docker socket 是否存在
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		return nil, err
	}

	localhost, _ := url.Parse("http://localhost")
	proxy := httputil.NewSingleHostReverseProxy(localhost)
	docketSockerDialer := &net.Dialer{}
	proxy.Transport = &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return (docketSockerDialer.Dial("unix", socketPath))
		},
	}

	handler := adaptor.HertzHandler(proxy)
	return &DockerProxy{
		prefix:     prefix,
		socketPath: socketPath,
		handler:    handler,
	}, nil
}

func (d *DockerProxy) HandleProxy(ctx context.Context, c *app.RequestContext) {
	// remove the prefix from the path
	path := strings.TrimPrefix(string(c.Request.Path()), d.prefix)
	c.Request.SetRequestURI(path)
	d.handler(ctx, c)
}
