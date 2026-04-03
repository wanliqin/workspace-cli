---
name: 长亭产品cli工具
description: 长亭 Workspace CLI (cws) 源码与编译环境。包含完整源码、已编译二进制文件及开发依赖，支持直接运行和二次开发。
---

# 长亭产品 CLI 开发技能

本目录包含 `github.com/chaitin/workspace-cli` 的完整源码和已编译好的 `cws` 二进制文件，可直接运行，也可继续二次开发。

## 目录说明

```
/Users/liqinwan/.config/agents/skills/长亭产品cli工具/
├── cws                      # 已编译的 macOS ARM64 二进制文件 (Go 1.25.0)
├── main.go                  # 入口文件
├── go.mod / go.sum          # Go 模块依赖
├── products/                # 各产品命令实现
│   ├── safeline/            # 雷池 WAF 企业版
│   ├── safeline-ce/         # 雷池 WAF 社区版
│   ├── cloudwalker/         # 牧云
│   ├── tanswer/             # 全悉 (T-Answer)
│   ├── xray/                # 洞鉴
│   └── chaitin/             # 通用命令
├── cmd/                     # 代码生成工具
├── config/                  # 配置加载
└── SKILL.md                 # 本文件
```

## 直接使用

### 配置方式

创建 `config.yaml`：

```yaml
safeline:
  url: https://192.168.220.11
  api_key: your-api-key
```

或临时通过命令行参数：

```bash
./cws --insecure safeline --url https://192.168.220.11 --api-key <token> site list
```

### 常用命令

```bash
# 雷池 WAF
./cws safeline site list
./cws safeline site get <id>
./cws safeline site enable <id>
./cws safeline site disable <id>
./cws safeline policy-rule list
./cws safeline policy-rule create --comment "xxx" --target urlpath --cmp infix --value "/admin" --action deny
./cws safeline acl template list
./cws safeline acl template create --name "Block IPs" --template-type manual --target-type cidr --action forbid --targets "1.1.1.1"
./cws safeline network workgroup list
./cws safeline stats overview

# 雷池 CE
./cws safeline-ce site list
./cws safeline-ce log attack list
./cws safeline-ce stat overview

# 全悉
./cws tanswer firewall search-white-list
./cws tanswer rules search-block-rules

# X-Ray
./cws xray --help
```

## 二次开发

### 环境要求

- Go 1.25.0+（当前二进制使用 `~/sdk/go1.25.0` 编译）

### 重新编译

```bash
cd /Users/liqinwan/.config/agents/skills/workspace-cli-main

# 使用本地 Go 1.25.0 编译
~/sdk/go1.25.0/bin/go build -o cws .

# 或使用系统默认 Go（需 >= 1.25）
go build -o cws .
```

### 开发任务

```bash
# 格式化
go fmt ./...

# 测试
go test ./...

# 运行示例
go run . safeline site list --insecure
```

### 添加新功能

- **新增产品命令**：在 `products/<product>/` 下实现，然后在 `main.go` 的 `newApp()` 中注册
- **新增 safeline 子命令**：在 `products/safeline/cmd/<module>/` 下实现，在 `products/safeline/command.go` 的 `RegisterModules()` 中注册
- **新增 safeline-ce 动态命令**：更新 `products/safeline-ce/openapi.json`，重新编译即可自动生成

## 重要提示

1. `cws` 二进制已包含完整的 `safeline-ce` OpenAPI 定义（嵌入 `openapi.json`）
2. 企业版 `safeline` 的 `site` 命令会根据 WAF 部署模式自动切换 API（软件反代/硬件透明代理等）
3. `network` 命令仅在硬件模式下可用
4. 由于源码放在 skills 目录，修改后建议重新执行 `go build -o cws .` 更新二进制
