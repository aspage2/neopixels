package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pnpleds/ledproto"

	"github.com/gerow/go-color"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received websocket request")
	conn, err := grpc.Dial("192.168.2.25:50051", grpc.WithInsecure())
	if err != nil {
		log.Println("dial:", err)
		return
	}
	defer conn.Close()
	client := ledproto.NewLedStripClient(conn)
	rpc, err := client.SetRealtime(context.Background())
	if err != nil {
		log.Println("ledserver:", err.Error())
		return
	}

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
		rgb := hsl.ToRGB()
		i := toInt(&rgb)
		err = rpc.Send(&ledproto.SetColorRequest{
			Color: i,
		})
		if err != nil {
			log.Println("send ledserver:", err.Error())
			return
		}
	}

}

func toInt(rgb *color.RGB) uint32 {
	return ((uint32(rgb.R*255) & 0xff) << 16) | (uint32(rgb.G*255) & 0xff << 8) | (uint32(rgb.B*255) & 0xff)
}

func main() {
	http.HandleFunc("/ws", serveWs)
	fmt.Println("Serving")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
