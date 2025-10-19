# pitanti
基于高性能分布式框架 Pitaya 构建的生产级 Go 语言服务端脚手架，一站式整合应用层设计与开箱即用的工具套件，旨在大幅提升开发效率。A production-ready Go backend scaffold built on Pitaya, designed to supercharge your development workflow with out-of-the-box application-layer abstractions and curated toolkits。

## 项目结构

```
.
├── app                 # 应用核心逻辑
├── client              # 客户端示例
├── common              # 通用组件
├── conf                # 配置文件
├── protos              # Protocol Buffers 定义
├── services            # 服务实现
│   ├── game            # 游戏服务
│   └── gateway         # 网关服务
├── utils               # 工具函数
└── README.md
```

## 开发

```bash
# 本地运行
go run ./services/gateway/main.go -port=3100 -env=dev -rpcsvport=3200
go run ./services/game/main.go -serverId=1 -env=dev
go run ./services/game/main.go -serverId=2 -env=dev
```
