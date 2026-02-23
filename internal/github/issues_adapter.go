package github

import (
	"context"
	"fmt"

	gh "github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"

	"orion-cli/internal/core"
)

type IssuesAdapter struct {
	client *gh.Client
}

func NewIssuesAdapter(token string) *IssuesAdapter {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	return &IssuesAdapter{client: gh.NewClient(httpClient)}
}

func (a *IssuesAdapter) CreateIssues(owner, repo string, drafts []core.IssueDraft) ([]core.CreatedIssue, error) {
	result := make([]core.CreatedIssue, 0, len(drafts))

	for _, draft := range drafts {
		issueReq := &gh.IssueRequest{
			Title:  gh.Ptr(draft.Title),
			Body:   gh.Ptr(draft.Body),
			Labels: &draft.Labels,
		}

		issue, _, err := a.client.Issues.Create(context.Background(), owner, repo, issueReq)
		if err != nil {
			return nil, fmt.Errorf("issue oluşturulamadı (%s): %w", draft.Title, err)
		}

		result = append(result, core.CreatedIssue{
			Number: issue.GetNumber(),
			URL:    issue.GetHTMLURL(),
			Title:  issue.GetTitle(),
		})
	}

	return result, nil
}
