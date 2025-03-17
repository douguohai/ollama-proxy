package main

import (
	"time"
)

// OpenAIModelResponse OpenAI风格的模型列表响应
type OpenAIModelResponse struct {
	Object string      `json:"object"`
	Data   []ModelData `json:"data"`
}

// ModelData 模型数据
type ModelData struct {
	ID         string            `json:"id"`
	Object     string            `json:"object"`
	Created    int64             `json:"created"`
	OwnedBy    string            `json:"owned_by"`
	Permission []ModelPermission `json:"permission"`
	Root       string            `json:"root"`
	Parent     interface{}       `json:"parent"`
}

// ModelPermission 模型权限
type ModelPermission struct {
	ID                 string      `json:"id"`
	Object             string      `json:"object"`
	Created            int64       `json:"created"`
	AllowCreateEngine  bool        `json:"allow_create_engine"`
	AllowSampling      bool        `json:"allow_sampling"`
	AllowLogprobs      bool        `json:"allow_logprobs"`
	AllowSearchIndices bool        `json:"allow_search_indices"`
	AllowView          bool        `json:"allow_view"`
	AllowFineTuning    bool        `json:"allow_fine_tuning"`
	Organization       string      `json:"organization"`
	Group              interface{} `json:"group"`
	IsBlocking         bool        `json:"is_blocking"`
}

// ConvertOllamaModelsResponse 将Ollama模型列表响应转换为OpenAI格式
func ConvertOllamaModelsResponse(ollamaResp map[string]interface{}) OpenAIModelResponse {
	var modelDataList []ModelData

	// 检查响应中是否包含models字段
	models, ok := ollamaResp["models"].([]interface{})
	if !ok || models == nil {
		return OpenAIModelResponse{
			Object: "list",
			Data:   modelDataList,
		}
	}

	// 处理所有模型数据
	for _, model := range models {
		modelInfo, ok := model.(map[string]interface{})
		if !ok {
			continue
		}

		modelName, ok := modelInfo["name"].(string)
		if !ok {
			continue
		}

		// 创建模型数据
		modelData := ModelData{
			ID:      modelName,
			Object:  "model",
			Created: time.Now().Unix(),
			OwnedBy: "organization-owner",
			Root:    modelName,
			Permission: []ModelPermission{
				{
					ID:                 "modelperm-" + time.Now().Format("20060102150405"),
					Object:             "model_permission",
					Created:            time.Now().Unix(),
					AllowCreateEngine:  false,
					AllowSampling:      true,
					AllowLogprobs:      true,
					AllowSearchIndices: false,
					AllowView:          true,
					AllowFineTuning:    false,
					Organization:       "*",
					IsBlocking:         false,
				},
			},
		}

		modelDataList = append(modelDataList, modelData)
	}

	return OpenAIModelResponse{
		Object: "list",
		Data:   modelDataList,
	}
}
