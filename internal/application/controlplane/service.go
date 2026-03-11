package controlplane

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/TwoEggDu/content-control-plane/internal/domain"
)

var (
	ErrInvalidRequest    = errors.New("invalid_request")
	ErrNotFound          = errors.New("not_found")
	ErrInvalidTransition = errors.New("invalid_transition")
)

type Repository interface {
	Ready(ctx context.Context) error
	CreateOrGetProjectByCode(ctx context.Context, code, name string, now time.Time) (*domain.Project, error)
	GetProject(ctx context.Context, id int64) (*domain.Project, error)
	FindScanTaskByImportKey(ctx context.Context, importKey string) (*domain.ScanTask, error)
	CreateScanTask(ctx context.Context, task *domain.ScanTask) (*domain.ScanTask, error)
	UpdateScanTask(ctx context.Context, task *domain.ScanTask) (*domain.ScanTask, error)
	ListScanTasks(ctx context.Context, filter domain.ScanTaskFilter) ([]*domain.ScanTask, error)
	GetScanTask(ctx context.Context, id int64) (*domain.ScanTask, error)
	CreateAttachment(ctx context.Context, attachment *domain.AttachmentFile) (*domain.AttachmentFile, error)
	ListAttachmentsByScanTask(ctx context.Context, scanTaskID int64) ([]*domain.AttachmentFile, error)
	CreateOrGetResource(ctx context.Context, resource *domain.ResourceItem) (*domain.ResourceItem, error)
	FindIssueByFingerprint(ctx context.Context, projectID int64, fingerprint string) (*domain.IssueItem, error)
	CreateIssue(ctx context.Context, issue *domain.IssueItem) (*domain.IssueItem, error)
	UpdateIssue(ctx context.Context, issue *domain.IssueItem) (*domain.IssueItem, error)
	ListIssues(ctx context.Context, filter domain.IssueFilter) ([]*domain.IssueItem, error)
	GetIssue(ctx context.Context, id int64) (*domain.IssueItem, error)
	GetResource(ctx context.Context, id int64) (*domain.ResourceItem, error)
	AppendIssueAction(ctx context.Context, action *domain.IssueActionLog) (*domain.IssueActionLog, error)
	ListIssueActions(ctx context.Context, issueID int64) ([]*domain.IssueActionLog, error)
}

type Service struct {
	repo Repository
	now  func() time.Time
}

type ImportAttachmentInput struct {
	FileType    string `json:"file_type"`
	FileName    string `json:"file_name"`
	StorageKey  string `json:"storage_key"`
	ContentHash string `json:"content_hash"`
	FileSize    int64  `json:"file_size"`
}

type ImportIssueInput struct {
	RuleCode      string `json:"rule_code"`
	Severity      string `json:"severity"`
	Message       string `json:"message"`
	ResourceGUID  string `json:"resource_guid"`
	ResourcePath  string `json:"resource_path"`
	ResourceName  string `json:"resource_name"`
	ResourceType  string `json:"resource_type"`
	ModuleName    string `json:"module_name"`
	OwnerName     string `json:"owner_name"`
	LocationKey   string `json:"location_key"`
	CurrentValue  string `json:"current_value"`
	ExpectedValue string `json:"expected_value"`
}

type ImportScanRequest struct {
	ProjectCode    string                  `json:"project_code"`
	ProjectName    string                  `json:"project_name"`
	SourceType     string                  `json:"source_type"`
	TaskNo         string                  `json:"task_no"`
	BranchName     string                  `json:"branch_name"`
	CommitSHA      string                  `json:"commit_sha"`
	ScannerVersion string                  `json:"scanner_version"`
	TriggeredBy    string                  `json:"triggered_by"`
	StartedAt      string                  `json:"started_at"`
	FinishedAt     string                  `json:"finished_at"`
	Attachments    []ImportAttachmentInput `json:"attachments"`
	Issues         []ImportIssueInput      `json:"issues"`
}

