package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"web_penetration/internal/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的 WebSocket 连接（生产环境需限制）
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()
	log.Println("WebSocket connected!")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		log.Printf("Received message: %s", string(msg))
	}
	log.Println("WebSocket disconnected!")

	// 定义 SSH 连接配置
	opts := utils.Options{
		Addr:     "127.0.0.1:22", // 替换为你的 SSH 服务器地址
		User:     "username",     // 替换为 SSH 用户名
		Password: "password",     // 替换为 SSH 密码
		Cols:     80,             // 初始终端列数
		Rows:     24,             // 初始终端行数
	}

	// 创建并运行终端会话
	terminal := utils.Terminal{
		Opts: opts,
		Ws:   conn,
	}
	terminal.Run()
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	serverAddr := "127.0.0.1:8080" // WebSocket 服务监听地址
	fmt.Printf("WebSocket 服务启动: ws://%s/ws\n", serverAddr)

	// 启动 HTTP 服务
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal("HTTP 服务启动失败:", err)
	}
}
