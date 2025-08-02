smart-dialog-ai/
├── cmd/                     # 主程序入口
│   └── main.go              # 主程序入口文件
├── internal/                # 项目核心代码
│   ├── api/                 # Gin 框架相关的 API 路由和控制器
│   │   ├── handlers.go      # API 请求处理逻辑
│   │   └── routes.go        # API 路由定义
│   ├── websocket/           # WebSocket 通信逻辑
│   │   ├── client.go        # WebSocket 客户端逻辑（与 Rust 通信）
│   │   ├── server.go        # WebSocket 服务器逻辑（与前端 Web 通信）
│   │   └── message.go       # WebSocket 消息处理逻辑
│   ├── service/             # 业务逻辑层
│   │   ├── asr_service.go   # ASR 服务调用逻辑
│   │   ├── llm_service.go   # LLM 服务调用逻辑
│   │   ├── rust_service.go  # Rust 服务器调用逻辑
│   │   └── tts_service.go   # TTS 服务调用逻辑
│   ├── model/               # 数据模型定义（Gorm）
│   │   ├── user.go          # 用户模型
│   │   └── conversation.go  # 对话记录模型
│   ├── repository/          # 数据访问层
│   │   ├── user_repository.go  # 用户数据访问
│   │   └── conversation_repository.go  # 对话记录数据访问
│   └── utils/               # 工具函数和通用逻辑
│       ├── http_client.go   # HTTP 客户端工具
│       └── log.go           # 日志工具
├── pkg/                     # 第三方库或工具包
│   └── rust_client/         # Rust 服务器调用的封装包
├── config/                  # 配置文件
│   └── config.yaml          # 项目配置文件（如数据库、服务地址等）
├── data/                    # 数据文件（如测试数据、配置模板等）
├── scripts/                 # 脚本文件（如数据库初始化脚本等）
├── tests/                   # 测试代码
│   ├── api_test.go          # API 测试
│   ├── service_test.go      # 业务逻辑测试
│   ├── model_test.go        # 数据模型测试
│   └── websocket_test.go    # WebSocket 测试
├── .env                     # 环境变量配置文件
├── go.mod                   # Go 模块依赖文件
├── go.sum                   # Go 模块依赖版本文件
└── README.md                # 项目说明文档