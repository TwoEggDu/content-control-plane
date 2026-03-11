package domain

import (
	"fmt"
	"strings"
	"time"
)

type Project struct {
	ID        int64
	Code      string
	Name      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ScanTaskStatus string

const (
	ScanTaskStatusReceived   ScanTaskStatus = "RECEIVED"
	ScanTaskStatusProcessing ScanTaskStatus = "PROCESSING"
	ScanTaskStatusImported   ScanTaskStatus = "IMPORTED"
	ScanTaskStatusFailed     ScanTaskStatus = "FAILED"
)

type ScanTask struct {
	ID                 int64
	ProjectID          int64
	ImportKey          string
	TaskNo             string
	SourceType         string
	BranchName         string
	CommitSHA          string
	ScannerVersion     string
	TriggeredBy        string
	StartedAt          time.Time
	FinishedAt         time.Time
	Status             ScanTaskStatus
	TotalIssueCount    int
	NewIssueCount      int
	ReopenedIssueCount int
	AttachmentCount    int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type ResourceItem struct {
	ID           int64
	ProjectID    int64
	ResourceKey  string
	ResourceGUID string
	ResourcePath string
	ResourceName string
	ResourceType string
	ModuleName   string
	OwnerName    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type IssueStatus string

const (
	IssueStatusNew      IssueStatus = "NEW"
	IssueStatusAssigned IssueStatus = "ASSIGNED"
	IssueStatusFixing   IssueStatus = "FIXING"
	IssueStatusResolved IssueStatus = "RESOLVED"
	IssueStatusVerified IssueStatus = "VERIFIED"
	IssueStatusIgnored  IssueStatus = "IGNORED"
)

type IssueItem struct {
	ID             int64
	ProjectID      int64
	ResourceID     int64
	LastScanTaskID int64
	Fingerprint    string
	RuleCode       string
	Severity       string
	Status         IssueStatus
	AssigneeName   string
	Message        string
	CurrentValue   string
	ExpectedValue  string
	LocationKey    string
	FirstSeenAt    time.Time
	LastSeenAt     time.Time
	ResolvedAt     *time.Time
	IgnoredAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type IssueActionLog struct {
	ID           int64
	IssueID      int64
	ActionType   string
	FromStatus   IssueStatus
	ToStatus     IssueStatus
	OperatorName string
	Comment      string
	CreatedAt    time.Time
}

type AttachmentFile struct {
	ID          int64
	ScanTaskID  int64
	FileType    string
	FileName    string
	StorageKey  string
	ContentHash string
	FileSize    int64
	CreatedAt   time.Time
}

func ParseIssueStatus(raw string) (IssueStatus, error) {
	normalized := IssueStatus(strings.ToUpper(strings.TrimSpace(raw)))
	switch normalized {
	case IssueStatusNew, IssueStatusAssigned, IssueStatusFixing, IssueStatusResolved, IssueStatusVerified, IssueStatusIgnored:
		return normalized, nil
	default:
		return "", fmt.Errorf("unknown issue status: %s", raw)
	}
}

func (s IssueStatus) CanTransitionTo(next IssueStatus) bool {
	if s == next {
		return true
	}

	switch s {
	case IssueStatusNew:
		return next == IssueStatusAssigned || next == IssueStatusFixing || next == IssueStatusResolved || next == IssueStatusIgnored
	case IssueStatusAssigned:
		return next == IssueStatusFixing || next == IssueStatusResolved || next == IssueStatusIgnored
	case IssueStatusFixing:
		return next == IssueStatusResolved || next == IssueStatusIgnored
	case IssueStatusResolved:
		return next == IssueStatusVerified || next == IssueStatusFixing
	case IssueStatusVerified:
		return next == IssueStatusFixing
	case IssueStatusIgnored:
		return next == IssueStatusNew || next == IssueStatusAssigned || next == IssueStatusFixing
	default:
		return false
	}
}
