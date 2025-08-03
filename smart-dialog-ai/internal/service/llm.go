package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

type SiliconFlowHandler struct {
	ctx    context.Context
	logger *logrus.Logger
	apiURL string
	apiKey string
	client *http.Client
}

// 创建新的 SiliconFlowHandler
func NewSiliconFlowHandler(apiURL, apiKey string) *SiliconFlowHandler {
	return &SiliconFlowHandler{
		ctx:    context.Background(),
		logger: logrus.New(),
		apiURL: apiURL,
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

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
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Tools    []Tool        `json:"tools"`
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

// 发送请求给大模型
func (s *SiliconFlowHandler) GenerateText(msg string) (string, error) {
	// 打印前端发来的消息
	log.Printf("web send message is :%s", msg)
	// 构造工具调用信息
	tools := []Tool{
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "DateMaster", // 工具名称
				Description: "一款多功能日期查询工具，它能提供用户指定日期的详细信息比如，星座，农历日期，生肖，岁次，黄历等",
				Parameters: map[string]interface{}{
					"type": "object", // 参数类型是 object
					"properties": map[string]interface{}{
						"data": map[string]interface{}{
							"type":        "string",              // 日期类型
							"description": "查询的日期，格式为YYYY-MM-DD", // 这个data参数的描述
						},
					},
					"required": []string{"data"}, // 必填参数
				},
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "jork",
				Description: "用来讲笑话的",
				Parameters: map[string]interface{}{
					"type": "object",
					"sort": map[string]interface{}{
						"type":        "int",
						"description": " 指定时间之前发布的，asc: 指定时间之后发布的",
					},
					"required": []string{"sort", "time"}, // 必填参数
				},
			},
		},
	}
	// 构造请求体
	requestBodyStruct := RequestBody{
		Model: "Qwen/QwQ-32B", // 使用的模型
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: msg, // 用户输入的问题
			},
			{
				Role:    "system",
				Content: "你是一位百科全书",
			},
		},
		Tools: tools, // 传入工具调用
	}

	// 将请求体转换为 JSON 格式
	jsonData, err := json.Marshal(requestBodyStruct)
	if err != nil {
		s.logger.Error("Error marshalling request body:", err)
		return "", fmt.Errorf("error marshalling request body: %v", err)
	}

	// 创建 HTTP 请求
	url := "https://api.siliconflow.cn/v1/chat/completions"        // 替换为你的 API 地址
	token := "sk-ogugyhoyqushnqplefczlsafysldjensioiucmqhwbbkcybs" // 替换为有效的 API Token

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resolvedIP := "47.103.87.49" // 替换为 curl 成功的 IP
	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				if addr == "api.siliconflow.cn:443" {
					addr = resolvedIP + ":443"
				}
				return (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext(ctx, network, addr)
			},
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// 解析响应
	var response ChatCompletionResponse
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}
	if len(response.Choices) <= 0 {
		return "", fmt.Errorf("llm response failed")
	} else {
		// 遍历 choices，获取详细信息
		for idx, choice := range response.Choices {
			log.Printf("== Choice %d ==", idx)
			log.Printf("推理内容: %s", choice.Message.ReasoningContent)
			log.Printf("结束原因: %s", choice.FinishReason)
			if len(choice.Message.ToolCalls) > 0 {
				for _, toolCall := range choice.Message.ToolCalls {
					log.Println("---- 工具调用 ----")
					log.Printf("toolCall.Function.Arguments:%s", toolCall.Function.Arguments)
					// function字段中既包含工具名也包含工具的参数
					// // 获取 function 字段，移除前后的空格
					function := toolCall.Function
					receviMessage, err := s.checkIfToolNeeded(function)
					if err != nil {
						log.Printf("调用工具失败:%v", err)
						break
					}

					// c.JSON(http.StatusOK, gin.H{
					// 	"response": receviMessage,
					// })
					return receviMessage, nil
				}

			} else {
				return "" + choice.Message.Content, nil
			}
		}
	}
	return "", fmt.Errorf("no response from the model")
}

func (s *SiliconFlowHandler) checkIfToolNeeded(functions Function) (string, error) {
	// get tool name
	name := functions.Name
	// get tool params
	arguments := functions.Arguments
	// 将参数序列化为json
	argumentsJSON, err := json.Marshal(arguments)
	if err != nil {
		return "", fmt.Errorf("error marshalling arguments: %v", err)
	}
	switch name {
	case "jork":
		result, err := s.tellJoke()
		fmt.Printf("Tool result:%s", string(result))
		return string(result), err
	case "DateMaster":
		argumentData := struct {
			Data string `json:"data"`
		}{}
		// get tool param
		if err := json.Unmarshal([]byte(argumentsJSON), &argumentData); err != nil {
			return "", fmt.Errorf("arguments become argumentData:%v", err)
		}
		return argumentData.Data, err
	default:
		return "", fmt.Errorf("UnKnown function:%s", name)
	}

}
func (s *SiliconFlowHandler) tellJoke() (string, error) {
	// 基本参数配置
	apiUrl := "http://v.juhe.cn/joke/content/list"
	apiKey := "ca9dee29ead19b5f7af99650f5a57299"
	// 接口请求入参配置
	requestParams := url.Values{}
	requestParams.Set("key", apiKey)
	requestParams.Set("sort", "desc")
	requestParams.Set("time", "1418816972")

	// 发起接口网络请求
	resp, err := http.Get(apiUrl + "?" + requestParams.Encode())
	if err != nil {
		return "", fmt.Errorf("网络请求异常:%v", err)
	}
	defer resp.Body.Close()

	var responseResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseResult)
	if err != nil {
		return "", fmt.Errorf("解析响应结果异常:%v", err)
	}
	// 用来判断是否还有限额
	reason := responseResult["reason"]
	// 判断是否限额
	if reason != "Success" {
		log.Fatalf("Error reason:%s", reason)
	}
	// 打印出reason
	log.Printf("reason:%s", reason)
	result, ok := responseResult["result"].(map[string]interface{})
	if !ok {
		log.Printf("Error result failed:")
	}
	data, ok := result["data"].([]interface{})
	if !ok {
		log.Fatalf("Error data failed ;")
	}
	//随机获取一个整数
	randomNum := rand.Intn(len(data))
	// 随机获取一个元素
	// .（map[string]string）是类型断言，告诉go这个数据的类型
	randomItem := data[randomNum].(map[string]interface{})
	content := randomItem["content"].(string)

	return content, nil
}
