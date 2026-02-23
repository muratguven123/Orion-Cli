package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"orion-cli/internal/core"
)

const systemPrompt = "Sen bir Technical Product Manager'sın. Verilen metni incele ve bana sadece aşağıdaki JSON formatında bir görev listesi dön. Markdown kullanma, sadece saf JSON dön."

type GeminiAdapter struct {
	apiKey string
	model  string
	client *http.Client
}

func NewGeminiAdapter(apiKey, model string, client *http.Client) *GeminiAdapter {
	if client == nil {
		client = http.DefaultClient
	}
	return &GeminiAdapter{apiKey: apiKey, model: model, client: client}
}

func (g *GeminiAdapter) Analyze(text string) ([]core.IssueDraft, error) {
	return g.AnalyzeWithContext(context.Background(), text)
}

func (g *GeminiAdapter) AnalyzeWithContext(ctx context.Context, text string) ([]core.IssueDraft, error) {
	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", g.model)

	payload := map[string]any{
		"system_instruction": map[string]any{
			"parts": []map[string]string{{"text": systemPrompt}},
		},
		"contents": []map[string]any{{
			"role":  "user",
			"parts": []map[string]string{{"text": text}},
		}},
		"generationConfig": map[string]any{
			"temperature":      0.2,
			"responseMimeType": "application/json",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("gemini isteği oluşturulamadı: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("gemini isteği hazırlanamadı: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", g.apiKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini API çağrısı başarısız: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gemini cevabı okunamadı: %w", err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("gemini API hata döndü (%d): %s", resp.StatusCode, string(respBody))
	}

	var parsed struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("gemini cevabı parse edilemedi: %w", err)
	}

	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini boş cevap döndü")
	}

	jsonText := sanitizeJSON(parsed.Candidates[0].Content.Parts[0].Text)

	var drafts []core.IssueDraft
	if err := json.Unmarshal([]byte(jsonText), &drafts); err != nil {
		return nil, fmt.Errorf("gemini çıktısı JSON değil: %w, içerik: %s", err, jsonText)
	}

	if len(drafts) == 0 {
		return nil, fmt.Errorf("gemini görev listesi üretmedi")
	}

	return drafts, nil
}

func sanitizeJSON(raw string) string {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.TrimPrefix(trimmed, "```json")
	trimmed = strings.TrimPrefix(trimmed, "```")
	trimmed = strings.TrimSuffix(trimmed, "```")
	return strings.TrimSpace(trimmed)
}
