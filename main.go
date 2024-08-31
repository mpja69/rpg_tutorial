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
	player      *entities.Player
	enemies     []*entities.Enemy
	potions     []*entities.Potion
	tilemapJSON *TilemapJSON
	tilemapImg  *ebiten.Image
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += 2.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= 2.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += 2.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= 2.0
	}

	for _, sprite := range g.enemies {
		if !sprite.CanFollow {
			continue
		}
		if sprite.X < g.player.X {
			sprite.X += 1.0
		}
		if sprite.X > g.player.X {
			sprite.X -= 1.0
		}
		if sprite.Y < g.player.Y {
			sprite.Y += 1.0
		}
		if sprite.Y > g.player.Y {
			sprite.Y -= 1.0
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	for _, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			x := index % layer.Width
			y := index / layer.Width
			x *= 16
			y *= 16

			srcX := (id - 1) % 22
			srcY := (id - 1) / 22
			srcX *= 16
			srcY *= 16

			opts.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(
				g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)
			opts.GeoM.Reset()
		}
	}

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

	tilemapJSON, err := NewTilemapJSON("./assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	tilemapImg, _, err := ebitenutil.NewImageFromFile("./assets/images/TilesetFloor.png")
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
		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImg,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

// vim: ts=4
