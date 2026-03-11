# Architecture

## Platform Architecture

```mermaid
flowchart LR
    A["Unity 检查器 / 命令行扫描器 / CI"] --> B["结果接入模块"]
    U["前端后台"] --> G["查询与操作 API"]

    B --> T["扫描任务模块"]
    B --> I["问题管理模块"]
    B --> R["资源索引模块"]

    G --> I
    G --> T
    G --> R
    G --> W["状态流转模块"]
    G --> S["统计报表模块"]
    G --> P["权限与用户模块"]
    G --> C["发布控制模块"]
    G --> A1["审计模块"]

    W --> N["通知与异步任务模块"]
    S --> N
    C --> N

    T --> DB[("PostgreSQL")]
    I --> DB
    R --> DB
    W --> DB
    P --> DB
    S --> DB
    C --> DB
    A1 --> DB

    B --> O[("MinIO / S3")]
    S --> O
    C --> O
```

## Module Boundaries

- 结果接入模块：负责接收扫描结果、报告、日志、截图，并做幂等导入
- 扫描任务模块：负责记录扫描来源、分支、提交号、扫描器版本、耗时和状态
- 问题管理模块：负责错误资源问题列表、筛选、分配、忽略、关闭和验证
- 资源索引模块：负责资源路径、GUID、类型、所属模块和责任人映射
- 状态流转模块：负责 `NEW / ASSIGNED / FIXING / RESOLVED / VERIFIED / IGNORED`
- 发布控制模块：负责版本登记、环境提升、审批、灰度和回滚元数据
- 统计报表模块：负责趋势、排行、汇总报表
- 权限与用户模块：负责角色、项目权限和访问控制
- 审计模块：负责关键操作留痕
- 通知与异步任务模块：负责提醒、定时报表、异步发布任务

## Data Model

```mermaid
erDiagram
    PROJECT ||--o{ SCAN_TASK : has
    PROJECT ||--o{ RESOURCE_ITEM : owns
    PROJECT ||--o{ RULE_DEF : uses
    USER_ACCOUNT ||--o{ RESOURCE_ITEM : responsible_for
    USER_ACCOUNT ||--o{ ISSUE_ITEM : assigned_to
    SCAN_TASK ||--o{ ISSUE_ITEM : produces
    RESOURCE_ITEM ||--o{ ISSUE_ITEM : hits
    RULE_DEF ||--o{ ISSUE_ITEM : triggered_by
    ISSUE_ITEM ||--o{ ISSUE_ACTION_LOG : has
    SCAN_TASK ||--o{ ATTACHMENT_FILE : outputs

    PROJECT {
        bigint id PK
        varchar code
        varchar name
        tinyint status
        datetime created_at
    }

    USER_ACCOUNT {
        bigint id PK
        varchar username
        varchar display_name
        varchar email
        varchar department
        tinyint status
        datetime created_at
    }

    RULE_DEF {
        bigint id PK
        bigint project_id FK
        varchar rule_code
        varchar rule_name
        varchar severity
        varchar category
        tinyint enabled
        varchar owner_team
        datetime updated_at
    }

    RESOURCE_ITEM {
        bigint id PK
        bigint project_id FK
        varchar resource_guid
        varchar resource_path
        varchar resource_name
        varchar resource_type
        varchar module_name
        bigint owner_user_id FK
        varchar asset_bundle
        datetime updated_at
    }

    SCAN_TASK {
        bigint id PK
        bigint project_id FK
        varchar task_no
        varchar source_type
        varchar branch_name
        varchar commit_sha
        varchar scanner_version
        varchar trigger_by
        datetime started_at
        datetime finished_at
        varchar status
        int total_issue_count
    }

    ISSUE_ITEM {
        bigint id PK
        bigint project_id FK
        bigint scan_task_id FK
        bigint resource_id FK
        bigint rule_id FK
        varchar fingerprint
        varchar severity
        varchar status
        bigint assignee_user_id FK
        text message
        text current_value
        text expected_value
        tinyint is_new
        tinyint is_ignored
        datetime first_seen_at
        datetime last_seen_at
        datetime resolved_at
    }

    ISSUE_ACTION_LOG {
        bigint id PK
        bigint issue_id FK
        bigint operator_user_id FK
        varchar action_type
        varchar from_status
        varchar to_status
        text comment
        datetime created_at
    }

    ATTACHMENT_FILE {
        bigint id PK
        bigint scan_task_id FK
        varchar file_type
        varchar file_name
        varchar storage_key
        varchar content_hash
        datetime created_at
    }
```

## Core Flow

1. 扫描器或 CI 生成结果并上传到结果接入模块。
2. 后端创建 `SCAN_TASK`，解析并写入 `ISSUE_ITEM`、`RESOURCE_ITEM`、附件记录。
3. 前端通过查询 API 查看问题列表、详情、趋势和责任归属。
4. 问题在状态流转模块中被分配、修复、验证或忽略。
5. 发布控制模块基于资源和版本信息执行审批、提升和回滚记录。
