package options

import (
	"os"
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

	return nil
}

type IOptions interface {
	DataPath() string

	StaticDataPath() string

	LogDataPath() string

	Agent() bool
}
