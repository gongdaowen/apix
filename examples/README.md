# Examples

本目录包含 Apix 的使用示例文件。

## 📁 目录结构

```
examples/
├── specs/          # OpenAPI 规范示例
│   ├── jsonplaceholder.yaml    # JSONPlaceholder API 示例
│   └── test-api.yaml           # 测试 API 示例
└── requests/       # 请求体示例
    ├── post.json               # 创建帖子的请求体
    └── new-post.json           # 新帖子的请求体（备用）
```

## 🚀 快速使用

### 使用示例规范文件

```bash
# 使用 jsonplaceholder 示例
apix -s examples/specs/jsonplaceholder.yaml listPosts

# 使用 test-api 示例
apix -s examples/specs/test-api.yaml listUsers --limit 10
```

### 使用示例请求体

```bash
# 使用 post.json 创建帖子
apix -s examples/specs/jsonplaceholder.yaml createPost -b examples/requests/post.json

# 预览请求
apix -s examples/specs/jsonplaceholder.yaml createPost \
  -b examples/requests/post.json \
  --dry-run
```

## 📝 示例说明

### OpenAPI 规范示例

#### jsonplaceholder.yaml
- **用途**：JSONPlaceholder fake API 的 OpenAPI 规范
- **服务器**：https://jsonplaceholder.typicode.com
- **操作**：
  - `listPosts` - 列出所有帖子
  - `getPost` - 获取单个帖子
  - `createPost` - 创建新帖子

#### test-api.yaml
- **用途**：简单的测试 API 规范
- **服务器**：https://api.example.com
- **操作**：
  - `listUsers` - 列出用户
  - `getUser` - 获取用户
  - `deleteUser` - 删除用户
  - `createPost` - 创建帖子

### 请求体示例

#### post.json
```json
{
  "title": "Hello World",
  "content": "This is a test post from apix CLI"
}
```

#### new-post.json
```json
{
  "title": "My Test Post",
  "body": "This is the content of my test post created via apix CLI",
  "userId": 1
}
```

## 💡 提示

1. **复制示例文件到项目根目录**可以自动检测，无需 `-s` 参数
2. **修改示例文件**以适应您的 API 需求
3. **参考这些示例**创建您自己的规范和请求体文件

## 🔗 相关链接

- [完整文档](../README.md)
- [快速开始](../QUICKSTART.md)
- [JSONPlaceholder](https://jsonplaceholder.typicode.com/)
