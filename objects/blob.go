package objects

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/barretodotcom/golang-git/hash"
)

func CreateBlob(file *os.File) (*os.File, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	h := hash.New(content)

	folder := h[0:2]
	fileName := h[2:]

	err = os.Mkdir(path.Join(".", ".gogit", "objects", folder), 0755)
	if os.IsExist(err) {
		err = os.RemoveAll(path.Join(".", ".gogit", "objects", folder))
		if err != nil {
			return nil, err
		}
		err = os.Mkdir(path.Join(".", ".gogit", "objects", folder), 0755)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	file, err = os.OpenFile(path.Join(".", ".gogit", "objects", folder, fileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil
	}
	defer file.Close()

	_, err = file.Write(content)

	return file, err
}

func GetBlobByHash(hash string) (*os.File, error) {
	folder := hash[0:2]
	fileName := hash[2:]

	path := filepath.Join(".", ".gogit", "objects", folder, fileName)

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)

	return file, err
}
