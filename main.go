package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"rpg-tutorial/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	SCREEN_WIDTH  = 320
	SCREEN_HEIGHT = 240
)

type Game struct {
	player      *entities.Player
	enemies     []*entities.Enemy
	potions     []*entities.Potion
	tilemapJSON *TilemapJSON
	tilesets    []Tileset
	tilemapImg  *ebiten.Image
	camera      *Camera
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

	// HACK:	Testar att få Potion om man är nära
	for _, potion := range g.potions {
		if potion.HealingPower > 0 &&
			g.player.X+16 > potion.X &&
			g.player.X < potion.X+16 &&
			g.player.Y+16 > potion.Y &&
			g.player.Y < potion.Y+16 {
			g.player.Health += potion.HealingPower
			potion.HealingPower = 0
			fmt.Println("Got Health:", g.player.Health)
		}

	}

	// HACK:	För att korrekt följa spelaren, måste justera för origo från TopLeft till Center
	//			Hårdkodat: spelaren är 16 stor, så lägg till hälften: 8
	g.camera.FollowTarget(g.player.X+8, g.player.Y+8, SCREEN_WIDTH, SCREEN_HEIGHT)
	g.camera.Contrain(
		float64(g.tilemapJSON.Layers[0].Width)*16.0,
		float64(g.tilemapJSON.Layers[0].Height)*16.0,
		SCREEN_WIDTH,
		SCREEN_HEIGHT,
	)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	for layerIndex, layer := range g.tilemapJSON.Layers {
		// fmt.Printf("LayerIndex: %d, \nLayer: %T\n%v\n\n", layerIndex, layer, layer)
		for index, id := range layer.Data {

			if id == 0 {
				continue
			}
			x := index % layer.Width
			y := index / layer.Width
			x *= 16
			y *= 16

			img := g.tilesets[layerIndex].Img(id)

			// Move according to the current tile
			opts.GeoM.Translate(float64(x), float64(y))
			// Move according to the tiles anchor point (Top-Left instead of Bottom-Left)
			opts.GeoM.Translate(0.0, -float64(img.Bounds().Dy()+16))
			// Move according to the Camera
			opts.GeoM.Translate(g.camera.X, g.camera.Y)

			screen.DrawImage(img, &opts)
			opts.GeoM.Reset()
		}
	}

	// Move according to the Player
	opts.GeoM.Translate(g.player.X, g.player.Y)
	// Move according to the Camera
	opts.GeoM.Translate(g.camera.X, g.camera.Y)

	screen.DrawImage(
		g.player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	for _, sprite := range g.enemies {
		opts.GeoM.Reset()
		// Move according to the Enemy
		opts.GeoM.Translate(sprite.X, sprite.Y)
		// Move according to the Camera
		opts.GeoM.Translate(g.camera.X, g.camera.Y)
		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
	}

	for _, sprite := range g.potions {
		// HACK: Om poition är slut...rita inte
		if sprite.HealingPower <= 0 {
			continue
		}
		opts.GeoM.Reset()
		// Move according to the Potion
		opts.GeoM.Translate(sprite.X, sprite.Y)
		// Move according to the Camera
		opts.GeoM.Translate(g.camera.X, g.camera.Y)
		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
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
	//fmt.Printf("Tilemap (spawn): %v\n\n", tilemapJSON)

	tilesets, err := tilemapJSON.GenTilesets()
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
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   50,
					Y:   300,
				},
				HealingPower: 40,
			},
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   310,
					Y:   150,
				},
				HealingPower: 20,
			},
		},
		tilemapJSON: tilemapJSON,
		tilesets:    tilesets,
		tilemapImg:  tilemapImg,
		camera:      NewCamera(0, 0),
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

// vim: ts=4
