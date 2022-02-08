package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

// 4kb で設定
const bufferSize = 4096

func handlerWebSocket(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Connection") != "Upgrade" || r.Header.Get("Upgrade") != "websocket" {
		fmt.Println("error")
		w.WriteHeader(400)
		return
	}

	// ここからのやりとりは、HTTP protocol ではなくなるので、 Hijacker を使う
	hijacker := w.(http.Hijacker)
	conn, readWriter, err := hijacker.Hijack()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// ブラウザからのキーを元にして、レスポンス用のキーを作成
	key := r.Header.Get("Sec-Websocket-Key")
	acceptKey := buildAcceptKey(key)

	readWriter.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	readWriter.WriteString("Upgrade: websocket\r\n")
	readWriter.WriteString("Connection: Upgrade\r\n")
	readWriter.WriteString("Sec-WebSocket-Accept: " + acceptKey + "\r\n")
	readWriter.WriteString("\r\n") // 空白行でステータスラインの終わりを示す
	readWriter.Flush()

	sendFrame := buildFrame("どやさ")
	readWriter.Write(sendFrame.toBytes())
	readWriter.Flush()

	data := make([]byte, bufferSize)
	for {
		frame := Frame{}
		n, err := readWriter.Read(data)
		if err != nil {
			panic(err)
		}

		frame.parse(data[:n])
		fmt.Println(string(frame.payloadData))
	}
}

func buildAcceptKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	h.Write([]byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func main() {
	var httpServer http.Server
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/websocket", handlerWebSocket)
	log.Println("start http listening :3000")
	httpServer.Addr = ":3000"
	log.Println(httpServer.ListenAndServe())
}
