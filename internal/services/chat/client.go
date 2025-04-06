package chat

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   int32
	username string
}

type Message struct {
	Username string    `json:"username"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
}

func NewClient(
	hub *Hub, conn *websocket.Conn, uID int32, uname string) *Client {
	return &Client{
		Hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   uID,
		username: uname,
	}
}

func (c *Client) HandleMessages(chatService IChatService) {
	defer func() {
		c.Hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var rawMsg struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(data, &rawMsg); err != nil {
			continue
		}

		msg := Message{
			Username: c.username,
			Content:  rawMsg.Content,
			Time:     time.Now(),
		}

		if err := chatService.SaveMessage(context.Background(), c.userID, msg.Content); err != nil {
			log.Printf("error saving message: %v", err)
			continue
		}

		messageJSON, err := json.Marshal(msg)
		if err != nil {
			log.Printf("error marshaling message: %v", err)
			continue
		}
		c.Hub.broadcast <- messageJSON
	}
}

func (c *Client) SendMessages() {
	defer c.conn.Close()

	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}
