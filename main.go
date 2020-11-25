package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	color "github.com/gerow/go-color"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received websocket request")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err.Error())
		return
	}
	defer ws.Close()

	var hsl color.HSL

	for {
		err := ws.ReadJSON(&hsl)
		if err != nil {
			log.Println("receive:", err.Error())
			return
		}
		fmt.Println(hsl)
		fmt.Println(hsl.ToHTML())
	}

}

func main () {
	http.HandleFunc("/ws", serveWs)
	fmt.Println("Serving")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
