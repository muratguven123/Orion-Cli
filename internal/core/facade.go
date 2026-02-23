package core

import "fmt"

type TaskFacade struct {
	Analyzer TaskAnalyzer
	Creator  IssueCreator
	Locator  RepoLocator
}

func NewTaskFacade(analyzer TaskAnalyzer, creator IssueCreator, locator RepoLocator) *TaskFacade {
	return &TaskFacade{Analyzer: analyzer, Creator: creator, Locator: locator}
}

func (f *TaskFacade) Execute(taskText string) ([]CreatedIssue, error) {
	if taskText == "" {
		return nil, fmt.Errorf("görev metni boş olamaz")
	}

	owner, repo, err := f.Locator.DetectOwnerRepo()
	if err != nil {
		return nil, err
	}

	drafts, err := f.Analyzer.Analyze(taskText)
	if err != nil {
		return nil, err
	}

	created, err := f.Creator.CreateIssues(owner, repo, drafts)
	if err != nil {
		return nil, err
	}

	return created, nil
}
