# API Contracts

## 文档目的

这个目录用于管理第一版控制面 API 契约，而不是把接口列表散落在 README 和聊天记录里。

## 当前冻结入口

第一版最小闭环契约已经收敛到：

- [01-mvp-api.md](01-mvp-api.md)
- [examples/import-scan-sample.json](examples/import-scan-sample.json)

## 第一版冻结接口

### 健康检查

- `GET /healthz`
- `GET /readyz`

### 结果接入

- `POST /api/scan-tasks/import`

### 扫描任务查询

- `GET /api/scan-tasks`
- `GET /api/scan-tasks/{id}`

### 问题查询与处理

- `GET /api/issues`
- `GET /api/issues/{id}`
- `POST /api/issues/{id}/status`

## 当前契约原则

### 1. 先冻结字段语义，再扩接口数量

第一版优先稳定：

- 导入输入字段
- 问题查询过滤字段
- 状态流转语义
- 错误响应结构

### 2. 第一版不拆多个动作接口

`assign`、`ignore` 和 `resolve` 都先通过统一状态流转接口表达。

### 3. 列表接口先做过滤，不先做复杂分页

第一版先验证闭环和字段语义，列表直接返回数组。

### 4. 当前本地实现允许先用内存仓储

当前契约服务于领域验证和 API 打通。

运行时实现可以先用内存仓储，但契约字段和语义要与后续 PostgreSQL 版本保持一致。

## 下一步

- 如果第一版字段稳定，再补 OpenAPI
- 如果查询维度稳定，再补分页、排序和批量操作
- 如果进入发布控制阶段，再新增 release 相关契约