type ImportScanResult struct {
	ScanTaskID         int64  `json:"scan_task_id"`
	ImportKey          string `json:"import_key"`
	Status             string `json:"status"`
	Created            bool   `json:"created"`
	Reused             bool   `json:"reused"`
	TotalIssueCount    int    `json:"total_issue_count"`
	NewIssueCount      int    `json:"new_issue_count"`
	ReopenedIssueCount int    `json:"reopened_issue_count"`
}

type AttachmentSummary struct {
	ID          int64     `json:"id"`
	FileType    string    `json:"file_type"`
	FileName    string    `json:"file_name"`
	StorageKey  string    `json:"storage_key"`
	ContentHash string    `json:"content_hash"`
	FileSize    int64     `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
}

type ScanTaskSummary struct {
	ID                 int64     `json:"id"`
	ProjectCode        string    `json:"project_code"`
	TaskNo             string    `json:"task_no"`
	SourceType         string    `json:"source_type"`
	BranchName         string    `json:"branch_name"`
	CommitSHA          string    `json:"commit_sha"`
	ScannerVersion     string    `json:"scanner_version"`
	Status             string    `json:"status"`
	TotalIssueCount    int       `json:"total_issue_count"`
	NewIssueCount      int       `json:"new_issue_count"`
	ReopenedIssueCount int       `json:"reopened_issue_count"`
	AttachmentCount    int       `json:"attachment_count"`
	StartedAt          time.Time `json:"started_at"`
	FinishedAt         time.Time `json:"finished_at"`
}

type ScanTaskDetail struct {
	ID                 int64               `json:"id"`
	ProjectCode        string              `json:"project_code"`
	TaskNo             string              `json:"task_no"`
	SourceType         string              `json:"source_type"`
	BranchName         string              `json:"branch_name"`
	CommitSHA          string              `json:"commit_sha"`
	ScannerVersion     string              `json:"scanner_version"`
	TriggeredBy        string              `json:"triggered_by"`
	Status             string              `json:"status"`
	TotalIssueCount    int                 `json:"total_issue_count"`
	NewIssueCount      int                 `json:"new_issue_count"`
	ReopenedIssueCount int                 `json:"reopened_issue_count"`
	AttachmentCount    int                 `json:"attachment_count"`
	StartedAt          time.Time           `json:"started_at"`
	FinishedAt         time.Time           `json:"finished_at"`
	Attachments        []AttachmentSummary `json:"attachments"`
}

type IssueSummary struct {
	ID           int64     `json:"id"`
	ProjectCode  string    `json:"project_code"`
	ScanTaskID   int64     `json:"scan_task_id"`
	ResourceID   int64     `json:"resource_id"`
	RuleCode     string    `json:"rule_code"`
	Severity     string    `json:"severity"`
	Status       string    `json:"status"`
	AssigneeName string    `json:"assignee_name"`
	Message      string    `json:"message"`
	ResourcePath string    `json:"resource_path"`
	ResourceGUID string    `json:"resource_guid"`
	FirstSeenAt  time.Time `json:"first_seen_at"`
	LastSeenAt   time.Time `json:"last_seen_at"`
}

type ResourceSummary struct {
	ID           int64  `json:"id"`
	ResourceGUID string `json:"resource_guid"`
	ResourcePath string `json:"resource_path"`
	ResourceName string `json:"resource_name"`
	ResourceType string `json:"resource_type"`
	ModuleName   string `json:"module_name"`
	OwnerName    string `json:"owner_name"`
}

type IssueActionRecord struct {
	ID           int64     `json:"id"`
	ActionType   string    `json:"action_type"`
	FromStatus   string    `json:"from_status"`
	ToStatus     string    `json:"to_status"`
	OperatorName string    `json:"operator_name"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
}

