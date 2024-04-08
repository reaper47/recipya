package models

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
	"time"
)

// Broker represents a client connection and its associated information.
type Broker struct {
	brokers map[int64]*Broker
	conn    *websocket.Conn
	mu      sync.Mutex
	userID  int64
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
	Toast    Toast  `json:"toast"`    // Toast to display to the user.
}

// NewBroker creates a new Broker instance for a specific user and adds it to the brokers map.
// The userID is used for identification and cleanup purposes.
func NewBroker(userID int64, brokers map[int64]*Broker, conn *websocket.Conn) *Broker {
	b := &Broker{
		brokers: brokers,
		conn:    conn,
		userID:  userID,
	}
	go b.ping()
	return b
}

// HideNotification hides the websocket's frontend notification.
func (b *Broker) HideNotification() {
	b.SendProgressStatus("", false, -1, -1)
}

// SendFile sends a file to the connected client.
func (b *Broker) SendFile(fileName string, data *bytes.Buffer) {
	fileNameAttr := slog.String("fileName", fileName)

	if b == nil {
		slog.Error("Websocket connection is nil", fileNameAttr)
		return
	}

	if data == nil {
		slog.Error("Data is nil", fileNameAttr)
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	err := b.conn.WriteJSON(Message{
		Type:     "file",
		FileName: fileName,
		Data:     base64.StdEncoding.EncodeToString(data.Bytes()),
	})
	if err != nil {
		slog.Error("Failed to send file", fileNameAttr, "error", err)
	}
}

// SendProgress sends a progress update with a title and value to the client.
// The isToastVisible parameter controls whether the progress bar is displayed in a toast notification.
func (b *Broker) SendProgress(title string, value, numValues int) {
	titleAttr := slog.String("title", title)
	valueAttr := slog.Int("value", value)
	numValuesAttr := slog.Int("numValues", numValues)

	if b == nil {
		slog.Error("Websocket connection is nil", titleAttr, valueAttr, numValuesAttr)
		return
	}

	content := fmt.Sprintf(`
		<div id="ws-notification-container" class="z-20 fixed bottom-0 right-0 p-6 cursor-default">
			<div class="bg-blue-500 text-white px-4 py-2 rounded shadow-md">
				<p class="font-medium text-center pb-1">%s</p>
				<div id="export-progress"><progress max="100" value="%f"></progress></div>
			</div>
		</div>`, title, float64(value)/float64(numValues)*100)

	b.mu.Lock()
	defer b.mu.Unlock()

	err := b.conn.WriteMessage(websocket.TextMessage, []byte(content))
	if err != nil {
		slog.Error("Broker.SendProgress failed", titleAttr, valueAttr, numValuesAttr, "error", err)
	}
}

// SendProgressStatus sends a progress update with a title and value, optionally hiding the toast notification.
func (b *Broker) SendProgressStatus(title string, isToastVisible bool, value, numValues int) {
	titleAttr := slog.String("title", title)
	valueAttr := slog.Int("value", value)
	numValuesAttr := slog.Int("numValues", numValues)

	if b == nil {
		slog.Error("Websocket connection is nil", titleAttr, valueAttr, numValuesAttr)
		return
	}

	content := `
		<div id="ws-notification-container" class="z-20 fixed bottom-0 right-0 p-6 cursor-default %s">
			<div class="bg-blue-500 text-white px-4 py-2 rounded shadow-md">
				<p class="font-medium text-center pb-1">%s</p>
				<div id="export-progress"><progress max="100" value="%f"></progress></div>
			</div>
		</div>`

	if isToastVisible {
		content = fmt.Sprintf(content, "", title, float64(value)/float64(numValues)*100)
	} else {
		content = fmt.Sprintf(content, "hidden", title, float64(value)/float64(numValues)*100)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	err := b.conn.WriteMessage(websocket.TextMessage, []byte(content))
	if err != nil {
		slog.Error("Broker.SendProgressStatus failed", titleAttr, valueAttr, numValuesAttr, "content", content, "error", err)
	}
}

func (b *Broker) ping() {
	defer func() {
		if b != nil && b.conn != nil {
			delete(b.brokers, b.userID)
			_ = b.conn.Close()
		}
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
	if b == nil || b.conn == nil {
		slog.Error("Broker or broker connection is nil")
		return
	}

	b.conn.SetPingHandler(func(_ string) error {
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
func (b *Broker) SendToast(toast Toast) {
	toastAttr := slog.Any("toast", toast)

	if b == nil {
		slog.Error("Websocket connection is nil", toastAttr)
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	err := b.conn.WriteJSON(Message{Type: "toast", Toast: toast})
	if err != nil {
		slog.Error("Broker.SendToast failed", toastAttr, "error", err)
	}
}
