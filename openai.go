package main

import (
	"time"
)

// OpenAIChatRequest OpenAI风格的聊天请求
type OpenAIChatRequest struct {
	Model    string          `json:"model"`
	Messages []ChatMessage   `json:"messages"`
	Stream   bool            `json:"stream,omitempty"`
	Options  *RequestOptions `json:"options,omitempty"`
}

// OpenAICompletionRequest OpenAI风格的生成请求
type OpenAICompletionRequest struct {
	Model   string          `json:"model"`
	Prompt  string          `json:"prompt"`
	Stream  bool            `json:"stream,omitempty"`
	Options *RequestOptions `json:"options,omitempty"`
}

// OpenAIEmbeddingRequest OpenAI风格的Embedding请求
type OpenAIEmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// ChatMessage 聊天消息结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// RequestOptions 请求选项
type RequestOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
}

// OllamaChatRequest Ollama聊天请求
type OllamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ChatMessage   `json:"messages"`
	Stream   bool            `json:"stream"`
	Options  *RequestOptions `json:"options,omitempty"`
}

// OllamaGenerateRequest Ollama生成请求
type OllamaGenerateRequest struct {
	Model   string          `json:"model"`
	Prompt  string          `json:"prompt"`
	Stream  bool            `json:"stream"`
	Options *RequestOptions `json:"options,omitempty"`
}

// OllamaEmbeddingRequest Ollama Embedding请求
type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OpenAIChatResponse OpenAI风格的聊天响应
type OpenAIChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

// OpenAICompletionResponse OpenAI风格的生成响应
type OpenAICompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

// OpenAIEmbeddingResponse OpenAI风格的Embedding响应
type OpenAIEmbeddingResponse struct {
	Object string            `json:"object"`
	Data   []EmbeddingResult `json:"data"`
	Model  string            `json:"model"`
}

// Choice 响应选项
type Choice struct {
	Message      ChatMessage `json:"message,omitempty"`
	Text         string      `json:"text,omitempty"`
	Index        int         `json:"index"`
	FinishReason string      `json:"finish_reason"`
}

// EmbeddingResult Embedding结果
type EmbeddingResult struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// ConvertOllamaChatResponse 将Ollama响应转换为OpenAI格式
func ConvertOllamaChatResponse(ollamaResp map[string]interface{}, model string) OpenAIChatResponse {
	return OpenAIChatResponse{
		ID:      generateResponseID(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []Choice{
			{
				Message: ChatMessage{
					Role:    "assistant",
					Content: ollamaResp["message"].(map[string]interface{})["content"].(string),
				},
				Index:        0,
				FinishReason: "stop",
			},
		},
	}
}

// ConvertOllamaGenerateResponse 将Ollama响应转换为OpenAI格式
func ConvertOllamaGenerateResponse(ollamaResp map[string]interface{}, model string) OpenAICompletionResponse {
	return OpenAICompletionResponse{
		ID:      generateResponseID(),
		Object:  "text_completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []Choice{
			{
				Text:         ollamaResp["response"].(string),
				Index:        0,
				FinishReason: "stop",
			},
		},
	}
}

// ConvertOllamaEmbeddingResponse 将Ollama响应转换为OpenAI格式
func ConvertOllamaEmbeddingResponse(ollamaResp map[string]interface{}, model string) OpenAIEmbeddingResponse {
	return OpenAIEmbeddingResponse{
		Object: "list",
		Data: []EmbeddingResult{
			{
				Object:    "embedding",
				Embedding: convertToFloat64Slice(ollamaResp["embedding"].([]interface{})),
				Index:     0,
			},
		},
		Model: model,
	}
}

// 辅助函数
func generateResponseID() string {
	return "chatcmpl-" + time.Now().Format("20060102150405")
}

func convertToFloat64Slice(input []interface{}) []float64 {
	result := make([]float64, len(input))
	for i, v := range input {
		result[i] = v.(float64)
	}
	return result
}
