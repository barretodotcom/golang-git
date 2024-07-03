package index

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type Index struct {
	Objects map[string]Entry `json:"objects"`
}

type Entry struct {
	Hash         string    `json:"hash"`
	LastModified time.Time `json:"lastModified"`
	Path         string    `json:"path"`
}

func AddBlob(originalFile *os.File, blob *os.File) error {
	index, err := BuildIndex()
	if err != nil {
		return err
	}

	fileInfo, err := originalFile.Stat()
	if err != nil {
		return err
	}

	parts := strings.Split(blob.Name(), "/")
	blobHash := parts[len(parts)-2] + parts[len(parts)-1]

	index.Objects[originalFile.Name()] = Entry{
		Hash:         blobHash,
		Path:         originalFile.Name(),
		LastModified: fileInfo.ModTime(),
	}

	err = Update(index)

	return err

}

func Update(index Index) error {
	file, err := os.OpenFile(path.Join(".", ".gogit", "index"), os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	idxContent, err := json.Marshal(index)
	if err != nil {
		return err
	}

	_, err = file.Write(idxContent)

	return err
}

func BuildIndex() (Index, error) {
	idxFile, err := os.OpenFile(path.Join(".", ".gogit", "index"), os.O_RDONLY, 0644)
	if err != nil {
		return Index{}, err
	}
	defer idxFile.Close()

	fileContent, err := io.ReadAll(idxFile)
	if err != nil {
		return Index{}, err
	}

	var index Index
	err = json.Unmarshal(fileContent, &index)

	return index, err
}
