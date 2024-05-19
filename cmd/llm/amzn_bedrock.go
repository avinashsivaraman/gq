package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/spf13/viper"
)

const (
	CLAUDE_MODEL_ID             = "anthropic.claude-v2"
	JURASSIC2_MODEL_ID          = "ai21.j2-mid-v1"
	LLAMA2_MODEL_ID             = "meta.llama2-13b-chat-v1"
	TITAN_IMAGE_MODEL_ID        = "amazon.titan-image-generator-v1"
	TITAN_TEXT_EXPRESS_MODEL_ID = "amazon.titan-text-express-v1"
)

type AmznBedrockAIProvider struct{}

func (_ AmznBedrockAIProvider) Chat(userQuery string, verbose bool) (string, error) {
	amznBedrock := viper.Sub("bedrock")

	modelName := amznBedrock.GetString("modelName")
	awsProfile := amznBedrock.GetString("awsProfile")
	awsRegion := amznBedrock.GetString("awsRegion")

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(awsProfile),
		config.WithRegion(awsRegion))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	// Create a new Bedrock Runtime client
	client := bedrockruntime.NewFromConfig(cfg)
	wrapper := InvokeModelWrapper{BedrockRuntimeClient: client}

	switch modelName {
	case CLAUDE_MODEL_ID:
		return wrapper.InvokeClaude(userQuery)
	case JURASSIC2_MODEL_ID:
		return wrapper.InvokeJurassic2(userQuery)
	case LLAMA2_MODEL_ID:
		return wrapper.InvokeLlama2(userQuery)
	case TITAN_IMAGE_MODEL_ID:
		return wrapper.InvokeTitanImage(userQuery, 0)
	case TITAN_TEXT_EXPRESS_MODEL_ID:
		return wrapper.InvokeTitanText(userQuery)
	default:
		return "", fmt.Errorf("modelID %s not found", modelName)
	}
}

// InvokeModelWrapper encapsulates Amazon Bedrock actions used in the examples.
// It contains a Bedrock Runtime client that is used to invoke foundation models.
type InvokeModelWrapper struct {
	BedrockRuntimeClient *bedrockruntime.Client
}

// Each model provider has their own individual request and response formats.
// For the format, ranges, and default values for Anthropic Claude, refer to:
// https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-claude.html

