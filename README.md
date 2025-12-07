# 分布式任务跑批系统 (Distributed Task Scheduler)

一个功能完善的分布式任务调度系统，支持任务管理、定时调度、分布式执行等核心功能。

## 🚀 技术栈

### 后端
- **Go 1.21+** - 高性能编程语言
- **Gin** - 轻量级HTTP Web框架
- **GORM** - ORM框架
- **MySQL 8.0** - 关系型数据库
- **Redis 7.0** - 缓存和分布式锁
- **JWT** - 身份认证
- **Zap** - 高性能日志库

### 前端
- **Vue 3** - 渐进式JavaScript框架
- **TypeScript** - 类型安全
- **Vite** - 下一代前端构建工具
- **Element Plus** - Vue 3组件库
- **Pinia** - 状态管理
- **Vue Router** - 路由管理
- **ECharts** - 图表库
- **Axios** - HTTP客户端

## ✨ 核心特性

- 🕐 **时间轮调度** - 高效的定时任务触发算法，O(1)时间复杂度
- 🔄 **多种路由策略** - 轮询、随机、一致性哈希、最少使用、故障转移
- 📊 **DAG工作流** - 支持任务依赖，按拓扑顺序执行
- 🚀 **任务分片** - 大任务自动拆分，并行执行
- 🔒 **分布式锁** - Redis实现，防止任务重复调度
- 📝 **实时日志** - 任务执行日志实时查看
- ⚡ **Goroutine池** - 高效的并发任务执行
- 🔔 **告警通知** - 任务失败自动告警

## 📁 项目结构

```
.
├── server/                 # Go后端
│   ├── cmd/               # 应用入口
│   │   └── admin/         # 管理服务
│   ├── internal/          # 内部包
│   │   ├── config/        # 配置
│   │   ├── handler/       # HTTP处理器
│   │   ├── middleware/    # 中间件
│   │   ├── model/         # 数据模型
│   │   ├── repository/    # 数据仓库
│   │   ├── router/        # 路由
│   │   ├── service/       # 业务逻辑
│   │   ├── scheduler/     # 调度器
│   │   │   ├── timewheel/ # 时间轮
│   │   │   ├── router/    # 路由策略
│   │   │   └── dag/       # DAG调度
│   │   └── executor/      # 执行器
│   │       └── pool/      # Goroutine池
│   └── pkg/               # 公共包
│       ├── logger/        # 日志
│       ├── mysql/         # MySQL
│       └── redis/         # Redis
├── web/                   # Vue3前端
│   ├── src/
│   │   ├── api/           # API接口
│   │   ├── components/    # 组件
│   │   ├── router/        # 路由
│   │   ├── store/         # 状态管理
│   │   ├── styles/        # 样式
│   │   ├── utils/         # 工具函数
│   │   └── views/         # 页面
│   └── package.json
└── sql/                   # SQL脚本
    └── init.sql           # 数据库初始化
```

## 🛠️ 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 7.0+

### 1. 初始化数据库

```bash
mysql -u root -p < sql/init.sql
```

### 2. 启动后端

```bash
cd server

# 安装依赖
go mod tidy

# 修改配置文件
vim config.yaml

# 运行
make run
# 或
go run cmd/admin/main.go
```

### 3. 启动前端

```bash
cd web

# 安装依赖
npm install

# 开发模式运行
npm run dev

# 构建生产版本
npm run build
```

### 4. 访问系统

- 前端地址: http://localhost:3000
- 后端地址: http://localhost:8080
- 默认账号: admin / admin123

## 📋 API接口

### 认证
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 退出登录

### 任务组
- `GET /api/v1/group` - 任务组列表
- `POST /api/v1/group` - 创建任务组
- `PUT /api/v1/group/:id` - 更新任务组
- `DELETE /api/v1/group/:id` - 删除任务组

### 任务
- `GET /api/v1/task` - 任务列表
- `POST /api/v1/task` - 创建任务
- `PUT /api/v1/task/:id` - 更新任务
- `DELETE /api/v1/task/:id` - 删除任务
- `POST /api/v1/task/:id/start` - 启动任务
- `POST /api/v1/task/:id/stop` - 停止任务
- `POST /api/v1/task/:id/trigger` - 手动触发

### 执行记录
- `GET /api/v1/instance` - 实例列表
- `GET /api/v1/instance/:id` - 实例详情
- `POST /api/v1/instance/:id/cancel` - 取消任务
- `POST /api/v1/instance/:id/retry` - 重试任务
- `GET /api/v1/instance/:id/logs` - 执行日志

### 执行器
- `POST /api/v1/executor/register` - 注册执行器
- `POST /api/v1/executor/heartbeat` - 心跳上报
- `GET /api/v1/executor` - 执行器列表

## 🎯 技术亮点

1. **时间轮算法** - 高效定时任务调度，O(1)复杂度
2. **一致性哈希** - 任务路由，最小化节点变化影响
3. **Goroutine Pool** - 控制并发，避免资源泄露
4. **分布式锁** - Redis实现，保证任务不重复执行
5. **DAG调度** - 拓扑排序实现任务依赖
6. **任务分片** - 大任务并行执行
7. **优雅停机** - Context + WaitGroup
8. **RBAC权限** - 基于角色的访问控制
9. **实时日志** - WebSocket推送
10. **限流中间件** - 令牌桶算法

## 📄 License

MIT License
