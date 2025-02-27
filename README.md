# Ollama Proxy

一个基于 Gin 框架的 Ollama API 代理服务，提供认证和请求转发功能。

## 功能特点

- 基于 Gin 框架开发
- 支持 Token 认证
- 支持跨域请求
- 完整的错误处理
- 支持所有 Ollama API 接口

## 配置文件

配置文件为 `config.yaml`，结构如下：

```yaml
auth:
  generate_tokens:    # 生成相关接口的token列表
    - "your-generate-token-1"
    - "your-generate-token-2"
  model_tokens:       # 模型管理相关接口的token列表
    - "your-model-token-1"
    - "your-model-token-2"
service:
  base_url: "http://localhost:11434"  # Ollama 服务地址
```

## API 接口文档

所有接口都需要在请求头中携带 `Authorization` Token 进行认证。模型管理接口需要使用 `model_tokens` 中的token，生成相关接口需要使用 `generate_tokens` 中的token。

### 模型管理接口

#### 1. 获取模型列表

- 请求方法：GET
- 请求路径：/sys/tags
- 请求头：

  ```json
  Authorization: your-model-token
  Content-Type: application/json
  ```

- 响应示例：

  ```json
  {
    "models": [
      {
        "name": "llama2",
        "modified_at": "2024-01-20T12:00:00Z",
        "size": 4000000000
      }
    ]
  }
  ```

#### 2. 拉取模型

- 请求方法：POST
- 请求路径：/sys/pull
- 请求头：

  ```json
  Authorization: your-model-token
  Content-Type: application/json
  ```

- 请求体：

  ```json
  {
    "name": "llama2"
  }
  ```

- 响应示例：

  ```json
  {
    "status": "downloading manifest",
    "completed": 0,
    "total": 100
  }
  ```

#### 3. 删除模型

- 请求方法：DELETE
- 请求路径：/sys/delete
- 请求头：

  ```json
  Authorization: your-model-token
  Content-Type: application/json
  ```

- 请求体：

  ```json
  {
    "name": "llama2"
  }
  ```

- 响应示例：

  ```json
  {
    "status": "success"
  }
  ```

#### 4. 复制模型

- 请求方法：POST
- 请求路径：/sys/copy
- 请求头：

  ```json
  Authorization: your-model-token
  Content-Type: application/json
  ```

- 请求体：

  ```json
  {
    "source": "llama2",
    "destination": "llama2-copy"
  }
  ```

- 响应示例：

  ```json
  {
    "status": "success"
  }
  ```

#### 5. 推送模型

- 请求方法：POST
- 请求路径：/sys/push
- 请求头：

  ```json
  Authorization: your-model-token
  Content-Type: application/json
  ```

- 请求体：

  ```json
  {
    "name": "llama2",
    "insecure": false
  }
  ```

- 响应示例：

  ```json
  {
    "status": "pushing manifest",
    "completed": 0,
    "total": 100
  }
  ```

#### 6. 查看模型详情

- 请求方法：POST
- 请求路径：/sys/show
- 请求头：

  ```json
  Authorization: your-model-token
  Content-Type: application/json
  ```

- 请求体：

  ```json
  {
    "name": "llama2"
  }
  ```

- 响应示例：

  ```json
  {
    "license": "MIT",
    "modelfile": "FROM llama2\nPARAMETER temperature 0.7",
    "parameters": "2.0B",
    "template": "{{ .Prompt }}",
    "size": 4000000000
  }
  ```

### 生成接口

### OpenAI 风格接口

所有接口都需要在请求头中携带 `Authorization` Token 进行认证，使用 `generate_tokens` 中的token。

#### 1. 聊天接口

- 请求方法：POST
- 请求路径：/v1/chat/completions
- 请求头：

  ```json
  Authorization: your-generate-token
  Content-Type: application/json
  ```

- 请求体：

  ```json
  {
    "model": "llama2",
    "messages": [
      {
        "role": "user",
        "content": "你好"
      }
    ],
    "stream": true,
    "temperature": 0.7,
    "top_p": 0.9
  }
  ```

