package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v48/github"
	"go.arsenm.dev/lure-repo-bot/internal/analyze"
	"go.arsenm.dev/lure-repo-bot/internal/shutils"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func startWebhookWorkers(ctx context.Context, jobQueue prQueue) {
	client := newClient(ctx, os.Getenv("LURE_BOT_GITHUB_TOKEN"))

	for i := 0; i < runtime.NumCPU(); i++ {
		go startWebhookWorker(ctx, jobQueue, client)
	}
}

func startWebhookWorker(ctx context.Context, jobQueue prQueue, client *github.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		case payload := <-jobQueue.Channel():
			if payload.Action == "opened" || payload.Action == "ready_for_review" || payload.Action == "review_requested" {
				if payload.PullRequest.Draft {
					continue
				}

				// Check if review was requested from the bot
				if payload.Action == "review_requested" {
					user, _, err := client.Users.Get(ctx, "")
					if err != nil {
						log.Println("Error getting github user:", err)
						continue
					}

					found := false
					for _, reviewer := range payload.PullRequest.RequestedReviewers {
						if reviewer.ID == *user.ID {
							found = true
							break
						}
					}

					if !found {
						continue
					}
				}

				fls, _, err := client.PullRequests.ListFiles(
					ctx,
					payload.PullRequest.Base.Repo.Owner.Login,
					payload.PullRequest.Base.Repo.Name,
					int(payload.PullRequest.Number),
					nil,
				)
				if err != nil {
					log.Println("Error listing PR files:", err)
					continue
				}

				head := payload.PullRequest.Head

				tmpdir, err := os.MkdirTemp("/tmp", "lure-repo-bot.*")
				if err != nil {
					log.Println("Error creating temporary directory:", err)
					continue
				}

				r, err := git.PlainClone(tmpdir, true, &git.CloneOptions{URL: head.Repo.HTMLURL})
				if err != nil {
					log.Println("Error cloning git repo:", err)
					continue
				}

				files, err := getFiles(r, head.Sha)
				if err != nil {
					log.Println("Error getting files in repo:", err)
					continue
				}

				err = files.ForEach(func(f *object.File) error {
					if !strings.Contains(f.Name, "lure.sh") {
						return nil
					}

					if !fileInPR(fls, f) {
						return nil
					}

					fmt.Println(f.Name, fileInPR(fls, f))

					r, err := f.Reader()
					if err != nil {
						return err
					}

					fl, err := syntax.NewParser().Parse(r, "lure.sh")
					if err != nil {
						return err
					}

					var nopRWC shutils.NopRWC
					runner, err := interp.New(
						interp.Env(expand.ListEnviron()),
						interp.StdIO(nopRWC, nopRWC, os.Stderr),
						interp.ExecHandler(shutils.NopExec),
						interp.ReadDirHandler(shutils.NopReadDir),
						interp.OpenHandler(shutils.NopOpen),
						interp.StatHandler(shutils.NopStat),
					)
					if err != nil {
						return err
					}

					err = runner.Run(ctx, fl)
					if err != nil {
						return err
					}

					findings, err := analyze.AnalyzeScript(runner, fl)
					if err != nil {
						return err
					}

					return writeFindings(ctx, client, findings, f.Name, &payload.PullRequest)
				})
				if err != nil {
					log.Println("Error analyzing files:", err)
					continue
				}

				err = os.RemoveAll(tmpdir)
				if err != nil {
					log.Println("Error removing temporary directory:", err)
					continue
				}
			}
		}
	}
}

func fileInPR(prFiles []*github.CommitFile, file *object.File) bool {
	for _, prFile := range prFiles {
		if *prFile.Filename == file.Name {
			return true
		}
	}
	return false
}

func getFiles(r *git.Repository, sha string) (*object.FileIter, error) {
	co, err := r.CommitObject(plumbing.NewHash(sha))
	if err != nil {
		return nil, err
	}

	tree, err := co.Tree()
	if err != nil {
		return nil, err
	}

	return tree.Files(), nil
}
