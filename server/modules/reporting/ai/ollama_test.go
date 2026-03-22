package ai

import (
	"strings"
	"testing"
)

func TestParseAIAnalysis(t *testing.T) {
	// Test valid JSON response
	validResponse := `{
		"attack_vectors": ["SQL Injection", "XSS"],
		"techniques": ["Input validation bypass", "Parameter tampering"],
		"commands": ["sqlmap -u 'http://target.com/search?q=test'", "nmap -sV target.com"],
		"priority": "alta",
		"confidence": 0.85
	}`

	result, err := ParseAIAnalysis(validResponse)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result.AttackVectors) != 2 {
		t.Errorf("Expected 2 attack vectors, got %d", len(result.AttackVectors))
	}

	if result.Priority != "alta" {
		t.Errorf("Expected priority 'alta', got '%s'", result.Priority)
	}

	if result.Confidence != 0.85 {
		t.Errorf("Expected confidence 0.85, got %f", result.Confidence)
	}
}

func TestParseAIAnalysis_InvalidJSON(t *testing.T) {
	invalidResponse := "This is not JSON at all"

	_, err := ParseAIAnalysis(invalidResponse)
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

func TestBuildAnalysisPrompt(t *testing.T) {
	service := NewOllamaService("http://localhost:11434")
	finding := "SQL injection vulnerability in login form"

	prompt := service.buildAnalysisPrompt(finding)

	if !strings.Contains(prompt, finding) {
		t.Errorf("Prompt should contain the finding text")
	}

	if !strings.Contains(prompt, "attack_vectors") {
		t.Errorf("Prompt should contain attack_vectors instruction")
	}

	if !strings.Contains(prompt, "JSON") {
		t.Errorf("Prompt should mention JSON format")
	}
}

func TestNewOllamaService(t *testing.T) {
	baseURL := "http://localhost:11434"
	service := NewOllamaService(baseURL)

	if service.BaseURL != baseURL {
		t.Errorf("Expected BaseURL %s, got %s", baseURL, service.BaseURL)
	}

	if service.HTTPClient == nil {
		t.Error("HTTPClient should not be nil")
	}

	if service.HTTPClient.Timeout != 30*1000000000 { // 30 seconds in nanoseconds
		t.Errorf("Expected timeout 30s, got %v", service.HTTPClient.Timeout)
	}
}
