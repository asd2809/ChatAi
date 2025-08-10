package model

import "time"

type ChatMessage struct {
	Role           string                 `json:"role"`                     // "user" / "assistant" / "system"
	Text           string                 `json:"text"`                     // 消息内容
	Timestamp      string                 `json:"timestamp"`                // 时间戳（建议 ISO8601 格式或 HH:MM:SS）
	MessageID      string                 `json:"message_id,omitempty"`     // 消息 ID（可选）
	ConversationID string                 `json:"conversation_id,omitempty"`// 会话 ID（可选）
	Extra          map[string]interface{} `json:"extra,omitempty"`          // 其他扩展信息（可选）
}
// type ChatMessage struct {
// 	Role    string `json:"role"`
// 	Content string `json:"content"`
// }
// ------以下是llm需要的结构体---------
type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type RequestBody struct {
	Model    string          `json:"model"`
	Messages []Message `json:"messages"`
	Tools    []Tool          `json:"tools"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"` // Arguments is a string containing JSON
	} `json:"function"`
}

// 响应体的结构体
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
		Message      struct {
			Role             string     `json:"role"`
			Content          string     `json:"content"`
			ReasoningContent string     `json:"reasoning_content,omitempty"`
			ToolCalls        []ToolCall `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
}

// 获取工具的姓名与参数的结构体
type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}
// ---------------


// ------------以下是聊天记录需的-------
// 聊天记录的内容
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}
// 存储到 MySQL 的结构
type ChatRecord struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"index"`
	Role      string    `gorm:"type:enum('user','assistant')"`
	Content   string
	CreatedAt time.Time
}
// ----------------------
// 用户表