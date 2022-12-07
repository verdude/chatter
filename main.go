package main

import (
  "bytes"
  "bufio"
  "os"
  "flag"
  "io"
  "net/http"
  "encoding/json"

  "github.com/verdude/zapr"
  "go.uber.org/zap"
)

type UserMessage struct {
  Sender string `json:"sender"`
  Message string `json:"message"`
}

type BotMessage struct {
  RecipientID string `json:"recipient_id"`
  Text string `json:"text"`
}

func sendMessage(msg string) (string, error) {
  zapr.V(5).I("Sending Message...")
  userMessage := UserMessage{Sender: "santi", Message: msg}
  payload, err := json.Marshal(userMessage)
  if err != nil {
    zapr.V(3).E("Failed to encode json string", zap.Any("error", err))
    return "", err
  }

  zapr.V(5).I("Sending msg to server", zap.Any("payload", payload))
  response, err := http.Post("http://localhost:5005/webhooks/rest/webhook", "application/json", bytes.NewBuffer(payload))
  if err != nil {
    zapr.V(3).E("Failed to send message", zap.Any("error", err))
    return "", err
  }

  defer response.Body.Close()
  body, err := io.ReadAll(response.Body)
  if err != nil {
    zapr.V(3).E("Failed to send message", zap.Any("error", err))
    return "", err
  }

  zapr.V(8).I("Somesing", zap.Any("sing", body))
  var botMsg []BotMessage
  err = json.Unmarshal(body, &botMsg)
  if err != nil {
    zapr.V(3).E("Failed to parse server response", zap.Any("error", err))
    return "", err
  }

  for i := 0; i < len(botMsg); i++ {
    zapr.V(3).I("Received bot messages", zap.Any("text", botMsg[i].Text))
  }
  return "haha", nil
}

func getMessage() (string, error) {
  zapr.V(5).I("Getting Message...")
  reader := bufio.NewReader(os.Stdin)
  msg, err := reader.ReadString('\n')
  if err != nil || len(msg) <= 0 {
    return "", err
  }
  return msg, nil
}

func main() {
  vLevel := flag.Int("v", 3, "Verbosity level.")
  flag.Parse()

  zapr.Init(uint8(*vLevel))
  defer zapr.Sync()

  for {
    msg, err := getMessage()
    if err != nil {
      zapr.V(3).I("Failure", zap.Any("error", err))
      break
    }
    zapr.V(5).I("Got message:", zap.Any("message", msg))

    response, err := sendMessage(msg)
    if err != nil {
      zapr.V(1).E("Error:", zap.Any("error", err))
      break
    }
    zapr.V(5).I("Response:", zap.Any("response", response))
  }
}
