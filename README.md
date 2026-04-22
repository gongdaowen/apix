# Apix - 通用 OpenAPI CLI 工具

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25+-00ADD8.svg)](https://golang.org/)
[![CI](https://github.com/gongdaowen/apix/actions/workflows/ci.yml/badge.svg)](https://github.com/gongdaowen/apix/actions/workflows/ci.yml)
[![Release](https://github.com/gongdaowen/apix/actions/workflows/release.yml/badge.svg)](https://github.com/gongdaowen/apix/actions/workflows/release.yml)

**Apix** 是一个功能强大的命令行工具，能够从 OpenAPI 3.x 规范自动生成 CLI 命令，让您无需编写任何代码即可调用任何 REST API。

## ✨ 特性

- 🚀 **自动命令生成** - 从 OpenAPI 规范自动生成所有 API 操作的 CLI 命令
- 🔍 **智能文件检测** - 自动识别当前目录的标准文件名（openapi.yaml、api.yaml 等）
- 🌍 **多环境支持** - 通过环境配置文件轻松切换 dev/staging/prod 环境
- 📝 **完整参数支持** - 支持路径参数、查询参数、请求头、请求体
- 🔐 **认证集成** - 内置 Bearer Token 和 API Key 认证支持
- 🎯 **多种输出格式** - 美化打印、原始响应、JSON 格式输出
- 🔬 **调试模式** - Dry-run 模式预览 curl 命令，Debug 模式查看请求详情
- 🌐 **国际化支持** - 支持中文和英文界面，自动检测系统语言
- 🔄 **多服务器支持** - 当 API 有多个服务器时，可灵活选择目标服务器

## 📦 安装

### 从 Release 页面下载（推荐）

访问 [Releases](https://github.com/apix-cli/apix/releases) 页面，下载适合您平台的二进制文件：

- 🐧 **Linux**: `apix-linux-amd64` 或 `apix-linux-arm64`
- 🍎 **macOS**: `apix-darwin-amd64` (Intel) 或 `apix-darwin-arm64` (Apple Silicon)
- 🪟 **Windows**: `apix-windows-amd64.exe`

**验证下载的文件：**
```bash
# 下载校验和文件
wget https://github.com/apix-cli/apix/releases/latest/download/checksums.txt

# 验证文件完整性
sha256sum -c checksums.txt --ignore-missing
```

### 从源码构建

**前置要求：**
- Go 1.25+
- Git

**使用 Makefile（Linux/macOS）：**
```bash
git clone https://github.com/apix-cli/apix.git
cd apix
make build        # 构建当前平台
make build-all    # 构建所有平台
make install      # 安装到 GOPATH/bin
```

**使用构建脚本（Windows）：**
```powershell
git clone https://github.com/apix-cli/apix.git
cd apix
.\build.bat build    # 构建
.\build.bat dev      # 构建并运行
```

**手动构建：**
```bash
git clone https://github.com/apix-cli/apix.git
cd apix
go build -o apix main.go
```

### 使用 Go install

```bash
go install github.com/apix-cli/apix@latest
```

## 🚀 快速开始

### 1. 准备 OpenAPI 规范

**选项 A：使用示例文件（推荐新手）**

项目提供了完整的示例文件，位于 `examples/` 目录：

```bash
# 查看示例
ls examples/specs/      # OpenAPI 规范示例
ls examples/requests/   # 请求体示例

# 使用示例规范
apix -s examples/specs/jsonplaceholder.yaml listPosts

# 复制示例到当前目录以启用自动检测
cp examples/specs/jsonplaceholder.yaml .
apix  # 现在可以自动检测了
```

**选项 B：使用自己的规范文件**

将您的 OpenAPI 3.x 规范文件放在项目目录中，支持以下标准命名：
- `openapi.yaml` / `openapi.yml` / `openapi.json`
- `api.yaml` / `api.yml` / `api.json`
- `swagger.yaml` / `swagger.yml` / `swagger.json`

### 2. 列出可用操作

```bash
apix
```

Apix 会自动检测规范文件并显示所有可用的 API 操作。

### 3. 调用 API

```bash
# 获取用户信息
apix getUser --id 123

# 创建新帖子（使用示例请求体）
apix createPost -b examples/requests/post.json

# 列出所有帖子
apix listPosts
```

## 📖 使用场景

### 场景 1：本地开发测试

在本地开发环境中快速测试 API 端点：

```bash
# 自动检测 openapi.yaml 并调用
apix listUsers --limit 10 --offset 0

# 创建新用户
apix createUser -b new-user.json

# 预览请求而不发送（dry-run）
apix updateUser --id 42 -b update.json --dry-run
```

### 场景 2：多环境管理

为不同环境维护独立的规范文件：

```
项目目录/
├── openapi-dev.yaml      # 开发环境
├── openapi-staging.yaml  # 预发布环境
└── openapi-prod.yaml     # 生产环境
```

```bash
# 使用开发环境
apix -P dev listUsers

# 使用预发布环境
apix -P staging listUsers

# 使用生产环境
apix -P prod listUsers
```

### 场景 3：带认证的 API 调用

```bash
# 使用 Bearer Token
apix getProtectedResource --id 123 -t YOUR_ACCESS_TOKEN

# 使用 API Key
apix getApiKeyResource -k YOUR_API_KEY

# 同时使用 Token 和自定义请求头
apix updateResource --id 42 \
  -b data.json \
  -t YOUR_TOKEN \
  -H "X-Custom-Header: value" \
  -H "X-Request-ID: abc123"
```

### 场景 4：调试和排查问题

```bash
# 预览生成的 curl 命令
apix getUser --id 123 --dry-run

# 输出示例：
# curl -X GET "https://api.example.com/users/123" \
#   -H "Authorization: Bearer xxx" \
#   -H "Content-Type: application/json"

# 启用调试模式查看详细信息
apix getUser --id 123 --debug

# 以 JSON 格式输出完整响应（包括状态码和响应头）
apix getUser --id 123 --json

# 仅输出响应体（适合脚本处理）
apix getUser --id 123 --raw
```

### 场景 5：CI/CD 自动化

在持续集成流程中自动化 API 测试：

```bash
#!/bin/bash
# test-api.sh

# 健康检查
apix healthCheck --raw | jq '.status'

# 创建测试数据
apix createTestUser -b test-user.json --json | jq '.id' > user-id.txt

# 验证创建
USER_ID=$(cat user-id.txt)
apix getUser --id $USER_ID --json | jq '.name'

# 清理测试数据
apix deleteUser --id $USER_ID
```

### 场景 6：远程规范文件

直接从 URL 加载 OpenAPI 规范：

```bash
# 从 URL 加载规范
apix -s https://api.example.com/openapi.json listResources

# 结合认证使用
apix -s https://api.internal.com/swagger.yaml \
  getResource --id 123 \
  -t $API_TOKEN
```

### 场景 7：多服务器选择

当 API 规范定义了多个服务器时：

```yaml
servers:
  - url: https://api-us.example.com
    description: US Region
  - url: https://api-eu.example.com
    description: EU Region
  - url: https://api-asia.example.com
    description: Asia Region
```

```bash
# 查看所有可用服务器
apix

# 使用美国服务器（默认，索引 0）
apix getData --server 0

# 使用欧洲服务器（索引 1）
apix getData --server 1

# 使用亚洲服务器（索引 2）
apix getData --server 2

# 完全覆盖服务器 URL
apix getData --base-url https://custom-server.example.com
```

### 场景 8：批量操作

结合 shell 脚本进行批量数据处理：

```bash
#!/bin/bash
# batch-update.sh

# 从文件读取 ID 列表并批量更新
while read -r id; do
  echo "Updating user $id..."
  apix updateUser --id "$id" -b "user-$id.json" --raw
done < user-ids.txt

# 批量删除
for id in 1 2 3 4 5; do
  apix deleteUser --id $id
  echo "Deleted user $id"
done
```

### 场景 9：API 文档探索

快速了解 API 的功能和参数：

```bash
# 查看特定操作的详细帮助
apix createPost --help

# 输出示例：
# Create a new post
#
# Operation ID: createPost
# Method: POST
# Path: /posts
#
# Parameters:
#   (此操作无路径/查询参数)
#
# Common Flags:
#   -b, --body string    请求体的JSON文件
#   -H, --header strings 请求头，格式为 '键: 值'（可多次指定）
#   -t, --token string   Bearer令牌用于认证
#   ...
#
# Examples:
#   # Call Create a new post
#   apix createPost 
#
#   # Call Create a new post with request body
#   apix createPost -b request.json 
#
#   # Preview Create a new post request
#   apix createPost  --dry-run
```

### 场景 10：与 jq 结合使用

强大的 JSON 处理能力：

```bash
# 提取特定字段
apix listUsers --raw | jq '.[].name'

# 过滤数据
apix listPosts --raw | jq '[.[] | select(.userId == 1)]'

# 格式化输出
apix getUser --id 123 --raw | jq '.'

# 链式处理
apix listUsers --raw | jq -r '.[].id' | while read id; do
  apix getUser --id $id --raw | jq '{name, email}'
done
```

## 📋 命令行参数

### 全局标志

| 标志 | 简写 | 说明 |
|------|------|------|
| `--spec` | `-s` | OpenAPI 规范文件的路径或 URL |
| `--profile` | `-P` | 使用环境名称，自动查找对应的规范文件 |
| `--base-url` | - | 覆盖服务器 URL（优先级高于规范中定义的 servers） |
| `--server` | - | 选择服务器索引（当有多个 servers 时），默认为 0 |
| `--raw` | - | 仅输出响应体，不包含头信息或状态 |
| `--json` | - | 以 JSON 格式输出完整响应，包括状态、头和体 |
| `--debug` | - | 显示调试信息，包括请求 URL 和方法 |

### 操作级别标志

| 标志 | 简写 | 说明 |
|------|------|------|
| `--body` | `-b` | 请求体的 JSON 文件路径 |
| `--header` | `-H` | 请求头，格式为 'Key: Value'（可多次指定） |
| `--token` | `-t` | Bearer Token 用于认证 |
| `--key` | `-k` | API Key 用于认证 |
| `--dry-run` | - | 打印 curl 命令而不发送请求 |

### 动态参数

每个 API 操作会根据其 OpenAPI 定义自动生成对应的参数标志。例如：

```yaml
parameters:
  - name: userId
    in: query
    required: true
    schema:
      type: integer
```

会生成：
```bash
apix someOperation --userId 123
```

## 🗂️ 项目结构

```
apix/
├── cmd/                    # 命令行入口
│   └── root.go            # 主命令和动态子命令注册
├── internal/              # 内部包
│   ├── builder/           # HTTP 请求构建器
│   ├── executor/          # HTTP 请求执行器
│   ├── formatter/         # 响应格式化器
│   ├── i18n/              # 国际化支持
│   │   ├── en.go         # 英文翻译
│   │   ├── i18n.go       # 翻译引擎
│   │   └── zh.go         # 中文翻译
│   ├── models/            # 数据模型
│   ├── parser/            # OpenAPI 规范解析器
│   └── resolver/          # 参数解析器
├── pkg/models/            # 公共模型
├── main.go                # 程序入口
├── go.mod                 # Go 模块定义
└── README.md              # 项目文档
```

## 🔧 配置示例

### 示例 1：基础 OpenAPI 规范

```yaml
openapi: 3.0.0
info:
  title: User API
  version: 1.0.0
servers:
  - url: https://api.example.com
paths:
  /users/{id}:
    get:
      operationId: getUser
      summary: Get user by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Success
```

使用方式：
```bash
apix getUser --id 123
```

### 示例 2：带请求体的 POST 请求

```yaml
/posts:
  post:
    operationId: createPost
    summary: Create a new post
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            properties:
              title:
                type: string
              content:
                type: string
```

**使用示例请求体文件：**

您可以直接使用项目提供的示例文件：
```bash
apix createPost -b examples/requests/post.json
```

或创建自己的 `post.json`：
```json
{
  "title": "My First Post",
  "content": "Hello, World!"
}
```

使用方式：
```bash
apix createPost -b post.json
```

### 示例 3：多环境配置

```
项目目录/
├── openapi-dev.yaml      # 开发环境配置
├── openapi-staging.yaml  # 预发布环境配置
└── openapi-prod.yaml     # 生产环境配置
```

**openapi-dev.yaml**:
```yaml
servers:
  - url: http://localhost:8080
```

**openapi-prod.yaml**:
```yaml
servers:
  - url: https://api.production.com
```

使用方式：
```bash
# 开发环境
apix -P dev listUsers

# 生产环境
apix -P prod listUsers
```

## 💡 最佳实践

### 1. 组织规范文件

```
project/
├── specs/
│   ├── openapi.yaml          # 主规范
│   ├── openapi-dev.yaml      # 开发环境
│   ├── openapi-staging.yaml  # 预发布环境
│   └── openapi-prod.yaml     # 生产环境
├── requests/                  # 请求体模板
│   ├── create-user.json
│   ├── update-post.json
│   └── ...
└── scripts/                   # 自动化脚本
    ├── test-api.sh
    └── deploy.sh
```

### 2. 使用环境变量存储敏感信息

```bash
# 不要硬编码 Token
export API_TOKEN="your-secret-token"
apix getProtectedData -t $API_TOKEN

# 或在 .env 文件中
source .env
apix getProtectedData -t $API_TOKEN
```

### 3. 利用 Dry-run 进行验证

在实际发送请求前，先用 `--dry-run` 验证：

```bash
# 检查生成的请求是否正确
apix updateUser --id 42 -b data.json --dry-run

# 确认无误后再发送
apix updateUser --id 42 -b data.json
```

### 4. 脚本中使用 --raw 和 --json

```bash
# 在脚本中提取数据
USER_ID=$(apix createUser -b user.json --raw | jq '.id')
echo "Created user with ID: $USER_ID"

# 完整的响应信息用于调试
apix createUser -b user.json --json | jq '.status'
```

### 5. 组合使用请求头

```bash
# 添加多个自定义请求头
apix updateResource --id 123 \
  -b data.json \
  -H "X-Request-ID: $(uuidgen)" \
  -H "X-Correlation-ID: abc123" \
  -H "Accept-Language: en-US"
```

## 🌍 国际化

Apix 支持中文和英文界面，会自动检测系统语言。您也可以手动指定：

```bash
# 强制使用英文
LANG=en apix

# 强制使用中文
LANG=zh apix
```

## 🐛 故障排查

### 问题 1：找不到规范文件

**错误信息**：`OpenAPI specification is required but not provided`

**解决方案**：
1. 确保规范文件在当前目录，并使用标准命名
2. 或使用 `--spec` 明确指定路径：
   ```bash
   apix -s path/to/your/spec.yaml listUsers
   ```
3. 或使用环境配置：
   ```bash
   apix -P dev listUsers
   ```

### 问题 2：规范文件格式错误

**错误信息**：`Failed to load OpenAPI specification`

**解决方案**：
1. 验证 YAML/JSON 格式是否正确
2. 确保是 OpenAPI 3.x 版本
3. 使用在线验证器检查规范：https://editor.swagger.io/

### 问题 3：缺少必需参数

**错误信息**：`required parameter "xxx" not provided`

**解决方案**：
查看操作的帮助信息，提供所有必需参数：
```bash
apix yourOperation --help
apix yourOperation --requiredParam value
```

### 问题 4：认证失败

**解决方案**：
1. 检查 Token 或 API Key 是否正确
2. 确认是否需要特定的请求头
3. 使用 `--debug` 查看实际发送的请求：
   ```bash
   apix getResource --id 123 -t YOUR_TOKEN --debug
   ```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

### 开发工作流

1. **Fork 仓库**
2. **创建特性分支**: `git checkout -b feature/amazing-feature`
3. **提交更改**: `git commit -m 'Add amazing feature'`
4. **推送到分支**: `git push origin feature/amazing-feature`
5. **创建 Pull Request**

### 本地开发

```bash
# 克隆仓库
git clone https://github.com/your-username/apix.git
cd apix

# 安装依赖
go mod download

# 运行测试
make test        # Linux/macOS
build.bat test   # Windows

# 构建并测试
make dev         # Linux/macOS
build.bat dev    # Windows
```

### 代码规范

- 遵循 [Effective Go](https://go.dev/doc/effective_go) 指南
- 使用 `gofmt` 格式化代码
- 为新功能编写单元测试
- 确保 CI 检查通过

## 🚀 CI/CD

本项目使用 GitHub Actions 进行持续集成和自动发布。

### 持续集成 (CI)

每次提交和 PR 都会自动运行：
- ✅ 单元测试和覆盖率检查
- ✅ 多平台构建验证（Linux、Windows、macOS）
- ✅ 代码质量检查（golangci-lint）
- ✅ OpenAPI 规范验证

查看 [CI 工作流](.github/workflows/ci.yml) 了解详情。

### 自动发布

创建 Git tag 时自动触发：

```bash
# 创建版本标签
git tag v1.0.0

# 推送标签
git push origin v1.0.0
```

GitHub Actions 会自动：
- 🏗️ 构建 5 个平台的二进制文件
- 🔐 生成 SHA256 校验和
- 📝 自动生成变更日志
- 📦 创建 GitHub Release 并上传所有文件

支持的平台：
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

查看 [Release 工作流](.github/workflows/release.yml) 和 [发布指南](.github/RELEASE_GUIDE.md) 了解详情。

### 手动触发 Release

您也可以通过 GitHub UI 手动触发 release：
1. 进入 **Actions** → **Release**
2. 点击 **Run workflow**
3. 输入版本号（如：`v1.0.0`）
4. 点击运行

## 📄 许可证

MIT License

## 🔗 相关链接

- [OpenAPI 规范](https://swagger.io/specification/)
- [Cobra CLI 框架](https://github.com/spf13/cobra)
- [kin-openapi](https://github.com/getkin/kin-openapi)

---

**享受使用 Apix 构建和测试您的 API！** 🚀
