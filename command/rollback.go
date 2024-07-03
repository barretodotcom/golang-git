package command

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/barretodotcom/golang-git/commit"
	"github.com/barretodotcom/golang-git/index"
	"github.com/barretodotcom/golang-git/objects"
	"github.com/spf13/cobra"
)

var Rollback = &cobra.Command{
	Use:   "rollback",
	Short: "rollback changes to previous commit",
	Run: func(cmd *cobra.Command, args []string) {
		head, err := commit.GetHEAD()
		if err != nil {
			log.Fatalf("error while building head: %s", err)
		}
		if head.ParentHash == "" {
			log.Fatalf("nothing to revert")
		}
		blob, err := objects.GetBlobByHash(head.ParentHash)
		if err != nil {
			log.Fatalf("error opening blob: %s", err)
		}
		defer blob.Close()

		previousHeadContent, err := io.ReadAll(blob)
		if err != nil {
			log.Fatalf("error reseting blob: %s", err)
		}

		var previousHead commit.Entry
		err = json.Unmarshal(previousHeadContent, &previousHead)
		if err != nil {
			log.Fatalf("error decoding blob: %s", err)
		}

		headFile, err := os.OpenFile(filepath.Join(".", ".gogit", "HEAD"), os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatalf("error reading HEAD file: %s", err)
		}
		defer headFile.Close()

		_, err = headFile.Write([]byte(previousHead.Hash))
		if err != nil {
			log.Fatalf("error writing HEAD file: %s", err)
		}

		idx, err := index.BuildIndex()
		if err != nil {
			log.Fatalf("error building index file: %s", err)
		}
		idx.Objects = previousHead.Index.Objects

		for _, obj := range idx.Objects {
			blob, err := objects.GetBlobByHash(obj.Hash)
			if err != nil {
				log.Fatalf("error getting blob file: %s", err)
			}
			defer blob.Close()

			fileContent, err := io.ReadAll(blob)
			if err != nil {
				log.Fatalf("error reading blob file: %s", err)
			}

			originalFile, err := os.OpenFile(filepath.Join(".", obj.Path), os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				log.Fatalf("error reading blob file: %s", err)
			}
			defer originalFile.Close()

			_, err = originalFile.Write(fileContent)
			if err != nil {
				log.Fatalf("error writing file: %s", err)
			}
		}

		err = index.Update(idx)
		if err != nil {
			log.Fatalf("error updating index file: %s", err)
		}

	},
}
