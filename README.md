# workspace-cli / 长亭产品命令行工具

[![CI](https://img.shields.io/github/actions/workflow/status/chaitin/workspace-cli/ci.yml?branch=main&label=CI)](https://github.com/chaitin/workspace-cli/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/chaitin/workspace-cli?label=Release)](https://github.com/chaitin/workspace-cli/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/chaitin/workspace-cli?label=Go)](https://github.com/chaitin/workspace-cli/blob/main/go.mod)
[![License](https://img.shields.io/github/license/chaitin/workspace-cli?label=License)](https://github.com/chaitin/workspace-cli/blob/main/LICENSE)

**Chaitin Workspace CLI for products**

长亭科技（Chaitin）统一产品命令行工具，一个二进制文件即可管理多个长亭安全产品：

- **雷池 WAF 企业版** (`safeline`)
- **雷池 WAF 社区版** (`safeline-ce`)
- **牧云** (`cloudwalker`)
- **全悉 T-Answer** (`tanswer`)
- **洞鉴 X-Ray** (`xray`)

---

## 快速开始

```bash
# 1. 查看帮助
./cws --help

# 2. 列出雷池 WAF 站点（需先配置 config.yaml）
./cws safeline site list

# 3. 查看统计概览
./cws safeline stats overview

# 4. 查看全悉阻断策略
./cws tanswer rules search-block-rules
```

---

## 配置说明 / Configuration

在 `./config.yaml` 中写入各产品的连接信息：

```yaml
safeline:
  url: https://192.168.220.11
  api_key: YOUR_API_KEY

safeline-ce:
  url: https://your-server:9443
  api_key: YOUR_API_KEY

cloudwalker:
  url: https://cloudwalker.example.com/rpc
  api_key: YOUR_API_KEY

tanswer:
  url: https://tanswer.example.com
  api_key: YOUR_API_KEY

xray:
  url: https://xray.example.com/api/v2
  api_key: YOUR_API_KEY
```

### 环境变量

你也可以通过环境变量或 `.env` 文件配置，变量命名规则为 `<PRODUCT>_<FIELD>`：

```text
cloudwalker.url      -> CLOUDWALKER_URL
cloudwalker.api_key  -> CLOUDWALKER_API_KEY
tanswer.url          -> TANSWER_URL
tanswer.api_key      -> TANSWER_API_KEY
xray.url             -> XRAY_URL
xray.api_key         -> XRAY_API_KEY
safeline-ce.url      -> SAFELINE_CE_URL
safeline-ce.api_key  -> SAFELINE_CE_API_KEY
safeline.url         -> SAFELINE_URL
safeline.api_key     -> SAFELINE_API_KEY
```

`.env` 示例：

```bash
SAFELINE_URL=https://safeline.example.com
SAFELINE_API_KEY=YOUR_API_KEY
XRAY_URL=https://xray.example.com/api/v2
XRAY_API_KEY=YOUR_API_KEY
```

### 优先级

配置优先级：**命令行 flags > 环境变量/.env > config.yaml**

### 多环境切换

使用 `-c` 或 `--config` 加载不同的配置文件，方便在多个 WAF 环境之间切换：

```bash
cws -c ./configs/safeline-prod.yaml safeline stats overview
cws -c ./configs/safeline-staging.yaml safeline stats overview
```

### 干跑模式 (Dry Run)

使用 `--dry-run` 可以预览命令会做什么，但不会真正发送请求：

```bash
cws --dry-run xray plan PostPlanFilter --filterPlan.limit=10
```

---

## 项目结构 / Project Structure

```text
main.go                # 程序入口和 CLI 路由
products/<name>/       # 每个产品独立的命令实现目录
Taskfile.yml           # 构建、运行、格式化任务定义
```

---

## 新增产品 / More Products

要添加新产品，请在 `products` 目录下实现，并完成以下步骤：

- 在 `main.go` 中导入产品包。
- 在 `newApp()` 中使用 `a.registerProductCommand(...)` 注册命令。
- 如果 `NewCommand()` 返回 `(*cobra.Command, error)`，需要在注册前处理错误。
- 如果产品需要读取 `config.yaml` 或根级别的运行时 flag，请在该产品包中实现 `ApplyRuntimeConfig(...)`，并在 `main.go` 的 `wrapProductCommand()` 中调用。
- 产品专属配置应在产品包内部通过 `config.Raw` 解码，**不要**把配置解析逻辑加到根命令中。

---

## 快捷调用方式 / BusyBox-Style Invocation

你可以通过创建软链接或重命名可执行文件的方式，直接以子命令名称调用：

```bash
task build
ln -s ./bin/cws ./chaitin
./chaitin
```

这等价于：

```bash
./bin/cws chaitin
```

---

## 常用任务 / Task

```bash
task build                          # 构建 cws 二进制文件
task run:chaitin                    # 运行 chaitin 演示命令
task fmt                            # 格式化 Go 代码
task lint                           # 运行 go vet 检查
task test                           # 运行单元测试
task package GOOS=linux GOARCH=amd64  # 打包指定平台的发行版
```

---

## Demo

### CloudWalker / 牧云

[![asciicast](https://asciinema.org/a/894643.svg)](https://asciinema.org/a/894643)

### T-Answer / 全悉

[![asciicast](https://asciinema.org/a/Pxabe3keAL0Z6PoJ.svg)](https://asciinema.org/a/Pxabe3keAL0Z6PoJ)

### SafeLine / 雷池企业版

[![asciicast](https://asciinema.org/a/ZDqZTHLD3nwXC27Z.svg)](https://asciinema.org/a/ZDqZTHLD3nwXC27Z)

### SafeLine-CE / 雷池社区版

[![asciicast](https://asciinema.org/a/dzJzibRTm8arWRmU.svg)](https://asciinema.org/a/dzJzibRTm8arWRmU)

### X-Ray / 洞鉴

[![asciicast](https://asciinema.org/a/XH6Hk9pWK0yp4VIt.svg)](https://asciinema.org/a/XH6Hk9pWK0yp4VIt)
