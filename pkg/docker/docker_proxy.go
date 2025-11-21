package docker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/docker/docker/client"
)

var dockerAPIVersionRegex = regexp.MustCompile(`^/v1\.\d+/`)

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

	localhost, _ := url.Parse("http://unix")
	u := fmt.Sprintf("%s%s", hostURL.Host, hostURL.Path)
	log.Printf("docker proxy host: %s\n, u: %s", host, u)
	proxy := httputil.NewSingleHostReverseProxy(localhost)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Scheme = "http"
		req.URL.Host = "unix"
		req.Host = "unix"
		if prefix != "" && len(req.URL.Path) >= len(prefix) &&
			req.URL.Path[:len(prefix)] == prefix {
			req.URL.Path = req.URL.Path[len(prefix):]
		}
	}
	proxy.Transport = &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			switch hostURL.Scheme {
			case "unix":
				return net.Dial("unix", u)
			case "tcp":
				return net.Dial("tcp", u)
			case "npipe":
				return net.Dial("npipe", u)
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
	d.handler(ctx, c)
}
