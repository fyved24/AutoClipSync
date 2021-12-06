package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)



var upgrader = websocket.Upgrader{} // use default options

var connStore = NewConnStore()

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	connStore.Set(c.RemoteAddr().String(), c)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	log.Printf("%s register:", c.RemoteAddr().String())
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()

		if err != nil {
			delete(connStore.connections, c.RemoteAddr().String())
			log.Printf("remove %s:", c.RemoteAddr().String())
			break
		}
		log.Printf("recv: %s from %s", message, c.RemoteAddr().String())
		go func() {
			for  key, conn := range connStore.connections {
				if key != c.RemoteAddr().String() {
					conn.Write(message)
					log.Printf("to %s", conn.wsconn.RemoteAddr().String())
				}
			}
		}()
	}
}

func StartWsServer(addr string, isServerStared chan  int) {
	http.HandleFunc("/echo", echo)
	isServerStared <- 1
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}


