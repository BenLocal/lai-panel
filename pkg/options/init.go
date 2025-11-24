package options

import (
	"os"
	"path"
)

const (
	LOG_BASE_PATH     = "log"
	STATIC_BASE_PATH  = "static"
	SERVICE_BASE_PATH = "service"
)

func InitOptions(options IOptions) error {
	dataPath := options.DataPath()
	err := os.MkdirAll(dataPath, 0755)
	if err != nil {
		return err
	}

	for _, v := range []string{LOG_BASE_PATH, STATIC_BASE_PATH, SERVICE_BASE_PATH} {
		p := path.Join(dataPath, v)
		err = os.MkdirAll(p, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

type IOptions interface {
	DataPath() string

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
