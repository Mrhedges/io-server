package bullet

import (
	"io-server/constants"
	"io-server/entity"
)

type Bullet struct {
	Entity *entity.Entity
	ParentID string
	X float64
	Y float64
	Dir string
}

func (b *Bullet) Update(dt float64) bool {
	b.Entity.Update(dt);
	return b.X < 0 || b.X > float64(constants.MapSize) || b.Y < 0 || b.Y > float64(constants.MapSize);
}