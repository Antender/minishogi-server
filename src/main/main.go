package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"shogi"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var blackPlayer *websocket.Conn
var whitePlayer *websocket.Conn
var moves []*shogi.Move
var moveMutex sync.Mutex

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
		err = conn.WriteMessage(websocket.TextMessage, []byte("b"))
	} else {
		whitePlayer = conn
		err = conn.WriteMessage(websocket.TextMessage, []byte("w"))
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	moveMutex.Lock()
	moveStr := ""
	for _, move := range moves {
		moveStr += move.String()
		moveStr += " "
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(moveStr))
	moveMutex.Unlock()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer clearPlayer(black)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		moveMutex.Lock()
		moves = append(moves, shogi.MoveFromString([]rune(string(p))))
		if black {
			if whitePlayer != nil {
				_ = whitePlayer.WriteMessage(websocket.TextMessage, p)
			}
			if blackPlayer != nil {
				err = blackPlayer.WriteMessage(websocket.TextMessage, []byte("ok"))
			}
		} else {
			if blackPlayer != nil {
				_ = blackPlayer.WriteMessage(websocket.TextMessage, p)
			}
			if whitePlayer != nil {
				err = whitePlayer.WriteMessage(websocket.TextMessage, []byte("ok"))
			}
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		moveMutex.Unlock()
	}
}

func serveApp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Application")
}

func main() {
	moves = make([]*shogi.Move, 0, 30)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	http.HandleFunc("/play", shogiHandler)
	http.HandleFunc("/", serveApp)
	http.ListenAndServe(":8080", nil)
}
