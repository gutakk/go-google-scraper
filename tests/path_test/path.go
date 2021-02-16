package path_test

import (
	"os"
	"path/filepath"
	"runtime"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/helpers/log"
)

func ChangeToRootDir() {
	err := os.Chdir(getRoot())
	if err != nil {
		log.Fatal(errorconf.ChangeToRootDirFailure, err)
	}
}

func getRoot() string {
	_, b, _, _ := runtime.Caller(0)

	// Root folder of the project
	root := filepath.Join(filepath.Dir(b), "../..")

	return root
}
