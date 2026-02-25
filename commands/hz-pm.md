---
description: 调用项目经理 agent 进行需求规划、任务编排或状态查询
argument-hint: <需求描述或查询>
---

# 项目经理调度

将需求或查询分发给 project-manager agent，完成需求分析、计划创建、任务编排或状态报告。

## 执行步骤

### 1. 收集项目上下文

使用 Bash 检查 docs 目录结构：

```bash
ls -la docs/ 2>/dev/null
ls docs/plans/ 2>/dev/null
```

使用 Glob 搜索所有计划的任务文件：

```
docs/plans/*/tasks.md
```

如找到 tasks.md 文件，使用 Read 读取每个文件内容，提取任务状态汇总。

使用 Read 读取以下文件（如存在）：
- `docs/specs.yaml` — 项目配置和编码规范
- `docs/project.md` — 项目概述

使用 Bash 获取 git 状态：

```bash
git status --short
git log --oneline -5
```

### 2. 判断请求类型

分析 `$ARGUMENTS`：

- 包含"进度"、"状态"、"怎么样"、"哪些" → **状态查询**
- 包含"添加"、"实现"、"开发"、"修复"、"重构"、"新建" → **新需求**
- 引用了已有计划路径 → **计划继续**
- 为空 → 默认**状态查询**，汇总最新计划进度

### 3. 调用 project-manager agent

使用 Task 工具（subagent_type: general-purpose），提供以下上下文：

```
## PM 任务

### 用户请求
$ARGUMENTS

### 请求类型
[状态查询 | 新需求 | 计划继续]

### 固定技术栈约束
本系统使用固定技术栈，brainstorming 和计划文档必须基于以下技术，不得替换：
- 后端: Go (Gin + GORM + Redis + Zap)，`server/` 目录，HZ 标记扩展
- 前端: React 19 + Vite + Tailwind CSS + Zustand，`web/` 目录
- 桌面端: Tauri 2
- 移动端: Capacitor
- 项目结构: `server/` + `web/` + `docs/`
如果用户的 PRD 中提到了其他技术选型，必须将需求映射到上述技术栈。

### 项目上下文
- 工作目录: [绝对路径]
- 已有计划: [docs/plans/ 下的目录列表]
- 项目规范: [specs.yaml 关键内容]
- 项目概述: [project.md 摘要]
- Git 状态: [分支 + 未提交变更]
- 最近提交: [5 条 commit]

### 当前任务状态
[各计划 tasks.md 的状态汇总]

### 期望产出
按 PM 报告格式输出：计划概览、编排决策、当前状态、下一步动作。
新需求 → 调用 brainstorming 后创建完整计划文档（架构和 specs.yaml 必须基于固定技术栈）。
状态查询 → 汇总所有相关任务完成情况。
```

### 4. 展示报告

将 PM agent 的产出完整展示给用户。

## 错误处理

- `docs/` 不存在 → 提示先初始化项目文档（`/project-docs`）
- `$ARGUMENTS` 为空且无计划 → 提示提供需求描述
