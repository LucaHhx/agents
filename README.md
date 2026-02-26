# HZ Multi-Agent 开发编排系统

通过 Claude Code 的 6 个专业 Agent 协作完成全栈开发（Go 后端 + React 19 前端）。

## 快速开始

```bash
# 1. 克隆此仓库
git clone https://github.com/LucaHhx/agents.git

# 2. 在你的项目中链接
cd /path/to/your-project
bash /path/to/agents/link.sh

# 3. 用 Claude Code 打开项目，使用斜杠命令
/hz-pm 我想添加用户认证功能
```

## Agent 团队

| 命令 | Agent | 模型 | 职责 |
|------|-------|------|------|
| `/hz-pm` | project-manager | opus | 需求规划、任务编排、文档管理 |
| `/hz-backend` | backend-dev | sonnet | Go 后端实现（API/Service/Model） |
| `/hz-frontend` | frontend-dev | sonnet | React 19 前端实现（组件/Store） |
| `/hz-design` | ui-designer | opus | 组件架构设计、交互规格定义 |
| `/hz-review` | code-reviewer | opus | 两阶段代码审查（只读） |
| `/hz-test` | qa-tester | sonnet | 测试执行、Bug 追踪 |

## 工作流

```
用户需求 → /hz-pm（brainstorming + 创建计划）
    → /hz-design（UI 规格，如有需要）
    → /hz-backend + /hz-frontend（实现）
    → /hz-review（代码审查）
    → /hz-test（测试验证）
    → PM 签收
```

每个里程碑完成后必须通过 `/hz-test` 门控测试才能进入下一阶段。

## 技术栈（固定）

| 层 | 技术 |
|---|---|
| 后端 | Go（Gin + GORM + Redis + Zap） |
| 前端 | React 19 + Vite + Tailwind CSS + Zustand |
| 桌面 | Tauri 2 |
| 移动 | Capacitor |
| 结构 | `server/` + `web/` + `docs/` |

## Skills

| Skill | 用途 |
|-------|------|
| brainstorming | 需求探索和方案设计 |
| project-init | 项目脚手架初始化 |
| project-docs | 文档结构管理 |
| subagent-driven-development | 顺序任务执行 + 审查 |
| dispatching-parallel-agents | 并行任务分发 |
| agent-browser | UI 自动化测试 |
| vercel-composition-patterns | React 组合模式 |
| vercel-react-best-practices | React 最佳实践 |

## 项目结构

```
agents/
├── CLAUDE.md          # 主配置（链接到目标项目的 .claude/）
├── link.sh            # 链接脚本
├── agents/            # Agent 定义
├── commands/          # 斜杠命令
└── skills/            # Skill 集合
```

## 已验证项目

- [LedgerX](https://github.com/LucaHhx/ledgerx) — 跨端记账应用（MVP 25/25 测试通过）
