package dashboard

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var distFS embed.FS

// GetDashboardFS returns the dist folder as the FS root.
func GetDashboardFS() fs.FS {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		// should not happen; return original FS to avoid panic
		return distFS
	}
	return sub
}
