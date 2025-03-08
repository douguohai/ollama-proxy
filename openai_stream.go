package main

import "time"

// ConvertOllamaChatStreamResponse 将Ollama流式响应转换为OpenAI格式
func ConvertOllamaChatStreamResponse(ollamaResp map[string]interface{}, model string) OpenAIChatResponse {

	return OpenAIChatResponse{
		ID:      generateResponseID(),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []ChatChoice{
			{
				Message: ChatMessage{
					Role:    "assistant",
					Content: ollamaResp["message"].(map[string]interface{})["content"].(string),
				},
				Index:        0,
				FinishReason: "",
			},
		},
	}
}

// ConvertOllamaGenerateStreamResponse 将Ollama流式响应转换为OpenAI格式
func ConvertOllamaGenerateStreamResponse(ollamaResp map[string]interface{}, model string) OpenAICompletionResponse {
	return OpenAICompletionResponse{
		ID:      generateResponseID(),
		Object:  "text_completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []Choice{
			{
				Text:         ollamaResp["response"].(string),
				Index:        0,
				FinishReason: "",
			},
		},
	}
}
