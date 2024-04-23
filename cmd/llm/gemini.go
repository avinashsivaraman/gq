package llm

import (
	"context"
	"encoding/json"
	"errors"
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

func (self *GeminiLLM) Chat(userQuery string) (string, error) {
	ctx := context.Background()
	apiKey := viper.GetString("api_key")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
		return "", errors.New("Gemini API Initialized failed")
	}

	defer client.Close()

	model := client.GenerativeModel("gemini-1.0-pro")

	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(512)

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
