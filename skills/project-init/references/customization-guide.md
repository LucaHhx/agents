# Customization Guide

## Adding a New Backend Module

Each business module follows the pattern: router → api → service → model.

### 1. Create model

```
server/models/{module}/
├── {entity}.go          # Domain model struct
├── request/{entity}.go  # Request DTOs
└── response/{entity}.go # Response DTOs
```

Model struct should embed `common.HZ_CRUD` (with soft delete: also `common.HZ_CRUD_DEL`).

### 2. Create service

```
server/service/{module}/{entity}.go
```

Add to `service/inlet.go`:
- Import: `"{{GO_MODULE}}/service/{module}"`
- Field: `{Module} {module}.ServiceGroup`

### 3. Create API handler

```
server/api/{module}/{entity}.go
```

Add to `api/inlet.go`:
- Import: `"{{GO_MODULE}}/api/{module}"`
- Field: `{Module} {module}.ApiGroup`

### 4. Create router

```
server/router/{module}/router.go
```

Implement `InitRouter(r *gin.RouterGroup)`.

Add to `router/inlet.go`:
- Import: `"{{GO_MODULE}}/router/{module}"`
- Field in V1RouterGroup
- Call in InitRouter

### 5. Register migration

In `entrance/migrate.go`, add model to `AutoMigrate()` call at `// HZ:MIGRATE:MODEL_LIST`.

## Adding Authentication

1. Add `config/jwt.go` with JWT config struct
2. Add JWT field to `config/config.go` Config struct
3. Add jwt section to `local.yaml`
4. Add `middleware/auth.go` with JWT validation logic
5. In `entrance/api.go`, create `authRouter` group with `middleware.Auth()`
6. Update `router/inlet.go` InitRouter to accept both public and auth router groups

## Adding Frontend Pages

### New page

Create `web/src/pages/{PageName}.tsx`, add route in `App.tsx`.

### New Zustand store

```typescript
// web/src/stores/{name}Store.ts
import { create } from 'zustand';

interface {Name}State {
  // state fields
  // action methods
}

export const use{Name}Store = create<{Name}State>((set, get) => ({
  // implementation
}));
```

### New API module

```typescript
// web/src/api/{name}.ts
import client from './client';

export const {name}Api = {
  getAll: () => client.get('/{name}'),
  getById: (id: number) => client.get(`/{name}/${id}`),
  create: (data: any) => client.post('/{name}', data),
  update: (id: number, data: any) => client.put(`/{name}/${id}`, data),
  delete: (id: number) => client.delete(`/{name}/${id}`),
};
```

## Multi-Platform Builds

### Desktop

```bash
cd web && npm run tauri:build
```

### iOS

```bash
cd web && npx tauri ios init    # First time only
npm run tauri:ios-dev            # Dev with simulator
npm run tauri:ios-build          # Production build
```

### Android

```bash
cd web && npx tauri android init  # First time only
npm run tauri:android-dev          # Dev build
npm run tauri:android-build        # APK build
```

## Database Configuration

Edit `server/local.yaml`:
- `mysql.host`, `mysql.port`, `mysql.user`, `mysql.password`, `mysql.dbname`
- `redis.host`, `redis.port`, `redis.password`
- `system.migrate: true` to auto-migrate on startup
- `system.redis: false` to disable Redis if not needed
