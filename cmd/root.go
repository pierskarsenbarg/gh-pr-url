package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
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
		localRepo, err := git.PlainOpen(currentDir)

		branches, err := localRepo.Branches()
		if err != nil {
			return err
		}

		_, err = branches.Next()
		if err != nil {
			return fmt.Errorf("No branches in repository")
		}

		headRef, err := localRepo.Head()
		if err != nil {
			return err
		}

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

		if headRef.Name().Short() == *restRepo.DefaultBranch {
			return err
		}

		var currentBranchName string
		err = branches.ForEach(func(branchRef *plumbing.Reference) error {
			if branchRef.Hash() == headRef.Hash() {
				currentBranchName = branchRef.Name().Short()
				return nil
			}

			return nil
		})
		if err != nil {
			return err
		}

		slog.Info(fmt.Sprintf("Current branch: %s", currentBranchName))

		var pulls []github.PullRequest

		err = client.Get(fmt.Sprintf("repos/%s/%s/pulls?head=%s:%s", repo.Owner, repo.Name, repo.Owner, currentBranchName), &pulls)
		if err != nil {
			return err
		}

		if len(pulls) == 0 {
			fmt.Printf("No Pull requests for this branch: %s", repo.Name)
		}

		if len(pulls) == 1 {
			fmt.Print(*pulls[0].HTMLURL)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
