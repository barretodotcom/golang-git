package command

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/andreyvit/diff"
	"github.com/barretodotcom/golang-git/commit"
	"github.com/barretodotcom/golang-git/objects"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var Show = &cobra.Command{
	Use:   "show",
	Short: "show commit changes",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("please provide a commit hash")
		}
		commitHash := args[0]

		currentCommit, err := commit.GetByHash(commitHash)
		if err != nil {
			log.Fatalf("couldn't get commit: %s", err)
		}

		if currentCommit.Hash == "" {
			log.Fatalf("commit not found")
		}

		parentCommit, err := commit.GetByHash(currentCommit.ParentHash)
		if err != nil {
			log.Fatalf("couldn't get parent commit: %s", err)
		}

		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		for path := range currentCommit.Index.Objects {
			if _, ok := parentCommit.Index.Objects[path]; !ok {
				newBlob, err := objects.GetBlobByHash(currentCommit.Index.Objects[path].Hash)
				if err != nil {
					log.Fatalf("current read blob: %s", err)
				}
				defer newBlob.Close()
				blobContent, err := io.ReadAll(newBlob)
				if err != nil {
					log.Fatalf("couldn't read blob: content %s", blobContent)
				}

				fmt.Println(green("file added: "))
				fmt.Println(green(path))
				fmt.Println(green(string(blobContent)))
				continue
			}
			if currentCommit.Index.Objects[path].Hash != parentCommit.Index.Objects[path].Hash {
				currentBlob, err := objects.GetBlobByHash(currentCommit.Index.Objects[path].Hash)
				if err != nil {
					log.Fatalf("current read blob: %s", err)
				}
				defer currentBlob.Close()
				currentContent, err := io.ReadAll(currentBlob)
				if err != nil {
					log.Fatalf("couldn't read blob: content %s", currentContent)
				}

				parentBlob, err := objects.GetBlobByHash(parentCommit.Index.Objects[path].Hash)
				if err != nil {
					log.Fatalf("current read blob: %s", err)
				}
				defer parentBlob.Close()

				parentContent, err := io.ReadAll(parentBlob)
				if err != nil {
					log.Fatalf("couldn't read blob content: %s", currentContent)
				}
				diffs := diff.LineDiffAsLines(string(currentContent), string(parentContent))
				diffs = diff.TrimLines(diffs)

				fmt.Printf("%s:\n", yellow(path))
				for _, line := range diffs {
					if strings.HasPrefix(line, "-") {
						fmt.Printf("%s\n", red(line))
					}
					if strings.HasPrefix(line, "+") {
						fmt.Printf("%s\n", green(line))
					}
				}
			}
		}

	},
}
