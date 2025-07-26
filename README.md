# Imsystem
即时通讯系统 (IM-System)
这是一个基于 Go 语言实现的轻量级即时通讯系统，支持公聊、私聊和用户管理功能。项目目前集成了 TCP 通信，提供命令行客户端和基础 Web 界面，并计划扩展为更强大的前后端应用。
项目概述

语言: Go
依赖: 内置标准库
功能:
公聊模式
私聊模式
用户名更新
在线用户查询


运行环境: Windows, Linux, macOS (需安装 Go)

安装步骤
依赖安装

安装 Go

确保安装 Go 1.16 或更高版本，下载地址: Go 官网
配置环境变量 GOROOT 和 GOPATH，并将 GOPATH/bin 添加到 PATH。


初始化项目

进入项目目录 (例如 D:\GO_code):cd D:\GO_code
go mod init im-system





编译项目

编译服务器和客户端:go build -o server.exe server.go main.go user.go
go build -o client.exe client.go



使用方法

启动服务器

运行编译后的服务器:server.exe -ip 127.0.0.1 -port 8888


服务器将在 127.0.0.1:8888 上运行（当前仅支持 TCP，后续将支持 WebSocket）。


运行客户端

启动客户端:client.exe -ip 127.0.0.1 -port 8888


按照提示选择模式 (0-退出, 1-公聊, 2-私聊, 3-更新用户名)。



文件结构

client.go: 客户端逻辑
server.go: 服务器逻辑
user.go: 用户管理逻辑
main.go: 入口文件
go.mod: Go 模块文件

开发与贡献

拉取请求

Fork 本仓库
创建新分支: git checkout -b feature/xxx
提交更改: git commit -m "描述"
推送并提交 PR


问题反馈

如遇问题，请在 Issues 页面提交。


许可证

默认 MIT 许可证，详情见 LICENSE 文件 (如需添加)。



未来展望
为了提升项目功能和用户体验，我们计划在未来进行以下开发：

前后端框架集成:

前端: 引入 React 框架，使用 JSX 和 Tailwind CSS 构建更美观、响应式的 Web 界面，支持实时消息通知和用户列表显示。
后端: 优化 Go 服务器，集成 Gin 框架以处理 RESTful API，支持 WebSocket 通信，提升并发性能。
目标: 实现一个完整的 Web 应用，允许用户通过浏览器轻松登录、聊天，并支持群组功能。


桌面应用扩展:

利用 webview/v2 库开发跨平台桌面客户端，结合 Electron 或原生 Go GUI 库，提供离线消息存储。


功能增强:

添加文件传输功能。
实现用户认证和权限管理。
支持多语言界面。


性能优化:

引入数据库（如 SQLite 或 PostgreSQL）存储聊天记录。
优化服务器架构，支持分布式部署。



我们欢迎社区成员参与这些功能的开发，共同打造一个功能强大、易用的即时通讯平台！
注意事项

当前版本仅支持 TCP 通信，后续将迁移至 WebSocket。
确保网络环境正常，防火墙未阻止 8888 端口。
