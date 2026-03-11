# Backlog

## 文档目的

本文档用于维护当前阶段的工作项池。

它和 `owner_decisions` 的区别是：

- `backlog` 管“要做什么”
- `owner_decisions` 管“哪些事情需要你拍板”

## 当前优先级

### Now

#### B-001 固定文档体系和阅读入口

- 状态：`DONE`
- 说明：完成项目背景、产品意图、流程、决策队列和文档地图。

#### B-002 拆第一版 SDD

- 状态：`IN_PROGRESS`
- 说明：把现有架构草图拆成主题化 SDD 文档。
- 已完成：`01-overview`、`02-ingestion-and-issue-domain`、`03-query-workflow-and-reporting`、`05-runtime-storage-and-ops`
- 剩余：`04-release-control-and-audit`、`06-testing-and-delivery-phases`

#### B-003 定义最小运维基线

- 状态：`DONE`
- 说明：明确运行方式、配置注入、日志、健康检查、迁移和失败恢复的当前阶段要求。
- 产出：`docs/sdd/05-runtime-storage-and-ops.md`

#### B-004 冻结第一版数据模型

- 状态：`DONE`
- 说明：把核心实体、唯一键、状态流转字段和索引建议提升为可实现设计。
- 产出：`docs/data_model.md`

#### B-005 拆查询、状态流转和报表 SDD

- 状态：`DONE`
- 说明：补完查询视图、状态流转、附件查看和最小统计设计。
- 产出：`docs/sdd/03-query-workflow-and-reporting.md`

#### B-006 搭 Go + Gin 项目骨架

- 状态：`DONE`
- 说明：建立 `cmd/api`、`internal` 和内存仓储版本的第一版服务骨架。
- 产出：`cmd/api`、`internal/api/http`、`internal/application/controlplane`、`internal/domain`、`internal/infrastructure/memory`

#### B-007 冻结第一版导入契约

- 状态：`DONE`
- 说明：定义最小导入 JSON 字段集、问题查询契约和状态流转契约。
- 产出：`docs/api/01-mvp-api.md`、`docs/api/examples/import-scan-sample.json`

#### B-008 打通导入到查询闭环

- 状态：`DONE`
- 说明：完成导入、扫描任务查询、问题列表、问题详情和状态流转最小闭环。
- 验证：`go test ./...`

### Next

#### B-009 引入 PostgreSQL 持久化和 migration

- 状态：`PLANNED`
- 说明：把内存仓储替换成 PostgreSQL 仓储，并补第一版 migration。

#### B-010 增加 API 集成测试和 HTTP 冒烟测试

- 状态：`PLANNED`
- 说明：补接口级测试，不只验证应用层用例。

#### B-011 补本地运行和导入验证 Runbook

- 状态：`DONE`
- 说明：已经补充本地启动、样例导入和基本验证步骤。
- 产出：`docs/runbook/README.md`

#### B-012 接入对象存储和真实附件上传

- 状态：`PLANNED`
- 说明：把附件从元信息模型推进到真实对象存储集成。

### Later

#### B-013 拆发布控制和审计 SDD

- 状态：`PLANNED`
- 说明：进入审批、环境提升、灰度、回滚和审计扩展设计。

#### B-014 增加问题观测历史和趋势统计

- 状态：`PLANNED`
- 说明：引入 `issue_occurrence` 或等价观测模型，为趋势统计打底。

## 使用规则

- 新工作优先进入 `backlog`，而不是直接堆进 `README`
- 真正影响方向的事项才进入 `owner_decisions`
- backlog 如果和仓库现状不一致，应优先修 backlog