package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// OpenAIChatRequest OpenAI风格的聊天请求
type OpenAIChatRequest struct {
	Model            string          `json:"model"`
	Messages         []ChatMessage   `json:"messages"`
	Stream           bool            `json:"stream"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	Temperature      float64         `json:"temperature,omitempty"`
	TopP             float64         `json:"top_p,omitempty"`
	N                int             `json:"n,omitempty"`
	Stop             []string        `json:"stop,omitempty"`
	PresencePenalty  float64         `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64         `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int  `json:"logit_bias,omitempty"`
	User             string          `json:"user,omitempty"`
	Options          *RequestOptions `json:"options,omitempty"`
}

// OpenAICompletionRequest OpenAI风格的生成请求
type OpenAICompletionRequest struct {
	Model            string          `json:"model"`
	Prompt           string          `json:"prompt"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	Temperature      float64         `json:"temperature,omitempty"`
	TopP             float64         `json:"top_p,omitempty"`
	N                int             `json:"n,omitempty"`
	Stream           bool            `json:"stream"`
	Logprobs         int             `json:"logprobs,omitempty"`
	Stop             []string        `json:"stop,omitempty"`
	PresencePenalty  float64         `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64         `json:"frequency_penalty,omitempty"`
	BestOf           int             `json:"best_of,omitempty"`
	User             string          `json:"user,omitempty"`
	Options          *RequestOptions `json:"options,omitempty"`
}

// OpenAIEmbeddingRequest OpenAI风格的Embedding请求
type OpenAIEmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	User  string `json:"user,omitempty"`
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
	Model string `json:"model"`
	Input string `json:"input"`
}

// OpenAIChatResponse OpenAI风格的聊天响应
type OpenAIChatResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Usage   Usage        `json:"usage"`
	Choices []ChatChoice `json:"choices"`
}

// OpenAICompletionResponse OpenAI风格的生成响应
type OpenAICompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}

// OpenAIEmbeddingResponse OpenAI风格的Embedding响应
type OpenAIEmbeddingResponse struct {
	Object string            `json:"object"`
	Data   []EmbeddingResult `json:"data"`
	Model  string            `json:"model"`
	Usage  Usage             `json:"usage"`
}

// Usage Token使用统计
type Usage struct {
	PromptTokens     float64 `json:"prompt_tokens"`
	CompletionTokens float64 `json:"completion_tokens"`
	TotalTokens      float64 `json:"total_tokens"`
}

// ChatChoice 聊天响应选项
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// Choice 响应选项
type Choice struct {
	Text         string    `json:"text"`
	Index        int       `json:"index"`
	Logprobs     *Logprobs `json:"logprobs,omitempty"`
	FinishReason string    `json:"finish_reason"`
}

// Logprobs 日志概率
type Logprobs struct {
	Tokens        []string             `json:"tokens"`
	TokenLogprobs []float64            `json:"token_logprobs"`
	TopLogprobs   []map[string]float64 `json:"top_logprobs"`
	TextOffset    []int                `json:"text_offset"`
}

// EmbeddingResult Embedding结果
type EmbeddingResult struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// ConvertOllamaChatResponse 将Ollama响应转换为OpenAI格式
func ConvertOllamaChatResponse(ollamaResp map[string]interface{}, model string) OpenAIChatResponse {
	// 计算token使用情况
	promptTokens := 0.0
	if prompt, ok := ollamaResp["prompt_eval_count"].(float64); ok {
		// 简单估算token数量，实际应该使用分词器
		promptTokens = prompt
	}

	evalCount := 0.0
	if eval, ok := ollamaResp["eval_count"].(float64); ok {
		// 简单估算token数量，实际应该使用分词器
		evalCount = eval
	}

	return OpenAIChatResponse{
		ID:      generateResponseID(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []ChatChoice{
			{
				Message: ChatMessage{
					Role:    "assistant",
					Content: ollamaResp["message"].(map[string]interface{})["content"].(string),
				},
				Index:        0,
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: evalCount,
			TotalTokens:      promptTokens + evalCount,
		},
	}
}

// ConvertOllamaGenerateResponse 将Ollama响应转换为OpenAI格式
func ConvertOllamaGenerateResponse(ollamaResp map[string]interface{}, model string) OpenAICompletionResponse {
	// 计算token使用情况
	promptTokens := 0.0
	if prompt, ok := ollamaResp["prompt_eval_count"].(float64); ok {
		// 简单估算token数量，实际应该使用分词器
		promptTokens = prompt
	}

	evalCount := 0.0
	if eval, ok := ollamaResp["eval_count"].(float64); ok {
		// 简单估算token数量，实际应该使用分词器
		evalCount = eval
	}

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
		Usage: Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: evalCount,
			TotalTokens:      promptTokens + evalCount,
		},
	}
}

// ConvertOllamaEmbeddingResponse 将Ollama响应转换为OpenAI格式
func ConvertOllamaEmbeddingResponse(ollamaResp map[string]interface{}, model string) OpenAIEmbeddingResponse {
	str, _ := json.Marshal(ollamaResp)
	fmt.Println(string(str))
	// 检查响应中是否包含embeddings字段
	embeddings, ok := ollamaResp["embeddings"]
	if !ok || embeddings == nil {
		fmt.Printf("Error: embeddings field missing in response: %+v\n", ollamaResp)
		return OpenAIEmbeddingResponse{
			Object: "error",
			Data:   []EmbeddingResult{},
			Model:  model,
		}
	}

	// 检查embeddings类型是否正确（应该是二维数组）
	embeddingsSlice, isSlice := embeddings.([]interface{})
	if !isSlice || len(embeddingsSlice) == 0 {
		fmt.Printf("Error: invalid embeddings type %T in response: %+v\n", embeddings, ollamaResp)
		return OpenAIEmbeddingResponse{
			Object: "error",
			Data:   []EmbeddingResult{},
			Model:  model,
		}
	}

	// 处理所有embedding向量
	var embeddingResults []EmbeddingResult
	for i, embedding := range embeddingsSlice {
		// 转换embedding数据为float64切片
		embeddingData := convertToFloat64Slice(embedding)

		// 创建EmbeddingResult对象
		embeddingResults = append(embeddingResults, EmbeddingResult{
			Object:    "embedding",
			Embedding: embeddingData,
			Index:     i,
		})
	}

	// 计算token使用情况
	promptTokens := 0.0
	if prompt, ok := ollamaResp["prompt_eval_count"].(float64); ok {
		// 简单估算token数量，实际应该使用分词器
		promptTokens = prompt
	}

	return OpenAIEmbeddingResponse{
		Object: "list",
		Data:   embeddingResults,
		Model:  model,
		Usage: Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: 0,
			TotalTokens:      promptTokens,
		},
	}
}

// 辅助函数
func generateResponseID() string {
	return "chatcmpl-" + time.Now().Format("20060102150405")
}

// convertToFloat64Slice 将interface{}类型的embedding数据转换为float64切片
func convertToFloat64Slice(data interface{}) []float64 {
	if data == nil {
		return []float64{}
	}

	// 尝试将数据转换为[]interface{}
	slice, ok := data.([]interface{})
	if !ok {
		fmt.Printf("Error: embedding data is not a slice: %+v\n", data)
		return []float64{}
	}

	// 转换每个元素为float64
	result := make([]float64, len(slice))
	for i, v := range slice {
		switch value := v.(type) {
		case float64:
			result[i] = value
		case float32:
			result[i] = float64(value)
		case int:
			result[i] = float64(value)
		case int64:
			result[i] = float64(value)
		default:
			fmt.Printf("Error: unexpected type in embedding data: %T\n", v)
			return []float64{}
		}
	}

	return result
}

// OpenAIError OpenAI标准错误格式
type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}
