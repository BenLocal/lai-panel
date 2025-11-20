package options

import (
	"os"
	"path"
)

func InitOptions(options IOptions) error {
	dataPath := options.DataPath()
	err := os.MkdirAll(dataPath, 0755)
	if err != nil {
		return err
	}

	// log
	logPath := options.LogDataPath()
	err = os.MkdirAll(logPath, 0755)
	if err != nil {
		return err
	}

	// static
	// add static file
	staticPath := options.StaticDataPath()
	err = os.MkdirAll(staticPath, 0755)
	if err != nil {
		return err
	}

	// service
	secretPath := options.ServicePath()
	err = os.MkdirAll(secretPath, 0755)
	if err != nil {
		return err
	}

	return nil
}

type IOptions interface {
	DataPath() string

	StaticDataPath() string

	LogDataPath() string

	ServicePath() string

	Agent() bool
}

func getDefaultDataPath(p string) string {
	d := DefaultDataPath

	if d == "" {
		home := os.Getenv("HOME")
		if home != "" {
			d = path.Join(os.Getenv("HOME"), ".lai-panel")
		}
	}

	if d == "" {
		tmp := os.TempDir()
		if tmp != "" {
			d = path.Join(tmp, "lai-panel")
		}
	}

	if d == "" {
		return p
	}

	return path.Join(d, p)
}
