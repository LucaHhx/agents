#!/bin/bash
# 将 HZ agents 框架链接到目标项目的 .claude/ 目录
# 用法: cd <目标项目> && bash <agents仓库路径>/link.sh

set -e

AGENTS_DIR="$(cd "$(dirname "$0")" && pwd)"
TARGET_DIR="$(pwd)"

if [ "$AGENTS_DIR" = "$TARGET_DIR" ]; then
  echo "❌ 请在目标项目目录下运行此脚本，而不是在 agents 仓库里"
  echo "用法: cd /path/to/your-project && bash $0"
  exit 1
fi

mkdir -p "$TARGET_DIR/.claude"

# 创建符号链接
for dir in skills agents commands; do
  if [ -L "$TARGET_DIR/.claude/$dir" ]; then
    echo "⚠️  $TARGET_DIR/.claude/$dir 已存在，跳过"
  elif [ -d "$TARGET_DIR/.claude/$dir" ]; then
    echo "⚠️  $TARGET_DIR/.claude/$dir 是普通目录，跳过（请手动处理）"
  else
    ln -s "$AGENTS_DIR/$dir" "$TARGET_DIR/.claude/$dir"
    echo "✅ 链接: .claude/$dir → $AGENTS_DIR/$dir"
  fi
done

# 复制 CLAUDE.md 模板（如果目标项目没有）
if [ ! -f "$TARGET_DIR/.claude/CLAUDE.md" ] && [ -f "$AGENTS_DIR/CLAUDE.md" ]; then
  cp "$AGENTS_DIR/CLAUDE.md" "$TARGET_DIR/.claude/CLAUDE.md"
  echo "✅ 复制: CLAUDE.md → .claude/CLAUDE.md（请根据项目修改）"
fi

# 复制 settings 模板（如果目标项目没有）
if [ ! -f "$TARGET_DIR/.claude/settings.local.json" ] && [ -f "$AGENTS_DIR/settings.template.json" ]; then
  cp "$AGENTS_DIR/settings.template.json" "$TARGET_DIR/.claude/settings.local.json"
  echo "✅ 复制: settings.local.json（预配置常用权限）"
fi

echo ""
echo "🎉 链接完成！在 $TARGET_DIR 中可以使用 /hz-pm, /hz-backend 等命令。"
