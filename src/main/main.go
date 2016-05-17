package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var blackPlayer *websocket.Conn
var whitePlayer *websocket.Conn

func clearPlayer(black bool) {
	if black {
		blackPlayer = nil
	} else {
		whitePlayer = nil
	}
}

func shogiHandler(w http.ResponseWriter, r *http.Request) {
	var black bool
	err := r.ParseForm()
	if r.FormValue("s") == "b" {
		if blackPlayer != nil {
			return
		}
		black = true
	} else if r.FormValue("s") == "w" {
		if whitePlayer != nil {
			return
		}
		black = false
	} else {
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if black {
		blackPlayer = conn
	} else {
		whitePlayer = conn
	}
	defer clearPlayer(black)
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if black {
			if whitePlayer != nil {
				err = whitePlayer.WriteMessage(messageType, p)
			}
		} else {
			if blackPlayer != nil {
				err = blackPlayer.WriteMessage(messageType, p)
			}
		}
		if err != nil {
			return
		}
	}
}

func serveApp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Application")
}

func main() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	http.HandleFunc("/play", shogiHandler)
	http.HandleFunc("/", serveApp)
	http.ListenAndServe(":8080", nil)
}
