# Runbook

## 文档目的

这个目录用于存放运行和操作手册。

Runbook 回答的不是“系统怎么设计”，而是：

- 这个系统怎么启动
- 常见操作怎么执行
- 出问题时先查什么
- 哪些步骤需要标准化

相关设计基线见 [../sdd/05-runtime-storage-and-ops.md](../sdd/05-runtime-storage-and-ops.md)。

## 当前可执行手册

### 1. 本地启动

当前仓库已经有第一版 `Gin` API 骨架，默认使用内存仓储。

本地启动：

```powershell
go run ./cmd/api
```

可选环境变量：

- `APP_PORT`：默认 `8080`

服务启动后可直接访问：

- `GET /healthz`
- `GET /readyz`
- `POST /api/scan-tasks/import`
- `GET /api/scan-tasks`
- `GET /api/issues`

### 2. 样例导入验证

样例输入文件：

- [../api/examples/import-scan-sample.json](../api/examples/import-scan-sample.json)

验证步骤：

1. 启动服务
2. 用样例 JSON 调用 `POST /api/scan-tasks/import`
3. 调用 `GET /api/scan-tasks`
4. 调用 `GET /api/issues`
5. 用 `POST /api/issues/{id}/status` 验证状态流转

### 3. 当前限制

- 当前实现使用内存仓储，重启进程后数据不会保留
- 还没有 PostgreSQL migration
- 还没有对象存储真实接入，附件目前只保留元信息

## 后续建议补充的 Runbook

### 1. PostgreSQL 启动和迁移执行手册

适用时机：数据库持久化引入后。

### 2. 对象存储检查手册

适用时机：附件和报告接入后。

### 3. 发布和回滚手册

适用时机：发布控制模块落地后。

## 维护规则

- 重复超过两次的手工操作，应考虑沉淀为 Runbook
- 出现事故或排障过程后，应优先补 Runbook
- Runbook 应偏操作步骤，不应写成大段设计说明