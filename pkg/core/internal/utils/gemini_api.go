package gemini_api

import (
	"context"
	"fmt"
	"os"

	genai "google.golang.org/genai"
)

// constants
const (
	// API version
	APIVersion = "v1beta"
	// Model name
	// ModelName = "gemini-2.0-flash-001"
	ModelName = "gemini-2.5-pro-preview-05-06"
	// MaxOutputTokens
	MaxOutputTokens = 128_000
)

func main() {
	fmt.Println("Paste a YouTube video URL to generate content based on it.")
	var url string
	fmt.Scanln(&url)
	resp, err := GenerateWithYTVideo(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating content: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated content:")
	fmt.Println(resp)
}

// get model name from environment variables or use defaults
func GetModelName() string {
	if modelName := os.Getenv("MODEL_NAME"); modelName != "" {
		return modelName
	}
	return ModelName
}

// get API version from environment variables or use defaults
func GetAPIVersion() string {
	if apiVersion := os.Getenv("API_VERSION"); apiVersion != "" {
		return apiVersion
	}
	return APIVersion
}

// GenerateWithYTVideo shows how to generate text using a YouTube video as input.
func GenerateWithYTVideo(url string) (string, error) {
	modelName := GetModelName()
	apiVersion := GetAPIVersion()
	return GenerateWithYTVideoAndModel(url, modelName, apiVersion)
}

// GenerateWithYTVideoAndModel allows specifying a custom model
func GenerateWithYTVideoAndModel(url, modelName, apiVersion string) (string, error) {

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: apiVersion},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create genai client: %w", err)
	}

	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "Write a short summary of the video using Markdown. Be as information dense as possible. Be thorough. Use bullet lists to break down complex ideas. Provide space between sections. Produce an overall summary, list key sections to listen to, then add a thoughtful critique of the video. Then include a 'Further Reading' section that connects ideas, expands on them, and provide further information with links."},
			{FileData: &genai.FileData{
				FileURI:  url,
				MIMEType: "video/mp4",
			}},
		}},
	}

	config := genai.GenerateContentConfig{
		MaxOutputTokens: MaxOutputTokens,
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: fmt.Sprintf("Keep your answer below %d tokens.", MaxOutputTokens)},
			},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, &config)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

	return respText, nil
}
