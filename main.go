package main

import (
	"image"
	"image/color"
	"log"
	"rpg-tutorial/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player  *entities.Player
	enemies []*entities.Enemy
	potions []*entities.Potion
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= 2
	}

	for _, sprite := range g.enemies {
		if !sprite.CanFollow {
			continue
		}
		if sprite.X < g.player.X {
			sprite.X += 1
		}
		if sprite.X > g.player.X {
			sprite.X -= 1
		}
		if sprite.Y < g.player.Y {
			sprite.Y += 1
		}
		if sprite.Y > g.player.Y {
			sprite.Y -= 1
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.player.X, g.player.Y)
	screen.DrawImage(
		g.player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	for _, sprite := range g.enemies {
		opts.GeoM.Reset()
		opts.GeoM.Translate(sprite.X, sprite.Y)
		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
	}

	for _, sprite := range g.potions {
		opts.GeoM.Reset()
		opts.GeoM.Translate(sprite.X, sprite.Y)
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
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   100,
				Y:   100,
			},
			Health: 0,
		},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   200,
					Y:   200,
				},
				CanFollow: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   300,
					Y:   100,
				},
				CanFollow: false,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   100,
					Y:   200,
				},
				HealingPower: 10,
			},
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

// vim: ts=4
