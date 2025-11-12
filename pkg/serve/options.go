package serve

type ServeOptions struct {
	DBPath         string
	MigrationsPath string
	Port           int
}

func NewServeOptions() *ServeOptions {
	return &ServeOptions{
		DBPath:         "lai-panel.db",
		MigrationsPath: "migrations",
		Port:           8080,
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

func WithPort(port int) func(o *ServeOptions) {
	return func(o *ServeOptions) {
		o.Port = port
	}
}
