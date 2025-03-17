# Ollama Proxy

一个基于 Gin 框架的 Ollama API 代理服务，提供认证和请求转发功能，同时支持 OpenAI 风格的 API 接口。

## 功能特点

- 基于 Gin 框架开发
- 支持 Token 认证
- 支持跨域请求
- 完整的错误处理
- 支持所有 Ollama API 接口
- 支持 OpenAI 风格的 API 接口
- 日志记录功能（记录请求信息、错误信息等）

## 配置文件

配置文件为 `config.yaml`，结构如下：

```yaml
auth:
  generate_tokens:    # 生成相关接口的token列表
    - "your-generate-token-1"
    - "your-generate-token-2"
service:
  base_url: "http://localhost:11434"  # Ollama 服务地址
```

## 运行方式

1. 确保已安装 Go 环境
2. 配置 `config.yaml` 文件
3. 运行服务：

```bash
go run .
```

服务默认运行在 8080 端口。

## 日志功能

系统会自动记录以下信息：

- 所有API请求的访问日志
- 认证失败的错误日志
- 代理转发过程中的错误日志
- 系统运行状态日志

日志文件存储在项目根目录下的 `logs` 文件夹中。

## API 接口文档

所有接口都需要在请求头中携带 `Authorization` Token 进行认证，使用 `generate_tokens` 中的token。

### Ollama 原生接口

#### 1. 获取模型列表

- 请求方法：GET
- 请求路径：/api/tags
- 请求头：

  ```json
  Authorization: your-generate-token
  Content-Type: application/json
  ```

- 响应示例：

  ```json
  {
    "models": [
        {
            "name": "deepseek-r1:32b",
            "model": "deepseek-r1:32b",
            "modified_at": "2025-02-26T22:08:41.599376152+08:00",
            "size": 19851337640,
            "digest": "38056bbcbb2d068501ecb2d5ea9cea9dd4847465f1ab88c4d4a412a9f7792717",
            "details": {
                "parent_model": "",
                "format": "gguf",
                "family": "qwen2",
                "families": [
                    "qwen2"
                ],
                "parameter_size": "32.8B",
                "quantization_level": "Q4_K_M"
            }
        }
    ]
  }
  ```

#### 2. 拉取模型

- 请求方法：POST
- 请求路径：/api/pull
- 请求头：

  ```json
  Authorization: your-generate-token
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
- 请求路径：/api/delete
- 请求头：

  ```json
  Authorization: your-generate-token
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
- 请求路径：/api/copy
- 请求头：

  ```json
  Authorization: your-generate-token
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
- 请求路径：/api/push
- 请求头：

  ```json
  Authorization: your-generate-token
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

- 请求方法：GET
- 请求路径：/api/show
- 请求头：

  ```json
  Authorization: your-generate-token
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

#### 7. 生成文本

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

#### 8. 聊天

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

#### 9. 嵌入向量

- 请求方法：POST
- 请求路径：/api/embed
- 请求头：

  ```json
  Authorization: your-generate-token
  Content-Type: application/json
  ```

- 请求体：

  ```json
  {
    "model": "llama2",
    "input": ["Hello World", "Another text to embed"]
  }
  ```

- 响应示例：

  ```json
  {
    "embeddings": [
      [0.1, 0.2, 0.3, ...],
      [0.2, 0.3, 0.4, ...]
    ],
    "model": "llama2",
    "prompt_eval_count": 8
  }
  ```

  或单个输入的响应：

  ```json
  {
    "embedding": [0.1, 0.2, 0.3, ...],
    "model": "llama2"
  }
  ```

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

- 请求参数结构：

  ```json
  {
    "model": "string",       // 模型名称
    "messages": [            // 对话消息列表
      {
        "role": "string",    // 角色：user/assistant
        "content": "string"  // 消息内容
      }
    ],
    "stream": boolean,       // 是否流式输出
    "temperature": number,   // 温度参数
    "top_p": number,         // Top-p采样参数
    "max_tokens": number,    // 最大生成token数
    "n": number,             // 生成数量
    "stop": ["string"],      // 停止词
    "presence_penalty": number,  // 存在惩罚
    "frequency_penalty": number, // 频率惩罚
    "logit_bias": {},        // 逻辑偏差
    "user": "string",        // 用户标识
    "options": {             // 选项
      "temperature": number,
      "top_p": number
    }
  }
  ```

- 请求示例：

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

- 响应参数结构：

  ```json
  {
    "id": "string",      // 响应ID
    "object": "string",  // 对象类型
    "created": number,   // 创建时间
    "model": "string",   // 模型名称
    "choices": [         // 生成的选项列表
      {
        "message": {      // 生成的消息
          "role": "string",
          "content": "string"
        },
        "index": number,        // 选项索引
        "finish_reason": "string" // 结束原因
      }
    ],
    "usage": {           // 使用统计
      "prompt_tokens": number,
      "completion_tokens": number,
      "total_tokens": number
    }
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
        "content": "你好！高兴见到你。"
      },
      "finish_reason": "stop"
    }],
    "usage": {
      "prompt_tokens": 10,
      "completion_tokens": 20,
      "total_tokens": 30
    }
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

