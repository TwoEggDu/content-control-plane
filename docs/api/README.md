# API Contracts

## 文档目的

这个目录用于管理 API 契约，而不是临时把接口列表散落在 README 和聊天记录里。

## 当前阶段的使用方式

项目还没有进入正式接口实现阶段，因此这里先定义“接口契约入口”和基本约定。

等第一版 SDD 稳定后，这里再补：

- OpenAPI 文档
- 请求和响应样例
- 错误码约定
- 导入契约样例

## 第一版优先接口族

### 结果接入

- `POST /api/scan-tasks/import`

### 扫描任务查询

- `GET /api/scan-tasks`
- `GET /api/scan-tasks/{id}`

### 问题查询与处理

- `GET /api/issues`
- `GET /api/issues/{id}`
- `POST /api/issues/{id}/assign`
- `POST /api/issues/{id}/status`
- `POST /api/issues/{id}/ignore`

### 资源与统计

- `GET /api/resources/{id}/issues`
- `GET /api/reports/summary`
- `GET /api/reports/trend`

## 当前契约原则

### 1. 输入契约优先明确幂等字段

导入类接口必须先定义：

- 哪些字段参与幂等判定
- 哪些字段用于来源追踪
- 哪些字段缺失时应拒绝导入

### 2. 查询接口优先支持治理维度过滤

问题列表接口至少应考虑这些维度：

- 项目
- 扫描任务
- 规则
- 资源
- 负责人
- 状态
- 严重级别

### 3. 状态流转接口必须保留操作痕迹

任何修改问题状态的接口，都不应只改主表，还应触发动作日志记录。

### 4. 第一版不要过度追求协议复杂度

在需求还没稳定前，先把字段语义和边界讲清楚，比一上来做非常完整的 OpenAPI 更重要。

## 下一步

- 在 SDD 稳定后，把第一版接口清单升级为正式契约
- 补充样例请求体和响应体
- 给导入接口准备 fixture 文件
