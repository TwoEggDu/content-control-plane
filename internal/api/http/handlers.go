package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/TwoEggDu/content-control-plane/internal/application/controlplane"
	"github.com/TwoEggDu/content-control-plane/internal/domain"
	"github.com/gin-gonic/gin"
)

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func RegisterRoutes(router *gin.Engine, service *controlplane.Service) {
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/readyz", func(c *gin.Context) {
		if err := service.Ready(c.Request.Context()); err != nil {
			renderError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	router.POST("/api/scan-tasks/import", func(c *gin.Context) {
		var request controlplane.ImportScanRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			renderError(c, errors.New("invalid_request: invalid json body"))
			return
		}

		result, err := service.ImportScan(c.Request.Context(), request)
		if err != nil {
			renderError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	})

	router.GET("/api/scan-tasks", func(c *gin.Context) {
		result, err := service.ListScanTasks(c.Request.Context(), domain.ScanTaskFilter{
			ProjectCode: strings.TrimSpace(c.Query("project_code")),
			Status:      strings.TrimSpace(c.Query("status")),
			SourceType:  strings.TrimSpace(c.Query("source_type")),
		})
		if err != nil {
			renderError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	})

	router.GET("/api/scan-tasks/:id", func(c *gin.Context) {
		id, err := parseID(c.Param("id"))
		if err != nil {
			renderError(c, err)
			return
		}
		result, err := service.GetScanTaskDetail(c.Request.Context(), id)
		if err != nil {
			renderError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	})

	router.GET("/api/issues", func(c *gin.Context) {
		filter := domain.IssueFilter{
			ProjectCode:  strings.TrimSpace(c.Query("project_code")),
			Status:       strings.TrimSpace(c.Query("status")),
			Severity:     strings.TrimSpace(c.Query("severity")),
			RuleCode:     strings.TrimSpace(c.Query("rule_code")),
			AssigneeName: strings.TrimSpace(c.Query("assignee_name")),
			ResourcePath: strings.TrimSpace(c.Query("resource_path")),
		}
		if rawScanTaskID := strings.TrimSpace(c.Query("scan_task_id")); rawScanTaskID != "" {
			scanTaskID, err := parseID(rawScanTaskID)
			if err != nil {
				renderError(c, err)
				return
			}
			filter.ScanTaskID = &scanTaskID
		}

		result, err := service.ListIssues(c.Request.Context(), filter)
		if err != nil {
			renderError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	})

	router.GET("/api/issues/:id", func(c *gin.Context) {
		id, err := parseID(c.Param("id"))
		if err != nil {
			renderError(c, err)
			return
		}
		result, err := service.GetIssueDetail(c.Request.Context(), id)
		if err != nil {
			renderError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	})

	router.POST("/api/issues/:id/status", func(c *gin.Context) {
		id, err := parseID(c.Param("id"))
		if err != nil {
			renderError(c, err)
			return
		}

		var request controlplane.UpdateIssueStatusRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			renderError(c, errors.New("invalid_request: invalid json body"))
			return
		}

		result, err := service.UpdateIssueStatus(c.Request.Context(), id, request)
		if err != nil {
			renderError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	})
}

func parseID(raw string) (int64, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || value <= 0 {
		return 0, errors.New("invalid_request: id must be a positive integer")
	}
	return value, nil
}

func renderError(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	code := "internal_error"
	message := err.Error()

	switch {
	case errors.Is(err, controlplane.ErrInvalidRequest):
		statusCode = http.StatusBadRequest
		code = "invalid_request"
		message = sanitizeMessage(message, "invalid_request:")
	case errors.Is(err, controlplane.ErrNotFound):
		statusCode = http.StatusNotFound
		code = "not_found"
		message = sanitizeMessage(message, "not_found:")
	case errors.Is(err, controlplane.ErrInvalidTransition):
		statusCode = http.StatusBadRequest
		code = "invalid_transition"
		message = sanitizeMessage(message, "invalid_transition:")
	default:
		if strings.HasPrefix(err.Error(), "invalid_request:") {
			statusCode = http.StatusBadRequest
			code = "invalid_request"
			message = sanitizeMessage(message, "invalid_request:")
		}
	}

	c.JSON(statusCode, errorEnvelope{
		Error: errorBody{
			Code:    code,
			Message: message,
		},
	})
}

func sanitizeMessage(message, prefix string) string {
	trimmed := strings.TrimSpace(strings.TrimPrefix(message, prefix))
	if trimmed == "" {
		return message
	}
	return trimmed
}
