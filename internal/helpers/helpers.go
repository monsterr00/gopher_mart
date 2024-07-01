package helpers

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
)

func AbsolutePath(pathStart string, pathEnd string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Server: cant get rooted path")
	}

	if path.Base(cwd) == "gopher_mart" { //проверяем, запущено из корня бинарником тестов или нет
		return path.Join(pathStart, cwd, pathEnd)
	} else {
		absPath, _ := url.JoinPath(pathStart, cwd, "../..", pathEnd)
		return absPath
	}
}

func WrapError(err error, message string) error {
	return fmt.Errorf("%v: %w", message, err)
}
