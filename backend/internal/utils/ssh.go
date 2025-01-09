package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"io"
	"sync"
	"time"
)

// wsBufferWriter
type wsBufferWriter struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (w *wsBufferWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

// 封装SSH连接
type SshConn struct {
	StdinPipe   io.WriteCloser  //SSH 会话的标准输入流
	ComboOutput *wsBufferWriter //捕获和存储 SSH 会话的输出流
	Session     *ssh.Session
}

func (s *SshConn) Close() {
	if s.Session != nil {
		s.Session.Close()
	}
}

// NewSshConn - 创建新的 SSH 会话,初始化终端
func NewSshConn(cols, rows int, sshClient *ssh.Client) (*SshConn, error) {
	sshSession, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	stdinP, err := sshSession.StdinPipe()
	if err != nil {
		return nil, err
	}

	comboWriter := new(wsBufferWriter)
	sshSession.Stdout = comboWriter
	sshSession.Stderr = comboWriter

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // 禁用回显
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
	}
	if err := sshSession.RequestPty("xterm", rows, cols, modes); err != nil {
		return nil, err
	}

	if err := sshSession.Shell(); err != nil {
		return nil, err
	}

	return &SshConn{StdinPipe: stdinP, ComboOutput: comboWriter, Session: sshSession}, nil
}

// NewSshClient - 创建 SSH 客户端
func NewSshClient(addr, user, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
	}
	return ssh.Dial("tcp", addr, config)
}

type Options struct {
	Addr     string `json:"addr"`
	User     string `json:"user"`
	Password string `json:"password"`
	Cols     int    `json:"cols"`
	Rows     int    `json:"rows"`
}

// Terminal - 封装 WebSocket 与 SSH 会话的管理
type Terminal struct {
	Opts     Options
	Ws       *websocket.Conn
	conn     *ssh.Client
	session  *SshConn
	cancelFn context.CancelFunc
}

// Run -运行Websocket 与 SSH 会话
func (t *Terminal) Run() {
	var err error

	// 创建 SSH 客户端
	t.conn, err = NewSshClient(t.Opts.Addr, t.Opts.User, t.Opts.Password)
	if t.handleWsError(err) {
		return
	}
	defer t.conn.Close()

	// 创建 SSH 会话
	t.session, err = NewSshConn(t.Opts.Cols, t.Opts.Rows, t.conn)
	if t.handleWsError(err) {
		return
	}
	defer t.session.Close()

	// 通道用于通知 goroutine 停止
	quitChan := make(chan bool, 3)

	// 日志缓冲区
	var logBuff bytes.Buffer

	go t.session.ReceiveWsMsg(t.Ws, &logBuff, quitChan)
	go t.session.SendComboOutput(t.Ws, quitChan)
	go t.session.SessionWait(quitChan)

	<-quitChan
}

// handleWsError - 处理 WebSocket 错误并发送错误信息
func (t *Terminal) handleWsError(err error) bool {
	if err != nil {
		if writeErr := t.Ws.WriteMessage(websocket.CloseMessage, []byte(err.Error())); writeErr != nil {
			// 记录日志（可换成实际日志工具）
			println("WebSocket CloseMessage 写入失败:", writeErr.Error())
		}
		return true
	}
	return false
}

// ReceiveWsMsg - 接收 WebSocket 消息并写入 SSH 会话
func (s *SshConn) ReceiveWsMsg(wsConn *websocket.Conn, logBuff *bytes.Buffer, exitCh chan bool) {
	defer close(exitCh)
	for {
		_, wsData, err := wsConn.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			return
		}
		if err != nil {
			println("WebSocket 读取失败:", err.Error())
			return
		}

		var msg wsMsg
		if err := json.Unmarshal(wsData, &msg); err != nil {
			println("WebSocket 消息解析失败:", err.Error())
			continue
		}

		switch msg.Type {
		case wsMsgResize:
			if msg.Cols > 0 && msg.Rows > 0 {
				if err := s.Session.WindowChange(msg.Rows, msg.Cols); err != nil {
					println("SSH 窗口大小更改失败:", err.Error())
				}
			}
		case wsMsgCmd:
			if _, err := s.StdinPipe.Write([]byte(msg.Cmd)); err != nil {
				println("WebSocket 数据写入 SSH stdin 失败:", err.Error())
			}
		}
	}
}

// SendComboOutput - 定期将 SSH 输出发送到 WebSocket
func (s *SshConn) SendComboOutput(wsConn *websocket.Conn, exitCh chan bool) {
	ticker := time.NewTicker(120 * time.Millisecond)
	defer ticker.Stop()
	defer close(exitCh)

	for {
		select {
		case <-ticker.C:
			if s.ComboOutput.buffer.Len() > 0 {
				if err := wsConn.WriteMessage(websocket.TextMessage, s.ComboOutput.buffer.Bytes()); err != nil {
					println("WebSocket 写入失败:", err.Error())
					return
				}
				s.ComboOutput.buffer.Reset()
			}
		case <-exitCh:
			return
		}
	}
}

// SessionWait - 等待 SSH 会话结束
func (s *SshConn) SessionWait(exitCh chan bool) {
	if err := s.Session.Wait(); err != nil {
		println("SSH 会话结束:", err.Error())
	}
	close(exitCh)
}

// wsMsg - WebSocket 消息类型
type wsMsg struct {
	Type string `json:"type"`
	Cmd  string `json:"cmd"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

const (
	wsMsgCmd    = "cmd"
	wsMsgResize = "resize"
)
