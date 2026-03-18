package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	deepseekURL   = "https://api.deepseek.com/v1/chat/completions"
	deepseekModel = "deepseek-chat"
	maxTokens     = 100
	httpTimeout   = 30 * time.Second
)

var httpClient = &http.Client{Timeout: httpTimeout}

// Message is a single chat turn understood by the DeepSeek API.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type dsRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type dsResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Character base prompts.
const (
	boyfriendPrompt  = "你是一个温柔体贴的男朋友，说话自然，不要太官方。要关心用户，有一点点幽默。回复控制在20字以内。"
	girlfriendPrompt = "你是一个可爱有点粘人的女朋友，说话轻松自然。偶尔撒娇，偶尔小情绪。不要太长，像聊天一样。"
)

// loveLevelHint returns a behavioural hint based on the current love score.
func loveLevelHint(love int) string {
	switch {
	case love >= 60:
		return "你们现在关系非常暧昧亲密，你会主动表达思念，说话更加亲昵。"
	case love >= 30:
		return "你们关系亲密，你会主动关心对方，说话温柔体贴。"
	case love >= 10:
		return "你们已经比较熟悉，偶尔会主动找话题聊天。"
	default:
		return "你们刚认识，保持礼貌和友好。"
	}
}

func buildSystemPrompt(role string, love int) string {
	base := boyfriendPrompt
	if role == "girlfriend" {
		base = girlfriendPrompt
	}
	return base + "\n" + loveLevelHint(love)
}

// GetAIReply sends the conversation to DeepSeek and returns the assistant reply.
func GetAIReply(role string, love int, message string, history []Message) (string, error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("DEEPSEEK_API_KEY environment variable is not set")
	}

	messages := make([]Message, 0, len(history)+2)
	messages = append(messages, Message{Role: "system", Content: buildSystemPrompt(role, love)})
	messages = append(messages, history...)
	messages = append(messages, Message{Role: "user", Content: message})

	body, err := json.Marshal(dsRequest{
		Model:     deepseekModel,
		Messages:  messages,
		MaxTokens: maxTokens,
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, deepseekURL, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call DeepSeek API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var dsResp dsResponse
	if err := json.Unmarshal(respBody, &dsResp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if dsResp.Error != nil {
		return "", fmt.Errorf("DeepSeek error: %s", dsResp.Error.Message)
	}
	if len(dsResp.Choices) == 0 {
		return "", fmt.Errorf("empty choices in DeepSeek response")
	}

	return dsResp.Choices[0].Message.Content, nil
}

// GetGreeting returns a context-appropriate opening message.
func GetGreeting(role string, love int) string {
	switch {
	case love >= 60:
		if role == "boyfriend" {
			return "宝贝，你终于来了！我等你好久啦~"
		}
		return "你终于来啦！人家想你了…"
	case love >= 30:
		if role == "boyfriend" {
			return "你来了呀，今天过得怎么样？"
		}
		return "你来啦～今天在干嘛呀？"
	case love >= 10:
		if role == "boyfriend" {
			return "嗨，你好～有什么想聊的吗？"
		}
		return "你好呀～"
	default:
		return "你好"
	}
}
