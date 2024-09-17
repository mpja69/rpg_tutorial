package entities

import (
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
	Direction  PlayerDirection
}

func (p *Player) Update() {
	p.Dx = 0
	p.Dy = 0
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.Dx = 2.0
		p.Direction = Right
		p.Animations[p.Direction].Update()
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.Dx = -2.0
		p.Direction = Left
		p.Animations[p.Direction].Update()
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.Dy = 2.0
		p.Direction = Down
		p.Animations[p.Direction].Update()
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.Dy = -2.0
		p.Direction = Up
		p.Animations[p.Direction].Update()
	}
}

func (p *Player) Draw(screen *ebiten.Image, translateFunc func(*ebiten.DrawImageOptions)) {
	opts := ebiten.DrawImageOptions{}
	// Move according to the Potion
	opts.GeoM.Translate(p.X, p.Y)
	// Move according to the Camera
	translateFunc(&opts)

	screen.DrawImage(
		p.Img.SubImage(
			p.SpriteSheet.Rect(p.Animations[p.Direction].Frame()),
		).(*ebiten.Image),
		&opts,
	)
}
