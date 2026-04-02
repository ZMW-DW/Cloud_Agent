package cloud_agent

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

type NLUResult struct {
	Intent string `json:"intent"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}

func TestInference(t *testing.T) {
	const envPath string = "/home/david/work/cloud_agent/.env"
	pm := NewPrompt_Manager()
	pm.RegisterFromYAML("/home/david/work/cloud_agent/configs/prompts.yaml")
	llmManager := NewLLM_Manager(pm)
	llmManager.AddLLM(getEnv(t, envPath))

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 使用 WaitGroup 等待所有并发任务完成
	var wg sync.WaitGroup
	wg.Add(2) // 我们有两个任务：Task1 (普通对话) 和 Task2 (结构化 NLU)

	// 1. 异步任务一：身份问询
	var got string
	var err1 error
	query1 := "你是由谁开发的" // 使用独立的变量防止竞态
	go func() {
		defer wg.Done()
		got, err1 = llmManager.Inference(ctx, "deepseek-chat", "default", query1)
	}()

	// 2. 异步任务二：NLU 解析
	var res NLUResult
	var err2 error
	query2 := "帮我把空调调到26度"
	go func() {
		defer wg.Done()
		// 直接在闭包里调用，不需要单独写 Async 函数，这样更灵活
		err2 = llmManager.InferenceStructured(ctx, "deepseek-chat", "nlu_parse", &res, query2)
	}()

	t.Log("🚀 两个 AI 任务已同时发出，正在并行处理...")

	// 3. 阻塞等待：直到两个 go func 都执行了 wg.Done()
	wg.Wait()

	// 4. 统一处理错误和结果
	if err1 != nil {
		t.Errorf("任务1出错: %v", err1)
	}
	if err2 != nil {
		t.Errorf("任务2出错: %v", err2)
	}

	// 5. 打印最终结果
	if err1 == nil && err2 == nil {
		t.Logf("NLU 结果 -> Intent: %s, Target: %s, Value: %d", res.Intent, res.Target, res.Value)
		t.Logf("身份回复 -> %s", got)
		t.Log("✅ 并发测试全部通过！")
	}
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
