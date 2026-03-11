# 01. MVP API

## 1. 目标

本文档冻结第一版最小闭环 API：

1. 导入扫描结果
2. 查询扫描任务
3. 查询问题列表
4. 查询问题详情
5. 更新问题状态

当前本地实现允许先使用内存仓储，但字段语义、状态语义和错误结构应与后续 PostgreSQL 版本一致。

## 2. 通用约定

- 协议：`HTTP + JSON`
- 编码：`UTF-8`
- 时间字段：`RFC3339`
- `id` 字段：`int64`
- 第一版列表接口直接返回数组，不提供分页对象

### 2.1 错误响应

第一版统一错误结构：

```json
{
  "error": {
    "code": "invalid_request",
    "message": "project_code is required"
  }
}
```

当前约定的错误码：

- `invalid_request`
- `not_found`
- `invalid_transition`
- `internal_error`

## 3. 健康检查

### `GET /healthz`

用途：进程存活检查。

响应示例：

```json
{
  "status": "ok"
}
```

### `GET /readyz`

用途：服务就绪检查。

响应示例：

```json
{
  "status": "ready"
}
```

## 4. 导入接口

### `POST /api/scan-tasks/import`

用途：导入一批扫描结果和附件元信息。

请求体：

```json
{
  "project_code": "twoegg-mobile",
  "project_name": "TwoEgg Mobile",
  "source_type": "CI",
  "task_no": "ci-20260311-001",
  "branch_name": "main",
  "commit_sha": "7d7f5afc",
  "scanner_version": "asset-checker@1.4.2",
  "triggered_by": "buildkite",
  "started_at": "2026-03-11T09:20:00Z",
  "finished_at": "2026-03-11T09:21:15Z",
  "attachments": [
    {
      "file_type": "REPORT",
      "file_name": "scan-report.json",
      "storage_key": "reports/twoegg-mobile/2026-03-11/scan-report.json",
      "content_hash": "sha256:report-001",
      "file_size": 12540
    }
  ],
  "issues": [
    {
      "rule_code": "TEXTURE_MAX_SIZE",
      "severity": "HIGH",
      "message": "Texture size exceeds project baseline",
      "resource_guid": "0f7f5d4d6a554f3cbef1c9b111111111",
      "resource_path": "Assets/Art/Hero/hero_diffuse.png",
      "resource_name": "hero_diffuse.png",
      "resource_type": "Texture2D",
      "module_name": "Hero",
      "owner_name": "art-team",
      "location_key": "import_settings.max_size",
      "current_value": "4096",
      "expected_value": "2048"
    }
  ]
}
```

字段要求：

- 必填：`project_code`、`source_type`、`branch_name`、`commit_sha`、`scanner_version`、`triggered_by`
- `task_no` 可空；为空时服务端使用派生幂等键
- `attachments` 可空数组
- `issues` 允许为空数组，但第一版更推荐由上游显式传入扫描结果
- `issues[*].rule_code`、`severity`、`message` 至少必须存在
- `issues[*].resource_guid` 和 `issues[*].resource_path` 不能同时为空

幂等语义：

- 如果命中已有导入批次，返回已有 `scan_task_id`
- 如果命中已有问题指纹，更新已有问题的最新观测信息，而不是重复创建长期问题实体

成功响应：

```json
{
  "scan_task_id": 1,
  "import_key": "task_no:ci-20260311-001",
  "status": "IMPORTED",
  "created": true,
  "reused": false,
  "total_issue_count": 1,
  "new_issue_count": 1,
  "reopened_issue_count": 0
}
```

重复导入响应示例：

```json
{
  "scan_task_id": 1,
  "import_key": "task_no:ci-20260311-001",
  "status": "IMPORTED",
  "created": false,
  "reused": true,
  "total_issue_count": 1,
  "new_issue_count": 1,
  "reopened_issue_count": 0
}
```

样例文件：

- [examples/import-scan-sample.json](examples/import-scan-sample.json)

## 5. 扫描任务查询

### `GET /api/scan-tasks`

支持的查询参数：

- `project_code`
- `status`
- `source_type`

成功响应：

```json
[
  {
    "id": 1,
    "project_code": "twoegg-mobile",
    "task_no": "ci-20260311-001",
    "source_type": "CI",
    "branch_name": "main",
    "commit_sha": "7d7f5afc",
    "scanner_version": "asset-checker@1.4.2",
    "status": "IMPORTED",
    "total_issue_count": 1,
    "new_issue_count": 1,
    "reopened_issue_count": 0,
    "attachment_count": 1,
    "started_at": "2026-03-11T09:20:00Z",
    "finished_at": "2026-03-11T09:21:15Z"
  }
]
```

