package command

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/barretodotcom/golang-git/commit"
	"github.com/barretodotcom/golang-git/index"
	"github.com/spf13/cobra"
)

var RollbackFile = &cobra.Command{
	Use:   "rollback-file",
	Short: "rollback a specific file to previous commit",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("provide a file to rollback")
		}

		idx, err := index.BuildIndex()
		if err != nil {
			log.Fatalf("error building index: %s", err)
		}

		head, err := commit.GetHEAD()
		if err != nil {
			log.Fatalf("couldn't get head: %s", err)
		}

		filePath := args[0]
		currentFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatalf("couldn't open file: %s", err)
		}
		defer currentFile.Close()

		if head.ParentHash == "" {
			initialFile, err := os.OpenFile(filepath.Join(".", ".gogit", "initial", filePath), os.O_RDONLY, 0755)
			if err != nil {
				log.Fatalf("couldn't read initial file: %s", err)
			}
			defer initialFile.Close()

			initialFileContent, err := io.ReadAll(initialFile)
			if err != nil {
				log.Fatalf("couldn't read initial file: %s", err)
			}

			_, err = currentFile.Write(initialFileContent)
			if err != nil {
				log.Fatalf("couldn't write file: %s", err)
			}
			return
		}

		if _, ok := idx.Objects[currentFile.Name()[2:]]; !ok {
			log.Fatalf("file not modified.")
		}

		previousHash := head.ParentHash

		head, err = commit.GetByHash(previousHash)
		if err != nil {
			log.Fatalf("couldn't get head: %s", err)
		}
		fileHash := head.Index.Objects[filePath].Hash
		folder := fileHash[0:2]
		fileName := fileHash[2:]

		file, err := os.OpenFile(filepath.Join(".", ".gogit", "objects", folder, fileName), os.O_RDONLY, 0655)
		if err != nil {
			log.Fatalf("couldn't open file: %s", err)
		}
		defer file.Close()

		fileContent, err := io.ReadAll(file)
		if err != nil {
			log.Fatalf("couldn't read file: %s", err)
		}

		_, err = currentFile.Write(fileContent)
		if err != nil {
			log.Fatalf("error while writing file: %s", err)
		}
	},
}
