---
description: 调用前端开发 agent 实现 React 19 前端功能（组件/Store/样式）
argument-hint: <实现任务描述>
---

# 前端开发调度

将实现任务分发给 frontend-dev agent，使用 React 19 + Tailwind + Zustand + Tauri 技术栈完成前端功能。

## 执行步骤

### 1. 定位前端目录

使用 Bash 查找前端根目录：

```bash
ls -d web/ 2>/dev/null || ls -d frontend/ 2>/dev/null || ls -d client/ 2>/dev/null
```

记录前端目录的绝对路径。如均不存在，报告错误并建议先初始化项目。

### 2. 查找 UI 设计规格

使用 Grep 在计划目录搜索 ui-designer 产出的设计规格：

```
搜索模式: UI 设计规格|组件架构|Component Architecture
路径: docs/plans/
```

如找到设计文档，使用 Read 读取完整内容。

使用 Grep 搜索 `docs/plans/*/decisions.md` 中的 UI 设计决策（搜索 "UI"、"组件"、"设计"）。

如未找到设计规格，标注"无上游设计规格"。

### 3. 查找后端 API 契约

使用 Grep 搜索后端路由定义：

```
搜索模式: \.GET\(|\.POST\(|\.PUT\(|\.DELETE\(
路径: server/router/ 或 backend/router/
输出模式: content
```

使用 Grep 搜索请求/响应结构体：

```
搜索模式: type.*Request struct|type.*Response struct
路径: server/ 或 backend/
```

如未找到，标注"API 未就绪，需创建模拟数据层"。

### 4. 扫描现有前端结构

使用 Glob 盘点文件：

```
[前端目录]/src/components/**/*.tsx
[前端目录]/src/pages/**/*.tsx
[前端目录]/src/stores/**/*.ts
[前端目录]/src/hooks/**/*.ts
[前端目录]/src/api/**/*.ts
[前端目录]/src/types/**/*.ts
```

使用 Read 读取以下文件（如存在）：
- Tailwind 配置文件
- `src/api/client.ts` — Axios 客户端配置

### 5. 调用 frontend-dev agent

使用 Task 工具（subagent_type: general-purpose），提供以下上下文：

```
## 前端实现任务

### 完整任务描述
$ARGUMENTS

### 工作目录
[前端目录绝对路径]

### UI 设计规格
[ui-designer 的完整设计产出，或标注"无上游设计规格"]

### 后端 API 契约
[API 端点列表 + 请求/响应类型，或标注"API 未就绪"]

### 现有前端结构
- 组件: [components/ 文件列表]
- 页面: [pages/ 文件列表]
- Store: [stores/ 文件列表]
- Hook: [hooks/ 文件列表]
- API: [api/ 文件列表]
- 类型: [types/ 文件列表]

### 设计 token
[tailwind.config theme.extend 内容]

### API 客户端
[client.ts 配置摘要]

### 期望产出
1. 实现内容（组件、Store、Provider 文件路径）
2. 组件架构（复合组件结构、Provider/Consumer 关系）
3. 测试结果（编写数量、通过数量）
4. 性能考量（懒加载、记忆化、包体积）
5. 自审发现
6. 疑问或顾虑

实现完成后运行 `npx tsc`（无错误）和 `npm test`（全部通过），然后提交 commit。
```

### 6. 展示实现报告

将 frontend-dev agent 的报告完整展示给用户。如标注了"API 未就绪"，提醒后续需集成。

## 错误处理

- 前端目录不存在 → 提示先初始化项目
- `$ARGUMENTS` 为空 → 搜索 tasks.md 中 frontend-dev 负责的未完成任务；如无则提示

## 文档更新（必须）

Agent 完成任务后，必须确认以下文档已更新：
1. `docs/plans/*/tasks.md` — 对应任务状态更新为「已完成」
2. `docs/plans/*/changelog.md` — 添加变更条目
3. Git commit message 引用任务 ID

如果 Agent 的报告中未包含文档更新确认，提醒 Agent 补充。
