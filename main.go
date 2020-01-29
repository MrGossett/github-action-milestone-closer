package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type config struct {
	Token      string `envconfig:"GITHUB_TOKEN" required:"true"`
	Repository string `envconfig:"GITHUB_REPOSITORY" required:"true"`
}

func main() {
	var c config
	if err := envconfig.Process("INPUT", &c); err != nil {
		log.Fatal(err)
	}

	client := ghClient(c.Repository, c.Token)
	if err := doTheThing(client); err != nil {
		log.Fatal(err)
	}
}

func ghClient(repo, token string) *gh {
	parts := strings.SplitN(repo, "/", 2)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return &gh{
		Client: github.NewClient(tc),
		owner:  parts[0],
		repo:   parts[1],
		ctx:    ctx,
	}
}

func doTheThing(client ifi) error {
	milestones, _, err := client.ListMilestones(&github.MilestoneListOptions{
		State:     "open",
		Direction: "desc",
	})
	if err != nil {
		return errors.Wrap(err, "could not list milestones")
	}

	for _, m := range milestones {
		shouldClose := (*m.OpenIssues == 0 && m.DueOn.Before(time.Now()) && m.ClosedAt == nil)
		if !shouldClose {
			continue
		}

		m.State = closed
		if _, _, err := client.EditMilestone(m); err != nil {
			return errors.Wrap(err, "could not close milestone")
		}
	}

	return nil
}

var closed = func(s string) *string { return &s }("closed")

type ifi interface {
	ListMilestones(*github.MilestoneListOptions) ([]*github.Milestone, *github.Response, error)
	EditMilestone(*github.Milestone) (*github.Milestone, *github.Response, error)
}

type gh struct {
	*github.Client
	owner, repo string
	ctx         context.Context
}

func (c *gh) ListMilestones(opts *github.MilestoneListOptions) ([]*github.Milestone, *github.Response, error) {
	return c.Issues.ListMilestones(c.ctx, c.owner, c.repo, opts)
}

func (c *gh) EditMilestone(m *github.Milestone) (*github.Milestone, *github.Response, error) {
	return c.Issues.EditMilestone(c.ctx, c.owner, c.repo, *m.Number, m)
}
