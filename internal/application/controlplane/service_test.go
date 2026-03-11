package controlplane

import (
	"context"
	"testing"
	"time"

	"github.com/TwoEggDu/content-control-plane/internal/domain"
	"github.com/TwoEggDu/content-control-plane/internal/infrastructure/memory"
)

func TestImportScanIsIdempotent(t *testing.T) {
	store := memory.NewStore()
	service := NewService(store)
	fixedNow := time.Date(2026, 3, 11, 9, 30, 0, 0, time.UTC)
	service.now = func() time.Time { return fixedNow }

	request := sampleImportRequest()

	first, err := service.ImportScan(context.Background(), request)
	if err != nil {
		t.Fatalf("first import failed: %v", err)
	}
	second, err := service.ImportScan(context.Background(), request)
	if err != nil {
		t.Fatalf("second import failed: %v", err)
	}

	if !first.Created || first.Reused {
		t.Fatalf("expected first import to create a new task, got %+v", first)
	}
	if second.Created || !second.Reused {
		t.Fatalf("expected second import to reuse the task, got %+v", second)
	}
	if first.ScanTaskID != second.ScanTaskID {
		t.Fatalf("expected same scan task id, got %d and %d", first.ScanTaskID, second.ScanTaskID)
	}

	tasks, err := service.ListScanTasks(context.Background(), domain.ScanTaskFilter{ProjectCode: "twoegg-mobile"})
	if err != nil {
		t.Fatalf("list scan tasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 scan task, got %d", len(tasks))
	}

	issues, err := service.ListIssues(context.Background(), domain.IssueFilter{ProjectCode: "twoegg-mobile"})
	if err != nil {
		t.Fatalf("list issues failed: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestUpdateIssueStatusWritesActionLog(t *testing.T) {
	store := memory.NewStore()
	service := NewService(store)
	fixedNow := time.Date(2026, 3, 11, 9, 30, 0, 0, time.UTC)
	service.now = func() time.Time { return fixedNow }

	if _, err := service.ImportScan(context.Background(), sampleImportRequest()); err != nil {
		t.Fatalf("import failed: %v", err)
	}

	issues, err := service.ListIssues(context.Background(), domain.IssueFilter{ProjectCode: "twoegg-mobile"})
	if err != nil {
		t.Fatalf("list issues failed: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}

	detail, err := service.UpdateIssueStatus(context.Background(), issues[0].ID, UpdateIssueStatusRequest{
		ToStatus:     "ASSIGNED",
		AssigneeName: "zhangsan",
		OperatorName: "lead-user",
		Comment:      "assign to art owner",
	})
	if err != nil {
		t.Fatalf("update status failed: %v", err)
	}

	if detail.Status != "ASSIGNED" {
		t.Fatalf("expected status ASSIGNED, got %s", detail.Status)
	}
	if detail.AssigneeName != "zhangsan" {
		t.Fatalf("expected assignee zhangsan, got %s", detail.AssigneeName)
	}
	if len(detail.Actions) != 1 {
		t.Fatalf("expected 1 action log, got %d", len(detail.Actions))
	}
	if detail.Actions[0].ActionType != "ASSIGNED" {
		t.Fatalf("expected action type ASSIGNED, got %s", detail.Actions[0].ActionType)
	}
}

func sampleImportRequest() ImportScanRequest {
	return ImportScanRequest{
		ProjectCode:    "twoegg-mobile",
		ProjectName:    "TwoEgg Mobile",
		SourceType:     "CI",
		TaskNo:         "ci-20260311-001",
		BranchName:     "main",
		CommitSHA:      "7d7f5afc",
		ScannerVersion: "asset-checker@1.4.2",
		TriggeredBy:    "buildkite",
		StartedAt:      "2026-03-11T09:20:00Z",
		FinishedAt:     "2026-03-11T09:21:15Z",
		Attachments: []ImportAttachmentInput{
			{
				FileType:    "REPORT",
				FileName:    "scan-report.json",
				StorageKey:  "reports/twoegg-mobile/2026-03-11/scan-report.json",
				ContentHash: "sha256:report-001",
				FileSize:    12540,
			},
		},
		Issues: []ImportIssueInput{
			{
				RuleCode:      "TEXTURE_MAX_SIZE",
				Severity:      "HIGH",
				Message:       "Texture size exceeds project baseline",
				ResourceGUID:  "0f7f5d4d6a554f3cbef1c9b111111111",
				ResourcePath:  "Assets/Art/Hero/hero_diffuse.png",
				ResourceName:  "hero_diffuse.png",
				ResourceType:  "Texture2D",
				ModuleName:    "Hero",
				OwnerName:     "art-team",
				LocationKey:   "import_settings.max_size",
				CurrentValue:  "4096",
				ExpectedValue: "2048",
			},
		},
	}
}
