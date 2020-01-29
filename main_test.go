package main

import (
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func TestGHClient(t *testing.T) {
	gh := ghClient("TestOrg/TestRepo", "abc123")
	assert.Equal(t, gh.owner, "TestOrg")
	assert.Equal(t, gh.repo, "TestRepo")
	assert.NotNil(t, gh.ctx)
	// *github.Client does not expose its underlying *http.Client, so it's not
	// possible to assert that it's using an *oauth2.Transport with the
	// appropriate token.
}

func TestDoTheThing(t *testing.T) {
	client := &testClient{}
	intPtr := func(i int) *int { return &i }

	// no err for no-op
	assert.NoError(t, doTheThing(client))

	// closes a milestone that is past due and has no open issues
	client.milestones = append(client.milestones, &github.Milestone{
		Number:     intPtr(1),
		OpenIssues: intPtr(0),
		DueOn:      func(t time.Time) *time.Time { return &t }(time.Now().Add(-1 * time.Hour)),
		State:      func(s string) *string { return &s }("open"),
	})
	assert.NoError(t, doTheThing(client))
	assert.Equal(t, "closed", *client.milestones[0].State)

	// does not close a milestone that has no open issues but is not past due
	client.milestones = append(client.milestones, &github.Milestone{
		Number:     intPtr(2),
		OpenIssues: intPtr(0),
		DueOn:      func(t time.Time) *time.Time { return &t }(time.Now().Add(time.Hour)),
		State:      func(s string) *string { return &s }("open"),
	})
	assert.NoError(t, doTheThing(client))
	assert.Equal(t, "open", *client.milestones[1].State)

	// does not close a milestone that is past due but has open issues
	client.milestones = append(client.milestones, &github.Milestone{
		Number:     intPtr(3),
		OpenIssues: intPtr(1),
		DueOn:      func(t time.Time) *time.Time { return &t }(time.Now().Add(-1 * time.Hour)),
		State:      func(s string) *string { return &s }("open"),
	})
	assert.NoError(t, doTheThing(client))
	assert.Equal(t, "open", *client.milestones[2].State)
}

var _ ifi = &testClient{}

type testClient struct {
	milestones []*github.Milestone
}

func (c *testClient) ListMilestones(_ *github.MilestoneListOptions) ([]*github.Milestone, *github.Response, error) {
	return c.milestones, new(github.Response), nil
}

func (c *testClient) EditMilestone(m *github.Milestone) (*github.Milestone, *github.Response, error) {
	for _, rec := range c.milestones {
		if rec.Number == nil || m.Number == nil || *rec.Number != *m.Number {
			continue
		}
		*rec = *m
	}
	return m, new(github.Response), nil
}
