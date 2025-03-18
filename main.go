package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/douguohai/ollama-proxy/models"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/douguohai/ollama-proxy/base"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// Config 配置结构体
type Config struct {
	Auth struct {
		GenerateTokens []string `yaml:"generate_tokens"`
		ModelTokens    []string `yaml:"model_tokens"`
	} `yaml:"auth"`
	Service struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"service"`
}

const (
	configFile = "config.yaml"
)

var logger *base.Logger

func main() {
	// 初始化日志记录器
	var err error
	logger, err = base.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// 读取配置文件
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(gin.Logger())
	// 添加日志中间件
	r.Use(logMiddleware())
	// 添加全局异常处理
	r.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": "服务器内部错误",
				})
				c.Abort()
			}
		}()
		c.Next()
	})

	// 配置CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API路由组，用于生成相关功能
	api := r.Group("/api", authMiddleware(*config))
	{
		// 生成相关接口
		api.POST("/generate", proxyOllama("/api/generate"))
		api.POST("/chat", proxyOllama("/api/chat"))
		api.POST("/embed", proxyOllama("/api/embed"))
		api.GET("/tags", proxyOllama("/api/tags"))
		api.POST("/pull", proxyOllama("/api/pull"))
		api.DELETE("/delete", proxyOllama("/api/delete"))
		api.POST("/copy", proxyOllama("/api/copy"))
		api.POST("/push", proxyOllama("/api/push"))
		api.GET("/show", proxyOllama("/api/show"))
	}

	// OpenAI风格的API路由组
	openai := r.Group("/v1", authMiddleware(*config))
	{
		// OpenAI风格的生成相关接口
		openai.GET("/models", handleOpenAIModels)
		openai.POST("/chat/completions", handleOpenAIChat)
		openai.POST("/completions", handleOpenAICompletion)
		openai.POST("/embeddings", handleOpenAIEmbedding)
	}

	r.Run(":8080")
}

func loadConfig() (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	// 解析YAML配置
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// authMiddleware 认证中间件
func authMiddleware(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的token
		token := c.GetHeader("Authorization")
		if token == "" {
			// 记录未提供token的情况
			logger.LogRequest(c.Request.Method, c.Request.URL.Path, nil, nil, fmt.Errorf("未提供认证token"), "", false)
			c.JSON(http.StatusOK, gin.H{
				"error": "未提供认证token",
			})
			c.Abort()
			return
		}

		// 检查token是否以Bearer开头，如果是则移除前缀
		token = strings.TrimPrefix(token, "Bearer ")

		// 根据路由组判断使用哪种token验证
		var validToken = false
		if strings.HasPrefix(c.Request.URL.Path, "/api") || strings.HasPrefix(c.Request.URL.Path, "/v1") {
			// 生成相关接口使用生成token
			for _, allowedToken := range config.Auth.GenerateTokens {
				if token == allowedToken {
					validToken = true
					break
				}
			}
		}

		// 记录token验证结果
		logger.LogRequest(c.Request.Method, c.Request.URL.Path, nil, nil, nil, token, validToken)

		if !validToken {
			c.JSON(http.StatusOK, gin.H{
				"error": "非授权访问",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// proxyOllama 创建一个代理处理函数
// OpenAI风格的API处理函数
func handleOpenAIChat(c *gin.Context) {
	var openAIReq models.OpenAIChatRequest
	if err := c.ShouldBindJSON(&openAIReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 转换为Ollama请求格式
	ollamaReq := models.OllamaChatRequest{
		Model:    openAIReq.Model,
		Messages: openAIReq.Messages,
		Stream:   openAIReq.Stream,
		Options:  openAIReq.Options,
	}

	if openAIReq.Stream {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		// 创建HTTP客户端和请求
		config, err := loadConfig()
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Error()})
			return
		}

		baseURL := "http://localhost:11434"
		if config.Service.BaseURL != "" {
			baseURL = config.Service.BaseURL
		}

		// 将请求数据转换为JSON
		jsonData, err := json.Marshal(ollamaReq)
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Error()})
			return
		}

		// 创建请求
		req, err := http.NewRequest("POST", baseURL+"/api/chat", strings.NewReader(string(jsonData)))
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Error()})
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Connection", "keep-alive")

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// 读取流式响应
		reader := bufio.NewReader(resp.Body)

		c.Stream(func(w io.Writer) bool {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					c.SSEvent("error", gin.H{"error": fmt.Sprintf("读取响应流出错: %v", err)})
				}
				return false
			}

			// 跳过空行
			if len(bytes.TrimSpace(line)) == 0 {
				return true
			}

			var result map[string]interface{}
			if err := json.Unmarshal(line, &result); err != nil {
				return true
			}

			// 转换为OpenAI流式响应格式
			openaiResp := models.ConvertOllamaChatStreamResponse(result, ollamaReq.Model)
			jsonData, err := json.Marshal(openaiResp)
			if err != nil {
				c.SSEvent("error", gin.H{"error": err.Error()})
				return false
			}
			c.SSEvent("message", string(jsonData))

			// 检查是否是最后一条消息
			if done, ok := result["done"].(bool); ok && done {
				return false
			}

			return true
		})
		return
	}

	// 非流式请求处理
	resp, err := sendToOllama("/api/chat", ollamaReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 检查响应中是否包含错误信息
	if errMsg, ok := resp["error"].(string); ok && errMsg != "" {
		c.JSON(http.StatusOK, gin.H{
			"error": errMsg,
		})
		return
	}

	// 转换为OpenAI响应格式
	openaiResp := models.ConvertOllamaChatResponse(resp, openAIReq.Model)
	c.JSON(http.StatusOK, openaiResp)
}

func handleOpenAICompletion(c *gin.Context) {
	var openAIReq models.OpenAICompletionRequest
	if err := c.ShouldBindJSON(&openAIReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 转换为Ollama请求格式
	ollamaReq := models.OllamaGenerateRequest{
		Model:   openAIReq.Model,
		Prompt:  openAIReq.Prompt,
		Stream:  openAIReq.Stream,
		Options: openAIReq.Options,
	}

	if openAIReq.Stream {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		// 创建HTTP客户端和请求
		config, err := loadConfig()
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Error()})
			return
		}

		baseURL := "http://localhost:11434"
		if config.Service.BaseURL != "" {
			baseURL = config.Service.BaseURL
		}

		// 将请求数据转换为JSON
		jsonData, err := json.Marshal(ollamaReq)
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Error()})
			return
		}

		// 创建请求
		req, err := http.NewRequest("POST", baseURL+"/api/generate", strings.NewReader(string(jsonData)))
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Error()})
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Connection", "keep-alive")

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.SSEvent("error", gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// 读取流式响应
		reader := bufio.NewReader(resp.Body)

		c.Stream(func(w io.Writer) bool {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					c.SSEvent("error", gin.H{"error": fmt.Sprintf("读取响应流出错: %v", err)})
				}
				return false
			}

			// 跳过空行
			if len(bytes.TrimSpace(line)) == 0 {
				return true
			}

			var result map[string]interface{}
			if err := json.Unmarshal(line, &result); err != nil {
				return true
			}

			// 转换为OpenAI流式响应格式
			openaiResp := models.ConvertOllamaGenerateStreamResponse(result, ollamaReq.Model)
			jsonData, err := json.Marshal(openaiResp)
			if err != nil {
				c.SSEvent("error", gin.H{"error": err.Error()})
				return false
			}
			c.SSEvent("message", string(jsonData))

			// 检查是否是最后一条消息
			if done, ok := result["done"].(bool); ok && done {
				return false
			}

			return true
		})
		return
	}

	// 非流式请求处理
	resp, err := sendToOllama("/api/generate", ollamaReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 检查响应中是否包含错误信息
	if errMsg, ok := resp["error"].(string); ok && errMsg != "" {
		c.JSON(http.StatusOK, gin.H{
			"error": errMsg,
		})
		return
	}

	// 转换为OpenAI响应格式
	openaiResp := models.ConvertOllamaGenerateResponse(resp, openAIReq.Model)
	c.JSON(http.StatusOK, openaiResp)
}

func handleOpenAIEmbedding(c *gin.Context) {
	var req models.OpenAIEmbeddingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 直接将OpenAI请求中的输入数组传递给Ollama
	ollamaReq := models.OllamaEmbeddingRequest{
		Model: req.Model,
		Input: req.Input,
	}

	// 发送批量请求到Ollama服务
	resp, err := sendToOllama("/api/embed", ollamaReq)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 检查响应中是否包含错误信息
	if errMsg, ok := resp["error"].(string); ok && errMsg != "" {
		c.JSON(http.StatusOK, gin.H{
			"error": errMsg,
		})
		return
	}
	// 获取token使用量
	totalTokens := 0.0
	if prompt, ok := resp["prompt_eval_count"].(float64); ok {
		totalTokens = prompt
	}

	// 获取embeddings数据
	var allEmbeddings []interface{}
	if embeddings, ok := resp["embeddings"].([]interface{}); ok {
		allEmbeddings = embeddings
	} else if embedding, ok := resp["embedding"].([]interface{}); ok {
		// 兼容单个embedding的情况
		allEmbeddings = []interface{}{embedding}
	}

	// 构造包含所有embeddings的响应
	combinedResp := map[string]interface{}{
		"embeddings":        allEmbeddings,
		"prompt_eval_count": totalTokens,
	}

	// 转换为OpenAI响应格式
	openaiResp := models.ConvertOllamaEmbeddingResponse(combinedResp, req.Model)
	c.JSON(http.StatusOK, openaiResp)
}

// handleOpenAIModels 处理OpenAI风格的模型列表请求
func handleOpenAIModels(c *gin.Context) {
	// 调用Ollama的tags接口获取模型列表
	resp, err := sendToOllamaGet("/api/tags")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 检查响应中是否包含错误信息
	if errMsg, ok := resp["error"].(string); ok && errMsg != "" {
		c.JSON(http.StatusOK, gin.H{
			"error": errMsg,
		})
		return
	}

	// 转换为OpenAI响应格式
	openaiResp := models.ConvertOllamaModelsResponse(resp)
	c.JSON(http.StatusOK, openaiResp)
}

// 发送GET请求到Ollama服务的函数
func sendToOllamaGet(path string) (map[string]interface{}, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	baseURL := "http://localhost:11434"
	if config.Service.BaseURL != "" {
		baseURL = config.Service.BaseURL
	}

	// 创建GET请求
	req, err := http.NewRequest("GET", baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取完整的响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// 发送请求到Ollama服务的通用函数
func sendToOllama(path string, data interface{}) (map[string]interface{}, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	baseURL := "http://localhost:11434"
	if config.Service.BaseURL != "" {
		baseURL = config.Service.BaseURL
	}

	// 将请求数据转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// 创建请求
	req, err := http.NewRequest("POST", baseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取完整的响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func proxyOllama(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		config, err := loadConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to load configuration",
			})
			return
		}

		// 检查是否为流式请求
		var requestBody map[string]interface{}
		if err := c.ShouldBindJSON(&requestBody); err == nil {
			if stream, ok := requestBody["stream"].(bool); ok && stream {
				c.Header("Content-Type", "text/event-stream")
				c.Header("Cache-Control", "no-cache")
				c.Header("Connection", "keep-alive")
			}
		}

		// 创建代理请求
		baseURL := "http://localhost:11434"
		if config.Service.BaseURL != "" {
			baseURL = config.Service.BaseURL
		}

		// 重新创建请求体
		jsonData, _ := json.Marshal(requestBody)
		req, err := http.NewRequest(c.Request.Method, baseURL+path, bytes.NewReader(jsonData))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create request",
			})
			return
		}

		// 复制原始请求的header
		for name, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}
		req.Header.Set("Content-Type", "application/json")

		// 发送请求到Ollama服务
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to connect to Ollama service",
			})
			return
		}
		defer resp.Body.Close()

		// 复制响应header
		for name, values := range resp.Header {
			for _, value := range values {
				c.Header(name, value)
			}
		}

		// 设置响应状态码
		c.Status(resp.StatusCode)

		// 检查是否为流式请求
		if stream, ok := requestBody["stream"].(bool); ok && stream {
			// 流式请求使用Stream方式返回
			c.Stream(func(w io.Writer) bool {
				_, err := io.Copy(w, resp.Body)
				return err == nil
			})
		} else {
			// 非流式请求处理
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to read response body",
				})
				return
			}
			// 其他接口直接返回原始响应
			c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
		}
	}
}

// logMiddleware 日志中间件
func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求体
		var requestBody map[string]interface{}
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			json.Unmarshal(bodyBytes, &requestBody)
		}

		// 检查是否为流式请求
		isStream := false
		if stream, ok := requestBody["stream"].(bool); ok {
			isStream = stream
		}

		// 如果是流式请求，只记录请求参数
		if isStream {
			logger.LogRequest(c.Request.Method, c.Request.URL.Path, requestBody, nil, nil, "", true)
			c.Next()
			return
		}

		// 创建自定义ResponseWriter以捕获响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// 解析响应体
		var response interface{}
		if err := json.Unmarshal(blw.body.Bytes(), &response); err == nil {
			logger.LogRequest(c.Request.Method, c.Request.URL.Path, requestBody, response, nil, "", true)
		}
	}
}

// bodyLogWriter 用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
