package options

import (
	"os"
	"strconv"
)

type ServeOptions struct {
	// current host ip, some time same as host ip
	// if use reverse proxy, it will be the reverse proxy ip
	masterHost string
	// current host port
	// if use reverse proxy, it will be the reverse proxy port
	// if not use reverse proxy, it will be the port of the current host
	masterPort int
	Port       int
	DBPath     string
	dataPath   string
}

func NewServeOptions() *ServeOptions {
	dataPath := getDefaultDataPath("data")
	port := 8080
	portEnv, ok := os.LookupEnv("PANEL_PORT")
	if ok {
		portInt, err := strconv.Atoi(portEnv)
		if err == nil {
			port = portInt
		}
	}

	masterHost, ok := os.LookupEnv("PANEL_MASTER_HOST")
	if !ok {
		masterHost = "127.0.0.1"
	}
	masterPort, ok := os.LookupEnv("PANEL_MASTER_PORT")
	if !ok {
		masterPort = ""
	}

	masterPortInt, err := strconv.Atoi(masterPort)
	if err != nil {
		masterPortInt = port
	}

	return &ServeOptions{
		DBPath:     "lai-panel.db",
		Port:       port,
		dataPath:   dataPath,
		masterHost: masterHost,
		masterPort: masterPortInt,
	}
}

func WithDBPath(dbPath string) func(o *ServeOptions) {
	return func(o *ServeOptions) {
		o.DBPath = dbPath
	}
}

func WithServePort(port int) func(o *ServeOptions) {
	return func(o *ServeOptions) {
		o.Port = port
	}
}

func (o *ServeOptions) DataPath() string {
	return o.dataPath
}

func (o *ServeOptions) Agent() bool {
	return false
}

func (o *ServeOptions) MasterHost() string {
	return o.masterHost
}

func (o *ServeOptions) MasterPort() int {
	return o.masterPort
}
