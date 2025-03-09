package utility

import (
	"go-crawl/logger"
	"os"
	"path/filepath"
	str "strings"
)

func Exit(err error) {
	logger.Log("[error] %s", err.Error())
	os.Exit(1)
}

func Finish(format string, args ...interface{}) {
	logger.Log(format, args...)
	os.Exit(0)
}

func CurrentDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func MakeDir(path string) (string, error) {
	var result string

	if str.HasPrefix(path, ".") {
		currentDir, err := CurrentDir()
		if err != nil {
			return "", err
		}

		result = currentDir + string(os.PathSeparator)
	}

	result = result + path
	result, err := filepath.Abs(result)
	if err != nil {
		return "", err
	}

	// check if dir exist
	if _, e := os.Stat(result); os.IsNotExist(e) {
		err = os.MkdirAll(result, 0775)
		if err != nil {
			return "", err
		}
		return result, nil
	}

	return result, err
}
