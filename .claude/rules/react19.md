# React 19 强制规范

IMPORTANT: 以下 React 19 规范必须遵守，违反会导致审查不通过：

- 使用 `use()` 消费 Context — 禁止使用 `useContext()`
- `ref` 作为普通 prop 传递 — 禁止使用 `forwardRef`
- 使用 Actions 处理表单 — `useActionState`、`useFormStatus`
- 异步操作设置 Suspense 边界
- 复合组件使用点号 API（如 `<Tabs.Root>`、`<Tabs.List>`）
- 状态管理: Zustand store 遵循 `{ state, actions, meta }` 模式
- 样式: Tailwind 工具类优先，用 `cn()` 合并条件类名
