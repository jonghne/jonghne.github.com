package ws

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type connection struct {
	conn *websocket.Conn
	id int64
	send chan []byte
}

type MsgProcessor interface {
	HandleMsg(id int64, msg []byte)
	Subscribe(send chan []byte) int64
	Unsubscribe(id int64)
}

func (link *connection) readMsg(processor MsgProcessor) {
	for {
		fmt.Println("begin")
		_, msg, err := link.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(msg)
		processor.HandleMsg(link.id, msg)
		//if string(msg) == "ping" {
		//	fmt.Println("ping")
		//	time.Sleep(time.Second * 2)
		//	err = conn.WriteMessage(msgType, []byte("pong"))
		//	if err != nil {
		//		fmt.Println(err)
		//		return
		//	}
		//} else {
		//	link.conn.Close()
		//	fmt.Println(string(msg))
		//	return
		//}
		//fmt.Println("after")
	}
	//fmt.Println("end")
}

func (link *connection) sendMsg() {
	for {
		select {
		case data := <- link.send:
			err := link.conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					fmt.Println(err)
					return
				}
		}
	}
}

func StartCommunicate(processor MsgProcessor) {
	indexFile, err := os.Open("ws/html/index.html")
	checkErr(err)

	index, err := ioutil.ReadAll(indexFile)
	checkErr(err)

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		send := make(chan []byte)
		id := processor.Subscribe(send)
		conn.SetCloseHandler(func(code int, text string) error {
			fmt.Println(code, text)
			processor.Unsubscribe(id)
			return nil
		})
		link := &connection{conn, id, send}

		go link.sendMsg()

		go link.readMsg(processor)

		//for {
		//	msgType, msg, err := conn.ReadMessage()
		//	if err != nil {
		//		fmt.Println(err)
		//		return
		//	}
		//	if string(msg) == "ping" {
		//		fmt.Println("ping")
		//		time.Sleep(time.Second * 2)
		//		err = conn.WriteMessage(msgType, []byte("pong"))
		//		if err != nil {
		//			fmt.Println(err)
		//			return
		//		}
		//	} else {
		//		conn.Close()
		//		fmt.Println(string(msg))
		//		return
		//	}
		//}

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	})

	http.ListenAndServe(":3000", nil)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
