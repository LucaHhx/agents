# HZ Multi-Agent 开发编排系统

本仓库是多 Agent 开发编排工具包，通过 `bash link.sh` 符号链接到目标项目的 `.claude/` 目录后使用。
6 个专业 Agent 通过斜杠命令协作完成全栈开发（Go 后端 + React 19 前端）。

## 语言

IMPORTANT: 所有 Claude 交互、任务分发和报告均使用中文。

## 开发工作流

完整生命周期：用户需求 → brainstorming → PM 创建计划 → 【每个里程碑：实现 → 里程碑验证】→ 全局审查 → 集成测试 → 签收

1. `/hz-pm <需求>` — PM 调用 brainstorming 探索需求，创建 `docs/plans/<name>/` 文档
2. PM 分析任务依赖，决定执行模式：
   - 有依赖 → 顺序执行（subagent-driven-development skill）
   - 独立任务 → 并行执行（dispatching-parallel-agents skill）
3. **每个里程碑内循环：**
   a. 如有 UI 需求 → `/hz-design` 先产出组件设计规格
   b. `/hz-backend` 和/或 `/hz-frontend` 执行实现
   c. **里程碑门控** → `/hz-test` 执行里程碑验证测试（启动服务，实际 HTTP 请求测试）
   d. 验证通过 → 进入下一里程碑；失败 → 修复后重新测试
4. 所有里程碑完成 → `/hz-review` 执行全局两阶段审查
5. `/hz-test` 执行端到端集成测试
6. PM 最终签收，更新 CHANGELOG

## 命令速查

| 命令 | Agent | 模型 | 用途 |
|------|-------|------|------|
| `/hz-init` | — | — | 初始化新项目并链接 HZ agents |
| `/hz-pm` | project-manager | opus | 需求规划、任务编排、文档管理 |
| `/hz-backend` | backend-dev | sonnet | Go 后端实现（API/Service/Model） |
| `/hz-frontend` | frontend-dev | sonnet | React 19 前端实现（组件/Store） |
| `/hz-design` | ui-designer | opus | 组件架构设计、交互规格定义 |
| `/hz-review` | code-reviewer | opus | 两阶段代码审查（只读，不修改代码） |
| `/hz-test` | qa-tester | sonnet | 测试执行、Bug 上报与追踪 |

## 技术栈约束（硬性门控）

IMPORTANT: 以下技术栈是本系统的固定约束，所有 Agent（包括 brainstorming 和 PM）在设计、规划和实现时必须严格遵守，不得自行替换为其他技术。

| 层 | 技术选型 | 备注 |
|---|---|---|
| 后端 | Go（Gin + GORM + Redis + Zap） | 通过 HZ 标记扩展 |
| 前端框架 | React 19 | 遵循 `.claude/rules/react19.md` |
| 桌面端 | Tauri 2 | Rust sidecar |
| 移动端 | Capacitor | iOS/Android |
| 构建工具 | Vite | 前端统一构建 |
| 状态管理 | Zustand | `{ state, actions, meta }` 模式 |
| 样式 | Tailwind CSS | 工具类优先，`cn()` 合并 |
| 项目结构 | `server/` + `web/` + `docs/` | 通过 `project-init` 初始化 |

**禁止行为：**
- brainstorming 时不得提议替换上述技术栈（如用 Expo 替换 Tauri、用 drizzle 替换 GORM、用 React Native 替换 React 19 + Capacitor）
- PM 创建计划时，架构文档和 specs.yaml 必须基于上述技术栈
- 如果 PRD 中的需求看似需要不同技术，brainstorming 应在上述技术栈约束内探索解决方案

## 关键规则

IMPORTANT: 以下规则必须严格遵守：

0. **文档实时更新（最高优先级）** — 任何 Agent 完成任务后，必须立即更新 `docs/plans/*/tasks.md` 中对应任务的状态（待办→已完成），并在 `docs/plans/*/changelog.md` 中添加变更记录。Git commit message 必须引用任务 ID。**不更新文档 = 任务未完成。** 这条规则适用于所有 Agent，无论是通过 PM 编排还是直接调用。
1. **技术栈不可替换** — 所有设计和实现必须使用上述固定技术栈，brainstorming 探索的是功能方案而非技术选型
2. **先 brainstorming 再实现** — 任何新功能必须先通过 brainstorming skill 探索需求，获得用户批准后才能开始实现
3. **Agent 通信全部经过 PM** — Agent 之间不直接分发任务，所有协调通过 project-manager
4. **审查链严格有序** — 规格合规审查通过 → 代码质量审查通过 → QA 测试，不可跳过或乱序
5. **分发时提供完整上下文** — 向 Agent 分发任务时内联完整任务描述和上下文，不要让 Agent 自己读文件
6. **code-reviewer 只读** — 审查 Agent 只报告问题不修改代码，修复由对应开发 Agent 完成
7. **里程碑门控测试** — 每个里程碑的所有任务完成后，必须启动服务并用实际 HTTP 请求/页面交互验证，通过后才能进入下一里程碑。不可跳过，不可推迟到最后统一测试

## 扩展系统

创建新组件时使用对应 skill 获取指导：

- 新 Agent → `agent-create` skill
- 新 Skill → `skill-creator` skill
- 新命令 → `command-creator` skill
- 发现社区 Skill → `find-skills` skill 或 `npx skills find <query>`

## 目标项目约定

通过 `project-init` skill 初始化的目标项目遵循以下结构：

- `server/` — Go 后端（Gin + GORM + Redis + Zap），通过 HZ 标记扩展
- `web/` — React 19 前端（Tauri/Capacitor + Vite + Tailwind + Zustand）
- `docs/` — 项目文档（通过 `project-docs` skill 管理计划、任务、决策记录）

HZ 标记扩展规范详见 `.claude/rules/hz-markers.md`。
React 19 强制规范详见 `.claude/rules/react19.md`。

## 详细参考

PM 编排架构: @agents/project-manager.md
需求探索流程: @skills/brainstorming/SKILL.md
顺序执行模式: @skills/subagent-driven-development/SKILL.md
并行执行模式: @skills/dispatching-parallel-agents/SKILL.md
文档管理结构: @skills/project-docs/SKILL.md
