package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/BrunodsLilly/Summarizer/pkg/core"
)

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
		expectExit1   bool
	}{
		{
			name:          "empty URL",
			input:         "\n",
			expectedError: "Error generating content:",
			expectExit1:   true,
		},
		{
			name:          "malformed URL",
			input:         "not-a-url\n",
			expectedError: "Error generating content:",
			expectExit1:   true,
		},
		{
			name:          "invalid youtube URL",
			input:         "https://example.com\n",
			expectedError: "Error generating content:",
			expectExit1:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if os.Getenv("BE_ERROR_TEST") == "1" {
				main()
				return
			}

			cmd := exec.Command(os.Args[0], "-test.run=TestErrorHandling")
			cmd.Env = append(os.Environ(), "BE_ERROR_TEST=1")

			stdin, err := cmd.StdinPipe()
			if err != nil {
				t.Fatal(err)
			}

			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}

			stdin.Write([]byte(tt.input))
			stdin.Close()

			err = cmd.Wait()
			
			if tt.expectExit1 {
				if exitError, ok := err.(*exec.ExitError); ok {
					if exitError.ExitCode() != 1 {
						t.Errorf("expected exit code 1, got %d", exitError.ExitCode())
					}
				} else {
					t.Errorf("expected exit error, got %v", err)
				}
			}

			stderrOutput := stderr.String()
			if !strings.Contains(stderrOutput, tt.expectedError) {
				t.Errorf("expected stderr to contain %q, got %q", tt.expectedError, stderrOutput)
			}
		})
	}
}

func TestCoreIntegrationErrors(t *testing.T) {
	t.Run("core.SummarizeURL error propagation", func(t *testing.T) {
		testURL := "https://youtube.com/watch?v=invalid"
		
		if os.Getenv("BE_CORE_ERROR_TEST") == "1" {
			var url string = testURL
			resp, err := core.SummarizeURL(url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating content: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Generated content:")
			fmt.Println(resp)
			return
		}

		cmd := exec.Command(os.Args[0], "-test.run=TestCoreIntegrationErrors")
		cmd.Env = append(os.Environ(), "BE_CORE_ERROR_TEST=1")

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err := cmd.Run()
		
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() != 1 {
				t.Errorf("expected exit code 1, got %d", exitError.ExitCode())
			}
		} else {
			t.Errorf("expected exit error, got %v", err)
		}

		stderrOutput := stderr.String()
		if !strings.Contains(stderrOutput, "Error generating content:") {
			t.Errorf("expected stderr to contain error message, got %q", stderrOutput)
		}
	})
}