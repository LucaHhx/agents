---
description: 调用后端开发 agent 实现 Go 后端功能（API/Service/Model/Cache）
argument-hint: <实现任务描述>
---

# 后端开发调度

将实现任务分发给 backend-dev agent，使用 Gin + GORM + Redis + Zap 技术栈完成 Go 后端功能。

## 执行步骤

### 1. 定位后端目录

使用 Bash 查找后端根目录：

```bash
ls -d server/ 2>/dev/null || ls -d backend/ 2>/dev/null
```

记录后端目录的绝对路径。如均不存在，报告错误并建议先初始化项目。

### 2. 扫描 HZ 标记位置

使用 Grep 搜索所有 HZ 标记：

```
搜索模式: // HZ:
路径: [后端目录]
输出模式: content（含行号）
```

重点提取以下入口文件的标记行号：
- `router/inlet.go` — HZ:ROUTER:*
- `api/inlet.go` — HZ:API:*
- `service/inlet.go` — HZ:SERVICE:*
- `entrance/migrate.go` — HZ:MIGRATE:*

### 3. 扫描现有模块

使用 Glob 列出现有模块目录：

```
[后端目录]/router/*/
[后端目录]/api/*/
[后端目录]/service/*/
[后端目录]/models/*/
```

### 4. 收集规范和决策

使用 Read 读取 `docs/specs.yaml`（如存在），提取 Go 编码规范。

使用 Grep 搜索 `docs/plans/*/decisions.md` 中的后端相关决策。

使用 Bash 获取 git 状态：

```bash
git status --short
git log --oneline -3
```

### 5. 调用 backend-dev agent

使用 Task 工具（subagent_type: general-purpose），提供以下上下文：

```
## 后端实现任务

### 完整任务描述
$ARGUMENTS

### 工作目录
[后端目录绝对路径]

### HZ 标记位置
[每个入口文件的标记及行号列表]

### 现有模块结构
[router/api/service/models 下的模块列表]

### 项目编码规范
[specs.yaml Go 规范摘要]

### 相关决策
[decisions.md 中的后端相关 DR]

### Git 状态
- 分支: [当前分支]
- 未提交变更: [有/无]

### 期望产出
1. 实现内容（创建/修改的文件列表）
2. 使用的 HZ 标记（file:line 引用）
3. 测试结果（编写数量、通过数量）
4. 自审发现
5. 疑问或顾虑

实现完成后运行 `go test ./...` 确保通过，然后提交 commit。
```

### 6. 展示实现报告

将 backend-dev agent 的报告完整展示给用户。如有疑问或顾虑，高亮展示。

## 错误处理

- 后端目录不存在 → 提示先初始化项目（project-init）
- HZ 标记不存在 → 标注"无 HZ 标记，按常规 Go 项目结构实现"
- `$ARGUMENTS` 为空 → 搜索 tasks.md 中 backend-dev 负责的未完成任务；如无则提示
