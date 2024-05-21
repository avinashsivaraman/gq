# gq

GQ: Generative (AI) Query

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
**Alternatively, you can specify a provider in real-time, overriding the default provider set in the YAML config file.**
```gq -q "Hi" -p openAI```


## API Key

To use a specific LLM model, create a `.gq.yaml` file in your $HOME/.config/gq/ directory and provide the API key and model specifications.

`touch ~/.config/gq/.gq.yaml`

```yaml
default: gemini
gemini:
  apiKey: <API_KEY>
  modelName: gemini-1.0-pro
  temperature: 0.7
  maxOutputTokens: 1024
openAI:
  apiKey: <API_KEY>
  modelName: gpt-3.5-turbo # supported params: gpt-3.5-turbo, gpt-4, gpt-4-turbo
  temperature: 0.5
  maxOutputTokens: 1024
azureOpenAI:
  apiKey: <API_KEY>
  modelDeploymentID: <DEPLOYMENT_NAME>
  modelEndpoint: <ENDPOINT>
  temperature: 0.5
  maxOutputTokens: 1024
bedrock:
  modelName: <MODEL_NAME>
  awsProfile: <AWS_PROFILE> # AWS Profile Name which has access to the model in ~/.aws/credentials
  awsRegion: <AWS_REGION>
```

## Supported Models

- Gemini
- OpenAI
- AzureOpenAI
- Amazon Bedrock Integration

## Future Plans

- Support for Ollama

## Amazon Bedrock Setup

To setup AWS profile for Amazon Bedrock, follow the steps below: 
1. AWS Profile: https://docs.aws.amazon.com/cli/v1/userguide/cli-configure-files.html. Profile should have access to invoke the model.
2. Enable Amazon Bedrock Model: https://docs.aws.amazon.com/bedrock/latest/userguide/model-access.html
