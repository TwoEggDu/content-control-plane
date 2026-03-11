# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- 建立项目文档地图，补充背景、计划、设计、交付四层文档结构。
- 增加 `roadmap`、`backlog`、`data_model`、`testing_strategy` 和 `runbook` 入口文档。
- 增加 `adr/README.md` 与 `ADR-0001`，固定文档驱动的推进方式。
- 增加 `api/README.md`，作为后续接口契约入口。
- 固定第一版 Web 框架为 `Gin`，并补充对应 ADR。
- 新增 `01-overview`、`02-ingestion-and-issue-domain`、`03-query-workflow-and-reporting` 与 `05-runtime-storage-and-ops`，把 SDD 从草图入口推进到主题文档阶段。
- 冻结第一版数据模型和 MVP API 契约，并补充导入样例 JSON。
- 建立 `Gin` API 服务骨架、内存仓储和最小闭环测试。
- 增加 `migrations/README.md` 和 `.gitignore`，为后续持久化实现留出入口。