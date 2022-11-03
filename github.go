package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v48/github"
	"go.arsenm.dev/lure-repo-bot/internal/analyze"
	"go.arsenm.dev/lure-repo-bot/internal/types"
	"golang.org/x/oauth2"
)

func newClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func writeFindings(ctx context.Context, c *github.Client, findings []analyze.Finding, path string, pr *types.PullRequest) error {
	comments := make([]*github.DraftReviewComment, len(findings))

	for i, finding := range findings {
		var name string
		if finding.Index != nil {
			name = fmt.Sprintf(
				"`%s[%v]` %s",
				finding.ItemName,
				finding.Index,
				finding.ItemType,
			)
		} else {
			name = fmt.Sprintf(
				"`%s` %s",
				finding.ItemName,
				finding.ItemType,
			)
		}

		msg := fmt.Sprintf(finding.Msg, name)

		if finding.ExtraMsg != "" {
			msg += "\n\n" + finding.ExtraMsg
		}

		if finding.Line == 0 {
			finding.Line = 1
		}

		fmt.Println(finding.Line)

		comments[i] = &github.DraftReviewComment{
			Line: github.Int(int(finding.Line)),
			Path: github.String(path),
			Body: github.String(msg),
			Side: github.String("RIGHT"),
		}
	}

	if len(findings) > 0 {
		rev, _, err := c.PullRequests.CreateReview(
			ctx,
			pr.Base.Repo.Owner.Login,
			pr.Base.Repo.Name,
			int(pr.Number),
			&github.PullRequestReviewRequest{
				Comments: comments,
			},
		)
		if err != nil {
			return err
		}

		_, _, err = c.PullRequests.SubmitReview(
			ctx,
			pr.Base.Repo.Owner.Login,
			pr.Base.Repo.Name,
			int(pr.Number),
			*rev.ID,
			&github.PullRequestReviewRequest{
				Body:  github.String("Please re-request review from the bot after applying these fixes"),
				Event: github.String("COMMENT"),
			},
		)
		if err != nil {
			return err
		}
	} else {
		rev, _, err := c.PullRequests.CreateReview(
			ctx,
			pr.Base.Repo.Owner.Login,
			pr.Base.Repo.Name,
			int(pr.Number),
			&github.PullRequestReviewRequest{},
		)
		if err != nil {
			return err
		}

		_, _, err = c.PullRequests.SubmitReview(
			ctx,
			pr.Base.Repo.Owner.Login,
			pr.Base.Repo.Name,
			int(pr.Number),
			*rev.ID,
			&github.PullRequestReviewRequest{
				Body:  github.String("No issues found!"),
				Event: github.String("APPROVE"),
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
