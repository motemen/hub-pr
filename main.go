//go:generate go run _tools/gen_commands.go

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/github/hub/Godeps/_workspace/src/github.com/octokit/go-octokit/octokit"
	"github.com/github/hub/git"
	"github.com/github/hub/github"
	"github.com/motemen/go-cli"
)

const defaultTemplate = `#{{.Number}}-{{.Head.Repo.Owner.Login}}/{{.Head.Ref}}`

// +command checkout - Checkout a branch for a pull request
//
// 	checkout [-f TEMPLATE] PULL_REQUEST_NUMBER
//
// Checks out a branch corresponding to given pull request
// number.  The branch name is based on the template,
// which defaults to
// '#{{.Number}}-{{.Head.Repo.Owner.Login}}/{{.Head.Ref}}',
// where the template context is *octokit.PullRequest.
func doCheckout(flags *flag.FlagSet, args []string) error {
	tmpl, err := git.GlobalConfig("hub-pr.checkoutBranch")
	if err != nil || tmpl == "" {
		tmpl = defaultTemplate
	}

	flags.StringVar(&tmpl, "f", tmpl, "branch name format")
	flags.Parse(args)

	tmplBranch := template.Must(template.New("branch").Parse(tmpl))

	prNumber := flags.Arg(0)
	if prNumber == "" {
		return cli.ErrUsage
	}

	cli, proj, err := setup()
	if err != nil {
		return err
	}

	pr, err := cli.PullRequest(proj, prNumber)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmplBranch.Execute(&buf, pr)
	if err != nil {
		return err
	}

	branchName := buf.String()

	remoteName := pr.Head.Repo.Owner.Login
	headBranch := pr.Head.Ref

	r := gitRunner{}

	remoteURL, err := git.Config("remote." + remoteName + ".url")

	if err == nil && remoteURL != "" {
		r.git("remote", "set-branches", "--add", remoteName, headBranch)
		r.git("fetch", remoteName)
	} else {
		r.git("remote", "add", "-t", headBranch, "-f", remoteName, pr.Head.Repo.GitURL)
	}

	if r.err != nil {
		return r.err
	}

	r.git("rev-parse", "--verify", branchName)
	if r.err != nil {
		r.resetError()
		r.git("checkout", "-b", branchName, "--track", fmt.Sprintf("remotes/%s/%s", remoteName, headBranch))
	}

	r.git("config", "--local", "branch."+branchName+".pushremote", "origin")
	r.git("config", "--local", "branch."+branchName+".prNumber", fmt.Sprintf("%d", pr.Number))

	return r.err
}

// +command list - List pull requests
//
// 	list
//
// Lists pull requests for current project.
func doList(flags *flag.FlagSet, args []string) error {
	flags.Parse(args)

	cli, proj, err := setup()
	if err != nil {
		return err
	}

	issues, err := cli.Issues(proj)
	if err != nil {
		return err
	}

	for _, issue := range issues {
		if issue.PullRequest.HTMLURL == "" {
			continue
		}

		fmt.Printf("%4d\t%s (@%s)\n", issue.Number, issue.Title, issue.User.Login)
	}

	return nil
}

// +command merge - Merge a branch of a pull request
//
// 	merge BRANCH
//
// Invokes 'git merge' for the branch created with 'hub-pr checkout' with a
// default commit message including Pull Request number and title, similar to
// the GitHub Merge Button.
func doMerge(flags *flag.FlagSet, args []string) error {
	flags.Parse(args)

	if flags.NArg() < 1 {
		return cli.ErrUsage
	}

	cli, proj, err := setup()
	if err != nil {
		return err
	}

	branch := flags.Arg(0)

	pr, err := corrPullRequest(cli, proj, branch)
	if err != nil {
		return err
	}

	mergeHead := fmt.Sprintf("%s/%s", pr.Head.Repo.Owner.Login, pr.Head.Ref)
	message := fmt.Sprintf("Merge pull request #%d from %s\n\n%s", pr.Number, mergeHead, pr.Title)

	return git.Run("merge", "--no-ff", "--edit", "-m", message, branch)
}

// +command browse - Open pull request page with browser
//
// 	browse
//
// Opens a web browser for the URL of Pull Request corresponding to current
// branch.
func doBrowse(flags *flag.FlagSet, args []string) error {
	flags.Parse(args)

	cli, proj, err := setup()
	if err != nil {
		return err
	}

	branch, err := git.Head()
	if err != nil {
		return err
	}

	branch = strings.TrimPrefix(branch, "refs/heads/")

	pr, err := corrPullRequest(cli, proj, branch)
	if err != nil {
		return err
	}

	return git.Run("web--browse", pr.HTMLURL)
}

// +command diff - Show pull request's diff
//
// 	diff
//
// Runs "git diff" between the base and head of Pull Request corresponding to
// current branch.
func doDiff(flags *flag.FlagSet, args []string) error {
	flags.Parse(args)

	cli, proj, err := setup()
	if err != nil {
		return err
	}

	branch, err := git.Head()
	if err != nil {
		return err
	}

	branch = strings.TrimPrefix(branch, "refs/heads/")

	pr, err := corrPullRequest(cli, proj, branch)
	if err != nil {
		return err
	}

	return git.Run("diff", fmt.Sprintf("%s...%s", pr.Base.Sha, pr.Head.Sha))
}

func corrPullRequest(cli *github.Client, proj *github.Project, branch string) (*octokit.PullRequest, error) {
	prNumber, err := git.Config("branch." + branch + ".prNumber")
	if err != nil {
		return nil, err
	}

	return cli.PullRequest(proj, prNumber)
}

type gitRunner struct {
	err error
}

func (r *gitRunner) git(commands ...string) {
	if r.err != nil {
		return
	}

	r.err = exec.Command("git", commands...).Run()
}

func (r *gitRunner) resetError() {
	r.err = nil
}

func setup() (*github.Client, *github.Project, error) {
	repo, err := github.LocalRepo()
	if err != nil {
		return nil, nil, err
	}

	proj, err := repo.MainProject()

	cli := github.NewClient(github.GitHubHost)

	return cli, proj, err
}

func main() {
	cli.Run(os.Args[1:])
}
