package entities

type Enemy struct {
	*Sprite
	CanFollow bool
	Player    *Player
}

// TODO: Add a FOV Field Of Vision, so the enemy can detect/loose the player
//			- 2 different threashholds

func (enemy *Enemy) Update() {
	enemy.Dx = 0
	enemy.Dy = 0

	if !enemy.CanFollow {
		return
		// continue
	}
	if enemy.X < enemy.Player.X {
		enemy.Dx += 1.0
	}
	if enemy.X > enemy.Player.X {
		enemy.Dx -= 1.0
	}
	if enemy.Y < enemy.Player.Y {
		enemy.Dy += 1.0
	}
	if enemy.Y > enemy.Player.Y {
		enemy.Dy -= 1.0
	}

}
