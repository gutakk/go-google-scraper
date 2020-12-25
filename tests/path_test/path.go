package path_test

import (
	"path/filepath"
	"runtime"
)

func GetRoot() string {
	_, b, _, _ := runtime.Caller(0)

	// Root folder of the project
	root := filepath.Join(filepath.Dir(b), "../..")

	return root
}
