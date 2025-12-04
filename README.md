# -
个人使用并发+分布式以及AI的功能完成处理海量数据的场景

```api
项目结构：
task-system/
├── cmd/
│   ├── api-server/      # API服务入口
│   └── worker/          # 工作节点入口
├── internal/
│   ├── api/             # HTTP接口层
│   ├── core/            # 核心业务逻辑
│   ├── models/          # 数据模型
│   ├── storage/         # 存储抽象
│   └── worker/          # 工作节点实现
└── pkg/
    └── utils/           # 工具函数
```
