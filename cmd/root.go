package cmd

import (
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/go-git/go-git/v6"
	"github.com/google/go-github/v42/github"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh-pr-url",
	Args:  cobra.ExactArgs(0),
	Short: "GitHub CLI extension for getting the PR URL associated with the current branch.",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentDir, err := os.Getwd()
		if err != nil {
			return err
		}

		localRepo, err := git.PlainOpenWithOptions(currentDir, &git.PlainOpenOptions{DetectDotGit: true})
		if err != nil {
			return fmt.Errorf("not a git repository: %w", err)
		}

		headRef, err := localRepo.Head()
		if err != nil {
			return fmt.Errorf("could not read HEAD: %w", err)
		}

		if !headRef.Name().IsBranch() {
			return fmt.Errorf("HEAD is detached — checkout a branch first")
		}

		currentBranchName := headRef.Name().Short()

		repo, err := repository.Current()
		if err != nil {
			return err
		}

		client, err := api.DefaultRESTClient()
		if err != nil {
			return err
		}

		var restRepo github.Repository

		err = client.Get(fmt.Sprintf("repos/%s/%s", repo.Owner, repo.Name), &restRepo)
		if err != nil {
			return err
		}

		if restRepo.DefaultBranch == nil {
			return fmt.Errorf("could not determine the default branch for this repository")
		}

		if currentBranchName == *restRepo.DefaultBranch {
			fmt.Println("You are on the default branch — there is no PR URL to show.")
			return nil
		}

		var pulls []github.PullRequest

		err = client.Get(fmt.Sprintf("repos/%s/%s/pulls?head=%s:%s", repo.Owner, repo.Name, repo.Owner, currentBranchName), &pulls)
		if err != nil {
			return err
		}

		if len(pulls) == 0 {
			fmt.Printf("No pull requests for branch: %s\n", currentBranchName)
			return nil
		}

		for _, pr := range pulls {
			if pr.HTMLURL != nil {
				fmt.Println(*pr.HTMLURL)
			}
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
