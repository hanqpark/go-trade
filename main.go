package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func onMessage(message []byte) {
	// 원하는 로직을 수행
	data := string(message)
	fmt.Println(data)
}

func onConnect(conn *websocket.Conn) {
	fmt.Println("connected!")

	// 연결 후 요청 전송
	req := []map[string]interface{}{
		{"ticket":"test"},
		{"type":"ticker","codes":[]string{"KRW-BTC"}},
		{"format":"SIMPLE"},
	}

	request, err := json.Marshal(req)
	if err != nil {
		log.Fatal("Encoding JSON Error!: ", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte(string(request))); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// WebSocket 연결 주소
	url := "wss://api.upbit.com/websocket/v1"

	// WebSocket 연결 설정
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 종료 시그널을 처리할 채널 생성
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 웹소켓 이벤트 처리
	go onConnect(conn)

	// 채널을 통해 종료 시그널을 받을 때까지 대기
	for {
		select {
		case <-interrupt:
			fmt.Println("Interrupt signal received, closing WebSocket connection...")
			return
		default:
			// 메시지 수신 대기
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
				return
			}

			// 수신한 메시지 처리
			onMessage(message)
		}
	}
}
