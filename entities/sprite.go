package entities

import (
	"rpg-tutorial/animations"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img              *ebiten.Image
	X, Y, Dx, Dy     float64
	SpriteSheet      *spritesheet.SpriteSheet
	RunningAnimation *animations.Animation
}

func (s *Sprite) Draw(screen *ebiten.Image, translateFunc func(*ebiten.DrawImageOptions)) {
	opts := ebiten.DrawImageOptions{}
	// Move according to the Potion
	opts.GeoM.Translate(s.X, s.Y)
	// Move according to the Camera
	translateFunc(&opts)

	screen.DrawImage(
		s.Img.SubImage(
			s.SpriteSheet.Rect(s.RunningAnimation.Frame()),
		).(*ebiten.Image),
		&opts,
	)

}
