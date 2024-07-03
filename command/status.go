package command

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/barretodotcom/golang-git/commit"
	"github.com/barretodotcom/golang-git/hash"
	"github.com/barretodotcom/golang-git/index"
	"github.com/barretodotcom/golang-git/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var idx index.Index
var changedFiles []string
var untrackedFiles []string
var stagedFiles []string
var err error

var Status = &cobra.Command{
	Use:   "status",
	Short: "see file changes",
	Run: func(cmd *cobra.Command, args []string) {

		idx, err = index.BuildIndex()
		if err != nil {
			log.Fatalf("error while building index: %s", err)
		}

		filepath.Walk(".", ClassifyFiles)

		red := color.New(color.FgRed).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		if len(untrackedFiles) > 0 {
			fmt.Println("Untracked files:  ")
			fmt.Println(red("   " + strings.Join(untrackedFiles, "\n   ")))
		}
		if len(stagedFiles) > 0 {
			fmt.Println("Staged files: ")
			fmt.Println(green("   " + strings.Join(stagedFiles, "\n   ")))
		}
		if len(changedFiles) > 0 {
			fmt.Println("Changed files: ")
			fmt.Println(red("   " + strings.Join(changedFiles, "\n   ")))
		}

		if len(changedFiles) == 0 && len(stagedFiles) == 0 && len(untrackedFiles) == 0 {
			fmt.Println("Nothing to commit, working tree clean")
		}

	},
}

func ClassifyFiles(path string, fileInfo fs.FileInfo, err error) error {
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

	idx, err := index.BuildIndex()
	if err != nil {
		return err
	}

	entry, ok := idx.Objects[path]
	if !ok {
		untrackedFiles = append(untrackedFiles, path)
		return nil
	}

	fileContent, err := utils.GetFileContent(path)
	if err != nil {
		return err
	}

	h := hash.New(fileContent)

	if entry.Hash != h {
		changedFiles = append(changedFiles, path)
		return nil
	}

	head, err := commit.GetHEAD()
	if err != nil {
		return err
	}

	if head.Index.Objects[path].Hash != idx.Objects[path].Hash {
		stagedFiles = append(stagedFiles, path)
		return nil
	}

	return nil
}
