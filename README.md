# Cloud-Agent 🚀

**Cloud-Agent** is a high-performance, Go-based Multi-Agent framework designed for seamless LLM integration and NLU (Natural Language Understanding) tasks. It is built with a focus on low-latency execution, structured outputs, and concurrent agent management, making it ideal for edge deployment (e.g., MTK G520) and cloud services.

## ✨ Key Features

  * **Concurrent Execution**: Built on Go's goroutines to handle multiple LLM requests in parallel, significantly reducing response latency.
  * **Structured NLU**: Native support for mapping LLM outputs directly into Go structs using the `InferenceStructured` method.
  * **Prompt Management**: Decoupled Prompt Management via `prompts.yaml`, allowing for dynamic persona setting and system instruction updates without code changes.
  * **Thread-Safe Model Management**: Safely manage multiple LLM instances (DeepSeek, OpenAI, etc.) using `sync.RWMutex`.
  * **Environment Agnostic**: Easy configuration via `.env` files and flexible directory structures.

## 📁 Project Structure

```text
.
├── cmd/                # Entry points for the application
├── configs/            # YAML configurations (Prompts, System Instructions)
├── internal/
│   └── llm/            # Core LLM & Prompt Manager logic
│       ├── llm.go      # Model & Client management
│       ├── prompt.go   # Template & Message building
│       └── llm_test.go # Concurrent & Structured testing
├── README.md
└── TODO.md             # Roadmap and pending features
```

## 🚀 Getting Started

### 1\. Prerequisites

  * Go 1.22.2+ (Recommended: Go 1.24 for latest Delve support)
  * An API Key from a supported Provider (OpenAI, DeepSeek, etc.)

### 2\. Configuration

Create a `.env` file in the root directory:

```env
OPENAI_API_KEY=your_api_key_here
OPENAI_BASE_URL=https://api.deepseek.com
MODEL_NAME=deepseek-chat
```

### 3\. Usage Example

```go
// Initialize Managers
pm := llm.NewPrompt_Manager()
pm.RegisterFromYAML("configs/prompts.yaml")
llmManager := llm.NewLLM_Manager(pm)

// Add a model
llmManager.AddLLM(apiKey, baseURL, "deepseek-chat")

// Perform structured NLU inference
var result NLUResult
err := llmManager.InferenceStructured(ctx, "deepseek-chat", "nlu_parse", &result, "Set AC to 26C")
```

## 🧪 Testing

The project includes a comprehensive test suite demonstrating concurrent agent calls:

```bash
cd internal/llm
go test -v
```

## 🛠 Roadmap (TODO)

  * [ ] Implement WebSocket support for real-time streaming.
  * [ ] Add Model Context Protocol (MCP) for tool calling.
  * [ ] Optimize token pruning for edge device deployment.

-----

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

**Developer:** [David (Dai Wei)](https://www.google.com/search?q=https://github.com/ZMW-DW)

