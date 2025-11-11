package server

type Options struct {
	DBPath         string
	MigrationsPath string
	Port           int
}

func NewOptions() *Options {
	return &Options{
		DBPath:         "lai-panel.db",
		MigrationsPath: "migrations",
		Port:           8080,
	}
}

func WithDBPath(dbPath string) func(o *Options) {
	return func(o *Options) {
		o.DBPath = dbPath
	}
}

func WithMigrationsPath(migrationsPath string) func(o *Options) {
	return func(o *Options) {
		o.MigrationsPath = migrationsPath
	}
}

func WithPort(port int) func(o *Options) {
	return func(o *Options) {
		o.Port = port
	}
}
