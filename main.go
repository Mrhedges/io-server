package main

import (
	"io-server/constants"
	"io-server/game"
	"io-server/handlers"

	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(so socketio.Conn) error {
		log.Println("on connection")

		so.Join("chat")

		return nil
	})

	g := new(game.Game)
	// g := game.Game

	server.OnEvent("/", constants.JoinGame, handlers.HandleJoinGame(g))
	server.OnEvent("/", constants.Input, handlers.HandleInput(g))

	server.OnEvent("/", "chat message", func(so socketio.Conn, msg string) {
		log.Println("emit:", msg)
		so.Emit("chat message", msg)

		br := socketio.NewBroadcast()
		br.Send("chat", "chat message", msg)
		//  so.BroadcastTo("chat", "chat message", msg)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Println("Closed", msg)
	})

	server.OnError("/", func(e error) {
		log.Println("Error:", err)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/assets/", http.FileServer(http.Dir("public/assets")))

	http.Handle("/stylesheets", http.FileServer(http.Dir("./static/stylesheets")))
	http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("./static/stylesheets"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// http.Handle("/stylesheets", http.FileServer(http.Dir("./static/stylesheets")))
	// http.Handle("/stylesheets", http.FileServer(http.Dir("./static/stylesheets")))
	// http.Handle("/stylesheets", http.FileServer(http.Dir("./static/stylesheets")))

	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
