package core

type IssueDraft struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

type CreatedIssue struct {
	Number int
	URL    string
	Title  string
}
