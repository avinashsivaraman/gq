# gq

Simple CLI tool to ask questions using LLMs (Large Language Models).

## Installation

```
chmod +x ./install.sh
./install.sh
```

## Usage

Run `gq` followed by your question.

Examples:

```
gq "What is the capital of France?"
gq "France" -q "What is the capital of this Country?"
gq -q "What is the capital of France?"
cat test_file.txt | gq -q "Explain this to me" 
```

## API Key

To use a specific LLM model, create a `.gq.yaml` file in your home directory and provide the API key and model specifications.

```yaml
gemini:
  apiKey: <API_KEY>
  modelName: gemini-1.0-pro
  temperature: 0.7
  maxOutputTokens: 256
```

## Supported Models

- Gemini

## Future Plans

- Support for OpenAI models
- Support for Ollama
