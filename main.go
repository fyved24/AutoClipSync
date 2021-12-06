package main

import (
	"AutoClipSync/server"
	"AutoClipSync/util"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
)

var addr = flag.String("addr", "localhost:9928", "http service address")

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	isServerStared :=  make(chan int)
	// 没有开启服务的话，开启一个服务
	if util.PortInUse(9928) == -1{
		go func() {
			server.StartWsServer(*addr, isServerStared)
		}()
		<-isServerStared
		fmt.Println("isServerStared")
	}

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	socket, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial", err)
	}

	defer socket.Close()

	done := make(chan struct{})
	go func() {
		for {
			_, message, err := socket.ReadMessage()
			fmt.Println(string(message))
			clipboard.WriteAll(string(message))
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}

}