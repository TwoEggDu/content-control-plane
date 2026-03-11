# ADR 目录

这个目录用于记录这个项目中已经做出的关键架构和流程决策。

## ADR 是什么

ADR 是 Architecture Decision Record。

它解决的问题不是“系统整体怎么设计”，而是：

- 某个重要选择为什么这么定
- 当时有哪些备选方案
- 这个选择会带来什么后果

## ADR 不是用来干什么的

ADR 不应该拿来记录：

- 还没想清楚的随手想法
- 可以轻易改回来的小默认值
- 已经在 SDD 中完整解释且没有分歧的常规设计

## 什么时候写 ADR

满足以下任一条件时，建议写 ADR：

- 会影响多个模块或较长时期实现
- 一旦开始实现，回滚成本明显变高
- 有多个合理方案，需要留下取舍依据
- 未来很容易被问“当时为什么这么选”

## 状态值

- `PROPOSED`
- `ACCEPTED`
- `DEPRECATED`
- `SUPERSEDED`

## 当前 ADR 列表

- [0001-document-driven-workflow.md](0001-document-driven-workflow.md)：采用文档驱动的项目推进方式

## 推荐阅读顺序

1. 先读 [../product_intent.md](../product_intent.md)
2. 再读 [../sdd/README.md](../sdd/README.md)
3. 最后按时间顺序阅读 ADR
