---
description: 初始化新项目并链接 HZ 多 Agent 编排系统
argument-hint: <项目名> [--module <go-module-path>]
---

# 项目初始化

快速创建一个新的全栈项目并链接 HZ agents 框架。

## 执行步骤

### 1. 解析参数

从 `$ARGUMENTS` 提取：
- `项目名`（必填，kebab-case）
- `--module`（可选，Go module 路径，默认 `<项目名>/server`）

如缺少项目名，询问用户。

### 2. 调用 project-init skill

使用 `project-init` skill 的脚本初始化项目：

```bash
python3 skills/project-init/scripts/init_project.py --name <项目名> --module <module-path>
```

### 3. 链接 agents 框架

```bash
cd <项目名>
bash <agents-repo>/link.sh
```

### 4. 初始化 Git

```bash
cd <项目名>
git init
git add -A
git commit -m "init: scaffold project with HZ agents"
```

### 5. 安装依赖

```bash
cd <项目名>/server && go mod tidy
cd <项目名>/web && npm install
```

### 6. 展示结果

```markdown
## 🎉 项目创建完成

- 项目: <项目名>
- 后端: server/ (Go + Gin + GORM)
- 前端: web/ (React 19 + Vite + Tailwind)
- 文档: docs/
- HZ Agents: 已链接 ✅

### 下一步
1. `cd <项目名>`
2. 提供 PRD 或需求描述
3. `/hz-pm <需求>` 开始开发
```
