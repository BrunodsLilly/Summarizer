package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestCLIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	t.Run("full workflow with mock input", func(t *testing.T) {
		if os.Getenv("BE_INTEGRATION_TEST") == "1" {
			main()
			return
		}

		cmd := exec.Command(os.Args[0], "-test.run=TestCLIIntegration")
		cmd.Env = append(os.Environ(), "BE_INTEGRATION_TEST=1")

		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}

		go func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Fprintln(stdin, "https://youtube.com/watch?v=test")
			stdin.Close()
		}()

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case err := <-done:
			output := stdout.String()
			if !strings.Contains(output, "Paste a YouTube video URL to generate content based on it.") {
				t.Errorf("expected prompt message in output, got: %s", output)
			}

			if err != nil {
				stderrOutput := stderr.String()
				if !strings.Contains(stderrOutput, "Error generating content:") {
					t.Errorf("unexpected error without proper error message: %v, stderr: %s", err, stderrOutput)
				}
			}
		case <-ctx.Done():
			cmd.Process.Kill()
			t.Fatal("test timed out")
		}
	})
}

func TestCLIBuildAndRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping build test in short mode")
	}

	tempDir := t.TempDir()
	binaryPath := fmt.Sprintf("%s/cli_test", tempDir)

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "/Users/brunodossantos/Code/Summarizer/cmd/cli"
	
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build CLI: %v", err)
	}

	t.Run("binary execution", func(t *testing.T) {
		cmd := exec.Command(binaryPath)
		
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}

		go func() {
			time.Sleep(100 * time.Millisecond)
			io.WriteString(stdin, "invalid-url\n")
			stdin.Close()
		}()

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case err := <-done:
			output := stdout.String()
			if !strings.Contains(output, "Paste a YouTube video URL to generate content based on it.") {
				t.Errorf("expected prompt message, got: %s", output)
			}

			if err == nil {
				t.Error("expected error for invalid URL but got none")
			}
		case <-ctx.Done():
			cmd.Process.Kill()
			t.Fatal("test timed out")
		}
	})
}