package entities

import (
	"log"
	"rpg-tutorial/animations"

	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerDirection uint8

const (
	Down PlayerDirection = iota
	Up
	Right
	Left
)

type Player struct {
	*Sprite
	Health     uint
	Animations map[PlayerDirection]*animations.Animation
}

func (p *Player) Draw(screen *ebiten.Image, translateFunc func(*ebiten.DrawImageOptions)) {
	opts := ebiten.DrawImageOptions{}
	// Move according to the Potion
	opts.GeoM.Translate(p.X, p.Y)
	// Move according to the Camera
	translateFunc(&opts)

	log.Println("Player.Draw()")

	screen.DrawImage(
		p.Img.SubImage(
			p.SpriteSheet.Rect(p.ActiveAnimation(int(p.Dx), int(p.Dy)).Frame()),
		).(*ebiten.Image),
		&opts,
	)

}
func (p *Player) ActiveAnimation(dx, dy int) *animations.Animation {
	if dx > 0 {
		return p.Animations[Right]
	}
	if dx < 0 {
		return p.Animations[Left]
	}
	if dy > 0 {
		return p.Animations[Down]
	}
	if dy < 0 {
		return p.Animations[Up]
	}
	return p.Animations[Down] // Maybe return DOWN as default...because it's facing forward?
}
