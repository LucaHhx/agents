---
description: 调用代码审查 agent 执行两阶段审查（规格合规 + 代码质量）
argument-hint: [审查范围或文件路径]
---

# 代码审查调度

将审查任务分发给 code-reviewer agent，执行两阶段审查：第一阶段验证规格合规，第二阶段评估代码质量。

## 执行步骤

### 1. 确定审查范围

如果 `$ARGUMENTS` 提供了具体文件路径：
- 使用指定文件作为审查范围

如果 `$ARGUMENTS` 为空或为通用描述：
- 使用 Bash 获取最近变更文件：

```bash
git diff --name-only HEAD~1..HEAD
```

- 如无结果，获取未提交变更：

```bash
git diff --name-only
git diff --name-only --cached
```

如果仍无变更文件 → 提示当前没有可审查的变更。

### 2. 获取 Git 信息

使用 Bash 获取 SHA 和差异：

```bash
git rev-parse HEAD
git rev-parse HEAD~1
git log --oneline -5
git diff HEAD~1..HEAD --stat
```

使用 Bash 获取详细 diff（变更量大时按文件分段）：

```bash
git diff HEAD~1..HEAD
```

### 3. 收集任务需求

使用 Grep 搜索最新计划 tasks.md 中与变更文件相关的任务描述：

```
搜索路径: docs/plans/*/tasks.md
```

使用 Read 读取 `docs/specs.yaml` — 编码规范全文。

### 4. 读取变更文件

使用 Read 逐个读取所有变更文件的当前内容。

按文件后缀分类：
- `.go` → 后端代码，检查 HZ 标记、Gin/GORM、Zap 日志
- `.tsx`/`.ts` → 前端代码，检查 React 19、复合组件、Tailwind、Zustand

### 5. 调用 code-reviewer agent

使用 Task 工具（subagent_type: general-purpose），提供以下上下文：

```
## 代码审查任务

### 审查范围
[变更文件列表 + 类型分类]

### 第一阶段：规格合规审查
#### 任务需求全文
[tasks.md 中的对应任务描述]
（如无明确任务需求，标注"无对应计划任务，仅进行代码质量审查"并跳至第二阶段）

#### 变更文件内容
[每个变更文件的完整内容]

### 第二阶段：代码质量审查
#### BASE_SHA
[commit hash]

#### HEAD_SHA
[commit hash]

#### Git Diff
[完整代码差异]

### 项目编码规范
[specs.yaml 全文]

### 审查要求
1. 第一阶段 — 逐项对照需求，验证完整性和正确性
2. 第一阶段通过后执行第二阶段 — 评估代码质量
3. 每个问题标注 file:line 引用
4. 按严重程度分类：严重 / 重要 / 次要
5. Go 代码: HZ 标记、Gin handler、GORM 规范、Zap 日志、错误包装
6. React 代码: 复合组件、React 19 API、Zustand、Tailwind
```

### 6. 展示审查报告

将两阶段审查报告展示给用户。

如第一阶段"不通过" → 高亮遗漏/误解项，建议分发给开发者修复。
如第二阶段"需要修改" → 按严重程度排序展示，严重问题必须修复。

## 错误处理

- git 历史不足（首次 commit）→ BASE_SHA 使用空树 hash `4b825dc642cb6eb9a060e54bf899d69f82ecf73`
- 变更文件列表为空 → 提示当前无可审查变更
- 指定文件不存在 → 报告未找到，列出可用文件
