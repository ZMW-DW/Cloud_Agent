package cloud_agent

import (
	"fmt"
	"os"
	"sync"

	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

// PromptConfig 对应 YAML 中的单个业务配置
type PromptConfig struct {
	System string `yaml:"system"`
	User   string `yaml:"user"` // 支持包含 %s 或 %d 的格式化字符串
}

type Prompt_Manager struct {
	prompts map[string]PromptConfig
	mu      sync.RWMutex
}

func NewPrompt_Manager() *Prompt_Manager {
	return &Prompt_Manager{
		prompts: make(map[string]PromptConfig),
	}
}

func (pm *Prompt_Manager) RegisterFromYAML(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取 YAML 失败: %v", err)
	}

	var tempMap map[string]PromptConfig
	if err := yaml.Unmarshal(data, &tempMap); err != nil {
		return fmt.Errorf("解析 YAML 格式错误: %v", err)
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()
	for k, v := range tempMap {
		pm.prompts[k] = v
	}
	return nil
}

func (pm *Prompt_Manager) BuildMessages(name string, args ...interface{}) ([]openai.ChatCompletionMessage, error) {
	pm.mu.RLock()
	config, ok := pm.prompts[name]
	pm.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("未找到 Prompt 配置: %s", name)
	}

	userContent := fmt.Sprintf(config.User, args...)

	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: config.System,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userContent,
		},
	}, nil
}
