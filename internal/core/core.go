import (
	"context"
	"fmt"
	"io"
	"os"

	genai "google.golang.org/genai"
)

// constants
const (
	// API version
	APIVersion = "v1beta"
	// Model name
    // ModelName = "gemini-2.5-pro-preview-05-06"
	ModelName = "gemini-2.0-flash-001"
)

func main() {
	fmt.Println("Paste a YouTube video URL to generate content based on it.")
	var url string
	fmt.Scanln(&url)
	if err := generateWithYTVideo(os.Stdout, url); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating content: %v\n", err)
		os.Exit(1)
	}
}

// get model name from environment variables or use defaults
func getModelName() string {
	if modelName := os.Getenv("MODEL_NAME"); modelName != "" {
		return modelName
	}
	return ModelName
}

// get API version from environment variables or use defaults
func getAPIVersion() string {
	if apiVersion := os.Getenv("API_VERSION"); apiVersion != "" {
		return apiVersion
	}
	return APIVersion
}

// generateWithYTVideo shows how to generate text using a YouTube video as input.
func generateWithYTVideo(w io.Writer, url string) (string, error) {
	modelName := getModelName()
	apiVersion := getAPIVersion()

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: apiVersion},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}


	contents := []*genai.Content{
		{Parts: []*genai.Part{
			{Text: "Write a short summary of the video. Be as information dense as possible. Be thorough. Produce an overall summary, then clearly identify video sections and summarize them. Add space between sections and use bullet points for key points."},
			{FileData: &genai.FileData{
				FileURI:  url,
				MIMEType: "video/mp4",
			}},
		}},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()

    return respText, nil
}
