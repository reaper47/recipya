package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"sync"
	"time"
)

// Broker represents a client connection and its associated information.
type Broker struct {
	brokers      map[int64]*Broker
	conn         *websocket.Conn
	mu           sync.Mutex
	notification *template.Template
	userID       int64
}

// Progress represents a current value and the total number of values.
// It is used to calculate progress as a percentage.
type Progress struct {
	Value int
	Total int
}

// Message represents the data format for file and progress updates sent to the client.
type Message struct {
	Type     string `json:"type"`     // Message type, e.g. file.
	FileName string `json:"fileName"` // File name (applicable for "file" type).
	Data     string `json:"data"`     // Message data to pass. Base64-encoded if type is "file".
}

type wsTemplateData struct {
	ContentHTML      template.HTML
	IsToastWSVisible bool
	Title            string
}

type toast struct {
	Message    string `json:"message"`
	Background string `json:"background"`
}

// NewBroker creates a new Broker instance for a specific user and adds it to the brokers map.
// The userID is used for identification and cleanup purposes.
func NewBroker(userID int64, brokers map[int64]*Broker, conn *websocket.Conn, notification *template.Template) *Broker {
	b := &Broker{
		brokers:      brokers,
		conn:         conn,
		notification: notification,
		userID:       userID,
	}
	go b.ping()
	return b
}

func (b *Broker) HideNotification() {
	_ = b.SendProgressStatus("", false, -1, -1)
}

// SendFile sends a file to the connected client.
func (b *Broker) SendFile(fileName string, data *bytes.Buffer) error {
	if b == nil {
		return errors.New("ws connection nil")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	return b.conn.WriteJSON(Message{
		Type:     "file",
		FileName: fileName,
		Data:     base64.StdEncoding.EncodeToString(data.Bytes()),
	})
}

// SendProgress sends a progress update with a title and value to the client.
// The isToastVisible parameter controls whether the progress bar is displayed in a toast notification.
func (b *Broker) SendProgress(title string, value, numValues int) error {
	if b == nil {
		return errors.New("ws connection nil")
	}

	var buf bytes.Buffer
	_ = b.notification.ExecuteTemplate(&buf, "toast-ws", wsTemplateData{
		IsToastWSVisible: true,
		ContentHTML:      template.HTML(fmt.Sprintf(`<div id="export-progress"><progress max="100" value="%f"></progress></div>`, float64(value)/float64(numValues)*100)),
		Title:            title,
	})

	b.mu.Lock()
	defer b.mu.Unlock()

	return b.conn.WriteMessage(websocket.TextMessage, buf.Bytes())
}

// SendProgressStatus sends a progress update with a title and value, optionally hiding the toast notification.
func (b *Broker) SendProgressStatus(title string, isToastVisible bool, value, numValues int) error {
	if b == nil {
		return errors.New("ws connection nil")
	}

	var buf bytes.Buffer
	_ = b.notification.ExecuteTemplate(&buf, "toast-ws", wsTemplateData{
		IsToastWSVisible: isToastVisible,
		ContentHTML:      template.HTML(fmt.Sprintf(`<div id="export-progress"><progress max="100" value="%f"></progress></div>`, float64(value)/float64(numValues)*100)),
		Title:            title,
	})

	b.mu.Lock()
	defer b.mu.Unlock()

	return b.conn.WriteMessage(websocket.TextMessage, buf.Bytes())
}

func (b *Broker) ping() {
	defer func() {
		delete(b.brokers, b.userID)
		_ = b.conn.Close()
	}()

	b.setPingPongHandlers()

	for {
		_, _, err := b.conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func (b *Broker) setPingPongHandlers() {
	b.conn.SetPingHandler(func(message string) error {
		return b.conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
	})

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			b.mu.Lock()
			err := b.conn.WriteMessage(websocket.PingMessage, []byte{})
			b.mu.Unlock()
			if err != nil {
				return
			}
		}
	}()
}

// SendToast sends a toast notification to the user.
func (b *Broker) SendToast(message, background string) error {
	if b == nil {
		return errors.New("ws connection nil")
	}

	xb, err := json.Marshal(toast{
		Message:    message,
		Background: background,
	})
	if err != nil {
		return err
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	return b.conn.WriteJSON(Message{
		Type: "toast",
		Data: string(xb),
	})
}
