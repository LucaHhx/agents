---
description: 调用 QA 测试 agent 执行全面测试验证并记录结果
argument-hint: [测试范围或计划路径]
---

# QA 测试调度

将测试任务分发给 qa-tester agent，执行测试计划、运行测试套件、验证集成功能并记录结果。

## 执行步骤

### 1. 定位测试计划

如果 `$ARGUMENTS` 提供了计划路径（如 `docs/plans/xxx/`）：
- 使用 Read 读取该目录下的 `testing.md`

如果 `$ARGUMENTS` 为通用描述或为空：
- 使用 Glob 搜索：

```
docs/plans/*/testing.md
```

- 取最新修改的 testing.md
- 如无 testing.md → 标注"无测试计划，从任务描述创建基础测试计划"

使用 Read 读取测试计划全文。

### 2. 收集实现摘要

使用 Read 读取对应计划的：
- `tasks.md` — 已完成任务描述
- `changelog.md` — 最近实现记录

使用 Bash 获取最近变更文件：

```bash
git log --name-only --oneline -10
```

### 3. 检测测试环境

使用 Bash 收集环境信息：

```bash
go version 2>/dev/null || echo "Go 未安装"
node --version 2>/dev/null || echo "Node 未安装"
uname -s
```

使用 Bash 确认后端和前端目录：

```bash
ls server/main.go 2>/dev/null || ls backend/main.go 2>/dev/null || echo "无后端"
ls web/package.json 2>/dev/null || ls frontend/package.json 2>/dev/null || echo "无前端"
```

### 4. 管理测试服务（pm-mcp）

**IMPORTANT: 测试前必须确保服务运行。使用 pm-mcp MCP 工具管理服务，不要用 Bash 直接启动。**

使用 `mcp__pm-mcp__list_processes` 检查已运行的服务。

如果后端服务未运行：
1. 使用 `mcp__pm-mcp__start_process` 启动后端：
   - name: `<项目名>-backend`
   - command: `go run main.go`
   - cwd: `<项目根>/server`
2. 使用 `mcp__pm-mcp__grep_logs` 搜索 `listening|started|:8080` 确认就绪

如果前端服务未运行且需要 UI 测试：
1. 使用 `mcp__pm-mcp__start_process` 启动前端：
   - name: `<项目名>-frontend`
   - command: `npm run dev`
   - cwd: `<项目根>/web`
2. 使用 `mcp__pm-mcp__grep_logs` 搜索 `ready in|Local:|localhost` 确认就绪

记录服务器进程 ID 和端口，传递给 QA agent。

### 5. 收集端点信息

使用 Grep 搜索后端 API 端点：

```
搜索模式: \.GET\(|\.POST\(|\.PUT\(|\.DELETE\(
路径: server/router/ 或 backend/router/
输出模式: content
```

使用 Grep 搜索前端路由/页面：

```
搜索模式: path:|Route|route
路径: web/src/ 或 frontend/src/
```

### 6. 调用 qa-tester agent

使用 Task 工具（subagent_type: general-purpose），提供以下上下文：

```
## QA 测试任务

### 测试范围
$ARGUMENTS
（如为空，执行测试计划中的全部项目）

### 测试计划
[testing.md 完整内容]
（如无测试计划 → "无测试计划，从已完成任务创建基础测试计划后执行"）

### 实现摘要
[已完成任务列表及描述]
[最近 changelog 条目]
[最近变更文件列表]

### 测试环境
- 操作系统: [uname 输出]
- Go 版本: [go version 输出]
- Node 版本: [node --version 输出]
- 后端目录: [路径或"无"]
- 前端目录: [路径或"无"]

### 服务状态
- 后端服务: [运行中/未运行] — 进程 ID: [id] — 端口: [port]
- 前端服务: [运行中/未运行] — 进程 ID: [id] — 端口: [port]
（如果服务已通过 pm-mcp 启动，提供进程 ID 供 QA agent 查看日志）

### API 端点列表
[从路由文件提取的端点]

### UI 页面列表
[从前端路由提取的页面]

### 计划路径
[docs/plans/xxx/]

### 执行要求
1. 按 testing.md 测试项依次执行
2. 后端单元测试: `cd [server] && go test ./... -v`
3. 前端单元测试: `cd [web] && npm test`（如有前端）
4. API 集成测试: 使用 curl 测试每个端点（服务器已在 [port] 运行）
5. UI 集成测试: 使用 agent-browser 测试前端页面（如有前端）
   - `agent-browser open http://localhost:[前端端口]`
   - `agent-browser snapshot -i` 获取交互元素
   - 验证页面渲染和交互
   - `agent-browser screenshot` 截图取证
6. 补充边界测试
7. 使用 TodoWrite 跟踪进度 — 每完成一项立即标记
8. 在 testing.md 记录结果
9. Bug → tasks.md 创建待办 + testing.md 添加已知问题
10. 如需查看服务器日志: `mcp__pm-mcp__get_logs(进程ID, lines=30)` 或 `mcp__pm-mcp__grep_logs(进程ID, pattern="error")`

### 期望产出
1. 环境信息
2. 测试执行摘要（总数/通过/失败/跳过）
3. 单元测试结果（后端 + 前端）
4. 集成测试结果表格
5. 发现的 Bug（严重程度 + 建议任务）
6. 已知问题更新
7. 总体结论: 通过 | 不通过
```

### 7. 展示测试报告

将 QA 报告完整展示给用户。

如"不通过" → 高亮 Bug 列表，建议通过 `/hz-pm` 协调修复。
如"通过" → 确认可进入最终签收，建议通知 PM。

## 错误处理

- 后端服务无法启动 → 通过 `mcp__pm-mcp__get_logs` 查看启动日志，跳过集成测试，仍执行单元测试
- 前端构建失败 → 报告为严重 Bug，跳过 UI 测试
- 端口被占用 → 通过 `mcp__pm-mcp__grep_logs` 搜索 `EADDRINUSE`，终止旧进程后重启
- agent-browser 未安装 → 降级为 curl 测试 API，跳过 UI 测试并在报告中注明
- `$ARGUMENTS` 和测试计划都为空 → 作为降级方案运行 `go test ./...` + `npm test`

## 文档更新（必须）

Agent 完成任务后，必须确认以下文档已更新：
1. `docs/plans/*/tasks.md` — 对应任务状态更新为「已完成」
2. `docs/plans/*/changelog.md` — 添加变更条目
3. Git commit message 引用任务 ID

如果 Agent 的报告中未包含文档更新确认，提醒 Agent 补充。