type IssueDetail struct {
	ID            int64               `json:"id"`
	ProjectCode   string              `json:"project_code"`
	ScanTaskID    int64               `json:"scan_task_id"`
	Resource      ResourceSummary     `json:"resource"`
	RuleCode      string              `json:"rule_code"`
	Severity      string              `json:"severity"`
	Status        string              `json:"status"`
	AssigneeName  string              `json:"assignee_name"`
	Message       string              `json:"message"`
	CurrentValue  string              `json:"current_value"`
	ExpectedValue string              `json:"expected_value"`
	LocationKey   string              `json:"location_key"`
	FirstSeenAt   time.Time           `json:"first_seen_at"`
	LastSeenAt    time.Time           `json:"last_seen_at"`
	ResolvedAt    *time.Time          `json:"resolved_at"`
	Actions       []IssueActionRecord `json:"actions"`
}

type UpdateIssueStatusRequest struct {
	ToStatus     string `json:"to_status"`
	AssigneeName string `json:"assignee_name"`
	OperatorName string `json:"operator_name"`
	Comment      string `json:"comment"`
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
}

func (s *Service) Ready(ctx context.Context) error {
	return s.repo.Ready(ctx)
}

func (s *Service) ImportScan(ctx context.Context, request ImportScanRequest) (*ImportScanResult, error) {
	if err := validateImportRequest(request); err != nil {
		return nil, err
	}

	importKey := buildImportKey(request)
	existingTask, err := s.repo.FindScanTaskByImportKey(ctx, importKey)
	if err != nil {
		return nil, err
	}
	if existingTask != nil {
		return &ImportScanResult{
			ScanTaskID:         existingTask.ID,
			ImportKey:          existingTask.ImportKey,
			Status:             string(existingTask.Status),
			Created:            false,
			Reused:             true,
			TotalIssueCount:    existingTask.TotalIssueCount,
			NewIssueCount:      existingTask.NewIssueCount,
			ReopenedIssueCount: existingTask.ReopenedIssueCount,
		}, nil
	}

	now := s.now().UTC()
	project, err := s.repo.CreateOrGetProjectByCode(ctx, strings.TrimSpace(request.ProjectCode), strings.TrimSpace(request.ProjectName), now)
	if err != nil {
		return nil, err
	}

	startedAt := parseTimeOrDefault(request.StartedAt, now)
	finishedAt := parseTimeOrDefault(request.FinishedAt, now)

	scanTask, err := s.repo.CreateScanTask(ctx, &domain.ScanTask{
		ProjectID:      project.ID,
		ImportKey:      importKey,
		TaskNo:         strings.TrimSpace(request.TaskNo),
		SourceType:     strings.TrimSpace(request.SourceType),
		BranchName:     strings.TrimSpace(request.BranchName),
		CommitSHA:      strings.TrimSpace(request.CommitSHA),
		ScannerVersion: strings.TrimSpace(request.ScannerVersion),
		TriggeredBy:    strings.TrimSpace(request.TriggeredBy),
		StartedAt:      startedAt,
		FinishedAt:     finishedAt,
		Status:         domain.ScanTaskStatusProcessing,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	if err != nil {
		return nil, err
	}

	for _, attachment := range request.Attachments {
		_, err := s.repo.CreateAttachment(ctx, &domain.AttachmentFile{
			ScanTaskID:  scanTask.ID,
			FileType:    strings.TrimSpace(attachment.FileType),
			FileName:    strings.TrimSpace(attachment.FileName),
			StorageKey:  strings.TrimSpace(attachment.StorageKey),
			ContentHash: strings.TrimSpace(attachment.ContentHash),
			FileSize:    attachment.FileSize,
			CreatedAt:   now,
		})
		if err != nil {
			return nil, err
		}
	}

	newIssueCount := 0
	reopenedIssueCount := 0
	observedAt := finishedAt
	if observedAt.IsZero() {
		observedAt = now
	}

	for _, item := range request.Issues {
		resourceKey := domain.StableResourceKey(item.ResourceGUID, item.ResourcePath)
		resource, err := s.repo.CreateOrGetResource(ctx, &domain.ResourceItem{
			ProjectID:    project.ID,
			ResourceKey:  resourceKey,
			ResourceGUID: strings.TrimSpace(item.ResourceGUID),
			ResourcePath: strings.TrimSpace(item.ResourcePath),
			ResourceName: strings.TrimSpace(item.ResourceName),
			ResourceType: strings.TrimSpace(item.ResourceType),
			ModuleName:   strings.TrimSpace(item.ModuleName),
			OwnerName:    strings.TrimSpace(item.OwnerName),
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		if err != nil {
			return nil, err
		}

		fingerprint := domain.GenerateIssueFingerprint(
			request.ProjectCode,
			item.RuleCode,
			resource.ResourceKey,
			item.LocationKey,
			item.CurrentValue,
			item.ExpectedValue,
			item.Message,
		)

		existingIssue, err := s.repo.FindIssueByFingerprint(ctx, project.ID, fingerprint)
		if err != nil {
			return nil, err
		}

		if existingIssue == nil {
			_, err = s.repo.CreateIssue(ctx, &domain.IssueItem{
				ProjectID:      project.ID,
				ResourceID:     resource.ID,
				LastScanTaskID: scanTask.ID,
				Fingerprint:    fingerprint,
				RuleCode:       strings.TrimSpace(item.RuleCode),
				Severity:       strings.ToUpper(strings.TrimSpace(item.Severity)),
				Status:         domain.IssueStatusNew,
				Message:        strings.TrimSpace(item.Message),
				CurrentValue:   strings.TrimSpace(item.CurrentValue),
				ExpectedValue:  strings.TrimSpace(item.ExpectedValue),
				LocationKey:    strings.TrimSpace(item.LocationKey),
				FirstSeenAt:    observedAt,
				LastSeenAt:     observedAt,
				CreatedAt:      now,
				UpdatedAt:      now,
			})
			if err != nil {
				return nil, err
			}
			newIssueCount++
			continue
		}

		previousStatus := existingIssue.Status
		existingIssue.ResourceID = resource.ID
		existingIssue.LastScanTaskID = scanTask.ID
		existingIssue.RuleCode = strings.TrimSpace(item.RuleCode)
		existingIssue.Severity = strings.ToUpper(strings.TrimSpace(item.Severity))
		existingIssue.Message = strings.TrimSpace(item.Message)
		existingIssue.CurrentValue = strings.TrimSpace(item.CurrentValue)
		existingIssue.ExpectedValue = strings.TrimSpace(item.ExpectedValue)
		existingIssue.LocationKey = strings.TrimSpace(item.LocationKey)
		existingIssue.LastSeenAt = observedAt
		existingIssue.UpdatedAt = now

		if previousStatus == domain.IssueStatusResolved || previousStatus == domain.IssueStatusVerified {
			existingIssue.Status = domain.IssueStatusNew
			existingIssue.ResolvedAt = nil
			reopenedIssueCount++
			_, err = s.repo.AppendIssueAction(ctx, &domain.IssueActionLog{
				IssueID:      existingIssue.ID,
				ActionType:   "REOPENED_BY_IMPORT",
				FromStatus:   previousStatus,
				ToStatus:     domain.IssueStatusNew,
				OperatorName: "system/import",
				Comment:      "issue observed again during import",
				CreatedAt:    now,
			})
			if err != nil {
				return nil, err
			}
		}

		if _, err := s.repo.UpdateIssue(ctx, existingIssue); err != nil {
			return nil, err
		}
	}

	scanTask.Status = domain.ScanTaskStatusImported
	scanTask.TotalIssueCount = len(request.Issues)
	scanTask.NewIssueCount = newIssueCount
	scanTask.ReopenedIssueCount = reopenedIssueCount
	scanTask.AttachmentCount = len(request.Attachments)
	scanTask.FinishedAt = finishedAt
	scanTask.UpdatedAt = now

	if _, err := s.repo.UpdateScanTask(ctx, scanTask); err != nil {
		return nil, err
	}

	return &ImportScanResult{
		ScanTaskID:         scanTask.ID,
		ImportKey:          scanTask.ImportKey,
		Status:             string(scanTask.Status),
		Created:            true,
		Reused:             false,
		TotalIssueCount:    scanTask.TotalIssueCount,
		NewIssueCount:      scanTask.NewIssueCount,
		ReopenedIssueCount: scanTask.ReopenedIssueCount,
	}, nil
}

func (s *Service) ListScanTasks(ctx context.Context, filter domain.ScanTaskFilter) ([]ScanTaskSummary, error) {
	tasks, err := s.repo.ListScanTasks(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make([]ScanTaskSummary, 0, len(tasks))
	for _, task := range tasks {
		project, err := s.repo.GetProject(ctx, task.ProjectID)
		if err != nil {
			return nil, err
		}
		result = append(result, toScanTaskSummary(project, task))
	}
	return result, nil
}

func (s *Service) GetScanTaskDetail(ctx context.Context, id int64) (*ScanTaskDetail, error) {
	task, err := s.repo.GetScanTask(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrNotFound
	}

	project, err := s.repo.GetProject(ctx, task.ProjectID)
	if err != nil {
		return nil, err
	}

	attachments, err := s.repo.ListAttachmentsByScanTask(ctx, task.ID)
	if err != nil {
		return nil, err
	}

	detail := &ScanTaskDetail{
		ID:                 task.ID,
		ProjectCode:        project.Code,
		TaskNo:             task.TaskNo,
		SourceType:         task.SourceType,
		BranchName:         task.BranchName,
		CommitSHA:          task.CommitSHA,
		ScannerVersion:     task.ScannerVersion,
		TriggeredBy:        task.TriggeredBy,
		Status:             string(task.Status),
		TotalIssueCount:    task.TotalIssueCount,
		NewIssueCount:      task.NewIssueCount,
		ReopenedIssueCount: task.ReopenedIssueCount,
		AttachmentCount:    task.AttachmentCount,
		StartedAt:          task.StartedAt,
		FinishedAt:         task.FinishedAt,
		Attachments:        make([]AttachmentSummary, 0, len(attachments)),
	}

	for _, attachment := range attachments {
		detail.Attachments = append(detail.Attachments, AttachmentSummary{
			ID:          attachment.ID,
			FileType:    attachment.FileType,
			FileName:    attachment.FileName,
			StorageKey:  attachment.StorageKey,
			ContentHash: attachment.ContentHash,
			FileSize:    attachment.FileSize,
			CreatedAt:   attachment.CreatedAt,
		})
	}

	return detail, nil
}

func (s *Service) ListIssues(ctx context.Context, filter domain.IssueFilter) ([]IssueSummary, error) {
	issues, err := s.repo.ListIssues(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make([]IssueSummary, 0, len(issues))
	for _, issue := range issues {
		project, err := s.repo.GetProject(ctx, issue.ProjectID)
		if err != nil {
			return nil, err
		}
		resource, err := s.repo.GetResource(ctx, issue.ResourceID)
		if err != nil {
			return nil, err
		}
		result = append(result, IssueSummary{
			ID:           issue.ID,
			ProjectCode:  project.Code,
			ScanTaskID:   issue.LastScanTaskID,
			ResourceID:   issue.ResourceID,
			RuleCode:     issue.RuleCode,
			Severity:     issue.Severity,
			Status:       string(issue.Status),
			AssigneeName: issue.AssigneeName,
			Message:      issue.Message,
			ResourcePath: resource.ResourcePath,
			ResourceGUID: resource.ResourceGUID,
			FirstSeenAt:  issue.FirstSeenAt,
			LastSeenAt:   issue.LastSeenAt,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].LastSeenAt.Equal(result[j].LastSeenAt) {
			return result[i].ID > result[j].ID
		}
		return result[i].LastSeenAt.After(result[j].LastSeenAt)
	})

	return result, nil
}

func (s *Service) GetIssueDetail(ctx context.Context, id int64) (*IssueDetail, error) {
	issue, err := s.repo.GetIssue(ctx, id)
	if err != nil {
		return nil, err
	}
	if issue == nil {
		return nil, ErrNotFound
	}
	return s.buildIssueDetail(ctx, issue)
}

func (s *Service) UpdateIssueStatus(ctx context.Context, issueID int64, request UpdateIssueStatusRequest) (*IssueDetail, error) {
	issue, err := s.repo.GetIssue(ctx, issueID)
	if err != nil {
		return nil, err
	}
	if issue == nil {
		return nil, ErrNotFound
	}

	nextStatus, err := domain.ParseIssueStatus(request.ToStatus)
	if err != nil {
		return nil, fmt.Errorf("%w: to_status is invalid", ErrInvalidRequest)
	}
	operatorName := strings.TrimSpace(request.OperatorName)
	if operatorName == "" {
		return nil, fmt.Errorf("%w: operator_name is required", ErrInvalidRequest)
	}
	assigneeName := strings.TrimSpace(request.AssigneeName)
	if nextStatus == domain.IssueStatusAssigned && assigneeName == "" {
		return nil, fmt.Errorf("%w: assignee_name is required when assigning", ErrInvalidRequest)
	}
	if !issue.Status.CanTransitionTo(nextStatus) {
		return nil, fmt.Errorf("%w: cannot move from %s to %s", ErrInvalidTransition, issue.Status, nextStatus)
	}

	previousStatus := issue.Status
	issue.Status = nextStatus
	if assigneeName != "" {
		issue.AssigneeName = assigneeName
	}

	now := s.now().UTC()
	switch nextStatus {
	case domain.IssueStatusResolved:
		issue.ResolvedAt = &now
		issue.IgnoredAt = nil
	case domain.IssueStatusIgnored:
		issue.IgnoredAt = &now
	case domain.IssueStatusNew, domain.IssueStatusAssigned, domain.IssueStatusFixing, domain.IssueStatusVerified:
		if nextStatus != domain.IssueStatusResolved {
			issue.ResolvedAt = nil
		}
		if nextStatus != domain.IssueStatusIgnored {
			issue.IgnoredAt = nil
		}
	}
	issue.UpdatedAt = now

	if _, err := s.repo.UpdateIssue(ctx, issue); err != nil {
		return nil, err
	}

	actionType := "STATUS_CHANGED"
	if nextStatus == domain.IssueStatusAssigned && assigneeName != "" {
		actionType = "ASSIGNED"
	}
	if _, err := s.repo.AppendIssueAction(ctx, &domain.IssueActionLog{
		IssueID:      issue.ID,
		ActionType:   actionType,
		FromStatus:   previousStatus,
		ToStatus:     nextStatus,
		OperatorName: operatorName,
		Comment:      strings.TrimSpace(request.Comment),
		CreatedAt:    now,
	}); err != nil {
		return nil, err
	}

	return s.buildIssueDetail(ctx, issue)
}

func (s *Service) buildIssueDetail(ctx context.Context, issue *domain.IssueItem) (*IssueDetail, error) {
	project, err := s.repo.GetProject(ctx, issue.ProjectID)
	if err != nil {
		return nil, err
	}
	resource, err := s.repo.GetResource(ctx, issue.ResourceID)
	if err != nil {
		return nil, err
	}
	actions, err := s.repo.ListIssueActions(ctx, issue.ID)
	if err != nil {
		return nil, err
	}

	detail := &IssueDetail{
		ID:          issue.ID,
		ProjectCode: project.Code,
		ScanTaskID:  issue.LastScanTaskID,
		Resource: ResourceSummary{
			ID:           resource.ID,
			ResourceGUID: resource.ResourceGUID,
			ResourcePath: resource.ResourcePath,
			ResourceName: resource.ResourceName,
			ResourceType: resource.ResourceType,
			ModuleName:   resource.ModuleName,
			OwnerName:    resource.OwnerName,
		},
		RuleCode:      issue.RuleCode,
		Severity:      issue.Severity,
		Status:        string(issue.Status),
		AssigneeName:  issue.AssigneeName,
		Message:       issue.Message,
		CurrentValue:  issue.CurrentValue,
		ExpectedValue: issue.ExpectedValue,
		LocationKey:   issue.LocationKey,
		FirstSeenAt:   issue.FirstSeenAt,
		LastSeenAt:    issue.LastSeenAt,
		ResolvedAt:    issue.ResolvedAt,
		Actions:       make([]IssueActionRecord, 0, len(actions)),
	}

	for _, action := range actions {
		detail.Actions = append(detail.Actions, IssueActionRecord{
			ID:           action.ID,
			ActionType:   action.ActionType,
			FromStatus:   string(action.FromStatus),
			ToStatus:     string(action.ToStatus),
			OperatorName: action.OperatorName,
			Comment:      action.Comment,
			CreatedAt:    action.CreatedAt,
		})
	}

	return detail, nil
}

func validateImportRequest(request ImportScanRequest) error {
	requiredFields := map[string]string{
		"project_code":    request.ProjectCode,
		"source_type":     request.SourceType,
		"branch_name":     request.BranchName,
		"commit_sha":      request.CommitSHA,
		"scanner_version": request.ScannerVersion,
		"triggered_by":    request.TriggeredBy,
	}

	for name, value := range requiredFields {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%w: %s is required", ErrInvalidRequest, name)
		}
	}

	for index, issue := range request.Issues {
		if strings.TrimSpace(issue.RuleCode) == "" {
			return fmt.Errorf("%w: issues[%d].rule_code is required", ErrInvalidRequest, index)
		}
		if strings.TrimSpace(issue.Severity) == "" {
			return fmt.Errorf("%w: issues[%d].severity is required", ErrInvalidRequest, index)
		}
		if strings.TrimSpace(issue.Message) == "" {
			return fmt.Errorf("%w: issues[%d].message is required", ErrInvalidRequest, index)
		}
		if domain.StableResourceKey(issue.ResourceGUID, issue.ResourcePath) == "" {
			return fmt.Errorf("%w: issues[%d] requires resource_guid or resource_path", ErrInvalidRequest, index)
		}
	}

	return nil
}

func buildImportKey(request ImportScanRequest) string {
	if taskNo := strings.TrimSpace(request.TaskNo); taskNo != "" {
		return "task_no:" + taskNo
	}

	parts := []string{
		strings.TrimSpace(request.ProjectCode),
		strings.TrimSpace(request.SourceType),
		strings.TrimSpace(request.BranchName),
		strings.TrimSpace(request.CommitSHA),
		strings.TrimSpace(request.ScannerVersion),
		strconv.Itoa(len(request.Attachments)),
		strconv.Itoa(len(request.Issues)),
	}

	for _, issue := range request.Issues {
		parts = append(parts,
			strings.TrimSpace(issue.RuleCode),
			domain.StableResourceKey(issue.ResourceGUID, issue.ResourcePath),
			strings.TrimSpace(issue.LocationKey),
			strings.TrimSpace(issue.CurrentValue),
			strings.TrimSpace(issue.ExpectedValue),
			strings.TrimSpace(issue.Message),
		)
	}

	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return "derived:" + hex.EncodeToString(sum[:])
}

func parseTimeOrDefault(value string, fallback time.Time) time.Time {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return fallback
	}
	parsed, err := time.Parse(time.RFC3339, normalized)
	if err != nil {
		return fallback
	}
	return parsed.UTC()
}

func toScanTaskSummary(project *domain.Project, task *domain.ScanTask) ScanTaskSummary {
	return ScanTaskSummary{
		ID:                 task.ID,
		ProjectCode:        project.Code,
		TaskNo:             task.TaskNo,
		SourceType:         task.SourceType,
		BranchName:         task.BranchName,
		CommitSHA:          task.CommitSHA,
		ScannerVersion:     task.ScannerVersion,
		Status:             string(task.Status),
		TotalIssueCount:    task.TotalIssueCount,
		NewIssueCount:      task.NewIssueCount,
		ReopenedIssueCount: task.ReopenedIssueCount,
		AttachmentCount:    task.AttachmentCount,
		StartedAt:          task.StartedAt,
		FinishedAt:         task.FinishedAt,
	}
}
