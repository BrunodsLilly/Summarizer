package core

import (
	"fmt"

    gemini_api "github.com/BrunodsLilly/Summarizer/pkg/core/internal/utils"
)

const (
	APIVersion = "v1beta"
	ModelName  = "gemini-2.5-pro-preview-05-06"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

// func getModelName() string {
// 	if modelName := os.Getenv("MODEL_NAME"); modelName != "" {
// 		return modelName
// 	}
// 	return ModelName
// }

// func getAPIVersion() string {
// 	if apiVersion := os.Getenv("API_VERSION"); apiVersion != "" {
// 		return apiVersion
// 	}
// 	return APIVersion
// }

// func generateWithYTVideo(url string) (string, error) {
// 	modelName := getModelName()
// 	apiVersion := getAPIVersion()

// 	ctx := context.Background()
// 	client, err := genai.NewClient(ctx, &genai.ClientConfig{
// 		HTTPOptions: genai.HTTPOptions{APIVersion: apiVersion},
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create genai client: %w", err)
// 	}

// 	contents := []*genai.Content{
// 		{Parts: []*genai.Part{
// 			{Text: "Write a short summary of the video. Be as information dense as possible. Be thorough. Produce an overall summary, then clearly identify video sections and summarize them."},
// 			{FileData: &genai.FileData{
// 				FileURI:  url,
// 				MIMEType: "video/mp4",
// 			}},
// 		}},
// 	}

// 	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate content: %w", err)
// 	}

// 	respText := resp.Text()
// 	return respText, nil
// }

func SummarizeURL(url string) (string, error) {
	resp, err := gemini_api.GenerateWithYTVideo(url)
	if err != nil {
		fmt.Println("Error generating summary:", err)
		return "", err
	}
	return resp, nil
}

// SummarizeURLWithModel allows specifying a custom model
func SummarizeURLWithModel(url, modelName string) (string, error) {
	apiVersion := gemini_api.GetAPIVersion()
	resp, err := gemini_api.GenerateWithYTVideoAndModel(url, modelName, apiVersion)
	if err != nil {
		fmt.Println("Error generating summary:", err)
		return "", err
	}
	return resp, nil
}

// GetModelInfo returns information about the current model being used
func GetModelInfo() (string, string) {
	return gemini_api.GetModelName(), gemini_api.GetAPIVersion()
}
