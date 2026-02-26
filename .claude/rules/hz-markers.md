# HZ 标记扩展系统

添加新后端模块时，按顺序在以下入口文件的 HZ 标记处扩展：

| 标记前缀 | 文件 | 子标记 |
|----------|------|--------|
| HZ:ROUTER | router/inlet.go | PACKAGE_IMPORTS, PACKAGE_FIELDS, PACKAGE_INIT |
| HZ:API | api/inlet.go | PACKAGE_IMPORTS, PACKAGE_FIELDS, PACKAGE_INIT |
| HZ:SERVICE | service/inlet.go | PACKAGE_IMPORTS, PACKAGE_FIELDS |
| HZ:MIGRATE | entrance/migrate.go | PACKAGE_IMPORTS, MODEL_LIST |

IMPORTANT: 如果目标文件中 HZ 标记不存在，报告问题并向 PM 请求指导，不要自行猜测插入位置。

全局变量前缀: `HZ_`（HZ_CONFIG, HZ_DB, HZ_REDIS, HZ_LOG, HZ_JSON）
