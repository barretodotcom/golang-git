package command

import (
	"io"
	"log"
	"os"

	"github.com/barretodotcom/golang-git/commit"
	"github.com/barretodotcom/golang-git/objects"
	"github.com/spf13/cobra"
)

var Reset = &cobra.Command{
	Use:   "reset",
	Short: "reset changes to the last commit",
	Run: func(cmd *cobra.Command, args []string) {
		head, err := commit.GetHEAD()
		if err != nil {
			log.Fatalf("error while building head: %s", err)
		}
		for filePath, obj := range head.Index.Objects {
			blob, err := objects.GetBlobByHash(obj.Hash)
			if err != nil {
				log.Fatalf("error opening blob: %s", err)
			}
			defer blob.Close()

			currentFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				log.Fatalf("error opening current file: %s", err)
			}
			defer currentFile.Close()

			blobContent, err := io.ReadAll(blob)
			if err != nil {
				log.Fatalf("error reading blob: %s", err)
			}

			_, err = currentFile.Write(blobContent)
			if err != nil {
				log.Fatalf("error writing in file: %s", err)
			}
		}

	},
}
