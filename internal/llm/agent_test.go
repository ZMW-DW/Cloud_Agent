package cloud_agent

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestInference(t *testing.T) {
	const envPath string = "/home/david/work/cloud_agent/.env"
	prompt_manager := NewPrompt_Manager()
	prompt_manager.RegisterFromYAML("/home/david/work/cloud_agent/configs/prompts.yaml")
	llmManager := NewLLM_Manager(prompt_manager)
	llmManager.AddLLM(getEnv(t, envPath))

	query := "你是由谁开发的"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	got, err := llmManager.Inference(ctx, "deepseek-chat", "default", query)

	if err != nil {
		t.Fatalf("Inference 运行出错: %v", err)
	}

	if got == "" {
		t.Error("期望得到回复，但结果为空")
	}

	t.Logf("✅ 测试通过！AI 回复: %s", got)
}

func getEnv(t *testing.T, file_path string) (apiKey string, baseUrl string, modelName string) {
	_ = godotenv.Load(file_path)

	apiKey = os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("跳过测试：未检测到 OPENAI_API_KEY")
	}

	baseUrl = os.Getenv("OPENAI_BASE_URL")
	if baseUrl == "" {
		t.Skip("跳过测试：未检测到 OPENAI_BASE_URL")
	}

	modelName = os.Getenv("MODEL_NAME")
	if modelName == "" {
		t.Skip("跳过测试：未检测到 MODEL_NAME")
	}

	return apiKey,
		baseUrl,
		modelName

}
