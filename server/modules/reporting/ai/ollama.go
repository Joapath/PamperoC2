package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OllamaService handles communication with Ollama API
type OllamaService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// OllamaRequest represents the request structure for Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents the response structure from Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// AIAnalysisResult represents the structured AI analysis output
type AIAnalysisResult struct {
	AttackVectors []string `json:"attack_vectors"`
	Techniques    []string `json:"techniques"`
	Commands      []string `json:"commands"`
	Priority      string   `json:"priority"`
	Confidence    float64  `json:"confidence"`
}

// NewOllamaService creates a new Ollama service instance
func NewOllamaService(baseURL string) *OllamaService {
	return &OllamaService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AnalyzeFinding analyzes a security finding using Ollama AI
func (s *OllamaService) AnalyzeFinding(finding string, model string) (string, error) {
	prompt := s.buildAnalysisPrompt(finding)

	req := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", s.BaseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	return ollamaResp.Response, nil
}

// buildAnalysisPrompt creates a structured prompt for security analysis
func (s *OllamaService) buildAnalysisPrompt(finding string) string {
	return fmt.Sprintf(`Analiza el siguiente hallazgo de seguridad y proporciona un análisis estructurado en formato JSON:

Hallazgo: %s

Por favor, proporciona tu respuesta ÚNICAMENTE en formato JSON válido con la siguiente estructura:
{
  "attack_vectors": ["vector1", "vector2", ...],
  "techniques": ["técnica1", "técnica2", ...],
  "commands": ["comando1", "comando2", ...],
  "priority": "alta|media|baja",
  "confidence": 0.0-1.0
}

Instrucciones:
- attack_vectors: Lista de posibles vectores de ataque que podrían explotar este hallazgo
- techniques: Técnicas específicas de ataque o mitigación recomendadas
- commands: Comandos técnicos concretos para explotación o remediación
- priority: Nivel de prioridad basado en criticidad (alta/media/baja)
- confidence: Nivel de confianza en el análisis (0.0 a 1.0)

Responde SOLO con el JSON, sin texto adicional.`, finding)
}

// ParseAIAnalysis parses the AI response into structured data
func ParseAIAnalysis(aiResponse string) (*AIAnalysisResult, error) {
	// Clean the response to extract JSON
	jsonStr := strings.TrimSpace(aiResponse)

	// Find JSON boundaries (in case there's extra text)
	start := strings.Index(jsonStr, "{")
	end := strings.LastIndex(jsonStr, "}")

	if start == -1 || end == -1 || end <= start {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr = jsonStr[start : end+1]

	var result AIAnalysisResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("error parsing AI analysis JSON: %w", err)
	}

	return &result, nil
}
