package entity

import (
	"io-server/point"

	"math"
)

type Entity struct {
	Id string
	X float64
	Y float64
	Dir float64
	Speed float64
}

func (e *Entity) Update(dt float64) {
	e.X += dt * e.Speed * math.Sin(e.Dir);
	e.Y -= dt * e.Speed * math.Cos(e.Dir);
}

func (e Entity) DistanceTo(other *Entity) float64 {
	dx := e.X - other.X;
	dy := e.Y - other.Y;
	return math.Sqrt(dx * dx + dy * dy);
}

func (e *Entity) SetDirection(dir float64) {
	e.Dir = dir;
}

func (e *Entity) SerializeForUpdate() *point.Point {
	p := point.Point{Id: e.Id, X: e.X, Y: e.Y}
	return &p;
}