package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)


type GeminiProvider struct{}

func (_ GeminiProvider) Chat(userQuery string, verbose bool) (string, error) {
	ctx := context.Background()

    geminiConfig := viper.Sub("gemini")

	apiKey := geminiConfig.GetString("apiKey")
    modelName := geminiConfig.GetString("modelName")
    temperature := float32(geminiConfig.GetFloat64("temperature"))
    maxOutputTokens := geminiConfig.GetInt32("maxOutputTokens")

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
		return "", errors.New("Gemini API Initialized failed")
	}

	defer client.Close()

    if verbose {
      fmt.Println("\033[33mModel Params:\033[0m")
      fmt.Println("\033[36mModel Name: ", modelName)
      fmt.Println("Temperature: ", temperature)
      fmt.Println("Max Output Tokens: ", maxOutputTokens)
      fmt.Println("\033[0m")
    }

	model := client.GenerativeModel(modelName)

	model.SetTemperature(temperature)
	model.SetMaxOutputTokens(maxOutputTokens)

	resp, err := model.GenerateContent(ctx, genai.Text(userQuery))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	content := resp.Candidates[0].Content

	outputResponse := ""
	if content != nil {
		r := content.Parts[0]
		rb, _ := json.MarshalIndent(r, "", "  ")
		outputResponse = string(rb)
	} else {
		outputResponse = "Failed to generate message. Try again"
	}

	return outputResponse, nil
}
