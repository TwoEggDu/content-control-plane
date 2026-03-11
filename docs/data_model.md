# Data Model

## 文档目的

本文档用于冻结第一版最小可运行版本的数据模型。

这里的“冻结”不是直接等于最终 DDL，而是表示这些对象、字段语义、唯一约束和关系方向已经进入实现基线，不应再在聊天里反复漂移。

当前冻结范围服务于第一版闭环：

1. 导入扫描结果
2. 查询扫描任务
3. 查询问题列表
4. 查看问题详情
5. 更新问题状态

## 与其他文档的关系

- `product_intent`：回答为什么要管理这些对象
- `docs/sdd/01-overview.md`：回答这些对象落在什么系统边界里
- `docs/sdd/02-ingestion-and-issue-domain.md`：回答导入、幂等、资源识别和问题指纹
- `docs/sdd/03-query-workflow-and-reporting.md`：回答这些对象如何被查询和流转
- `docs/api/01-mvp-api.md`：回答这些对象如何通过 HTTP 暴露

## 第一版冻结对象

### `project`

用途：定义治理边界。

冻结字段：

- `id`：主键
- `code`：项目编码，唯一且稳定
- `name`：项目显示名称
- `status`：当前项目状态，第一版建议 `ACTIVE / ARCHIVED`
- `created_at`
- `updated_at`

约束：

- 唯一键：`code`

说明：

- 第一版权限体系未落地前，`project` 仍然是所有查询的最高过滤维度

### `scan_task`

用途：表示一次扫描结果导入批次。

冻结字段：

- `id`：主键
- `project_id`：所属项目
- `import_key`：导入幂等键
- `task_no`：上游任务号，可空
- `source_type`：来源类型，例如 `UNITY_EDITOR / CLI / CI`
- `branch_name`
- `commit_sha`
- `scanner_version`
- `triggered_by`
- `started_at`
- `finished_at`
- `status`：`RECEIVED / PROCESSING / IMPORTED / FAILED`
- `total_issue_count`
- `new_issue_count`
- `reopened_issue_count`
- `attachment_count`
- `created_at`
- `updated_at`

约束：

- 唯一键：`import_key`
- 可选唯一约束：`(project_id, task_no)`，仅在 `task_no` 非空时生效

说明：

- `scan_task` 是批次对象，不是长期治理对象
- 批次状态和批次统计都应挂在这里，而不是分散到附件或问题表里

### `resource_item`

用途：表示被治理的资源对象。

冻结字段：

- `id`：主键
- `project_id`：所属项目
- `resource_key`：资源稳定键，优先 GUID，否则使用规范化路径
- `resource_guid`：可空
- `resource_path`
- `resource_name`
- `resource_type`
- `module_name`
- `owner_name`：第一版先用字符串，不冻结用户外键
- `created_at`
- `updated_at`

约束：

- 唯一键：`(project_id, resource_key)`

说明：

- `resource_key` 是实现层字段，用于保证跨扫描复用同一资源对象
- 第一版不冻结 `owner_user_id`，避免在权限体系未稳定前把用户模型提前做重

### `issue_item`

用途：表示可持续治理的问题实体。

冻结字段：

- `id`：主键
- `project_id`：所属项目
- `resource_id`：关联资源
- `last_scan_task_id`：最近一次观察到该问题的扫描任务
- `fingerprint`：问题稳定指纹
- `rule_code`
- `severity`
- `status`：`NEW / ASSIGNED / FIXING / RESOLVED / VERIFIED / IGNORED`
- `assignee_name`：第一版先用字符串，可空
- `message`
- `current_value`
- `expected_value`
- `location_key`：资源内定位键，可空
- `first_seen_at`
- `last_seen_at`
- `resolved_at`
- `ignored_at`
- `created_at`
- `updated_at`

约束：

- 唯一键：`(project_id, fingerprint)`
- 索引建议：
  - `(project_id, status, severity)`
  - `(project_id, rule_code)`
  - `(project_id, resource_id)`
  - `(last_scan_task_id)`

说明：

- `issue_item` 是长期治理对象，不是单次扫描快照
- 如果后续需要完整保留每次扫描观测，第二阶段再增加 `issue_occurrence`

### `issue_action_log`

用途：记录问题处理历史。

冻结字段：

- `id`：主键
- `issue_id`
- `action_type`：第一版至少支持 `STATUS_CHANGED` 和系统重开动作
- `from_status`
- `to_status`
- `operator_name`
- `comment`
- `created_at`

约束：

- 索引建议：`(issue_id, created_at DESC)`

说明：

- 不要只保留当前状态，任何重要动作都必须有历史

### `attachment_file`

用途：记录扫描批次附件元信息。

冻结字段：

- `id`：主键
- `scan_task_id`
- `file_type`：例如 `REPORT / LOG / SCREENSHOT / RAW_RESULT`
- `file_name`
- `storage_key`
- `content_hash`
- `file_size`
- `created_at`

约束：

- 索引建议：`(scan_task_id, created_at DESC)`

说明：

- 大文件走对象存储，数据库只保存定位和校验元信息

## 当前未冻结对象

这些对象概念上存在，但不进入第一版冻结范围：

- `rule_def`
- `user_account`
- `issue_occurrence`
- `release_record`
- `approval_flow`

原因不是它们不重要，而是它们不该在第一版最小闭环里抢占建模主导权。

## 核心关系

第一版关系方向固定如下：

- `project 1 -> n scan_task`
- `project 1 -> n resource_item`
- `project 1 -> n issue_item`
- `scan_task 1 -> n attachment_file`
- `scan_task 1 -> n issue_item` 通过 `last_scan_task_id` 表示最近观察关系
- `resource_item 1 -> n issue_item`
- `issue_item 1 -> n issue_action_log`

## 建模原则

### 1. 幂等键和指纹都必须可追溯

- `import_key` 负责识别导入批次
- `fingerprint` 负责识别长期问题实体

两者都不能只存结果，不存业务语义来源。

### 2. 责任信息第一版先轻量建模

- `owner_name` 和 `assignee_name` 先以字符串进入模型
- 用户体系成熟后再替换成外键，不影响第一版 API 和闭环

### 3. 批次统计字段放在 `scan_task`

第一版不做复杂统计仓库，因此新问题数、重开问题数和附件数都直接落在 `scan_task`。

### 4. 问题历史靠日志，不靠覆盖描述

当前状态、当前责任人和最近一次消息可以被更新，但动作历史必须保留在 `issue_action_log`。

## 当前结论

第一版冻结模型已经足够支撑导入、查询和状态流转闭环。

后续如果要新增字段，应优先保证不破坏：

- `project -> scan_task`
- `project -> resource_item`
- `project/resource -> issue_item`
- `issue_item -> issue_action_log`

这四条主关系。