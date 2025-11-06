package server

type Options struct {
	DBPath         string
	MigrationsPath string
}

func NewOptions() *Options {
	return &Options{
		DBPath:         "lai-panel.db",
		MigrationsPath: "migrations",
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
