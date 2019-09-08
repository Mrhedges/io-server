package game

import (
	"io-server/bullet"
	"io-server/constants"
	"io-server/helpers"
	"io-server/player"
	"io-server/point"

	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

type UpdatedState struct {
	T           float64
	Me          *player.SerializedPlayer
	Others      []*player.SerializedPlayer
	Bullets     []*point.Point
	Leaderboard []Leader
}

type Game struct {
	sockets          map[string]socketio.Conn
	players          map[string]*player.Player
	bullets          []*bullet.Bullet
	lastUpdateTime   int
	shouldSendUpdate bool
}

func (g Game) AddPlayer(s socketio.Conn, username string) {
	g.sockets[s.ID()] = s

	// Generate a position to start this player at.
	x := float64(constants.MapSize) * (0.25 + rand.Float64()*0.5)
	y := float64(constants.MapSize) * (0.25 + rand.Float64()*0.5)

	log.Println(x, y)
	// todo need to add an entity to this
	g.players[s.ID()] = &player.Player{Id: s.ID(), Username: username, X: x, Y: y}
}

func (g Game) RemovePlayer(s socketio.Conn) {

	// remove socket from game
	_, socket := g.sockets[s.ID()]
	if socket {
		delete(g.sockets, s.ID())
	}

	// remove player from game
	_, player := g.players[s.ID()]
	if player {
		delete(g.players, s.ID())
	}
}

func (g *Game) HandleInput(s socketio.Conn, dir float64) {

	_, player := g.players[s.ID()]
	if player {
		g.players[s.ID()].Entity.SetDirection(dir)
	}
}

func (g Game) Update() {
	// Calculate time elapsed
	t := int(time.Now().UnixNano() / int64(time.Millisecond))
	dt := (t - g.lastUpdateTime) / 1000
	g.lastUpdateTime = t

	// Update each bullet
	var bulletsToRemove []*bullet.Bullet

	for _, bullet := range g.bullets {
		if bullet.Update(float64(dt)) {
			bulletsToRemove = append(bulletsToRemove, bullet)
		}
	}

	for i, bullet := range g.bullets {
		for _, bulletToRemove := range bulletsToRemove {
			if bullet.Entity.Id == bulletToRemove.Entity.Id {
				g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
			}
		}
	}

	// // Update each player
	for k, _ := range g.sockets {
		player := g.players[k]
		newBullet := player.Update(float64(dt))

		if newBullet != nil {
			g.bullets = append(g.bullets, newBullet)
		}
	}

	// Apply collisions, give players score for hitting bullets
	var players []*player.Player

	for _, player := range g.players {
		players = append(players, player)
	}

	var destroyedBullets = helpers.ApplyCollisions(players, g.bullets)

	for _, destroyedBullet := range destroyedBullets {
		_, player := g.players[destroyedBullet.ParentID]

		if player {
			g.players[destroyedBullet.ParentID].OnDealtDamage()
		}

		// if g.players[destroyedBullet.ParentID] {
		// 	g.players[destroyedBullet.ParentID].OnDealtDamage()
		// }
	}

	// // Check if any players are dead
	for i, socket := range g.sockets {
		player := g.players[i]

		if player.Hp <= 0 {
			socket.Emit(constants.GameOver)
			g.RemovePlayer(socket)
		}
	}

	// // Check if any players are dead
	// Object.keys(g.sockets).forEach(playerID => {
	// 	const socket = g.sockets[playerID];
	// 	const player = g.players[playerID];
	// 	if (player.hp <= 0) {
	// 		socket.emit(Constants.MSG_TYPES.GAME_OVER);
	// 		g.removePlayer(socket);
	// 	}
	// });

	// Send a game update to each player every other time
	if g.shouldSendUpdate {
		leaderboard := g.GetLeaderboard()

		for _, socket := range g.sockets {
			player := g.players[socket.ID()]
			socket.Emit(constants.GameUpdate, g.CreateUpdate(player, leaderboard))
		}

		g.shouldSendUpdate = false
	} else {
		g.shouldSendUpdate = true
	}

	// // Send a game update to each player every other time
	// if (g.shouldSendUpdate) {
	// 	const leaderboard = g.getLeaderboard();
	// 	Object.keys(g.sockets).forEach(playerID => {
	// 		const socket = g.sockets[playerID];
	// 		const player = g.players[playerID];
	// 		socket.emit(Constants.MSG_TYPES.GAME_UPDATE, g.createUpdate(player, leaderboard));
	// 	});
	// 	g.shouldSendUpdate = false;
	// } else {
	// 	g.shouldSendUpdate = true;
	// }
}

type Leader struct {
	Username string
	Score    float64
}

func (g Game) GetLeaderboard() []Leader {

	var topFive []*player.Player
	var leaderboard []Leader
	var sortedPlayers []*player.Player

	for _, player := range g.players {
		sortedPlayers = append(sortedPlayers, player)
	}

	sort.Slice(sortedPlayers, func(i, j int) bool {
		return sortedPlayers[i].Score < sortedPlayers[j].Score
	})

	topFive = sortedPlayers[:5]

	for _, player := range topFive {
		leaderboard = append(leaderboard, Leader{Username: player.Username, Score: math.Round(player.Score)})
	}

	return leaderboard
	// return Object.values(g.players)
	// .sort((p1, p2) => p2.score - p1.score)
	// .slice(0, 5)
	// .map(p => ({ username: p.username, score: Math.round(p.score) }));
}

func (g Game) CreateUpdate(p *player.Player, leaderboard []Leader) *UpdatedState {

	var nearbyPlayers []*player.Player
	var nearbyBullets []*bullet.Bullet
	var serializedPlayers []*player.SerializedPlayer
	var serialiedBullets []*point.Point

	for _, player := range g.players {
		if p != player && p.Entity.DistanceTo(player.Entity) <= float64(constants.MapSize)/2 {
			nearbyPlayers = append(nearbyPlayers, player)
		}
	}

	for _, player := range nearbyPlayers {
		serializedPlayers = append(serializedPlayers, player.SerializeForUpdate())
	}

	for _, bullet := range g.bullets {
		if bullet.Entity.DistanceTo(p.Entity) <= float64(constants.MapSize)/2 {
			nearbyBullets = append(nearbyBullets, bullet)
		}
	}

	for _, bullet := range nearbyBullets {
		serialiedBullets = append(serialiedBullets, bullet.Entity.SerializeForUpdate())
	}

	return &UpdatedState{
		T:           float64(time.Now().UnixNano() / int64(time.Millisecond)),
		Me:          p.SerializeForUpdate(),
		Others:      serializedPlayers,
		Bullets:     serialiedBullets,
		Leaderboard: leaderboard,
	}
	// const nearbyPlayers = Object.values(g.players).filter(
	// 	p => p !== player && p.distanceTo(player) <= Constants.MAP_SIZE / 2,
	// );
	// const nearbyBullets = g.bullets.filter(
	// 	b => b.distanceTo(player) <= Constants.MAP_SIZE / 2,
	// );

	// return {
	// t: Date.now(),
	// me: player.serializeForUpdate(),
	// others: nearbyPlayers.map(p => p.serializeForUpdate()),
	// bullets: nearbyBullets.map(b => b.serializeForUpdate()),
	// leaderboard,
	// };
}
