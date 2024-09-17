package entities

import (
	"image"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img          *ebiten.Image
	X, Y, Dx, Dy float64
	SpriteSheet  *spritesheet.SpriteSheet
}

func (s *Sprite) Draw(screen *ebiten.Image, translateFunc func(*ebiten.DrawImageOptions)) {
	opts := ebiten.DrawImageOptions{}
	// Move according to the Potion
	opts.GeoM.Translate(s.X, s.Y)
	// Move according to the Camera
	translateFunc(&opts)

	screen.DrawImage(
		s.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

}

func (s *Sprite) Rect() image.Rectangle {
	return image.Rect(int(s.X), int(s.Y), int(s.X)+16, int(s.Y)+16)
}

func (sprite *Sprite) AdjustForColliders(colliders []image.Rectangle) {
	sprite.X += sprite.Dx
	sprite.X = sprite.CheckCollisionX(colliders)

	sprite.Y += sprite.Dy
	sprite.Y = sprite.CheckCollisionY(colliders)
}

func (sprite *Sprite) CheckCollisionX(colliders []image.Rectangle) float64 {
	for _, collider := range colliders {
		if collider.Overlaps(sprite.Rect()) {
			// if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+16, int(sprite.Y)+16)) {
			if sprite.Dx > 0 {
				return float64(collider.Min.X) - 16
			} else if sprite.Dx < 0 {
				return float64(collider.Max.X)
			}
		}
	}
	return sprite.X
}
func (sprite *Sprite) CheckCollisionY(colliders []image.Rectangle) float64 {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+16, int(sprite.Y)+16)) {
			if sprite.Dy > 0 {
				return float64(collider.Min.Y) - 16
			} else if sprite.Dy < 0 {
				return float64(collider.Max.Y)
			}
		}
	}
	return sprite.Y
}
