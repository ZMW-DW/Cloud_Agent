package cloud_agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sashabaranov/go-openai"
)

type llm struct {
	client *openai.Client
	model  string
}

type LLM_Manager struct {
	models map[string]*llm
	mu     sync.RWMutex
	pm     *Prompt_Manager
}

func NewLLM_Manager(pm *Prompt_Manager) *LLM_Manager {
	return &LLM_Manager{
		models: make(map[string]*llm),
		pm:     pm,
	}
}

func (m *LLM_Manager) AddLLM(apiKey, baseURL, modelName string) {
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.models[modelName] = &llm{
		client: openai.NewClientWithConfig(config),
		model:  modelName,
	}
}

func (m *LLM_Manager) Inference(ctx context.Context, modelName string, promptName string, query ...interface{}) (string, error) {
	messages, err := m.pm.BuildMessages(promptName, query...)
	if err != nil {
		return "", err
	}

	m.mu.RLock()
	inst, ok := m.models[modelName]
	m.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("model [%s] not found", modelName)
	}

	resp, err := inst.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    inst.model,
			Messages: messages,
		},
	)

	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (m *LLM_Manager) InferenceStructured(ctx context.Context, modelName string, promptName string, target interface{}, query ...interface{}) error {
	messages, err := m.pm.BuildMessages(promptName, query...)
	if err != nil {
		return err
	}

	m.mu.RLock()
	inst, ok := m.models[modelName]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("model not found")
	}

	resp, err := inst.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    inst.model,
			Messages: messages,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)

	if err != nil {
		return err
	}

	content := resp.Choices[0].Message.Content
	err = json.Unmarshal([]byte(content), target)
	if err != nil {
		return fmt.Errorf("结构化解析失败: %v, 原始输出: %s", err, content)
	}

	return nil
}