- 响应示例：

  ```json
  {
    "id": "chatcmpl-123",
    "object": "chat.completion",
    "created": 1677652288,
    "model": "llama2",
    "choices": [{
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "你好！很高兴见到你。"
      },
      "finish_reason": "stop"
    }]
  }
  ```

#### 2. 生成接口

- 请求方法：POST
- 请求路径：/v1/completions
- 请求头：

  ```json
  Authorization: your-generate-token
  Content-Type: application/json
  ```

- 请求体：

  ```json
  {
    "model": "llama2",
    "prompt": "讲个故事",
    "stream": true,
    "temperature": 0.7,
    "top_p": 0.9
  }
  ```

- 响应示例：

  ```json
  {
    "id": "cmpl-123",
    "object": "text_completion",
    "created": 1677652288,
    "model": "llama2",
    "choices": [{
      "text": "从前有座山...",
      "index": 0,
      "finish_reason": "stop"
    }]
  }
  ```

#### 3. Embeddings接口

- 请求方法：POST
- 请求路径：/v1/embeddings
- 请求头：

  ```json
  Authorization: your-generate-token
  Content-Type: application/json
  ```

- 请求体：

  ```json
  {
    "model": "llama2",
    "input": "Hello World"
  }
  ```

- 响应示例：

  ```json
  {
    "object": "list",
    "data": [
      {
        "object": "embedding",
        "embedding": [0.1, 0.2, 0.3, ...],
        "index": 0
      }
    ],
    "model": "llama2"
  }
  ```

- 响应示例：

  ```json
  {
    "embedding": [0.1, 0.2, 0.3, ...],
    "model": "llama2"
  }
  ```

## 错误处理

### 错误响应格式

所有的错误响应都遵循以下统一格式：

```json
{
  "code": 500,
  "message": "错误信息描述",
  "data": null
}
```

### 错误码列表

| 错误码 | 说明 | 处理建议 |
|--------|------|----------|
| 401 | 认证失败 | 检查 Token 是否有效或是否已提供 |
| 500 | 服务器内部错误 | 检查服务器日志，确认 Ollama 服务是否正常运行 |

### 常见错误处理

1. **认证失败**
   - 确保请求头中包含 `Authorization` 字段
   - 验证 Token 是否在 config.yaml 中正确配置
   - 区分使用场景，模型管理接口使用 model_tokens，生成接口使用 generate_tokens

2. **模型相关错误**
   - 拉取模型失败时，检查网络连接和磁盘空间
   - 模型不存在时，确认模型名称是否正确，必要时重新拉取
   - 模型生成超时，考虑调整请求参数或检查服务器负载

3. **服务异常处理**
   - 定期检查服务日志
   - 监控服务资源使用情况
   - 配置适当的重试策略
   - 在客户端实现错误重试机制

### 最佳实践

1. **错误重试**
   - 对于网络类错误（502、504），建议实现指数退避重试
   - 避免对客户端错误（400、401、403）进行重试
   - 设置合理的重试次数和超时时间

2. **日志记录**
   - 记录详细的错误信息和请求上下文
   - 包含时间戳和请求ID
   - 定期检查和分析错误日志

3. **监控告警**
   - 监控服务可用性和响应时间
   - 设置错误率阈值告警
   - 关注异常流量和资源使用

## 项目结构

```txt
.
├── main.go         # 主程序入口
├── config.yaml     # 配置文件
└── README.md       # 项目文档

```

## Docker 部署

### 构建镜像

```bash
docker build -t ollama-proxy .
```

### 运行容器

```bash
docker run -d \
  -p 8080:8080 \
  -v /path/to/your/config.yaml:/app/config.yaml \
  --name ollama-proxy \
  ollama-proxy
```

### 环境变量说明

容器内已配置以下环境变量：

- `GO111MODULE=on`: 启用Go模块
- `CGO_ENABLED=0`: 禁用CGO
- `GOOS=linux`: 目标操作系统
- `GOARCH=amd64`: 目标架构

### 注意事项

1. 运行容器时需要将本地的`config.yaml`文件挂载到容器的`/app/config.yaml`
2. 确保`config.yaml`中的`base_url`配置正确指向Ollama服务地址
3. 容器默认暴露8080端口，可以根据需要修改端口映射
