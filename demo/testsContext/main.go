package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gorm.io/gorm"
)

const apiURL = "https://api.siliconflow.cn/v1/chat/completions"
const apiToken = "sk-ogugyhoyqushnqplefczlsafysldjensioiucmqhwbbkcybs" // 替换为你的真实 token

func main() {
	db := InitDB()
	userID := "user123" // 可动态生成或从登录用户获取

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("💬 AI 聊天助手已启动，输入 exit 退出")

	for {
		fmt.Print("🧑‍💻 你：")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		if userInput == "exit" {
			break
		}

		// 读取历史
		history := LoadHistory(db, userID, 10)

		// 加入本轮用户输入
		userMsg := Message{Role: "user", Content: userInput}
		history = append(history, userMsg)
		SaveMessage(db, userID, userMsg)

		// 请求大模型
		reqBody := ChatRequest{
			Model:    "Qwen/QwQ-32B",
			Messages: history,
		}

		jsonData, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+apiToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("❌ 请求失败:", err)
			continue
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var aiResp ChatResponse
		if err := json.Unmarshal(body, &aiResp); err != nil || len(aiResp.Choices) == 0 {
			fmt.Println("❌ 解析失败:", string(body))
			continue
		}

		// 打印回复
		aiReply := aiResp.Choices[0].Message.Content
		fmt.Println("🤖 AI：", aiReply)

		// 保存回复
		aiMsg := Message{Role: "assistant", Content: aiReply}
		SaveMessage(db, userID, aiMsg)
	}
}

func SaveMessage(db *gorm.DB, userID string, msg Message) {
	record := ChatRecord{
		UserID:  userID,
		Role:    msg.Role,
		Content: msg.Content,
	}
	
	if err := db.Create(&record).Error; err != nil {
		fmt.Println("❌ 保存失败:", err)
	}
}

func LoadHistory(db *gorm.DB, userID string, limit int) []Message {
	var records []ChatRecord
	var history []Message

	db.Where("user_id = ?", userID).Order("created_at asc").Limit(limit).Find(&records)

	for _, r := range records {
		history = append(history, Message{
			Role:    r.Role,
			Content: r.Content,
		})
	}
	return history
}
