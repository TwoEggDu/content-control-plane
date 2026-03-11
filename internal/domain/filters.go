package domain

type ScanTaskFilter struct {
	ProjectCode string
	Status      string
	SourceType  string
}

type IssueFilter struct {
	ProjectCode  string
	ScanTaskID   *int64
	Status       string
	Severity     string
	RuleCode     string
	AssigneeName string
	ResourcePath string
}
