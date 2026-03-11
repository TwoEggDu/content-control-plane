# Docs Map

这个目录存放这个应用从“为什么做”到“怎么做”、再到“怎么交付和运行”的完整文档体系。

## 文档分层

### 1. 背景与边界

这层回答“为什么要做”和“第一版做什么”。

- [project_context.md](project_context.md)：项目背景、个人能力来源、仓库定位
- [product_intent.md](product_intent.md)：产品意图、边界、MVP、成功标准

### 2. 计划与决策

这层回答“先做什么”和“哪些问题需要拍板”。

- [autonomous_workflow.md](autonomous_workflow.md)：默认推进规则和上提机制
- [roadmap.md](roadmap.md)：阶段目标和推进顺序
- [backlog.md](backlog.md)：当前工作项池和优先级
- [owner_decisions/current.md](owner_decisions/current.md)：当前待你决定的问题

### 3. 设计与契约

这层回答“系统怎么设计”和“对外怎么约定”。

- [sdd/README.md](sdd/README.md)：系统设计入口
- [architecture.md](architecture.md)：当前架构草图基线
- [data_model.md](data_model.md)：核心数据对象与建模原则
- [api/README.md](api/README.md)：API 契约入口和约定
- [adr/README.md](adr/README.md)：关键架构决策记录

### 4. 交付与运行

这层回答“怎么验证、上线和维护”。

- [testing_strategy.md](testing_strategy.md)：测试分层和验证策略
- [runbook/README.md](runbook/README.md)：运行和操作手册入口
- [../CHANGELOG.md](../CHANGELOG.md)：变更历史

## 推荐阅读顺序

第一次进入仓库，建议按下面顺序读：

1. [project_context.md](project_context.md)
2. [product_intent.md](product_intent.md)
3. [autonomous_workflow.md](autonomous_workflow.md)
4. [owner_decisions/current.md](owner_decisions/current.md)
5. [roadmap.md](roadmap.md)
6. [sdd/README.md](sdd/README.md)
7. [adr/README.md](adr/README.md)
8. [testing_strategy.md](testing_strategy.md)
9. [runbook/README.md](runbook/README.md)

## 当前项目阶段

当前项目仍处在“文档驱动的 M0/M1 过渡阶段”：

- M0：固定文档体系、开发流程和关键待决策事项
- M1：进入质量门禁控制面的第一个可运行闭环

## 更新原则

- 新增需求前，先判断它属于哪个层级
- 可逆的日常假设，不要写 ADR
- 会影响实现顺序的事项，先进入 backlog
- 会影响产品边界或核心设计的事项，进入 owner decisions
- 当设计开始稳定时，再把草图沉淀进 SDD 和 API/Data Model 文档
