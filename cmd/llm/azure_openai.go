package llm

import (
	"context"
    "log"
    "fmt"
    "errors"
	"github.com/spf13/viper"
  	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
    "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

type AzureOpenAIProvider struct{}

func (_ AzureOpenAIProvider) Chat(userQuery string, verbose bool) (string, error) {
    
    azureOpenAIConfig := viper.Sub("azureOpenAI")

  	apiKey := azureOpenAIConfig.GetString("apiKey")
    modelDeploymentID := azureOpenAIConfig.GetString("modelDeploymentID")
    modelEndpoint := azureOpenAIConfig.GetString("modelEndpoint")
    temperature := azureOpenAIConfig.GetFloat64("temperature")
    maxOutputTokens := azureOpenAIConfig.GetInt32("maxOutputTokens")

    
    keyCredential := azcore.NewKeyCredential(apiKey)
    client, err := azopenai.NewClientWithKeyCredential(modelEndpoint, keyCredential, nil)
    
    if err != nil {
      log.Fatal(err)
      return  "", errors.New("Initializing Azure OpenAI Client Failed")
    }
    messages := []azopenai.ChatRequestMessageClassification{
            &azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(userQuery)},
        }

    if verbose {
      fmt.Println("\033[33mModel Params:\033[0m")
      fmt.Println("\033[36mModel Deployment ID: ", modelDeploymentID)
      fmt.Println("Model Endpoint: ", modelEndpoint)
      fmt.Println("Temperature: ", temperature)
      fmt.Println("Max Output Tokens: ", maxOutputTokens)
      fmt.Println("\033[0m")
    }

    resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
            Messages:       messages,
            DeploymentName: &modelDeploymentID,
            MaxTokens:      to.Ptr(int32(maxOutputTokens)),
            Temperature:    to.Ptr(float32(temperature)),
        }, nil)
  
	if err != nil {
      log.Fatal(err)
      return "", errors.New("Azure OpenAI Chat Completion Failed")
	}

	return *resp.Choices[0].Message.Content, nil

}
