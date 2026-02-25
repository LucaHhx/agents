---
description: 调用 UI 设计师 agent 进行组件架构设计和交互规格产出
argument-hint: <设计任务描述>
---

# UI 设计调度

将设计任务分发给 ui-designer agent，产出可直接用于前端实现的组件架构和交互规格。

## 执行步骤

### 1. 盘点现有前端组件

使用 Glob 扫描现有组件文件：

```
web/src/components/**/*.tsx
web/src/pages/**/*.tsx
web/src/stores/**/*.ts
web/src/hooks/**/*.ts
```

如 `web/` 不存在，尝试 `src/components/**/*.tsx` 等备选路径。

记录找到的组件列表（名称 + 路径）。

### 2. 获取设计 token

使用 Glob 搜索 Tailwind 配置文件：

```
**/tailwind.config.*
```

如找到，使用 Read 读取其中 `theme.extend` 的内容。

### 3. 收集计划上下文

使用 Glob 搜索计划文档：

```
docs/plans/*/decisions.md
docs/plans/*/plan.md
```

使用 Grep 搜索 decisions.md 中的 UI 设计相关决策（搜索 "UI"、"设计"、"组件"）。

使用 Read 读取 `docs/specs.yaml`（如存在），提取前端规范。

### 4. 调用 ui-designer agent

使用 Task 工具（subagent_type: general-purpose），提供以下上下文：

```
## UI 设计任务

### 任务描述
$ARGUMENTS

### 现有组件清单
[按目录分组的文件列表]
- components/: [...]
- pages/: [...]
- stores/: [...]
- hooks/: [...]

### 设计系统现状
- Tailwind 配置: [theme.extend 内容]
- 已有复合组件: [如有]

### 计划上下文
- 相关需求: [tasks.md 中的 UI 任务]
- 已有设计决策: [decisions.md 中的 UI 相关 DR]

### 技术栈约束
- React 19 + Tailwind CSS + Zustand
- 桌面端: Tauri; 移动端: Capacitor
- 组件路径: src/components/<feature>/<Name>.tsx
- Store 路径: src/stores/<name>.ts

### 期望产出
1. 组件架构（组件树、Provider 结构、复合组件 API）
2. 各组件规格（名称、路径、Props 接口、状态形态、Context 契约）
3. 交互流程（用户操作 → 状态变化 → 视觉反馈）
4. 响应式策略（桌面 → 平板 → 移动）
5. 设计决策记录（DR 条目）
```

### 5. 展示设计规格

将设计产出展示给用户，标注该规格可传递给 `/hz-frontend` 实现。

## 错误处理

- 前端目录不存在 → 提示先初始化项目
- `$ARGUMENTS` 为空 → 搜索 tasks.md 中 design 类型的未完成任务；如无则提示提供描述
