package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "sync"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // 允许所有来源
    },
}

type Message struct {
    Content string `json:"content"`
    From    string `json:"from"`
    Time    string `json:"time"`
}

type Server struct {
    messages []Message
    mu       sync.Mutex
    clients  map[*websocket.Conn]bool
}

func (s *Server) handleWebSocket(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    s.mu.Lock()
    s.clients[conn] = true
    s.mu.Unlock()

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            break
        }

        // 处理用户消息
        msg := Message{
            Content: string(message),
            From:    "user",
            Time:    fmt.Sprintf("%02d:%02d", 12, 30),
        }
        s.mu.Lock()
        s.messages = append(s.messages, msg)
        s.mu.Unlock()

        // 模拟机器人的回复
        botMsg := Message{
            Content: "This is a response from the bot.",
            From:    "bot",
            Time:    fmt.Sprintf("%02d:%02d", 12, 32),
        }
        s.mu.Lock()
        s.messages = append(s.messages, botMsg)
        s.mu.Unlock()

        // 向所有客户端发送机器人消息
        for client := range s.clients {
            err := client.WriteMessage(websocket.TextMessage, []byte(botMsg.Content))
            if err != nil {
                log.Println(err)
            }
        }
    }
}

func (s *Server) getMessages(c *gin.Context) {
    c.Header("Content-Type", "application/json")
    s.mu.Lock()
    defer s.mu.Unlock()
    err := json.NewEncoder(c.Writer).Encode(struct {
        Messages []Message `json:"messages"`
    }{
        Messages: s.messages,
    })
    if err != nil {
        log.Println("Error sending messages:", err)
    }
}

func main() {
    server := &Server{
        clients: make(map[*websocket.Conn]bool),
    }

    r := gin.Default()

    // CORS 处理中间件
    r.Use(cors.Default()) // 允许所有来源，您也可以根据需要配置来源

    r.GET("/get-messages", server.getMessages)
    r.GET("/ws", server.handleWebSocket)

    log.Println("Server started on :8083")
    if err := r.Run(":8083"); err != nil {
        log.Fatal("ListenAndServe failed:", err)
    }
}
