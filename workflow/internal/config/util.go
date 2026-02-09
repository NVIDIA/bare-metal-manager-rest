package config

import (
	"path/filepath"
	"runtime"
)

var (
	_, cur, _, _ = runtime.Caller(0)

	// ProjectRoot describes the folder path of this project
	ProjectRoot = filepath.Join(filepath.Dir(cur), "../..")
)
