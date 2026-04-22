# GitHub Actions 使用指南

## 📋 概述

本项目配置了两个 GitHub Actions 工作流：

1. **CI (ci.yml)** - 持续集成，在每次提交和 PR 时运行
2. **Release (release.yml)** - 自动发布，在打标签时触发

## 🚀 快速开始

### 方式 1：通过 Git Tag 触发 Release（推荐）

```bash
# 1. 确保代码已提交
git add .
git commit -m "Prepare for release v1.0.0"

# 2. 创建标签
git tag v1.0.0

# 3. 推送标签到 GitHub
git push origin v1.0.0
```

推送标签后，GitHub Actions 会自动：
- ✅ 构建 5 个平台的二进制文件
- ✅ 生成 SHA256 校验和
- ✅ 自动生成变更日志
- ✅ 创建 GitHub Release 并上传所有文件

### 方式 2：通过 GitHub UI 手动触发

1. 进入仓库的 **Actions** 标签页
2. 选择 **Release** 工作流
3. 点击 **Run workflow**
4. 输入版本号（如：`v1.0.0`）
5. 点击 **Run workflow** 按钮

## 📦 生成的文件

每次 Release 会生成以下文件：

```
apix-linux-amd64              # Linux 64位
apix-linux-amd64.sha256       # 校验和
apix-linux-arm64              # Linux ARM64
apix-linux-arm64.sha256
apix-darwin-amd64             # macOS Intel
apix-darwin-amd64.sha256
apix-darwin-arm64             # macOS Apple Silicon
apix-darwin-arm64.sha256
apix-windows-amd64.exe        # Windows 64位
apix-windows-amd64.exe.sha256
checksums.txt                 # 所有文件的综合校验和
```

## 🔍 验证下载的文件

```bash
# 下载二进制文件和对应的 sha256 文件
wget https://github.com/gongdaowen/apix/releases/download/v1.0.0/apix-linux-amd64
wget https://github.com/gongdaowen/apix/releases/download/v1.0.0/apix-linux-amd64.sha256

# 验证校验和
sha256sum -c apix-linux-amd64.sha256

# 或使用综合校验和文件
sha256sum -c checksums.txt --ignore-missing
```

## 🎯 CI 工作流说明

CI 工作流在以下情况自动运行：
- 推送到 `main` 或 `master` 分支
- 创建 Pull Request

### CI 包含的检查：

1. **测试 (Test)**
   - 运行所有单元测试
   - 生成代码覆盖率报告
   - 上传到 Codecov

2. **构建验证 (Build Validation)**
   - 在 Linux、Windows、macOS 上编译
   - 验证二进制文件可以正常运行
   - 测试 `--help` 命令

3. **代码检查 (Lint)**
   - 运行 golangci-lint
   - 检查代码质量和风格

4. **规范验证 (Validate OpenAPI Specs)**
   - 验证 YAML 文件格式
   - 确保 OpenAPI 规范语法正确

## 🏷️ 版本命名规范

遵循语义化版本规范（Semantic Versioning）：

- **主版本**：不兼容的 API 修改 `v1.0.0` → `v2.0.0`
- **次版本**：向后兼容的功能性新增 `v1.0.0` → `v1.1.0`
- **修订版本**：向后兼容的问题修正 `v1.0.0` → `v1.0.1`

预发布版本：
- `v1.0.0-alpha.1`
- `v1.0.0-beta.1`
- `v1.0.0-rc.1`

## 📝 发布流程最佳实践

### 1. 准备发布

```bash
# 更新 CHANGELOG.md（如果有）
# 更新 README.md 中的版本号
# 确保所有测试通过

git add .
git commit -m "Prepare for release v1.0.0"
```

### 2. 创建并发布标签

```bash
# 创建带注释的标签（推荐）
git tag -a v1.0.0 -m "Release version 1.0.0

Features:
- Add multi-platform support
- Improve error handling
- Add debug mode

Bug Fixes:
- Fix parameter parsing issue
"

# 推送标签
git push origin v1.0.0
```

### 3. 监控构建进度

1. 访问 GitHub Actions 页面
2. 查看 Release 工作流运行状态
3. 等待所有任务完成（通常 2-5 分钟）

### 4. 验证 Release

1. 访问 Releases 页面
2. 检查所有二进制文件是否已上传
3. 下载并测试一个平台
4. 验证校验和

## 🔧 故障排查

### 问题 1：Release 工作流未触发

**检查项：**
- 标签格式是否正确（必须是 `v*` 格式）
- 标签是否已推送到远程仓库
- GitHub Actions 是否启用

**解决方案：**
```bash
# 检查标签
git tag -l

# 重新推送标签
git push origin v1.0.0

# 或删除后重新创建
git tag -d v1.0.0
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### 问题 2：构建失败

**检查项：**
- Go 代码是否有编译错误
- 依赖是否正确
- CI 工作流是否通过

**解决方案：**
```bash
# 本地测试构建
go build -v

# 测试多平台构建
GOOS=linux GOARCH=amd64 go build
GOOS=darwin GOARCH=amd64 go build
GOOS=windows GOARCH=amd64 go build

# 查看 Actions 日志
# 访问 GitHub -> Actions -> 失败的工作流 -> 查看详细日志
```

### 问题 3：版本号不正确

**解决方案：**
```bash
# 删除错误的标签
git tag -d v1.0.0
git push --delete origin v1.0.0

# 创建正确的标签
git tag -a v1.0.1 -m "Release v1.0.1"
git push origin v1.0.1
```

## 📊 自定义工作流

### 修改支持的平台

编辑 `.github/workflows/release.yml` 中的 matrix：

```yaml
strategy:
  matrix:
    include:
      - goos: linux
        goarch: amd64
        artifact_name: apix-linux-amd64
      # 添加新平台
      - goos: linux
        goarch: arm
        artifact_name: apix-linux-arm
```

### 添加额外的检查

编辑 `.github/workflows/ci.yml` 添加新的 job：

```yaml
security-scan:
  name: Security Scan
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - name: Run security scan
      run: gosec ./...
```

### 配置通知

在 workflow 中添加通知步骤：

```yaml
- name: Notify on success
  if: success()
  run: |
    curl -X POST -H 'Content-type: application/json' \
      --data '{"text":"Release ${{ steps.version.outputs.version }} successful!"}' \
      ${{ secrets.SLACK_WEBHOOK }}
```

## 🔐 安全建议

1. **不要提交敏感信息**
   - 使用 GitHub Secrets 存储敏感数据
   - 不要在 workflow 文件中硬编码 token

2. **限制权限**
   - 使用最小权限原则
   - 明确声明需要的 permissions

3. **定期更新 Actions**
   - 保持使用的 actions 版本最新
   - 关注安全公告

## 📚 相关资源

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [语义化版本规范](https://semver.org/)
- [Go 交叉编译指南](https://go.dev/doc/install/source#environment)
- [softprops/action-gh-release](https://github.com/softprops/action-gh-release)

---

**祝发布顺利！** 🎉
