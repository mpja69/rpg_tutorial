package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	PlayerImage *ebiten.Image
	x, y        float64
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.x += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.x -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.y -= 1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.x, g.y)
	screen.DrawImage(
		g.PlayerImage.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	playerImage, _, err := ebitenutil.NewImageFromFile("./assets/images/ninja.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(&Game{PlayerImage: playerImage, x: 100, y: 100}); err != nil {
		log.Fatal(err)
	}
}

// vim: ts=4
