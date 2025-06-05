package main

import (
	"bytes"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BrunodsLilly/Summarizer/cmd/web/templates"
	"github.com/BrunodsLilly/Summarizer/pkg/core"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func init() {
	// Register MIME types for JavaScript
	mime.AddExtensionType(".js", "application/javascript")
}

func main() {
	log.Println("Starting web server...")
	
	// Get the directory where the executable is located
	execDir, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get executable directory:", err)
	}
	webDir := filepath.Dir(execDir)
	
	// For 'go run', determine the web directory based on where we're running from
	if strings.Contains(execDir, "go-build") {
		wd, _ := os.Getwd()
		// Check if we're running from the parent directory (Summarizer)
		if filepath.Base(wd) == "Summarizer" {
			webDir = filepath.Join(wd, "cmd", "web")
		} else {
			// Assume we're running from cmd/web directory
			webDir = wd
		}
	}
	
	staticDir := filepath.Join(webDir, "static")
	log.Printf("Web directory: %s", webDir)
	log.Printf("Static directory: %s", staticDir)
	
	// Check if static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("Warning: Static directory does not exist: %s", staticDir)
		log.Printf("Trying current directory as fallback...")
		staticDir = "./static"
		if _, err := os.Stat(staticDir); os.IsNotExist(err) {
			log.Printf("Warning: Fallback static directory also not found: %s", staticDir)
		}
	}
	
	// Create custom file server with proper MIME types
	http.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the file path
		filePath := r.URL.Path
		log.Printf("Serving static file: %s", filePath)
		
		// Determine MIME type from file extension
		ext := filepath.Ext(filePath)
		mimeType := mime.TypeByExtension(ext)
		
		if mimeType != "" {
			w.Header().Set("Content-Type", mimeType)
			log.Printf("Set MIME type for %s file %s: %s", ext, filePath, mimeType)
		} else if strings.HasSuffix(filePath, ".js") {
			// Fallback for JavaScript files
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			log.Printf("Fallback MIME type for JS file: %s", filePath)
		}
		
		// Serve the file using absolute path
		fullPath := filepath.Join(staticDir, filePath)
		log.Printf("Serving file from: %s", fullPath)
		
		// Check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			log.Printf("File not found: %s", fullPath)
			http.NotFound(w, r)
			return
		}
		
		http.ServeFile(w, r, fullPath)
	})))
	
	// Application routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/summarize", summarizeHandler)
	http.HandleFunc("/test-summary", testSummaryHandler)
	http.HandleFunc("/health", healthHandler)

	// Get port from environment variable (for Cloud Run) or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("Static files served from: %s\n", staticDir)
	fmt.Printf("Test summary available at http://localhost:%s/test-summary\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	component := templates.Index()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template rendering error: %v", err)
		return
	}
}

func summarizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	url := r.FormValue("url")
	if strings.TrimSpace(url) == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	selectedModel := r.FormValue("model")
	if strings.TrimSpace(selectedModel) == "" {
		selectedModel = "gemini-2.5-pro-preview-05-06" // Default model
	}

	// Generate summary with selected model
	summary := generateSummaryWithModel(url, selectedModel)

	component := templates.SummaryResult(summary)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template rendering error: %v", err)
		return
	}
}

func generateSummary(url string) string {
	res, err := core.SummarizeURL(url)
	if err != nil {
		log.Printf("Error summarizing URL: %v", err)
		// return "Error generating summary"
		return fmt.Sprintf("Error generating summary for URL: %s\n%s", url, err.Error())
	}
	return markdownToHTML(res)
}

func generateSummaryWithModel(url, modelName string) string {
	res, err := core.SummarizeURLWithModel(url, modelName)
	if err != nil {
		log.Printf("Error summarizing URL with model %s: %v", modelName, err)
		return fmt.Sprintf("Error generating summary for URL: %s using model %s\n%s", url, modelName, err.Error())
	}
	return markdownToHTML(res)
}

func markdownToHTML(markdown string) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		log.Printf("Error converting markdown to HTML: %v", err)
		return markdown // fallback to original text
	}
	
	return buf.String()
}

func testSummaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Sample summary content for testing
	sampleSummary := `# The Future of Artificial Intelligence: A Comprehensive Overview

## Introduction

Artificial Intelligence (AI) has emerged as one of the most transformative technologies of the 21st century, revolutionizing industries from healthcare to transportation. This comprehensive analysis explores the current state of AI, its applications, and the profound implications for society.

## Current State of AI Technology

### Machine Learning Breakthroughs

The field has witnessed remarkable progress in **machine learning algorithms**, particularly in:

- **Deep Learning Networks**: Neural networks with multiple layers that can process complex data patterns
- **Natural Language Processing**: Systems that understand and generate human language with unprecedented accuracy
- **Computer Vision**: AI systems that can interpret and analyze visual information better than ever before
- **Reinforcement Learning**: Algorithms that learn through trial and error, achieving superhuman performance in games and simulations

### Real-World Applications

AI is already transforming numerous sectors:

1. **Healthcare**: AI-powered diagnostic tools can detect diseases earlier and more accurately than traditional methods
2. **Finance**: Algorithmic trading and fraud detection systems process millions of transactions in real-time
3. **Transportation**: Autonomous vehicles are being tested on roads worldwide, promising safer and more efficient travel
4. **Education**: Personalized learning platforms adapt to individual student needs and learning styles

## Challenges and Ethical Considerations

### Technical Limitations

Despite remarkable progress, AI systems face several significant challenges:

- **Data Quality**: AI models are only as good as the data they're trained on
- **Interpretability**: Many AI systems operate as "black boxes," making it difficult to understand their decision-making process
- **Generalization**: AI systems often struggle to apply knowledge from one domain to another
- **Computational Requirements**: Training state-of-the-art AI models requires enormous computational resources

### Ethical Implications

The rapid advancement of AI raises important ethical questions:

> "With great power comes great responsibility. As AI becomes more powerful, we must ensure it's developed and deployed in ways that benefit all of humanity."

Key ethical considerations include:

- **Bias and Fairness**: Ensuring AI systems don't perpetuate or amplify existing societal biases
- **Privacy**: Protecting individual privacy in an age of ubiquitous data collection
- **Job Displacement**: Addressing the potential for AI to automate human jobs
- **Autonomy**: Maintaining human agency and control over important decisions

## Future Directions

### Emerging Technologies

Several promising areas of AI research are poised to drive the next wave of innovation:

- **Quantum Machine Learning**: Leveraging quantum computers to solve complex AI problems
- **Neuromorphic Computing**: Computer architectures inspired by the human brain
- **Federated Learning**: Training AI models across distributed data sources while preserving privacy
- **Explainable AI**: Developing AI systems that can provide clear explanations for their decisions

### Societal Integration

As AI becomes more prevalent, society must adapt:

- **Education Reform**: Preparing the workforce for an AI-driven economy
- **Regulatory Frameworks**: Developing appropriate governance structures for AI technologies
- **International Cooperation**: Fostering global collaboration on AI safety and ethics
- **Public Engagement**: Ensuring public understanding and participation in AI development

## Economic Impact

The economic implications of AI are profound and far-reaching:

### Market Growth

The global AI market is experiencing explosive growth:

- Current market size: $150 billion (2023)
- Projected market size: $1.5 trillion by 2030
- Annual growth rate: 25-30%

### Industry Transformation

AI is reshaping entire industries:

1. **Manufacturing**: Smart factories with predictive maintenance and quality control
2. **Retail**: Personalized shopping experiences and supply chain optimization
3. **Media**: Automated content creation and personalized recommendations
4. **Agriculture**: Precision farming techniques that optimize crop yields

## Conclusion

Artificial Intelligence represents both an unprecedented opportunity and a significant challenge for humanity. As we stand at the threshold of an AI-driven future, it's crucial that we approach this technology thoughtfully, ensuring that its benefits are widely shared while mitigating potential risks.

The path forward requires collaboration between technologists, policymakers, ethicists, and society at large. By working together, we can harness the power of AI to solve some of humanity's greatest challenges while preserving the values and principles that define our civilization.

The future of AI is not predetermined—it's a future we're actively creating through the choices we make today. Let's ensure it's a future that benefits everyone.

---

*This summary provides a comprehensive overview of artificial intelligence, covering technical aspects, applications, challenges, and future directions. The content is designed to be informative yet accessible to a broad audience.*`

	// Convert markdown to HTML
	htmlSummary := markdownToHTML(sampleSummary)

	// Render the test page with full layout
	component := templates.TestSummaryPage(htmlSummary)
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template rendering error: %v", err)
		return
	}
}


func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
