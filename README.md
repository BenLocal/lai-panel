# lai-panel

一个基于 Go 和 Vue 3 的应用管理面板，支持 Docker 容器管理、节点管理和应用部署。

## 功能特性

- 🚀 **应用管理** - 创建、编辑和管理应用程序
- 🖥️ **节点管理** - 管理本地和远程服务器节点
- 🐳 **Docker 管理** - 容器、镜像、网络和卷管理
- 📦 **服务部署** - 通过 Docker Compose 部署服务
- 💻 **终端访问** - 通过 Web 终端访问节点和容器
- 📁 **工作区管理** - 管理应用的工作区文件
- 🔐 **用户认证** - 基于 JWT 的用户认证系统

## 技术栈

### 后端
- Go 1.25+
- CloudWeGo Hertz (HTTP 框架)
- SQLite (数据库)
- Docker API
- SignalR (实时通信)

### 前端
- Vue 3
- PrimeVue 4
- TypeScript
- Vite
- TailwindCSS

## 快速开始

### 前置要求

- Go 1.25 或更高版本
- Node.js 18+ 和 npm
- Docker (可选，用于容器管理)

### 安装

1. 克隆仓库
```bash
git clone https://github.com/benlocal/lai-panel.git
cd lai-panel
```

2. 构建项目
```bash
make all
```

3. 构建前端
```bash
cd dashboard
npm install
npm run build
cd ..
```

4. 运行服务
```bash
make run-serve
```

服务默认运行在 `http://localhost:8080`

### 开发模式

运行前端开发服务器：
```bash
cd dashboard
npm run dev
```

前端开发服务器运行在 `http://localhost:5173`

## 项目结构

```
lai-panel/
├── cmd/
│   ├── serve/      # 主服务程序
│   └── agent/       # 代理程序
├── dashboard/       # 前端 Vue 应用
├── pkg/            # Go 包
│   ├── handler/    # HTTP 处理器
│   ├── repository/ # 数据访问层
│   ├── model/      # 数据模型
│   └── ...
└── migrations/     # 数据库迁移文件
```

## 使用说明

1. 访问 `http://localhost:8080` 打开管理面板
2. 使用默认凭据登录（首次运行需要创建用户）
3. 添加节点、创建应用并开始部署服务

## 许可证

查看 [LICENSE](LICENSE) 文件了解详情。
