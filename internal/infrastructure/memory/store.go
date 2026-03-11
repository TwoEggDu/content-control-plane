package memory

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/TwoEggDu/content-control-plane/internal/domain"
)

type Store struct {
	mu sync.RWMutex

	nextID int64

	projects      map[int64]*domain.Project
	projectByCode map[string]int64

	scanTasks           map[int64]*domain.ScanTask
	scanTaskByImportKey map[string]int64

	resources     map[int64]*domain.ResourceItem
	resourceByKey map[string]int64

	issues             map[int64]*domain.IssueItem
	issueByFingerprint map[string]int64

	attachments           map[int64]*domain.AttachmentFile
	attachmentsByScanTask map[int64][]int64

	actions        map[int64]*domain.IssueActionLog
	actionsByIssue map[int64][]int64
}

func NewStore() *Store {
	return &Store{
		nextID:                1,
		projects:              make(map[int64]*domain.Project),
		projectByCode:         make(map[string]int64),
		scanTasks:             make(map[int64]*domain.ScanTask),
		scanTaskByImportKey:   make(map[string]int64),
		resources:             make(map[int64]*domain.ResourceItem),
		resourceByKey:         make(map[string]int64),
		issues:                make(map[int64]*domain.IssueItem),
		issueByFingerprint:    make(map[string]int64),
		attachments:           make(map[int64]*domain.AttachmentFile),
		attachmentsByScanTask: make(map[int64][]int64),
		actions:               make(map[int64]*domain.IssueActionLog),
		actionsByIssue:        make(map[int64][]int64),
	}
}

func (s *Store) Ready(context.Context) error {
	return nil
}

func (s *Store) CreateOrGetProjectByCode(_ context.Context, code, name string, now time.Time) (*domain.Project, error) {
	normalizedCode := strings.TrimSpace(strings.ToLower(code))
	if normalizedCode == "" {
		return nil, fmt.Errorf("project code is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if id, ok := s.projectByCode[normalizedCode]; ok {
		project := s.projects[id]
		if strings.TrimSpace(name) != "" && project.Name != strings.TrimSpace(name) {
			project.Name = strings.TrimSpace(name)
			project.UpdatedAt = now
		}
		return cloneProject(project), nil
	}

	project := &domain.Project{
		ID:        s.allocateIDLocked(),
		Code:      strings.TrimSpace(code),
		Name:      strings.TrimSpace(name),
		Status:    "ACTIVE",
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.projects[project.ID] = project
	s.projectByCode[normalizedCode] = project.ID
	return cloneProject(project), nil
}

func (s *Store) GetProject(_ context.Context, id int64) (*domain.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	project, ok := s.projects[id]
	if !ok {
		return nil, fmt.Errorf("project not found: %d", id)
	}
	return cloneProject(project), nil
}

func (s *Store) FindScanTaskByImportKey(_ context.Context, importKey string) (*domain.ScanTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.scanTaskByImportKey[strings.TrimSpace(importKey)]
	if !ok {
		return nil, nil
	}
	return cloneScanTask(s.scanTasks[id]), nil
}

func (s *Store) CreateScanTask(_ context.Context, task *domain.ScanTask) (*domain.ScanTask, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cloned := cloneScanTask(task)
	cloned.ID = s.allocateIDLocked()
	s.scanTasks[cloned.ID] = cloned
	s.scanTaskByImportKey[strings.TrimSpace(cloned.ImportKey)] = cloned.ID
	return cloneScanTask(cloned), nil
}

func (s *Store) UpdateScanTask(_ context.Context, task *domain.ScanTask) (*domain.ScanTask, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.scanTasks[task.ID]; !ok {
		return nil, fmt.Errorf("scan task not found: %d", task.ID)
	}
	cloned := cloneScanTask(task)
	s.scanTasks[cloned.ID] = cloned
	s.scanTaskByImportKey[strings.TrimSpace(cloned.ImportKey)] = cloned.ID
	return cloneScanTask(cloned), nil
}

func (s *Store) ListScanTasks(_ context.Context, filter domain.ScanTaskFilter) ([]*domain.ScanTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var projectID int64
	if filter.ProjectCode != "" {
		id, ok := s.projectByCode[strings.TrimSpace(strings.ToLower(filter.ProjectCode))]
		if !ok {
			return []*domain.ScanTask{}, nil
		}
		projectID = id
	}

	result := make([]*domain.ScanTask, 0, len(s.scanTasks))
	for _, task := range s.scanTasks {
		if projectID != 0 && task.ProjectID != projectID {
			continue
		}
		if filter.Status != "" && !strings.EqualFold(string(task.Status), strings.TrimSpace(filter.Status)) {
			continue
		}
		if filter.SourceType != "" && !strings.EqualFold(task.SourceType, strings.TrimSpace(filter.SourceType)) {
			continue
		}
		result = append(result, cloneScanTask(task))
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].StartedAt.Equal(result[j].StartedAt) {
			return result[i].ID > result[j].ID
		}
		return result[i].StartedAt.After(result[j].StartedAt)
	})

	return result, nil
}

