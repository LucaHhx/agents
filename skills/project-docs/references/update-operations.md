# Update Operations Guide

## Updating Task Status (tasks.md)

Change the status column value:
- 待办 → 进行中: Set start date to today
- 进行中 → 已完成: Set completion date to today
- Any → 已取消: Add reason in notes column

When multiple tasks change, update all in a single edit.

## Recording Decisions (decisions.md)

Add a new decision record:
1. Copy the DR-NNN template block
2. Increment the number (DR-001, DR-002, ...)
3. Fill in all fields; leave "决定" and "理由" empty if status is 待定
4. When a decision is made, update status to 已决定 and fill remaining fields

## Adding Changelog Entries (changelog.md per plan)

Add under today's date header (create if not exists):
- Format: `- [类型] 变更内容 (原因)`
- Types: 新增 / 变更 / 修复 / 移除
- Group entries by date, newest at top

Also update the root `docs/CHANGELOG.md` if the change is project-significant.

## Recording Test Results (testing.md)

1. Add a new date section under "测试结果"
2. Fill in environment description
3. Add result rows for each test item
4. Save screenshots to the `testing/` directory with descriptive names
5. Move failed tests to "已知问题" table with severity and linked task number

## Cross-File Consistency

When updating any file, check if related files need updates:
- Task completed → Add changelog entry
- Decision made → Check if tasks need updating
- Test failed → Create task in tasks.md, add to known issues
