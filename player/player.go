package player

import (
	"io-server/bullet"
	"io-server/constants"
	"io-server/entity"

	"math"
)

type Player struct {
	Entity       *entity.Entity
	Id           string
	Username     string
	X            float64
	Y            float64
	Score        float64
	FireCooldown float64
	Hp           int
}

type SerializedPlayer struct {
	Id  string
	X   float64
	Y   float64
	Hp  int
	Dir float64
}

func (p *Player) Update(dt float64) *bullet.Bullet {
	// Update score
	p.Score += dt * float64(constants.ScorePerSecond)

	// Make sure the player stays in bounds
	p.X = math.Max(0, math.Min(float64(constants.MapSize), p.X))
	p.Y = math.Max(0, math.Min(float64(constants.MapSize), p.Y))

	// Fire a bullet, if needed
	p.FireCooldown -= dt

	if p.FireCooldown <= 0 {
		p.FireCooldown += float64(constants.PlayerFireCooldown)

		b := bullet.Bullet{Entity: p.Entity, ParentID: p.Entity.Id, X: p.Entity.X, Y: p.Entity.Y}
		return &b
	}

	return nil
}

func (p *Player) TakeBulletDamage() {
	p.Hp -= constants.BulletDamage
}

func (p *Player) OnDealtDamage() {
	p.Score += float64(constants.ScoreBulletHit)
}

func (p *Player) SerializeForUpdate() *SerializedPlayer {
	a := p.Entity.SerializeForUpdate()

	s := SerializedPlayer{Id: a.Id, X: a.X, Y: a.Y}
	s.Dir = p.Entity.Dir
	s.Hp = p.Hp

	return &s
}
