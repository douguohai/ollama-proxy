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

#### 1. 聊天接口

- 请求方法：POST
- 请求路径：/api/chat
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
    "options": {
      "temperature": 0.7,
      "top_p": 0.9
    }
  }
  ```

- 响应示例：

  ```json
  {
    "model": "llama2",
    "created_at": "2024-01-20T12:00:00Z",
    "message": {
      "role": "assistant",
      "content": "你好！很高兴见到你。"
    },
    "done": false
  }
  ```

#### 2. 生成接口

- 请求方法：POST
- 请求路径：/api/generate
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
    "options": {
      "temperature": 0.7,
      "top_p": 0.9
    }
  }
  ```

- 响应示例：

  ```json
  {
    "model": "llama2",
    "created_at": "2024-01-20T12:00:00Z",
    "response": "从前有座山...",
    "done": false
  }
  ```

#### 3. Embeddings接口

- 请求方法：POST
- 请求路径：/api/embeddings
- 请求头：

  ```json
  Authorization: your-generate-token
  Content-Type: application/json
  ```

- 请求体：

  ```json
  {
    "model": "llama2",
    "prompt": "Hello World",
    "options": {
      "temperature": 0.7
    }
  }
  ```

- 响应示例：

  ```json
  {
    "embedding": [0.1, 0.2, 0.3, ...],
    "model": "llama2"
  }
  ```

## 错误码说明

- 401: 认证失败（Token 无效或未提供）
- 500: 服务器内部错误

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
