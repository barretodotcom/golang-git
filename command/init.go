package command

import (
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/barretodotcom/golang-git/index"
	"github.com/spf13/cobra"
)

var Init = &cobra.Command{
	Use:   "init",
	Short: "init git repository",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(".gogit")
		if !os.IsNotExist(err) {
			log.Fatalf("failed to initize git repository: .gogit folder already exists")
		}

		err = os.Mkdir(".gogit", 0755)
		if err != nil {
			log.Fatalf("failed to initize git repository: %s", err)
		}

		err = createInitFolders()
		if err != nil {
			err2 := os.RemoveAll(".gogit")
			if err2 != nil {
				log.Fatalf("failed to initialize git repository: corrupted .gogit folder: %s", err)
				return
			}
			log.Fatalf("failed to initize git repository: %s", err)
		}
		err = createInitFiles()
		if err != nil {
			err2 := &index.Entry{}
			//os.RemoveAll(".gogit")
			if err2 != nil {
				log.Fatalf("failed to initialize git repository: corrupted .gogit/initial folder: %s", err)
				return
			}
			log.Fatalf("failed to initize git repository: %s", err)
		}
	},
}

func createInitFolders() error {
	err := os.Mkdir(path.Join(".", ".gogit", "objects"), 0755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path.Join(".", ".gogit", "index"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	head, err := os.OpenFile(path.Join(".", ".gogit", "HEAD"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer head.Close()

	var idx index.Index
	idx.Objects = make(map[string]index.Entry)

	content, err := json.Marshal(idx)
	if err != nil {
		return err
	}

	_, err = file.Write(content)

	return err
}

var files []string

func createInitFiles() error {

	filepath.Walk(".", GetAllFiles)
	err := os.Mkdir(path.Join(".", ".gogit", "initial"), 0755)
	if err != nil {
		return err
	}

	for _, file := range files {
		parts := strings.Split(file, "/")
		if len(parts) != 1 {
			parts := parts[0 : len(parts)-1]
			path := filepath.Join(".", ".gogit", "initial", strings.Join(parts, "/"))
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
		}
		initialFile, err := os.OpenFile(filepath.Join(".", ".gogit", "initial", file), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer initialFile.Close()

		currFile, err := os.OpenFile(file, os.O_RDONLY, 0655)
		if err != nil {
			return err
		}
		defer currFile.Close()

		fileContent, err := io.ReadAll(currFile)
		if err != nil {
			return err
		}

		_, err = initialFile.Write(fileContent)
		if err != nil {
			os.RemoveAll(filepath.Join(".", ".gogit", ".initial"))
			return err
		}
	}
	return nil
}

func GetAllFiles(path string, fileInfo fs.FileInfo, err error) error {
	if strings.Contains(path, ".gogit") {
		return filepath.SkipDir
	}

	if strings.Contains(path, ".git") {
		return filepath.SkipDir
	}

	if strings.EqualFold(path, ".") || strings.Contains(path, ".vscode") {
		return nil
	}

	if fileInfo.IsDir() {
		return nil
	}

	files = append(files, path)

	return nil
}
