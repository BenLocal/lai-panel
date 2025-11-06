package handler

import (
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/benlocal/lai-panel/pkg/di"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type DockerProxy struct {
	prefix     string
	socketPath string

	proxyClient *fasthttp.HostClient
}

func HandleDockerProxyWithDI(ctx *fasthttp.RequestCtx) {
	err := di.Invoke(func(dd *DockerProxy) {
		if string(ctx.Request.Header.Peek("Upgrade")) != "" ||
			string(ctx.Request.Header.Peek("Connection")) == "Upgrade" {
			dd.handleExecStartUpgrade(ctx)
			return
		}
		dd.handleHTTPProxy(ctx)
	})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
	}
}

func NewDockerProxy(socketPath string, prefix string) (*DockerProxy, error) {
	// 检查 Docker socket 是否存在
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		return nil, err
	}

	// 创建 Docker socket 的 HostClient
	proxyClient := &fasthttp.HostClient{
		StreamResponseBody: true,
		Addr:               "unix",
		Dial: func(addr string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	}

	return &DockerProxy{
		prefix:      prefix,
		socketPath:  socketPath,
		proxyClient: proxyClient,
	}, nil
}

func (d *DockerProxy) WithRoutes() func(router *router.Router) {
	return func(r *router.Router) {
		r.ANY("/docker.proxy/{any:*}", func(ctx *fasthttp.RequestCtx) {
			// 处理普通 HTTP 请求

			if string(ctx.Request.Header.Peek("Upgrade")) != "" ||
				string(ctx.Request.Header.Peek("Connection")) == "Upgrade" {
				d.handleExecStartUpgrade(ctx)
				return
			}
			d.handleHTTPProxy(ctx)
		})
	}
}

func (d *DockerProxy) handleHTTPProxy(ctx *fasthttp.RequestCtx) {
	req := &ctx.Request
	resp := &ctx.Response

	d.prepareRequest(req)

	if err := d.proxyClient.Do(req, resp); err != nil {
		log.Printf("error when proxying the request: %s", err)
		ctx.SetStatusCode(fasthttp.StatusBadGateway)
		ctx.SetBodyString("Docker proxy error: " + err.Error())
		return
	}

	d.postprocessResponse(resp)
}

func (d *DockerProxy) prepareRequest(req *fasthttp.Request) {
	req.Header.Del("Connection")
	req.Header.Del("Upgrade")
	req.Header.Set("Host", "localhost")
	uri := req.URI()
	path := uri.Path()
	prefix := d.prefix
	// 处理 Docker socket 的路径
	if len(path) >= len(prefix) && string(path[:len(prefix)]) == prefix {
		newPath := path[len(prefix):]
		if len(newPath) == 0 || newPath[0] != '/' {
			newPath = append([]byte("/"), newPath...)
		}
		uri.SetPathBytes(newPath)
	}
}

func (d *DockerProxy) postprocessResponse(resp *fasthttp.Response) {
	resp.Header.Del("Connection")
}

func (d *DockerProxy) handleExecStartUpgrade(ctx *fasthttp.RequestCtx) {
	log.Printf("Handling ExecStart upgrade for path: %s", ctx.Path())

	// 获取底层连接
	ctx.HijackSetNoResponse(true)
	ctx.Hijack(func(conn net.Conn) {
		defer conn.Close()

		// 连接到 Docker socket
		dockerConn, err := net.Dial("unix", d.socketPath)
		if err != nil {
			log.Printf("Failed to connect to Docker socket: %v", err)
			return
		}
		defer dockerConn.Close()

		// 构造请求
		req := &ctx.Request
		d.prepareExecRequest(req)

		log.Printf("Sending request to Docker: %s %s", req.Header.Method(), req.URI().Path())

		// 发送请求到 Docker socket
		if _, err := req.WriteTo(dockerConn); err != nil {
			log.Printf("Failed to write request to Docker socket: %v", err)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusSwitchingProtocols)
		log.Printf("Starting bidirectional copy for ExecStart")
		// 直接开始双向数据转发，让 HTTP 升级握手透传
		d.bidirectionalCopy(conn, dockerConn)

		log.Printf("ExecStart session ended")
	})

}

func (d *DockerProxy) prepareExecRequest(req *fasthttp.Request) {
	// 不要删除 Connection 和 Upgrade 头，ExecStart 需要这些
	req.Header.Set("Host", "localhost")

	uri := req.URI()
	path := uri.Path()
	prefix := d.prefix

	// 处理路径
	if len(path) >= len(prefix) && string(path[:len(prefix)]) == prefix {
		newPath := path[len(prefix):]
		if len(newPath) == 0 || newPath[0] != '/' {
			newPath = append([]byte("/"), newPath...)
		}
		uri.SetPathBytes(newPath)
	}
}

func (d *DockerProxy) bidirectionalCopy(client, docker net.Conn) {
	done := make(chan struct{}, 2)

	// 通用复制函数
	copyData := func(dst, src net.Conn, direction string) {
		defer func() {
			// 尝试关闭写端
			if closer, ok := dst.(interface{ CloseWrite() error }); ok {
				closer.CloseWrite()
			}
			// 非阻塞通知
			select {
			case done <- struct{}{}:
			default:
			}
		}()

		if _, err := io.Copy(dst, src); err != nil && !isConnectionClosed(err) {
			log.Printf("Error copying %s: %v", direction, err)
		}
	}

	go copyData(docker, client, "from client to docker")
	go copyData(client, docker, "from docker to client")

	<-done // 等待任一方向完成
}

func isConnectionClosed(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "use of closed network connection") ||
		strings.Contains(errStr, "connection reset by peer") ||
		strings.Contains(errStr, "broken pipe") ||
		err == io.EOF
}
