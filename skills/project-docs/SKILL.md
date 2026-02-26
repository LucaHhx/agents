---
name: project-docs
description: "Manage structured project documentation throughout the full lifecycle. Use when the user wants to: (1) Initialize project docs (triggers: '初始化项目文档', '创建文档结构', 'init docs'), (2) Create a new plan (triggers: '新建计划', '创建计划', 'new plan'), (3) Update plan status, record decisions, log changes, or manage test records in the docs/ structure."
---

# Project Docs

Manage a structured `docs/` directory for project documentation and plan tracking.

## Directory Structure

### Root (`docs/`)

```
docs/
├── project.md       # 项目概览
├── specs.yaml       # 项目规范
├── CHANGELOG.md     # 项目变更日志
├── glossary.md      # 术语表
└── plans/           # 计划目录
```

### Each Plan (`docs/plans/<plan-name>/`)

```
docs/plans/<plan-name>/
├── plan.md          # 总计划 + 目标 + 时间线 (What/Why/When)
├── tasks.md         # 任务分解 + 进度 (How/Who/Progress)
├── decisions.md     # 决策记录 + 风险 (ADR 格式)
├── changelog.md     # 变更日志
├── testing.md       # 测试计划 + 结果
└── testing/         # 测试资产 (截图等)
```

## Workflows

### 1. Initialize Docs

Trigger: User says "初始化项目文档", "创建文档结构", or similar.

1. Run `scripts/init_docs.py <project-root>`
2. Fill in `docs/project.md` based on existing project files (README, package.json, etc.) and user context
3. Fill in `docs/specs.yaml` based on existing configs (.editorconfig, linting, git conventions)
4. Leave CHANGELOG.md and glossary.md with initial entries

### 2. New Plan

Trigger: User says "新建计划", "创建计划", or similar.

1. Ask user for plan name (kebab-case) and brief description if not provided
2. Create plan directory with all required files. Two options:
   - **Option A (recommended):** Run `scripts/new_plan.py <project-root>/docs/plans <plan-name>` — auto-generates all files from templates
   - **Option B:** Manually create files — but you MUST create ALL of them
3. Fill in `plan.md` with goals, scope, and timeline based on user context
4. Fill in `tasks.md` with initial task breakdown
5. Fill in `decisions.md` with initial design decisions (DO NOT leave empty)
6. Fill in `changelog.md` with plan creation entry
7. Fill in `testing.md` with test plan outline

<HARD-GATE>
After creating a plan, VERIFY completeness: `ls docs/plans/<name>/` must show ALL of: plan.md, tasks.md, decisions.md, changelog.md, testing.md, testing/. Missing any file = plan creation incomplete.
</HARD-GATE>

### 3. Update Operations

Trigger: User asks to update status, record a decision, log a change, or add test results.

Read `references/update-operations.md` for detailed guidance on each operation type and cross-file consistency rules.

## Script Usage

```bash
# Initialize docs structure
python scripts/init_docs.py <project-root>

# Create a new plan
python scripts/new_plan.py <plans-directory> <plan-name>
```

## Conventions

- Plan names: kebab-case (e.g., `backend-api`, `user-auth`)
- Date format: YYYY-MM-DD
- Task status: 待办 / 进行中 / 已完成 / 已取消
- Decision status: 待定 / 已决定 / 已废弃
- Change types: 新增 / 变更 / 修复 / 移除
- Decision numbering: DR-001, DR-002, ...
