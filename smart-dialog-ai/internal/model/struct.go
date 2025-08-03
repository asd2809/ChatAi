package model

type ChatMessage struct {
	Role           string                 `json:"role"`                     // "user" / "assistant" / "system"
	Text           string                 `json:"text"`                     // 消息内容
	Timestamp      string                 `json:"timestamp"`                // 时间戳（建议 ISO8601 格式或 HH:MM:SS）
	MessageID      string                 `json:"message_id,omitempty"`     // 消息 ID（可选）
	ConversationID string                 `json:"conversation_id,omitempty"`// 会话 ID（可选）
	Extra          map[string]interface{} `json:"extra,omitempty"`          // 其他扩展信息（可选）
}
