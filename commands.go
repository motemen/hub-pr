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

	cli.Use(
		&cli.Command{
			Name:   "diff",
			Action: doDiff,
			Short:  "Show pull request's diff",
			Long:   "diff\n\nRuns \"git diff\" between the base and head of Pull Request corresponding to\ncurrent branch.",
		},
	)

	cli.Use(
		&cli.Command{
			Name:   "show",
			Action: doShow,
			Short:  "Show pull request's information",
			Long:   "show [-f TEMPLATE] [PULL_REQUEST_NUMBER]\n\nShows the information os a pull request (given # by argument or one current branch).\nThe context of the template is an octokit.Issue <http://godoc.org/github.com/github/hub/Godeps/_workspace/src/github.com/octokit/go-octokit/octokit#Issue>.",
		},
	)
}
