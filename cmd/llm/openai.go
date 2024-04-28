package llm

import (
	"context"
    "log"
    "fmt"
    "errors"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

type OpenAIProvider struct{}

func (_ OpenAIProvider) Chat(userQuery string, verbose bool) (string, error) {
    
    openAIConfig := viper.Sub("openAI")

  	apiKey := openAIConfig.GetString("apiKey")
    temperature := openAIConfig.GetFloat64("temperature")
    modelName := openAIConfig.GetString("modelName")
    maxOutputTokens := openAIConfig.GetInt32("maxOutputTokens")
	client := openai.NewClient(apiKey)

    var model string

    switch modelName {
      case "gpt-3.5-turbo": 
        model = openai.GPT3Dot5Turbo
      case "gpt-4-turbo":
        model = openai.GPT4Turbo
      case "gpt-4":
        model = openai.GPT4
      default:
        log.Fatal("Unsupported Model. Make sure the model name parsed is correct")
        return "", errors.New("Unsupported Model name found")
  }   
    
    if verbose {
      fmt.Println("\033[33mModel Params:\033[0m")
      fmt.Println("\033[36mModel Name: ", modelName)
      fmt.Println("Temperature: ", temperature)
      fmt.Println("Max Output Tokens: ", maxOutputTokens)
      fmt.Println("\033[0m")
    }
  
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
            Temperature: float32(temperature),
            MaxTokens: int(maxOutputTokens),
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userQuery,
				},
			},
		},
	)

	if err != nil {
      log.Fatal(err)
      return "", errors.New("OpenAI Chat Completion Failed")
	}

	return resp.Choices[0].Message.Content, nil
}
