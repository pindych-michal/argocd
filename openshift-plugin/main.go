package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: openshift-plugin <command> [args]")
		fmt.Fprintln(os.Stderr, "Commands: generate, discover")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "generate":
		if err := generate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "discover":
		discover()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// generate reads all YAML files from ARGOCD_APP_SOURCE_PATH and outputs them to stdout
func generate() error {
	sourcePath := os.Getenv("ARGOCD_APP_SOURCE_PATH")
	if sourcePath == "" {
		sourcePath = "."
	}

	files, err := filepath.Glob(filepath.Join(sourcePath, "*.yaml"))
	if err != nil {
		return fmt.Errorf("failed to glob yaml files: %w", err)
	}

	ymlFiles, err := filepath.Glob(filepath.Join(sourcePath, "*.yml"))
	if err != nil {
		return fmt.Errorf("failed to glob yml files: %w", err)
	}
	files = append(files, ymlFiles...)

	for i, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file, err)
		}

		// Output the manifest
		fmt.Print(string(content))

		// Add separator between documents if not the last file
		if i < len(files)-1 && !strings.HasSuffix(string(content), "\n---\n") {
			if !strings.HasSuffix(string(content), "\n") {
				fmt.Println()
			}
			fmt.Println("---")
		}
	}

	return nil
}

// discover checks if this plugin should handle the repository
func discover() {
	sourcePath := os.Getenv("ARGOCD_APP_SOURCE_PATH")
	if sourcePath == "" {
		sourcePath = "."
	}

	// Check for OpenShift-specific marker file
	markerFile := filepath.Join(sourcePath, ".openshift-plugin")
	if _, err := os.Stat(markerFile); err == nil {
		fmt.Println(`{"find": {"glob": "**/manifest*.yaml"}}`)
		return
	}

	// Check if any YAML contains OpenShift resources
	files, _ := filepath.Glob(filepath.Join(sourcePath, "*.yaml"))
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		content, _ := io.ReadAll(f)
		f.Close()

		// Look for OpenShift-specific resource kinds
		if strings.Contains(string(content), "kind: Route") ||
			strings.Contains(string(content), "kind: DeploymentConfig") ||
			strings.Contains(string(content), "kind: BuildConfig") ||
			strings.Contains(string(content), "kind: ImageStream") {
			fmt.Println(`{"find": {"glob": "**/*.yaml"}}`)
			return
		}
	}

	// Not an OpenShift project
	os.Exit(1)
}