func (s *Store) GetScanTask(_ context.Context, id int64) (*domain.ScanTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.scanTasks[id]
	if !ok {
		return nil, nil
	}
	return cloneScanTask(task), nil
}

func (s *Store) CreateAttachment(_ context.Context, attachment *domain.AttachmentFile) (*domain.AttachmentFile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cloned := cloneAttachment(attachment)
	cloned.ID = s.allocateIDLocked()
	s.attachments[cloned.ID] = cloned
	s.attachmentsByScanTask[cloned.ScanTaskID] = append(s.attachmentsByScanTask[cloned.ScanTaskID], cloned.ID)
	return cloneAttachment(cloned), nil
}

func (s *Store) ListAttachmentsByScanTask(_ context.Context, scanTaskID int64) ([]*domain.AttachmentFile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := s.attachmentsByScanTask[scanTaskID]
	result := make([]*domain.AttachmentFile, 0, len(ids))
	for _, id := range ids {
		result = append(result, cloneAttachment(s.attachments[id]))
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].CreatedAt.Equal(result[j].CreatedAt) {
			return result[i].ID > result[j].ID
		}
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

func (s *Store) CreateOrGetResource(_ context.Context, resource *domain.ResourceItem) (*domain.ResourceItem, error) {
	if strings.TrimSpace(resource.ResourceKey) == "" {
		return nil, fmt.Errorf("resource key is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	indexKey := resourceIndexKey(resource.ProjectID, resource.ResourceKey)
	if id, ok := s.resourceByKey[indexKey]; ok {
		current := s.resources[id]
		mergeResourceFields(current, resource)
		return cloneResource(current), nil
	}

	cloned := cloneResource(resource)
	cloned.ID = s.allocateIDLocked()
	s.resources[cloned.ID] = cloned
	s.resourceByKey[indexKey] = cloned.ID
	return cloneResource(cloned), nil
}

func (s *Store) FindIssueByFingerprint(_ context.Context, projectID int64, fingerprint string) (*domain.IssueItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.issueByFingerprint[issueIndexKey(projectID, fingerprint)]
	if !ok {
		return nil, nil
	}
	return cloneIssue(s.issues[id]), nil
}

func (s *Store) CreateIssue(_ context.Context, issue *domain.IssueItem) (*domain.IssueItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cloned := cloneIssue(issue)
	cloned.ID = s.allocateIDLocked()
	s.issues[cloned.ID] = cloned
	s.issueByFingerprint[issueIndexKey(cloned.ProjectID, cloned.Fingerprint)] = cloned.ID
	return cloneIssue(cloned), nil
}

func (s *Store) UpdateIssue(_ context.Context, issue *domain.IssueItem) (*domain.IssueItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.issues[issue.ID]; !ok {
		return nil, fmt.Errorf("issue not found: %d", issue.ID)
	}
	cloned := cloneIssue(issue)
	s.issues[cloned.ID] = cloned
	s.issueByFingerprint[issueIndexKey(cloned.ProjectID, cloned.Fingerprint)] = cloned.ID
	return cloneIssue(cloned), nil
}

func (s *Store) ListIssues(_ context.Context, filter domain.IssueFilter) ([]*domain.IssueItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var projectID int64
	if filter.ProjectCode != "" {
		id, ok := s.projectByCode[strings.TrimSpace(strings.ToLower(filter.ProjectCode))]
		if !ok {
			return []*domain.IssueItem{}, nil
		}
		projectID = id
	}

	resourcePathFilter := strings.TrimSpace(strings.ToLower(filter.ResourcePath))
	result := make([]*domain.IssueItem, 0, len(s.issues))

	for _, issue := range s.issues {
		if projectID != 0 && issue.ProjectID != projectID {
			continue
		}
		if filter.ScanTaskID != nil && issue.LastScanTaskID != *filter.ScanTaskID {
			continue
		}
		if filter.Status != "" && !strings.EqualFold(string(issue.Status), strings.TrimSpace(filter.Status)) {
			continue
		}
		if filter.Severity != "" && !strings.EqualFold(issue.Severity, strings.TrimSpace(filter.Severity)) {
			continue
		}
		if filter.RuleCode != "" && !strings.EqualFold(issue.RuleCode, strings.TrimSpace(filter.RuleCode)) {
			continue
		}
		if filter.AssigneeName != "" && !strings.EqualFold(issue.AssigneeName, strings.TrimSpace(filter.AssigneeName)) {
			continue
		}
		if resourcePathFilter != "" {
			resource := s.resources[issue.ResourceID]
			if resource == nil || !strings.Contains(strings.ToLower(resource.ResourcePath), resourcePathFilter) {
				continue
			}
		}
		result = append(result, cloneIssue(issue))
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].LastSeenAt.Equal(result[j].LastSeenAt) {
			return result[i].ID > result[j].ID
		}
		return result[i].LastSeenAt.After(result[j].LastSeenAt)
	})

	return result, nil
}

func (s *Store) GetIssue(_ context.Context, id int64) (*domain.IssueItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	issue, ok := s.issues[id]
	if !ok {
		return nil, nil
	}
	return cloneIssue(issue), nil
}

func (s *Store) GetResource(_ context.Context, id int64) (*domain.ResourceItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resource, ok := s.resources[id]
	if !ok {
		return nil, fmt.Errorf("resource not found: %d", id)
	}
	return cloneResource(resource), nil
}

func (s *Store) AppendIssueAction(_ context.Context, action *domain.IssueActionLog) (*domain.IssueActionLog, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cloned := cloneAction(action)
	cloned.ID = s.allocateIDLocked()
	s.actions[cloned.ID] = cloned
	s.actionsByIssue[cloned.IssueID] = append(s.actionsByIssue[cloned.IssueID], cloned.ID)
	return cloneAction(cloned), nil
}

func (s *Store) ListIssueActions(_ context.Context, issueID int64) ([]*domain.IssueActionLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := s.actionsByIssue[issueID]
	result := make([]*domain.IssueActionLog, 0, len(ids))
	for _, id := range ids {
		result = append(result, cloneAction(s.actions[id]))
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].CreatedAt.Equal(result[j].CreatedAt) {
			return result[i].ID > result[j].ID
		}
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

func (s *Store) allocateIDLocked() int64 {
	id := s.nextID
	s.nextID++
	return id
}

func resourceIndexKey(projectID int64, resourceKey string) string {
	return fmt.Sprintf("%d|%s", projectID, strings.TrimSpace(resourceKey))
}

func issueIndexKey(projectID int64, fingerprint string) string {
	return fmt.Sprintf("%d|%s", projectID, strings.TrimSpace(fingerprint))
}

func mergeResourceFields(current, incoming *domain.ResourceItem) {
	if incoming.ResourceGUID != "" {
		current.ResourceGUID = incoming.ResourceGUID
	}
	if incoming.ResourcePath != "" {
		current.ResourcePath = incoming.ResourcePath
	}
	if incoming.ResourceName != "" {
		current.ResourceName = incoming.ResourceName
	}
	if incoming.ResourceType != "" {
		current.ResourceType = incoming.ResourceType
	}
	if incoming.ModuleName != "" {
		current.ModuleName = incoming.ModuleName
	}
	if incoming.OwnerName != "" {
		current.OwnerName = incoming.OwnerName
	}
	if !incoming.UpdatedAt.IsZero() {
		current.UpdatedAt = incoming.UpdatedAt
	}
}

func cloneProject(value *domain.Project) *domain.Project {
	if value == nil {
		return nil
	}
	clone := *value
	return &clone
}

func cloneScanTask(value *domain.ScanTask) *domain.ScanTask {
	if value == nil {
		return nil
	}
	clone := *value
	return &clone
}

func cloneResource(value *domain.ResourceItem) *domain.ResourceItem {
	if value == nil {
		return nil
	}
	clone := *value
	return &clone
}

func cloneIssue(value *domain.IssueItem) *domain.IssueItem {
	if value == nil {
		return nil
	}
	clone := *value
	if value.ResolvedAt != nil {
		resolvedAt := *value.ResolvedAt
		clone.ResolvedAt = &resolvedAt
	}
	if value.IgnoredAt != nil {
		ignoredAt := *value.IgnoredAt
		clone.IgnoredAt = &ignoredAt
	}
	return &clone
}

func cloneAction(value *domain.IssueActionLog) *domain.IssueActionLog {
	if value == nil {
		return nil
	}
	clone := *value
	return &clone
}

func cloneAttachment(value *domain.AttachmentFile) *domain.AttachmentFile {
	if value == nil {
		return nil
	}
	clone := *value
	return &clone
}