### `GET /api/scan-tasks/{id}`

成功响应：

```json
{
  "id": 1,
  "project_code": "twoegg-mobile",
  "task_no": "ci-20260311-001",
  "source_type": "CI",
  "branch_name": "main",
  "commit_sha": "7d7f5afc",
  "scanner_version": "asset-checker@1.4.2",
  "triggered_by": "buildkite",
  "status": "IMPORTED",
  "total_issue_count": 1,
  "new_issue_count": 1,
  "reopened_issue_count": 0,
  "attachment_count": 1,
  "started_at": "2026-03-11T09:20:00Z",
  "finished_at": "2026-03-11T09:21:15Z",
  "attachments": [
    {
      "id": 1,
      "file_type": "REPORT",
      "file_name": "scan-report.json",
      "storage_key": "reports/twoegg-mobile/2026-03-11/scan-report.json",
      "content_hash": "sha256:report-001",
      "file_size": 12540,
      "created_at": "2026-03-11T09:21:15Z"
    }
  ]
}
```

## 6. 问题查询

### `GET /api/issues`

支持的查询参数：

- `project_code`
- `scan_task_id`
- `status`
- `severity`
- `rule_code`
- `assignee_name`
- `resource_path`

成功响应：

```json
[
  {
    "id": 1,
    "project_code": "twoegg-mobile",
    "scan_task_id": 1,
    "resource_id": 1,
    "rule_code": "TEXTURE_MAX_SIZE",
    "severity": "HIGH",
    "status": "NEW",
    "assignee_name": "",
    "message": "Texture size exceeds project baseline",
    "resource_path": "Assets/Art/Hero/hero_diffuse.png",
    "resource_guid": "0f7f5d4d6a554f3cbef1c9b111111111",
    "first_seen_at": "2026-03-11T09:21:15Z",
    "last_seen_at": "2026-03-11T09:21:15Z"
  }
]
```

### `GET /api/issues/{id}`

成功响应：

```json
{
  "id": 1,
  "project_code": "twoegg-mobile",
  "scan_task_id": 1,
  "resource": {
    "id": 1,
    "resource_guid": "0f7f5d4d6a554f3cbef1c9b111111111",
    "resource_path": "Assets/Art/Hero/hero_diffuse.png",
    "resource_name": "hero_diffuse.png",
    "resource_type": "Texture2D",
    "module_name": "Hero",
    "owner_name": "art-team"
  },
  "rule_code": "TEXTURE_MAX_SIZE",
  "severity": "HIGH",
  "status": "NEW",
  "assignee_name": "",
  "message": "Texture size exceeds project baseline",
  "current_value": "4096",
  "expected_value": "2048",
  "location_key": "import_settings.max_size",
  "first_seen_at": "2026-03-11T09:21:15Z",
  "last_seen_at": "2026-03-11T09:21:15Z",
  "resolved_at": null,
  "actions": []
}
```

## 7. 问题状态流转

### `POST /api/issues/{id}/status`

请求体：

```json
{
  "to_status": "ASSIGNED",
  "assignee_name": "zhangsan",
  "operator_name": "lead-user",
  "comment": "assign to art owner"
}
```

字段要求：

- `to_status` 必填
- `operator_name` 必填
- 当 `to_status=ASSIGNED` 时，`assignee_name` 必填
- 其他状态允许不传 `assignee_name`

允许的人工流转：

| Current | Allowed Next |
| --- | --- |
| `NEW` | `ASSIGNED`, `FIXING`, `RESOLVED`, `IGNORED` |
| `ASSIGNED` | `FIXING`, `RESOLVED`, `IGNORED` |
| `FIXING` | `RESOLVED`, `IGNORED` |
| `RESOLVED` | `VERIFIED`, `FIXING` |
| `VERIFIED` | `FIXING` |
| `IGNORED` | `NEW`, `ASSIGNED`, `FIXING` |

成功响应返回更新后的问题详情。

## 8. 当前范围外接口

这些接口概念上存在，但不进入第一版冻结范围：

- `POST /api/issues/{id}/assign`
- `POST /api/issues/{id}/ignore`
- `GET /api/resources/{id}/issues`
- `GET /api/reports/summary`
- `GET /api/reports/trend`

## 9. 当前结论

第一版 API 已经足够支撑导入到查询闭环。

后续新增接口时，应优先保证不破坏：

- 导入幂等语义
- 问题指纹语义
- 状态流转语义
- 错误结构一致性