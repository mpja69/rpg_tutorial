package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Img  *ebiten.Image
	x, y float64
}

type Enemy struct {
	*Sprite
	canFollow bool
}

type Potion struct {
	*Sprite
	healingPower uint
}

type Player struct {
	*Sprite
	Health uint
}

type Game struct {
	player  *Player
	enemies []*Enemy
	potions []*Potion
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.x += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.x -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.y -= 2
	}

	for _, sprite := range g.enemies {
		if !sprite.canFollow {
			continue
		}
		if sprite.x < g.player.x {
			sprite.x += 1
		}
		if sprite.x > g.player.x {
			sprite.x -= 1
		}
		if sprite.y < g.player.y {
			sprite.y += 1
		}
		if sprite.y > g.player.y {
			sprite.y -= 1
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.player.x, g.player.y)
	screen.DrawImage(
		g.player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	for _, sprite := range g.enemies {
		opts.GeoM.Reset()
		opts.GeoM.Translate(sprite.x, sprite.y)
		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
	}

	for _, sprite := range g.potions {
		opts.GeoM.Reset()
		opts.GeoM.Translate(sprite.x, sprite.y)
		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	playerImg, _, err := ebitenutil.NewImageFromFile("./assets/images/ninja.png")
	if err != nil {
		log.Fatal(err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("./assets/images/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("./assets/images/potion.png")
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: &Player{
			Sprite: &Sprite{
				Img: playerImg,
				x:   100,
				y:   100,
			},
			Health: 0,
		},
		enemies: []*Enemy{
			{
				&Sprite{
					Img: skeletonImg,
					x:   200,
					y:   200,
				},
				true,
			},
			{
				&Sprite{
					Img: skeletonImg,
					x:   300,
					y:   100,
				},
				false,
			},
		},
		potions: []*Potion{
			{
				&Sprite{
					Img: potionImg,
					x:   100,
					y:   200,
				},
				10,
			},
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

// vim: ts=4
