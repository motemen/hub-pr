// auto-generated file

package main

import "github.com/motemen/go-cli"

func init() {
	cli.Use(
		&cli.Command{
			Name:   "checkout",
			Action: doCheckout,
			Short:  "Checkout a branch for a pull request",
			Long:   "checkout [-f TEMPLATE] PULL_REQUEST_NUMBER\n\nChecks out a branch corresponding to given pull request\nnumber.  The branch name is based on the template,\nwhich defaults to\n'#{{.Number}}-{{.Head.Repo.Owner.Login}}/{{.Head.Ref}}',\nwhere the template context is *octokit.PullRequest.",
		},
	)

	cli.Use(
		&cli.Command{
			Name:   "list",
			Action: doList,
			Short:  "List pull requests",
			Long:   "list\n\nLists pull requests for current project.",
		},
	)

	cli.Use(
		&cli.Command{
			Name:   "merge",
			Action: doMerge,
			Short:  "Merge a branch of a pull request",
			Long:   "merge BRANCH\n\nInvokes 'git merge' for the branch created with 'hub-pr checkout' with a\ndefault commit message including Pull Request number and title, similar to\nthe GitHub Merge Button.",
		},
	)

	cli.Use(
		&cli.Command{
			Name:   "browse",
			Action: doBrowse,
			Short:  "Open pull request page with browser",
			Long:   "browse\n\nOpens a web browser for the URL of Pull Request corresponding to current\nbranch.",
		},
	)
}
