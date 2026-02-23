package core

type TaskAnalyzer interface {
	Analyze(text string) ([]IssueDraft, error)
}

type IssueCreator interface {
	CreateIssues(owner, repo string, drafts []IssueDraft) ([]CreatedIssue, error)
}

type RepoLocator interface {
	DetectOwnerRepo() (string, string, error)
}
