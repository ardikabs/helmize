package util

import (
	"os"
	"path"
	"path/filepath"
)

func WriteTempValues(name string, data []byte) (string, error) {
	out, err := os.MkdirTemp("", name)
	if err != nil {
		return "", err
	}

	return WriteToFile(out, "values.yaml", data, false)
}

func WriteToFile(outputDir string, name string, data []byte, append bool) (filename string, err error) {
	outfileName := filepath.Join(outputDir, name)
	if err := ensureDirectoryForFile(outfileName); err != nil {
		return "", err
	}

	f, err := createOrOpenFile(outfileName, append)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err = f.Write(data); err != nil {
		return "", err
	}

	return f.Name(), nil
}

// check if the directory exists to create file. creates if don't exists
func ensureDirectoryForFile(file string) error {
	baseDir := path.Dir(file)
	_, err := os.Stat(baseDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return os.MkdirAll(baseDir, 0755)
}

func createOrOpenFile(filename string, append bool) (*os.File, error) {
	if append {
		return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	}

	return os.Create(filename)
}
