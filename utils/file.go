package utils

import (
	"io"
	"os"
)

func GetFileContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileContent, err := io.ReadAll(file)

	return fileContent, err
}
