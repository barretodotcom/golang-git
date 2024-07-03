package command

import (
	"fmt"
	"log"
	"strconv"

	"github.com/barretodotcom/golang-git/commit"
	"github.com/spf13/cobra"
)

var Log = &cobra.Command{
	Use:   "log",
	Short: "show the most recent commits",
	Run: func(cmd *cobra.Command, args []string) {
		head, err := commit.GetHEAD()
		if err != nil {
			log.Fatalf("error while opening head: %s", err)
		}

		limit := int(5)
		if len(args) > 0 {
			count, err := strconv.ParseInt(args[0], 10, 64)
			if err == nil {
				limit = int(count)
			}
		}
		currentHash := head.ParentHash
		commits := []commit.Entry{head}
		for i := 1; i < limit; i++ {
			if currentHash == "" {
				break
			}
			commit, err := commit.GetByHash(currentHash)
			if err != nil {
				log.Fatalf("error while opening head: %s", err)
			}
			commits = append(commits, commit)
			currentHash = commit.ParentHash
		}

		for i := range commits {
			fmt.Println(commits[i].String())
		}

	},
}
