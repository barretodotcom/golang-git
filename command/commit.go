package command

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/barretodotcom/golang-git/commit"
	"github.com/barretodotcom/golang-git/index"
	"github.com/spf13/cobra"
)

var Commit = &cobra.Command{
	Use:   "commit",
	Short: "apply changes",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("provide at least one argument      ")
		}

		idx, err := index.BuildIndex()
		if err != nil {
			log.Fatalf("error while building index: %s", err)
		}

		authorName := getEnvOrDefault("GIT_USER_NAME", "Jane Doe")
		authorEmail := getEnvOrDefault("GIT_USER_EMAIL", "janedoe@gmail.com")
		c := commit.Entry{
			AuthorName:  authorName,
			AuthorEmail: authorEmail,
			AuthorDate:  time.Now(),
			Index:       idx,
			Message:     args[0],
		}

		head, err := commit.GetHEAD()
		if err != nil {
			log.Fatalf("could not get head file: %s", err)
		}

		if head.Hash != "" {
			c.ParentHash = head.Hash
		}

		hash, err := commit.Write(c)
		if err != nil {
			log.Fatalf("could not write commit: %s", err)
		}
		fmt.Printf("[%s] %s\n", hash[0:6], c.Message)
	},
}

func getEnvOrDefault(envVariable string, defaultVariable string) string {
	env := os.Getenv(envVariable)
	if env == "" {
		return defaultVariable
	}
	return env

}
