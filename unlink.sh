#!/bin/bash
# 从目标项目移除 HZ agents 框架链接
# 用法: cd <目标项目> && bash <agents仓库路径>/unlink.sh

set -e
TARGET_DIR="$(pwd)"

for dir in skills agents commands; do
  if [ -L "$TARGET_DIR/.claude/$dir" ]; then
    rm "$TARGET_DIR/.claude/$dir"
    echo "✅ 移除链接: .claude/$dir"
  fi
done

echo "🎉 清理完成"
