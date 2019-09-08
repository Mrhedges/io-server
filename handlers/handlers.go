package handlers

import (
	"io-server/game"

	socketio "github.com/googollee/go-socket.io"
)

// when to use '*'?
func HandleJoinGame(g *game.Game) func(socketio.Conn, string) {
	return func(s socketio.Conn, username string) {
		g.AddPlayer(s, username)
	}
}

func HandleInput(g *game.Game) func(socketio.Conn, float64) {
	return func(s socketio.Conn, dir float64) {
		g.HandleInput(s, dir)
	}
}

func RemovePlayer(g *game.Game) func(socketio.Conn) {
	return func(s socketio.Conn) {
		g.RemovePlayer(s)
	}
}