- 请求参数结构：

  ```json
  {
    "model": "string",       // 模型名称
    "prompt": "string",      // 提示文本
    "stream": boolean,       // 是否流式输出
    "temperature": number,   // 温度参数
    "top_p": number,         // Top-p采样参数
    "max_tokens": number,    // 最大生成token数
    "n": number,             // 生成数量
    "stop": ["string"],      // 停止词
    "presence_penalty": number,  // 存在惩罚
    "frequency_penalty": number, // 频率惩罚
    "logprobs": number,      // 日志概率
    "best_of": number,       // 最佳数量
    "user": "string",        // 用户标识
    "options": {             // 选项
      "temperature": number,
      "top_p": number
    }
  }
  ```

- 请求示例：

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
      "text": "从前有一座山，山上有一座庙，庙里有一个老和尚在讲故事...",
      "index": 0,
      "finish_reason": "stop"
    }],
    "usage": {
      "prompt_tokens": 5,
      "completion_tokens": 25,
      "total_tokens": 30
    }
  }
  ```

#### 3. 嵌入向量接口

- 请求方法：POST
- 请求路径：/v1/embeddings
- 请求头：

  ```json
  Authorization: your-generate-token
  Content-Type: application/json
  ```

- 请求参数结构：

  ```json
  {
    "model": "string",    // 模型名称
    "input": ["string"]   // 输入文本数组，支持批量处理
  }
  ```

- 请求示例：

  ```json
  {
    "model": "llama2",
    "input": ["Hello World", "Another text to embed"]
  }
  ```

  或单个输入：

  ```json
  {
    "model": "llama2",
    "input": ["Hello World"]
  }
  ```

- 响应参数结构：

  ```json
  {
    "object": "string",   // 对象类型
    "data": [             // 嵌入向量数据
      {
        "object": "string",  // 对象类型
        "embedding": [number],  // 嵌入向量
        "index": number      // 索引
      }
    ],
    "model": "string",    // 模型名称
    "usage": {           // 使用统计
      "prompt_tokens": number,
      "total_tokens": number
    }
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
    "model": "llama2",
    "usage": {
      "prompt_tokens": 2,
      "total_tokens": 2
    }
  }
  ```

#### 4. 模型列表接口

- 请求方法：GET
- 请求路径：/v1/models
- 请求头：

  ```json
  Authorization: your-generate-token
  Content-Type: application/json
  ```

- 响应参数结构：

  ```json
  {
    "object": "string",   // 对象类型，固定为"list"
    "data": [             // 模型数据列表
      {
        "id": "string",      // 模型ID
        "object": "string",  // 对象类型，固定为"model"
        "created": number,   // 创建时间戳
        "owned_by": "string" // 模型所有者
      }
    ]
  }
  ```

- 响应示例：

  ```json
  {
    "object": "list",
    "data": [
      {
        "id": "llama2",
        "object": "model",
        "created": 1677610602,
        "owned_by": "ollama"
      },
      {
        "id": "deepseek-r1:32b",
        "object": "model",
        "created": 1677649963,
        "owned_by": "ollama"
      }
    ]
  }
  ```

## 部署方式

### Docker 部署

1. 构建 Docker 镜像：

```bash
docker build -t ollama-proxy .
```

2. 运行容器：

```bash
docker run -d -p 8080:8080 -v $(pwd)/config.yaml:/app/config.yaml ollama-proxy
```

### 二进制部署

1. 编译项目：

```bash
go build -o ollama-proxy
```

2. 运行服务：

```bash
./ollama-proxy
```

## 常见问题

1. **认证失败**
   - 检查 `config.yaml` 中的 `generate_tokens` 配置是否正确
   - 确保请求头中包含正确的 `Authorization` 信息（不带Bearer前缀）
   - 系统会返回 "未提供认证token" 或 "非授权访问" 的错误信息

2. **无法连接到 Ollama 服务**
   - 检查 Ollama 服务是否正常运行
   - 确认 `config.yaml` 中的 `base_url` 配置是否正确（默认为 "<http://localhost:11434"）>
   - 系统会返回 "Failed to connect to Ollama service" 的错误信息

3. **流式输出不正常**
   - 确保客户端支持 SSE (Server-Sent Events) 格式
   - 检查请求体中 `stream` 参数是否设置为 `true`
   - 确保网络连接稳定，避免连接中断
   - 如果使用OpenAI风格API，确保正确处理流式响应格式
