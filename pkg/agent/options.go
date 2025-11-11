package agent

type Options struct {
	Port       int
	MasterHost string
	MasterPort int
}

func NewOptions() *Options {
	return &Options{
		Port:       8081,
		MasterHost: "127.0.0.1",
		MasterPort: 8080,
	}
}

func WithPort(port int) func(o *Options) {
	return func(o *Options) {
		o.Port = port
	}
}

func WithMasterHost(masterHost string) func(o *Options) {
	return func(o *Options) {
		o.MasterHost = masterHost
	}
}

func WithMasterPort(masterPort int) func(o *Options) {
	return func(o *Options) {
		o.MasterPort = masterPort
	}
}
