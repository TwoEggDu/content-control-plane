# 03. Query, Workflow, and Reporting

## 1. 文档目的

本文档用于定义第一阶段最小闭环里的查询能力、问题状态流转和最小统计视图。

它回答的是：

- 扫描任务和问题要如何被查询
- 问题状态流转第一版如何约束
- 附件和报告在查询侧如何暴露
- 哪些统计能力在第一阶段必须进入系统

## 2. 范围

第一阶段本主题只覆盖这些内容：

- 扫描任务列表和详情查询
- 问题列表和详情查询
- 统一的问题状态流转接口
- 扫描任务附件查看
- 支撑治理闭环的最小统计字段

第一阶段不覆盖这些内容：

- 复杂报表中心
- 自定义仪表盘
- 导出中心
- 专门的资源视图和资源主页
- 发布控制视图

## 3. 查询视图设计

### 3.1 扫描任务列表

扫描任务列表的目标不是展示所有原始输入，而是回答：

- 这次导入来自哪里
- 导入是否成功
- 这次导入带来了多少问题
- 是否有可追踪的附件和上下文

第一版建议支持这些过滤维度：

- `project_code`
- `status`
- `source_type`
- `branch_name`
- `commit_sha`

第一版排序默认按：

1. `started_at DESC`
2. `id DESC`

扫描任务列表项至少应返回：

- `id`
- `project_code`
- `task_no`
- `source_type`
- `branch_name`
- `commit_sha`
- `scanner_version`
- `status`
- `total_issue_count`
- `new_issue_count`
- `reopened_issue_count`
- `attachment_count`
- `started_at`
- `finished_at`

### 3.2 扫描任务详情

扫描任务详情用于给治理者提供一次导入批次的完整上下文。

第一版详情至少应包含：

- 扫描任务元信息
- 批次统计字段
- 关联附件列表
- 本次批次关联问题的摘要视图

第一版不强制在扫描任务详情里做完整报表聚合，但要保证：

- 能看到导入是否成功
- 能看到本次导入影响了哪些问题
- 能定位到原始附件

### 3.3 问题列表

问题列表是第一阶段最重要的治理视图。

它的目标不是“原样复读扫描器输出”，而是把问题组织成可筛选、可分派、可跟踪的治理对象。

第一版建议支持这些过滤维度：

- `project_code`
- `scan_task_id`
- `status`
- `severity`
- `rule_code`
- `assignee_name`
- `resource_path`

第一版排序默认按：

1. `last_seen_at DESC`
2. `id DESC`

问题列表项至少应返回：

- `id`
- `project_code`
- `scan_task_id`
- `resource_id`
- `rule_code`
- `severity`
- `status`
- `assignee_name`
- `message`
- `resource_path`
- `resource_guid`
- `first_seen_at`
- `last_seen_at`

### 3.4 问题详情

问题详情用于回答：

- 这个问题到底是什么
- 它挂在哪个资源上
- 最近一次在哪次扫描里被看到
- 当前处于什么治理状态
- 历史上被谁如何处理过

第一版详情至少应包含：

- 问题核心字段
- 资源上下文
- 最近一次扫描任务引用
- 状态流转历史
- 最近一次消息、期望值和当前值

第一版暂不引入单独的问题级附件，因此附件仍通过关联扫描任务查看。

## 4. 状态流转模型

### 4.1 统一操作入口

第一版不拆分 `assign`、`ignore`、`resolve` 多个独立接口，而是统一收敛到：

- `POST /api/issues/{id}/status`

这样做的原因是：

- 第一阶段优先验证领域模型，而不是扩展操作面
- 所有状态变化都必须进入同一条审计链路
- 可以减少 Web 层和前端契约分叉

### 4.2 第一版允许的状态流转

| Current | Allowed Next |
| --- | --- |
| `NEW` | `ASSIGNED`, `FIXING`, `RESOLVED`, `IGNORED` |
| `ASSIGNED` | `FIXING`, `RESOLVED`, `IGNORED` |
| `FIXING` | `RESOLVED`, `IGNORED` |
| `RESOLVED` | `VERIFIED`, `FIXING` |
| `VERIFIED` | `FIXING` |
| `IGNORED` | `NEW`, `ASSIGNED`, `FIXING` |

约束说明：

- 进入 `ASSIGNED` 时应提供 `assignee_name`
- 所有人工流转都应提供 `operator_name`
- `comment` 可选，但建议在 `IGNORED`、`RESOLVED` 和回退场景中填写
- 导入链路可把 `RESOLVED` 或 `VERIFIED` 的问题重新打开为 `NEW`，这属于系统行为，不走人工接口

### 4.3 状态流转日志

每次成功流转都必须写入 `issue_action_log`。

第一版动作日志至少记录：

- `issue_id`
- `action_type`
- `from_status`
- `to_status`
- `operator_name`
- `comment`
- `created_at`

如果后续要补审批、通知或 SLA，动作日志会成为直接输入。

## 5. 附件和报告视图

### 5.1 附件归属

第一版附件只归属到 `scan_task`。

原因是：

- 报告、日志和截图通常描述的是整批扫描
- 第一阶段先保证导入链路可追踪
- 问题级附件会引入额外权限、存储和界面复杂度

### 5.2 扫描任务详情中的附件暴露方式

扫描任务详情至少应返回附件数组，每个附件项应包含：

- `id`
- `file_type`
- `file_name`
- `storage_key`
- `content_hash`
- `file_size`
- `created_at`

### 5.3 最小报告策略

第一阶段不单独建设 `reports` API。

第一版报告能力先体现在：

- 扫描任务级计数字段
- 问题列表过滤能力
- 问题详情和动作日志

如果这些基础视图都还不稳定，提前做趋势报表只会放大脏数据问题。

## 6. 最小统计视图

第一阶段必须进入系统的不是复杂 BI，而是治理闭环必须依赖的统计字段。

第一版最小统计包括：

- 每次扫描任务的 `total_issue_count`
- 每次扫描任务的 `new_issue_count`
- 每次扫描任务的 `reopened_issue_count`
- 每个问题的 `first_seen_at`
- 每个问题的 `last_seen_at`
- 每个问题的当前 `status`
- 每个问题的当前 `severity`

这些字段足以支持：

- 导入批次对比
- 问题池筛选
- 基础治理节奏跟踪

## 7. API 设计约束

为了先验证领域闭环，第一版 API 做这些收敛：

- 列表接口先返回数组，不在第一版引入复杂分页对象
- 查询接口优先做过滤语义稳定，而不是一次做完所有排序和统计参数
- 状态流转先收敛成一个统一接口
- 问题详情必须能带出动作日志
- 扫描任务详情必须能带出附件

这些约束的目的不是偷简化，而是避免第一版把接口面做得比领域模型还复杂。

## 8. 对实现的直接要求

这份设计对后续实现提出这些直接要求：

- 应用层必须提供扫描任务列表、扫描任务详情、问题列表、问题详情和状态流转用例
- 持久化层必须支持按治理维度过滤问题
- 状态流转成功后必须立即写动作日志
- 扫描任务详情必须能返回附件列表
- 导入结果必须更新扫描任务统计字段

## 9. 当前结论

第一阶段查询与工作流设计的重点不是“做很多页面”，而是把治理者最需要的视图和动作稳定下来。

只要扫描任务可查、问题可查、状态可流转、动作可追溯，这个系统就已经具备了第一版控制面的骨架。