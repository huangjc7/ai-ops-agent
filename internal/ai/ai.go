package ai

import (
	"ai-ops-agent/pkg/env"
	"fmt"
	"log"
)

type ClientProvider struct {
	TextGenTextModelClient   Controller
	ImageGenImageModelClient Controller
	ImageGenTextModelClient  Controller
}

const (
	TextGenTextModel   = "generate_text"
	ImageGenImageModel = "generate_image"
	ImageGenTextModel  = "vision_image"
)

func NewAIClientProvider() (*ClientProvider, error) {
	config := &Config{
		ApiKey:  env.Get("API_KEY", ""),
		BaseURL: env.Get("BASE_URL", "https://dashscope.aliyuncs.com/compatible-mode/v1"),
		Model:   env.Get("MODEL", "qwen3-max"),
	}

	textGenClient, err := InitModel(TextGenTextModel, config)
	if err != nil {
		return nil, fmt.Errorf("初始化 TextGenClient 失败: %w", err)
	}

	imageGenClient, err := InitModel(ImageGenImageModel, config)
	if err != nil {
		return nil, fmt.Errorf("初始化 ImageGenClient 失败: %w", err)
	}

	visionClient, err := InitModel(ImageGenTextModel, config)
	if err != nil {
		return nil, fmt.Errorf("初始化 VisionClient 失败: %w", err)
	}

	return &ClientProvider{
		TextGenTextModelClient:   textGenClient,
		ImageGenImageModelClient: imageGenClient,
		ImageGenTextModelClient:  visionClient,
	}, nil
}

func InitModel(modelType string, cfg *Config) (Controller, error) {
	aiclient := NewAIClient(&Config{
		ApiKey:  cfg.ApiKey,
		BaseURL: cfg.BaseURL,
		Model:   cfg.Model,
	})

	return aiclient, nil
}

func GetAIModel() *ClientProvider {
	AIClientP, err := NewAIClientProvider()
	if err != nil {
		log.Println(err)
		return nil
	}
	return AIClientP

}
