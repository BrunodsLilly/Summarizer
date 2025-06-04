package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainFunction(t *testing.T) {
	if os.Getenv("BE_MAIN") == "1" {
		main()
		return
	}

	tests := []struct {
		name           string
		input          string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "empty input",
			input:          "\n",
			expectedOutput: "Paste a YouTube video URL to generate content based on it.",
			expectError:    true,
		},
		{
			name:           "invalid URL",
			input:          "invalid-url\n",
			expectedOutput: "Paste a YouTube video URL to generate content based on it.",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=TestMainFunction")
			cmd.Env = append(os.Environ(), "BE_MAIN=1")
			
			stdin, err := cmd.StdinPipe()
			if err != nil {
				t.Fatal(err)
			}

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}

			if _, err := io.WriteString(stdin, tt.input); err != nil {
				t.Fatal(err)
			}
			stdin.Close()

			err = cmd.Wait()
			
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !strings.Contains(stdout.String(), tt.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tt.expectedOutput, stdout.String())
			}
		})
	}
}

func TestMainWithMockInput(t *testing.T) {
	t.Run("test with subprocess", func(t *testing.T) {
		if os.Getenv("BE_MOCK_TEST") == "1" {
			main()
			return
		}

		cmd := exec.Command(os.Args[0], "-test.run=TestMainWithMockInput")
		cmd.Env = append(os.Environ(), "BE_MOCK_TEST=1")

		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}

		go func() {
			fmt.Fprintln(stdin, "https://youtube.com/watch?v=test")
			stdin.Close()
		}()

		err = cmd.Wait()

		output := stdout.String()
		if !strings.Contains(output, "Paste a YouTube video URL to generate content based on it.") {
			t.Errorf("expected prompt message, got: %s", output)
		}
	})
}