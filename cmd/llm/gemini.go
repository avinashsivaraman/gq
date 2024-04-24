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

type GeminiLLM struct {
	chatMessage string
}

func NewGeminiLLM() *GeminiLLM {
	return &GeminiLLM{}
}

func (self *GeminiLLM) Chat(userQuery string, verbose bool) (string, error) {
	ctx := context.Background()
	apiKey := viper.GetString("api_key")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
		return "", errors.New("Gemini API Initialized failed")
	}

	defer client.Close()

    modelName := "gemini-1.0-pro"
    var temperature float32 = 0.7
    var maxOutputTokens int32 = 512

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
