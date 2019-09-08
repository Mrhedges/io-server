package helpers

import (
	"io-server/bullet"
	"io-server/constants"
	"io-server/player"
)

func ApplyCollisions(players []*player.Player, bullets []*bullet.Bullet) []*bullet.Bullet {
	var destroyedBullets []*bullet.Bullet

	for _, bullet := range bullets {
		// Look for a player (who didn't create the bullet) to collide each bullet with.
		// As soon as we find one, break out of the loop to prevent double counting a bullet.
		for _, player := range players {

			// const bullet = bullets[i];
			// const player = players[j];
			if bullet.ParentID != player.Id && player.Entity.DistanceTo(bullet.Entity) <= float64(constants.PlayerRadius)+float64(constants.BulletRadius) {
				destroyedBullets = append(destroyedBullets, bullet)
				player.TakeBulletDamage()
				break
			}
		}
	}

	return destroyedBullets
}
