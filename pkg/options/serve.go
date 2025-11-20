package options

import (
	"path"
)

type ServeOptions struct {
	DBPath         string
	MigrationsPath string
	Port           int
	dataPath       string
}

func NewServeOptions() *ServeOptions {
	dataPath := getDefaultDataPath("serve")
	return &ServeOptions{
		DBPath:         "lai-panel.db",
		MigrationsPath: "migrations",
		Port:           8080,
		dataPath:       dataPath,
	}
}

func WithDBPath(dbPath string) func(o *ServeOptions) {
	return func(o *ServeOptions) {
		o.DBPath = dbPath
	}
}

func WithMigrationsPath(migrationsPath string) func(o *ServeOptions) {
	return func(o *ServeOptions) {
		o.MigrationsPath = migrationsPath
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

func (o *ServeOptions) StaticDataPath() string {
	return path.Join(o.dataPath, "static")
}

func (o *ServeOptions) LogDataPath() string {
	return path.Join(o.dataPath, "logs")
}

func (o *ServeOptions) ServicePath() string {
	return path.Join(o.dataPath, "services")
}
