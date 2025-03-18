package main

import "time"

// ConvertOllamaChatStreamResponse 将Ollama流式响应转换为OpenAI格式
func ConvertOllamaChatStreamResponse(ollamaResp map[string]interface{}, model string) StreamOpenAIChatResponse {

	return StreamOpenAIChatResponse{
		ID:      generateResponseID(),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []StreamChatChoice{
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
func ConvertOllamaGenerateStreamResponse(ollamaResp map[string]interface{}, model string) StreamOpenAICompletionResponse {
	return StreamOpenAICompletionResponse{
		ID:      generateResponseID(),
		Object:  "text_completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []StreamChoice{
			{
				Text:         ollamaResp["response"].(string),
				Index:        0,
				FinishReason: "",
			},
		},
	}
}
