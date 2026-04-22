package i18n

func loadChineseMessages() map[string]string {
	return map[string]string{
		// Root command
		"root.short": "通用CLI工具，从OpenAPI规范构建并发送HTTP请求",
		"root.long": `apix 读取 OpenAPI 3.x 规范，让您能够从命令行调用任何端点。
支持交互模式、直接操作选择和灵活的参数输入。`,
		
		// Features
		"feature.load_specs": "从本地文件或URL加载OpenAPI规范",
		"feature.auto_detect": "自动检测当前目录的标准文件名（openapi.yaml, api.yaml等）",
		"feature.auto_generate": "为每个API操作自动生成CLI命令",
		"feature.support_params": "支持路径参数、查询参数、请求头和请求体",
		"feature.output_formats": "多种输出格式（美化打印、原始、JSON）",
		"feature.dry_run": "干跑模式预览curl命令",
		"feature.debug_mode": "调试模式检查请求详情",
		
		// Examples
		"example.load_spec": "自动检测规范并列出可用操作",
		"example.call_operation": "调用特定的API操作",
		"example.use_env_profile": "使用环境配置（dev/prod/staging）",
		"example.auth_request": "发送带认证的请求",
		"example.preview_curl": "预览curl命令而不发送",
		
		// Flags - Global
		"flag.spec": "OpenAPI规范文件的路径或URL",
		"flag.profile": "使用环境名称，自动查找对应的规范文件（如：dev -> openapi-dev.yaml）",
		"flag.base_url": "覆盖服务器URL（优先级高于文档中定义的servers）",
		"flag.server": "选择服务器索引（当有多个servers时），默认为0（第一个服务器）",
		"flag.lang": "设置语言偏好（en/zh）。未指定时自动检测",
		"flag.full": "显示完整响应，包括状态、头和响应体",
		"flag.raw": "以美化JSON格式输出响应体",
		"flag.json": "以JSON格式输出响应体（与--full组合输出完整响应的JSON）",
		"flag.debug": "显示调试信息，包括请求URL和方法",
		
		// Flags - Operation
		"flag.body": "请求体的JSON文件",
		"flag.header": "请求头，格式为 '键: 值'（可多次指定）",
		"flag.token": "Bearer令牌用于认证",
		"flag.key": "API密钥用于认证",
		"flag.dry_run": "打印curl命令而不发送请求",
		
		// Error messages
		"error.no_spec": "未提供必需的OpenAPI规范",
		"error.no_spec.hint": "您可以通过以下方式之一提供OpenAPI规范：",
		"error.no_spec.option1": "使用--spec标志指定文件路径或URL：",
		"error.no_spec.option2": "使用环境配置（-P dev/prod/staging）：",
		"error.no_spec.option3": "将规范文件放在当前目录，使用标准命名：",
		"error.no_spec.examples": "示例：",
		"error.no_spec.more_info": "运行 'apix --help' 获取更多信息。",
		
		"error.load_spec_failed": "无法从 %q 加载OpenAPI规范",
		"error.load_spec.check": "请检查：",
		"error.load_spec.check.path": "文件路径或URL是否正确",
		"error.load_spec.check.exists": "文件是否存在且可访问",
		"error.load_spec.check.valid": "文件是否包含有效的OpenAPI 3.x规范",
		
		"error.env_profile_not_found": "未找到环境 %q 的规范文件",
		"error.env_profile.hint": "提示：请确保存在 openapi-{env}.yaml、api-{env}.yaml 或 swagger-{env}.yaml 文件",
		
		// Help text
		"help.usage": "用法：",
		"help.examples": "示例：",
		"help.flags": "标志：",
		"help.global_flags": "全局标志：",
		"help.available_commands": "可用命令：",
		"help.more_info": "使用 \"%s [command] --help\" 获取更多关于命令的信息。",
		
		// Operation help
		"op.operation_id": "操作ID：",
		"op.method": "方法：",
		"op.path": "路径：",
		"op.servers": "可用服务器：",
		"op.default": "（默认）",
		"op.server_flag_hint": "使用 --server 标志选择服务器：",
		"op.parameters": "参数：",
		"op.common_flags": "通用标志：",
		"op.required": "（必需）",
		"op.example.call": "调用 %s",
		"op.example.with_body": "调用 %s（带请求体）",
		"op.example.preview": "预览 %s 请求",
	}
}