type ClaudeRequest struct {
	Prompt            string   `json:"prompt"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	Temperature       float64  `json:"temperature,omitempty"`
	StopSequences     []string `json:"stop_sequences,omitempty"`
}

type ClaudeResponse struct {
	Completion string `json:"completion"`
}

// Invokes Anthropic Claude on Amazon Bedrock to run an inference using the input
// provided in the request body.
func (wrapper InvokeModelWrapper) InvokeClaude(prompt string) (string, error) {
	modelId := "anthropic.claude-v2"

	// Anthropic Claude requires enclosing the prompt as follows:
	enclosedPrompt := "Human: " + prompt + "\n\nAssistant:"

	body, err := json.Marshal(ClaudeRequest{
		Prompt:            enclosedPrompt,
		MaxTokensToSample: 200,
		Temperature:       0.5,
		StopSequences:     []string{"\n\nHuman:"},
	})
	if err != nil {
		log.Fatal("failed to marshal", err)
	}

	output, err := wrapper.BedrockRuntimeClient.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		ProcessError(err, modelId)
	}

	var response ClaudeResponse
	if err := json.Unmarshal(output.Body, &response); err != nil {
		log.Fatal("failed to unmarshal", err)
	}

	return response.Completion, nil
}

// Each model provider has their own individual request and response formats.
// For the format, ranges, and default values for AI21 Labs Jurassic-2, refer to:
// https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-jurassic2.html

type Jurassic2Request struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"maxTokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

type Jurassic2Response struct {
	Completions []Completion `json:"completions"`
}
type Completion struct {
	Data Data `json:"data"`
}
type Data struct {
	Text string `json:"text"`
}

// Invokes AI21 Labs Jurassic-2 on Amazon Bedrock to run an inference using the input
// provided in the request body.
func (wrapper InvokeModelWrapper) InvokeJurassic2(prompt string) (string, error) {
	modelId := "ai21.j2-mid-v1"

	body, err := json.Marshal(Jurassic2Request{
		Prompt:      prompt,
		MaxTokens:   200,
		Temperature: 0.5,
	})
	if err != nil {
		log.Fatal("failed to marshal", err)
	}

	output, err := wrapper.BedrockRuntimeClient.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		ProcessError(err, modelId)
	}

	var response Jurassic2Response
	if err := json.Unmarshal(output.Body, &response); err != nil {
		log.Fatal("failed to unmarshal", err)
	}

	return response.Completions[0].Data.Text, nil
}

// Each model provider has their own individual request and response formats.
// For the format, ranges, and default values for Meta Llama 2 Chat, refer to:
// https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-meta.html

type Llama2Request struct {
	Prompt       string  `json:"prompt"`
	MaxGenLength int     `json:"max_gen_len,omitempty"`
	Temperature  float64 `json:"temperature,omitempty"`
}

type Llama2Response struct {
	Generation string `json:"generation"`
}

// Invokes Meta Llama 2 Chat on Amazon Bedrock to run an inference using the input
// provided in the request body.
func (wrapper InvokeModelWrapper) InvokeLlama2(prompt string) (string, error) {
	modelId := "meta.llama2-13b-chat-v1"

	body, err := json.Marshal(Llama2Request{
		Prompt:       prompt,
		MaxGenLength: 512,
		Temperature:  0.5,
	})
	if err != nil {
		log.Fatal("failed to marshal", err)
	}

	output, err := wrapper.BedrockRuntimeClient.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		ProcessError(err, modelId)
	}

	var response Llama2Response
	if err := json.Unmarshal(output.Body, &response); err != nil {
		log.Fatal("failed to unmarshal", err)
	}

	return response.Generation, nil
}

type TitanImageRequest struct {
	TaskType              string                `json:"taskType"`
	TextToImageParams     TextToImageParams     `json:"textToImageParams"`
	ImageGenerationConfig ImageGenerationConfig `json:"imageGenerationConfig"`
}
type TextToImageParams struct {
	Text string `json:"text"`
}
type ImageGenerationConfig struct {
	NumberOfImages int     `json:"numberOfImages"`
	Quality        string  `json:"quality"`
	CfgScale       float64 `json:"cfgScale"`
	Height         int     `json:"height"`
	Width          int     `json:"width"`
	Seed           int64   `json:"seed"`
}

type TitanImageResponse struct {
	Images []string `json:"images"`
}

// Invokes the Titan Image model to create an image using the input provided
// in the request body.
func (wrapper InvokeModelWrapper) InvokeTitanImage(prompt string, seed int64) (string, error) {
	modelId := "amazon.titan-image-generator-v1"

	body, err := json.Marshal(TitanImageRequest{
		TaskType: "TEXT_IMAGE",
		TextToImageParams: TextToImageParams{
			Text: prompt,
		},
		ImageGenerationConfig: ImageGenerationConfig{
			NumberOfImages: 1,
			Quality:        "standard",
			CfgScale:       8.0,
			Height:         512,
			Width:          512,
			Seed:           seed,
		},
	})
	if err != nil {
		log.Fatal("failed to marshal", err)
	}

	output, err := wrapper.BedrockRuntimeClient.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		ProcessError(err, modelId)
	}

	var response TitanImageResponse
	if err := json.Unmarshal(output.Body, &response); err != nil {
		log.Fatal("failed to unmarshal", err)
	}

	base64ImageData := response.Images[0]

	return base64ImageData, nil
}

// Each model provider has their own individual request and response formats.
// For the format, ranges, and default values for Amazon Titan Text, refer to:
// https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-titan-text.html
type TitanTextRequest struct {
	InputText            string               `json:"inputText"`
	TextGenerationConfig TextGenerationConfig `json:"textGenerationConfig"`
}

type TextGenerationConfig struct {
	Temperature   float64  `json:"temperature"`
	TopP          float64  `json:"topP"`
	MaxTokenCount int      `json:"maxTokenCount"`
	StopSequences []string `json:"stopSequences,omitempty"`
}

type TitanTextResponse struct {
	InputTextTokenCount int      `json:"inputTextTokenCount"`
	Results             []Result `json:"results"`
}

type Result struct {
	TokenCount       int    `json:"tokenCount"`
	OutputText       string `json:"outputText"`
	CompletionReason string `json:"completionReason"`
}

func (wrapper InvokeModelWrapper) InvokeTitanText(prompt string) (string, error) {
	modelId := "amazon.titan-text-express-v1"

	body, err := json.Marshal(TitanTextRequest{
		InputText: prompt,
		TextGenerationConfig: TextGenerationConfig{
			Temperature:   0,
			TopP:          1,
			MaxTokenCount: 4096,
		},
	})
	if err != nil {
		log.Fatal("failed to marshal", err)
	}

	output, err := wrapper.BedrockRuntimeClient.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		ProcessError(err, modelId)
	}

	var response TitanTextResponse
	if err := json.Unmarshal(output.Body, &response); err != nil {
		log.Fatal("failed to unmarshal", err)
	}

	return response.Results[0].OutputText, nil
}

func ProcessError(err error, modelId string) {
	errMsg := err.Error()
	if strings.Contains(errMsg, "no such host") {
		log.Printf(`The Bedrock service is not available in the selected region.
                    Please double-check the service availability for your region at
                    https://aws.amazon.com/about-aws/global-infrastructure/regional-product-services/.\n`)
	} else if strings.Contains(errMsg, "Could not resolve the foundation model") {
		log.Printf(`Could not resolve the foundation model from model identifier: \"%v\".
                    Please verify that the requested model exists and is accessible
                    within the specified region.\n
                    `, modelId)
	} else {
		log.Printf("Couldn't invoke model: \"%v\". Here's why: %v\n", modelId, err)
	}
}
