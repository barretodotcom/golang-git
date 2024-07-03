package command

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/barretodotcom/golang-git/hash"
	"github.com/barretodotcom/golang-git/index"
	"github.com/barretodotcom/golang-git/objects"
	"github.com/barretodotcom/golang-git/utils"
	"github.com/spf13/cobra"
)

var Add = &cobra.Command{
	Use:   "add",
	Short: "add file to staging area",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("please provide at least one file argument      ")
		}
		if len(args) == 1 && args[0] == "." {
			filepath.Walk(".", ClassifyFiles)
			changedFiles = append(changedFiles, untrackedFiles...)
			args = changedFiles
		}

		for _, filePath := range args {
			idx, err := index.BuildIndex()
			if err != nil {
				log.Fatalf("error while reading index   : %s", err)
			}

			fileContent, err := utils.GetFileContent(filePath)
			if err != nil {
				log.Fatalf("error while opening file: %s", err)
			}

			h := hash.New(fileContent)
			if file, ok := idx.Objects[filePath]; ok {
				if file.Hash == h {
					continue
				}
			}

			file, err := os.OpenFile(filePath, os.O_RDONLY, 0655)
			if err != nil {
				log.Fatalf("error   while ssopening file: %s", err)
			}
			defer file.Close()

			blob, err := objects.CreateBlob(file)
			if err != nil {
				log.Fatalf("error while creating blob: %s", err)
			}
			defer blob.Close()

			err = index.AddBlob(file, blob)
			if err != nil {
				log.Fatalf("cannot add blob: %s", err)
			}

		}
	},
}

func compareFiles(file, blob *os.File) (bool, error) {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return false, err
	}

	blobContent, err := io.ReadAll(blob)
	if err != nil {
		return false, err
	}

	return string(fileContent) != string(blobContent), nil
}
