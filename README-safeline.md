# 雷池 WAF CLI 使用手册

`cws safeline` 是长亭雷池 WAF（企业版 Skyview）的命令行管理工具，支持站点管理、策略规则、ACL 封禁、IP 组、日志查询、网络配置、系统管理等全生命周期操作。

---

## 目录

1. [快速开始](#快速开始)
2. [配置方式](#配置方式)
3. [命令总览](#命令总览)
4. [详细命令说明](#详细命令说明)
   - [stats - 统计概览](#stats---统计概览)
   - [site - 站点管理](#site---站点管理)
   - [ip-group - IP 组管理](#ip-group---ip-组管理)
   - [acl - ACL 管理](#acl---acl-管理)
   - [policy-group - 策略组](#policy-group---策略组)
   - [policy-rule - 自定义规则](#policy-rule---自定义规则)
   - [log - 日志查询](#log---日志查询)
   - [network - 网络配置（硬件模式）](#network---网络配置硬件模式)
   - [system - 系统管理](#system---系统管理)
5. [部署模式说明](#部署模式说明)
6. [常用操作示例](#常用操作示例)
7. [注意事项](#注意事项)

---

## 快速开始

```bash
# 查看帮助
./cws safeline --help

# 列出所有站点
./cws safeline site list

# 查看统计概览
./cws safeline stats overview --duration h

# 封禁一个 IP
./cws safeline acl template create --name "Block-IP" --template-type manual --target-type cidr --action forbid --targets "10.0.0.1"
```

---

## 配置方式

### 方式 1：config.yaml（推荐）

在 `cws` 同级目录创建 `config.yaml`：

```yaml
safeline:
  url: https://192.168.220.11
  api_key: YOUR_API_KEY
```

### 方式 2：环境变量

```bash
export SAFELINE_URL=https://192.168.220.11
export SAFELINE_API_KEY=YOUR_API_KEY
```

### 方式 3：命令行参数

```bash
./cws --insecure safeline --url https://192.168.220.11 --api-key YOUR_API_KEY site list
```

### 配置优先级

**命令行 flags > 环境变量 > config.yaml**

### 常用全局参数

| 参数 | 说明 |
|------|------|
| `--url` | Skyview API 地址 |
| `--api-key` | API Token |
| `--insecure` | 跳过 TLS 证书验证（自签证书必需，默认 true） |
| `--indent` | 输出格式化 JSON |
| `--dry-run` | 干跑模式，仅预览不执行 |
| `-c, --config` | 指定其他配置文件 |

---

## 命令总览

| 命令 | 功能 |
|------|------|
| `stats` | 查看 WAF 统计概览 |
| `site` | 站点（网站）管理 |
| `ip-group` | IP 组管理 |
| `acl` | ACL 访问控制 / 限速 / 封禁 |
| `policy-group` | 策略组（检测引擎开关） |
| `policy-rule` | 自定义策略规则 |
| `log` | 攻击日志、访问日志、限速日志查询 |
| `network` | 网络配置（仅硬件模式） |
| `system` | 许可证、机器 ID、系统日志 |

---

## 详细命令说明

### stats - 统计概览

```bash
# 查看 24 小时统计
./cws safeline stats overview --duration h

# 查看 30 天统计
./cws safeline stats overview --duration d
```

---

### site - 站点管理

管理 SafeLine 站点（网站）。**自动适配当前 WAF 部署模式**。

```bash
# 列出所有站点
./cws safeline site list

# 查看站点详情
./cws safeline site get <id>

# 更新站点策略组
./cws safeline site update <id> --policy-group <group-id>
./cws safeline site update <id> --policy-group 0  # 取消关联

# 启用/禁用站点
./cws safeline site enable <id>
./cws safeline site disable <id>
```

**注意**：当前 CLI 不支持 `site create`，如需创建站点请直接调用 API 或 fork 后二次开发。

---

### ip-group - IP 组管理

别名：`ipgroup`

```bash
# 列出 IP 组
./cws safeline ip-group list
./cws safeline ip-group list --name "office"
./cws safeline ip-group list --count 50

# 查看详情
./cws safeline ip-group get <id>

# 创建 IP 组
./cws safeline ip-group create --name "Office" --ips "192.168.1.0/24,10.0.0.1"
./cws safeline ip-group create --name "DC" --ips "172.16.0.0/16" --comment "数据中心"

# 删除 IP 组（支持批量）
./cws safeline ip-group delete <id>
./cws safeline ip-group delete 1 2 3

# 添加/移除 IP
./cws safeline ip-group add-ip <id> --ips "192.168.3.0/24,10.0.1.0/24"
./cws safeline ip-group remove-ip <id> --ips "192.168.3.0/24"
```

---

### acl - ACL 管理

分为 **template（模板）** 和 **rule（规则/被封禁用户）** 两层管理。

#### acl template - ACL 模板

```bash
# 列出模板
./cws safeline acl template list
./cws safeline acl template list --name "limit"

# 查看模板详情
./cws safeline acl template get <id>

# 手动规则：封禁指定 IP
./cws safeline acl template create --name "Block IPs" \
  --template-type manual --target-type cidr --action forbid \
  --targets "192.168.1.100,10.0.0.50"

# 手动规则：使用 IP 组封禁
./cws safeline acl template create --name "Block Group" \
  --template-type manual --target-type cidr --action forbid \
  --ip-groups 1,2

# 自动规则：限速
./cws safeline acl template create --name "Rate Limit" \
  --template-type auto --period 60 --limit 100 --action forbid

# 自动规则：限速（限流动作）
./cws safeline acl template create --name "Throttle" \
  --template-type auto --period 60 --limit 100 --action limit_rate \
  --limit-rate-limit 10 --limit-rate-period 60

# 启用/禁用/删除模板
./cws safeline acl template enable <id>
./cws safeline acl template disable <id>
./cws safeline acl template delete <id>
```

#### acl rule - ACL 规则（被封禁用户）

```bash
# 列出被封禁的规则（需指定模板 ID）
./cws safeline acl rule list --template-id <id>

# 删除规则（解封）
./cws safeline acl rule delete <id>
./cws safeline acl rule delete <id> --add-to-whitelist  # 同时加入白名单

# 清空所有规则（批量解封）
./cws safeline acl rule clear --template-id <id>
./cws safeline acl rule clear --template-id <id> --add-to-whitelist
```

---

### policy-group - 策略组

管理策略组（检测引擎配置），控制各类攻击检测模块的开关。

```bash
# 列出策略组
./cws safeline policy-group list

# 查看策略组详情
./cws safeline policy-group get <id>

# 更新模块状态：启用/禁用检测模块
./cws safeline policy-group update <id> --module m_sqli,m_xss --state enabled
./cws safeline policy-group update <id> --module m_cmd_injection,m_ssrf --state disabled
```

**可用模块键名**：

| 模块键 | 说明 |
|--------|------|
| `m_sqli` | SQL 注入检测 |
| `m_xss` | XSS 检测 |
| `m_cmd_injection` | 命令注入检测 |
| `m_file_include` | 文件包含检测 |
| `m_file_upload` | 文件上传检测 |
| `m_php_code_injection` | PHP 代码注入检测 |
| `m_php_unserialize` | PHP 反序列化检测 |
| `m_java` | Java 检测 |
| `m_java_unserialize` | Java 反序列化检测 |
| `m_ssrf` | SSRF 检测 |
| `m_ssti` | SSTI 检测 |
| `m_csrf` | CSRF 检测 |
| `m_scanner` | 扫描器检测 |
| `m_response` | 响应检测 |
| `m_rule` | 内置规则 |

---

### policy-rule - 自定义规则

管理自定义策略规则，支持简单模式和复杂 JSON 模式。

```bash
# 列出规则
./cws safeline policy-rule list              # 全局规则（默认）
./cws safeline policy-rule list --global=false  # 自定义规则

# 查看规则详情
./cws safeline policy-rule get <id>

# 简单模式：拦截包含 /admin 的 URL
./cws safeline policy-rule create --comment "Block admin" \
  --target urlpath --cmp infix --value "/admin" \
  --action deny

# JSON 模式：复杂条件
./cws safeline policy-rule create --comment "Complex rule" \
  --pattern-json '{"$AND":[{"infix":{"urlpath":"admin"}}]}' \
  --action deny

# 查看可用目标字段
./cws safeline policy-rule targets

# 查看特定字段支持的操作符
./cws safeline policy-rule targets --cmp urlpath

# 启用/禁用/删除规则
./cws safeline policy-rule enable <id>
./cws safeline policy-rule disable <id>
./cws safeline policy-rule delete <id>
```

**常用参数说明**：

| 参数 | 说明 |
|------|------|
| `--comment` | 规则描述（必填） |
| `--action` | `deny`（拦截）、`dry_run`（仅记录）、`allow`（放行） |
| `--risk-level` | 风险等级：0=无, 1=低, 2=中, 3=高 |
| `--enabled` | 是否启用（默认 true） |
| `--expire-time` | 过期时间戳（0=永不过期） |
| `--target` | 目标字段（简单模式） |
| `--cmp` | 比较操作符（简单模式） |
| `--value` | 匹配值（简单模式） |
| `--pattern-json` | 完整 JSON 条件（复杂模式） |

**常用 target 字段**：`urlpath`, `decoded_path`, `decoded_query`, `method`, `host`, `user_agent`, `request_body`, `response_status`

---

### log - 日志查询

#### log detect - 检测日志（攻击日志）

```bash
./cws safeline log detect list --count 50
./cws safeline log detect list --current-page 1 --target-page 2 --tail-sort '[1743628800, "abc123"]'
./cws safeline log detect get --event-id "xxx" --timestamp "1774857841"
```

#### log access - 访问日志

```bash
./cws safeline log access list --count 50
./cws safeline log access get --event-id "xxx" --req-start-time "1775117700"
```

#### log rate-limit - 限速日志（别名：rl）

```bash
./cws safeline log rate-limit list --count 50 --offset 20
./cws safeline log rl list --count 50
```

---

### network - 网络配置（硬件模式）

> ⚠️ **注意：仅硬件模式支持。** 软件模式下执行会报错：`该功能在软件模式下不支持`

```bash
# 工作组（别名：wg）
./cws safeline network workgroup list
./cws safeline network workgroup get <name>

# 网络接口（别名：if）
./cws safeline network interface list
./cws safeline network interface ip <name>

# 默认网关（别名：gw）
./cws safeline network gateway get

# 静态路由（别名：sr）
./cws safeline network route list
```

---

### system - 系统管理

```bash
# 查看许可证
./cws safeline system license

# 获取机器 ID
./cws safeline system machine-id

# 系统日志
./cws safeline system log list
./cws safeline system log list --count 50 --offset 0
```

---

## 部署模式说明

`cws safeline site` 命令会**自动探测**当前 WAF 的部署模式，并调用对应的 API：

| 部署模式 | site 命令自动调用的 API |
|---------|------------------------|
| Software Reverse Proxy | `/api/SoftwareReverseProxyWebsiteAPI` |
| Software Cluster Reverse Proxy | `/api/SoftwareClusterReverseProxyWebsiteAPI` |
| Software Port Mirroring | `/api/SoftwarePortMirroringWebsiteAPI` |
| Hardware Reverse Proxy | `/api/HardwareReverseProxyWebsiteAPI` |
| Hardware Transparent Proxy | `/api/HardwareTransparentProxyWebsiteAPI` |
| Hardware Transparent Bridging | `/api/HardwareTransparentBridgingWebsiteAPI` |
| Hardware Port Mirroring | `/api/HardwarePortMirroringWebsiteAPI` |
| Hardware Traffic Detection | `/api/HardwareTrafficDetectionWebsiteAPI` |
| Hardware Router Proxy | `/api/HardwareReverseProxyWebsiteAPI` |

探测接口：`GET /api/ServerControlledConfigAPI?type=operation_mode`

---

## 常用操作示例

### 场景 1：快速封禁恶意 IP

```bash
./cws safeline acl template create --name "紧急封禁" \
  --template-type manual --target-type cidr --action forbid \
  --targets "192.168.1.100,10.0.0.50"
```

### 场景 2：批量解封 IP

```bash
# 1. 找到模板 ID
./cws safeline acl template list

# 2. 清空该模板下所有封禁规则
./cws safeline acl rule clear --template-id 1 --add-to-whitelist
```

### 场景 3：开启 SQL 注入和 XSS 检测

```bash
./cws safeline policy-group update 1 --module m_sqli,m_xss --state enabled
```

### 场景 4：拦截特定 URL 路径

```bash
./cws safeline policy-rule create --comment "Block /admin path" \
  --target urlpath --cmp infix --value "/admin" \
  --action deny --risk-level 3
```

### 场景 5：查看最近攻击日志

```bash
./cws safeline log detect list --count 20
```

---

## 注意事项

1. **证书问题**：由于 WAF 通常使用自签证书，建议始终使用 `--insecure` 参数（该参数默认值已为 `true`）。
2. **模式限制**：`network` 相关命令仅在硬件模式下可用，软件模式会报错。
3. **站点创建**：当前 CLI 暂未实现 `site create` 子命令，如需批量创建站点建议直接调用 API 或二次开发。
4. **dry-run**：修改类命令（`create`, `update`, `delete`, `enable`, `disable`）支持 `--dry-run`，建议先预览再执行。
5. **输出格式**：默认输出表格格式，加 `--indent` 可输出格式化 JSON。
