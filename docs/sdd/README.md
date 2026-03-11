# SDD 目录

这个目录存放“内容控制平台”的系统设计文档集合。

SDD 负责回答“系统怎么设计”，但它不负责替代所有其他文档。

## SDD 和其他文档的边界

- `product_intent`：回答为什么做、第一版做什么
- `roadmap / backlog`：回答先做什么
- `owner_decisions`：回答哪些问题还需要拍板
- `ADR`：回答关键设计为什么这么选
- `API`：回答对外接口怎么约定
- `data_model`：回答核心对象如何建模
- `testing_strategy`：回答如何验证设计是否成立
- `runbook`：回答系统如何运行和操作
- `SDD`：把这些输入收束成可实现的系统设计

## 当前直接入口

- 文档总览：[../README.md](../README.md)
- 项目背景：[../project_context.md](../project_context.md)
- 产品意图：[../product_intent.md](../product_intent.md)
- 路线图：[../roadmap.md](../roadmap.md)
- 工作项池：[../backlog.md](../backlog.md)
- 数据模型：[../data_model.md](../data_model.md)
- API 契约入口：[../api/README.md](../api/README.md)
- 当前待决策：[../owner_decisions/current.md](../owner_decisions/current.md)

## 当前阅读顺序

建议按下面顺序阅读：

1. [../project_context.md](../project_context.md)
2. [../product_intent.md](../product_intent.md)
3. [../roadmap.md](../roadmap.md)
4. [../owner_decisions/current.md](../owner_decisions/current.md)
5. [../data_model.md](../data_model.md)
6. [../api/README.md](../api/README.md)
7. [../architecture.md](../architecture.md)

## 接下来要拆出的 SDD 主题

后续会逐步拆成这些主题文档：

- `01-overview.md`：系统总览、边界和阶段目标
- `02-ingestion-and-issue-domain.md`：结果接入、问题模型和问题指纹
- `03-query-workflow-and-reporting.md`：查询、状态流转、统计和附件
- `04-release-control-and-audit.md`：发布控制、审批、回滚和审计
- `05-runtime-storage-and-ops.md`：运行方式、存储、异步任务和通知
- `06-testing-and-delivery-phases.md`：测试策略和分阶段交付

## 当前默认假设

在你还没有正式拍板前，系统设计先采用以下默认假设：

- 第一版实现范围优先收敛到质量门禁闭环
- 第一可运行版本优先打通 API 和导入链路
- 架构采用 Go 模块化单体
- Web 层框架优先 Gin，领域和持久化层保持解耦
- 数据层优先 PostgreSQL，对象产物优先 MinIO / S3

## 当前结论

这个目录已经作为后续 SDD 的入口建立完成。

下一步默认是把当前的架构草图拆成真正的主题化 SDD 文档，而不是继续把所有设计都堆在 `README.md` 里。
