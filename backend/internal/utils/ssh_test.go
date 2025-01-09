package utils

import (
	"bytes"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Mock WebSocket Server for Testing
func startMockWebSocketServer(t *testing.T, messageHandler func(msgType int, msg []byte)) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			messageHandler(msgType, msg)
		}
	}))
	return server
}

// Mock SSH Server (Simulated)
func startMockSSHServer() (*ssh.Client, error) {
	// Replace this with a real SSH server address or local testing SSH server.
	return NewSshClient("127.0.0.1:22", "testuser", "password")
}

// Test SSH and WebSocket Integration
func TestSSHWebSocketIntegration(t *testing.T) {
	// Mock WebSocket Server
	mockServer := startMockWebSocketServer(t, func(msgType int, msg []byte) {
		t.Logf("Received message: %s", string(msg))
	})
	defer mockServer.Close()

	// Establish WebSocket connection to the mock server
	wsURL := strings.Replace(mockServer.URL, "http", "ws", 1)
	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err, "Failed to connect to mock WebSocket server")
	defer wsConn.Close()

	// Setup SSH connection (replace with your SSH server)
	client, err := NewSshClient("43.136.179.241:22", "ubuntu", "CTct288513832##")
	assert.NoError(t, err, "Failed to connect to SSH server")
	defer client.Close()

	// Initialize SSH Terminal
	opts := Options{
		Addr:     "127.0.0.1:22",
		User:     "username",
		Password: "password",
		Cols:     80,
		Rows:     24,
	}
	terminal := NewTerminal(wsConn, opts)

	// Run the terminal in a separate goroutine
	go terminal.Run()

	// Simulate WebSocket messages
	sendMessages := []string{
		`{"type": "cmd", "cmd": "ls\n"}`,
		`{"type": "cmd", "cmd": "pwd\n"}`,
		`{"type": "resize", "cols": 100, "rows": 30}`,
	}

	for _, msg := range sendMessages {
		err = wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
		assert.NoError(t, err, "Failed to send WebSocket message")
	}

	// Allow some time for processing
	time.Sleep(2 * time.Second)
}

// Test SSH Output Handling
func TestSSHOutputHandling(t *testing.T) {
	// Mock WebSocket
	mockBuffer := &bytes.Buffer{}

	mockWriter := &SshConn{
		ComboOutput: &wsBufferWriter{
			buffer: *mockBuffer,
		},
	}

	// Write test data
	_, err := mockWriter.ComboOutput.Write([]byte("Test output"))
	assert.NoError(t, err, "Failed to write to buffer")

	assert.Equal(t, "Test output", mockBuffer.String(), "Output buffer does not match expected content")
}

// Test WebSocket Resize Handling
func TestWebSocketResizeHandling(t *testing.T) {
	// Mock SSH Session
	mockSession := &SshConn{}

	// Simulate Resize Message
	resizeMsg := wsMsg{
		Type: "resize",
		Cols: 120,
		Rows: 40,
	}

	err := mockSession.Session.WindowChange(resizeMsg.Rows, resizeMsg.Cols)
	assert.NoError(t, err, "Failed to handle resize message")
}
