package commit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/barretodotcom/golang-git/hash"
	"github.com/barretodotcom/golang-git/index"
)

type Entry struct {
	AuthorName  string      `json:"authorName"`
	AuthorEmail string      `json:"authorEmail"`
	AuthorDate  time.Time   `json:"authorDate"`
	Message     string      `json:"message"`
	Hash        string      `json:"hash"`
	ParentHash  string      `json:"parentHash"`
	Index       index.Index `json:"index"`
}

func (e Entry) String() string {
	commitString := fmt.Sprintf(
		`
commit %s
	Author: %s <%s>
	Date: %s
	Message: %s
`,
		e.Hash, e.AuthorName, e.AuthorEmail, e.AuthorDate.String(), e.Message)

	return commitString
}

func GetHEAD() (Entry, error) {
	file, err := os.OpenFile(path.Join(".", ".gogit", "HEAD"), os.O_RDONLY, 0655)
	if err != nil {
		return Entry{}, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		if err == io.EOF {
			return Entry{}, nil
		}
		return Entry{}, err
	}

	return GetByHash(string(content))

}

func GetByHash(hash string) (Entry, error) {
	if hash == "" {
		return Entry{}, nil
	}

	folder := hash[0:2]
	fileName := hash[2:]
	path := filepath.Join(".", ".gogit", "objects", folder, fileName)

	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return Entry{}, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return Entry{}, err
	}

	var c Entry
	if err = json.Unmarshal(content, &c); err != nil {
		return Entry{}, fmt.Errorf("not possible to unmarshal commit file: %w", err)
	}

	return c, err
}

func Write(c Entry) (string, error) {
	idx := c.Index

	hashes := make([]string, len(idx.Objects))

	var i int
	for _, obj := range idx.Objects {
		hashes[i] = obj.Hash
		i++
	}
	sort.Strings(hashes)
	allHashes := strings.Join(hashes, ".")
	c.Hash = hash.New([]byte(allHashes))

	head, err := GetHEAD()
	if err != nil {
		return "", err
	}
	if head.Hash == c.Hash {
		return "", fmt.Errorf("nothing to commit, working tree clean.")
	}

	folder := c.Hash[0:2]
	fileName := c.Hash[2:]
	path := filepath.Join(".", ".gogit", "objects", folder)
	err = os.Mkdir(path, 0755)
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	file, err := os.OpenFile(filepath.Join(path, fileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	_, err = file.Write(content)
	if err != nil {
		return "", err
	}

	headFile, err := os.OpenFile(filepath.Join(".", ".gogit", "HEAD"), os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return "", err
	}
	defer headFile.Close()

	_, err = headFile.Write([]byte(c.Hash))

	return c.Hash, err
}
