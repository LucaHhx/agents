---
name: project-init
description: "Scaffold new full-stack projects from templates. Use when the user wants to: (1) Create a new project (triggers: '新建项目', '初始化项目', 'new project', 'init project', 'scaffold project', 'create project'), (2) Start a new Go + Tauri/React application. Generates a project with server/ (Go backend with Gin/GORM/Redis/Zap), web/ (Tauri + React + Vite + Tailwind + Zustand), and docs/ (via project-docs skill)."
---

# Project Init

Scaffold a new full-stack project with Go backend + Tauri/React frontend.

## Generated Structure

```
{project-name}/
├── server/     # Go (Gin + GORM + Redis + Zap)
├── web/        # Tauri + React (Vite + Tailwind + Zustand)
└── docs/       # Project docs (via project-docs skill)
```

## Usage

```bash
python scripts/init_project.py --name <project-name> --module <go-module-path> [--output <dir>]
```

### Parameters

- `--name` (required): kebab-case project name (e.g., `my-app`)
- `--module` (required): Go module path (e.g., `github.com/luca/my-app`)
- `--output` (optional): parent directory, defaults to cwd

## Workflow

1. Ask user for project name and Go module path (if not provided)
2. Run `scripts/init_project.py` with parameters
3. Invoke project-docs skill to initialize `{project-name}/docs/`
4. Run `go mod tidy` in server/ and `npm install` in web/
5. For mobile: `cd web && npx tauri ios init` / `npx tauri android init`

## Placeholder System

| Placeholder | Example |
|-------------|---------|
| `{{PROJECT_NAME}}` | `my-app` |
| `{{PROJECT_TITLE}}` | `My App` |
| `{{GO_MODULE}}` | `github.com/luca/my-app` |
| `{{APP_IDENTIFIER}}` | `com.hz.my-app` |

Go import paths: template uses `"server/..."`, replaced with `"{GO_MODULE}/..."` at init time.

## Server Template

Global variable prefix: `HZ_` (HZ_CONFIG, HZ_DB, HZ_REDIS, HZ_LOG, HZ_JSON).

Extension markers in inlet files:
- `// HZ:MIGRATE:*` — entrance/migrate.go
- `// HZ:ROUTER:*` — router/inlet.go
- `// HZ:API:*` — api/inlet.go
- `// HZ:SERVICE:*` — service/inlet.go

## Web Template

Minimal skeleton: React Router + Axios client + Loading component + empty HomePage.
No auth, no layout components. Add as needed per project.

## Post-Init Customization

See `references/customization-guide.md` for:
- Adding new backend modules
- Adding authentication
- Adding frontend pages and stores
- Multi-platform builds
