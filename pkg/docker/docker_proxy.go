package docker

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/docker/docker/client"
)

type DockerProxy struct {
	prefix  string
	hostURL *url.URL

	handler app.HandlerFunc
}

func NewDockerProxy(host string, prefix string) (*DockerProxy, error) {
	hostURL, err := client.ParseHostURL(host)
	if err != nil {
		return nil, err
	}

	localhost, _ := url.Parse("http://localhost")
	proxy := httputil.NewSingleHostReverseProxy(localhost)
	docketSockerDialer := &net.Dialer{}
	proxy.Transport = &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			switch hostURL.Scheme {
			case "unix":
				return docketSockerDialer.Dial("unix", fmt.Sprintf("%s%s", hostURL.Host, hostURL.Path))
			case "tcp":
				return docketSockerDialer.Dial("tcp", fmt.Sprintf("%s%s", hostURL.Host, hostURL.Path))
			case "npipe":
				return docketSockerDialer.Dial("npipe", fmt.Sprintf("%s%s", hostURL.Host, hostURL.Path))
			default:
				return nil, errors.New("invalid host scheme")
			}
		},
	}

	handler := adaptor.HertzHandler(proxy)
	return &DockerProxy{
		prefix:  prefix,
		hostURL: hostURL,
		handler: handler,
	}, nil
}

func (d *DockerProxy) HandleProxy(ctx context.Context, c *app.RequestContext) {
	// remove the prefix from the path
	path := strings.TrimPrefix(string(c.Request.Path()), d.prefix)
	c.Request.SetRequestURI(path)
	d.handler(ctx, c)
}
