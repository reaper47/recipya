package models

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"maps"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"slices"
	"sync"
	"time"
)

// Broker represents a message broker that manages WebSocket connections for subscribers.
type Broker struct {
	subscribers map[int64][]*websocket.Conn
	mu          sync.Mutex
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

// NewBroker creates a subscribes a new Broker for a specific user.
func NewBroker() *Broker {
	b := &Broker{
		subscribers: make(map[int64][]*websocket.Conn),
		mu:          sync.Mutex{},
	}
	b.Monitor()
	return b
}

// Add adds a connection to the subscriber.
func (b *Broker) Add(userID int64, c *websocket.Conn) {
	b.mu.Lock()
	b.subscribers[userID] = append(b.subscribers[userID], c)
	b.mu.Unlock()

	go func(_ int64, conn *websocket.Conn) {
		for {
			_, _, err := conn.Read(context.Background())
			if err != nil {
				return
			}
		}
	}(userID, c)
}

// Clone creates a deep copy of the Broker.
func (b *Broker) Clone() *Broker {
	return &Broker{
		subscribers: maps.Clone(b.subscribers),
		mu:          sync.Mutex{},
	}
}

// Has checks whether a subscriber with the given userID exists in the Broker's subscribers map.
func (b *Broker) Has(userID int64) bool {
	if b == nil {
		return false
	}

	_, ok := b.subscribers[userID]
	return ok
}

// HideNotification hides the websocket's frontend notification.
func (b *Broker) HideNotification(userID int64) {
	b.SendProgressStatus("", false, -1, -1, userID)
}

// Monitor monitors the websocket connections to clean those closed.
func (b *Broker) Monitor() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)

		defer func() {
			ticker.Stop()
			err := recover()
			if err != nil {
				slog.Error("Websocket ping pong panic", "error", err)
			}
		}()

		for range ticker.C {
			for userID, connections := range b.subscribers {
				connectionsCopy := make([]*websocket.Conn, len(connections))
				copy(connectionsCopy, connections)

				for _, c := range connectionsCopy {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

					err := c.Ping(ctx)
					if err != nil {
						b.mu.Lock()
						b.subscribers[userID] = slices.DeleteFunc(b.subscribers[userID], func(conn *websocket.Conn) bool {
							return conn == c
						})
						b.mu.Unlock()
					}

					cancel()
				}
			}
		}
	}()
}

// SendToast sends a toast notification to the user.
func (b *Broker) SendToast(toast Toast, userID int64) {
	userIDAttr := slog.Int64("userID", userID)
	toastAttr := slog.Any("toast", toast)

	xc, ok := b.subscribers[userID]
	if !ok || len(xc) == 0 {
		slog.Warn("User does not have any websocket connections", userIDAttr, toastAttr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for i, c := range xc {
		err := wsjson.Write(ctx, c, Message{Type: "toast", Toast: toast})
		if err != nil {
			slog.Error("Failed to send toast", userIDAttr, toastAttr, "i", i, "error", err)
		}
	}
}

// SendFile sends a file to the connected client.
func (b *Broker) SendFile(fileName string, data *bytes.Buffer, userID int64) {
	fileNameAttr := slog.String("fileName", fileName)
	userIDAttr := slog.Int64("userID", userID)

	xc, ok := b.subscribers[userID]
	if !ok || len(xc) == 0 {
		slog.Warn("User does not have any websocket connections", userIDAttr, fileNameAttr)
		return
	}

	if data == nil {
		slog.Error("Data is nil", userIDAttr, fileNameAttr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for i, c := range xc {
		err := wsjson.Write(ctx, c, Message{
			Type:     "file",
			FileName: fileName,
			Data:     base64.StdEncoding.EncodeToString(data.Bytes()),
		})
		if err != nil {
			slog.Error("Failed to send file through websocket", userIDAttr, fileNameAttr, "i", i, "file", fileName, "error", err)
		}
	}
}

// SendProgress sends a progress update with a title and value to the client.
// The isToastVisible parameter controls whether the progress bar is displayed in a toast notification.
func (b *Broker) SendProgress(title string, value, numValues int, userID int64) {
	userIDAttr := slog.Int64("userID", userID)
	titleAttr := slog.String("title", title)
	valueAttr := slog.Int("value", value)
	numValuesAttr := slog.Int("numValues", numValues)

	xc, ok := b.subscribers[userID]
	if !ok || len(xc) == 0 {
		slog.Warn("User does not have any websocket connections", userIDAttr, titleAttr, valueAttr, numValuesAttr)
		return
	}

	content := fmt.Sprintf(`
		<div id="ws-notification-container" class="z-20 fixed bottom-0 right-0 p-6 cursor-default">
			<div class="bg-blue-500 text-white px-4 py-2 rounded shadow-md">
				<p class="font-medium text-center pb-1">%s</p>
				<div id="export-progress"><progress max="100" value="%f"></progress></div>
			</div>
		</div>`, title, float64(value)/float64(numValues)*100)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for i, c := range xc {
		err := c.Write(ctx, websocket.MessageText, []byte(content))
		if err != nil {
			slog.Error("Failed to send progress through websocket", titleAttr, valueAttr, numValuesAttr, "content", content, "i", i, "error", err)
		}
	}
}

// SendProgressStatus sends a progress update with a title and value, optionally hiding the toast notification.
func (b *Broker) SendProgressStatus(title string, isToastVisible bool, value, numValues int, userID int64) {
	userIDAttr := slog.Int64("userID", userID)
	titleAttr := slog.String("title", title)
	valueAttr := slog.Int("value", value)
	numValuesAttr := slog.Int("numValues", numValues)

	xc, ok := b.subscribers[userID]
	if !ok || len(xc) == 0 {
		slog.Warn("User does not have any websocket connections", userIDAttr, titleAttr, valueAttr, numValuesAttr)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for i, c := range xc {
		err := c.Write(ctx, websocket.MessageText, []byte(content))
		if err != nil {
			slog.Error("Failed to send progress status through websocket", userIDAttr, titleAttr, valueAttr, numValuesAttr, "i", i, "content", content, "error", err)
		}
	}
}
