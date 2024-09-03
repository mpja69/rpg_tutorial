package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"rpg-tutorial/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	SCREEN_WIDTH  = 320
	SCREEN_HEIGHT = 240
)

func CheckCollisionX(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+16, int(sprite.Y)+16)) {
			if sprite.Dx > 0 {
				sprite.X = float64(collider.Min.X) - 16
			} else if sprite.Dx < 0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}
func CheckCollisionY(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+16, int(sprite.Y)+16)) {
			if sprite.Dy > 0 {
				sprite.Y = float64(collider.Min.Y) - 16
			} else if sprite.Dy < 0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}

type Game struct {
	player      *entities.Player
	enemies     []*entities.Enemy
	potions     []*entities.Potion
	tilemapJSON *TilemapJSON
	tilesets    []Tileset
	// tilemapImg  *ebiten.Image
	camera    *Camera
	colliders []image.Rectangle
}

func (g *Game) Update() error {

	g.player.Dx = 0
	g.player.Dy = 0
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Dx = 2.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Dx = -2.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Dy = 2.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Dy = -2.0
	}

	g.player.X += g.player.Dx
	CheckCollisionX(g.player.Sprite, g.colliders)

	g.player.Y += g.player.Dy
	CheckCollisionY(g.player.Sprite, g.colliders)

	// HACK: Få enemies att följa mig

	for _, enemy := range g.enemies {
		enemy.Dx = 0
		enemy.Dy = 0

		if !enemy.CanFollow {
			continue
		}
		if enemy.X < g.player.X {
			enemy.Dx += 1.0
		}
		if enemy.X > g.player.X {
			enemy.Dx -= 1.0
		}
		if enemy.Y < g.player.Y {
			enemy.Dy += 1.0
		}
		if enemy.Y > g.player.Y {
			enemy.Dy -= 1.0
		}

		enemy.X += enemy.Dx
		CheckCollisionX(enemy.Sprite, g.colliders)

		enemy.Y += enemy.Dy
		CheckCollisionY(enemy.Sprite, g.colliders)
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

	// Draw layers in the map
	for layerIndex, layer := range g.tilemapJSON.Layers {
		layer.Draw(screen, g.tilesets[layerIndex], g.camera)
		// layer.Draw(screen, g.tilesets[layerIndex], func(o *ebiten.DrawImageOptions) {
		// 	o.GeoM.Translate(g.camera.X, g.camera.Y)
		// })
	}

	// Draw Player
	g.player.Draw(screen, func(o *ebiten.DrawImageOptions) {
		o.GeoM.Translate(g.camera.X, g.camera.Y)
	})

	// Draw Enemies
	for _, enemy := range g.enemies {
		enemy.Draw(screen, func(o *ebiten.DrawImageOptions) {
			o.GeoM.Translate(g.camera.X, g.camera.Y)
		})
	}

	// Draw Potions
	for _, potion := range g.potions {
		if potion.HealingPower <= 0 {
			continue
		}
		potion.Draw(screen, func(o *ebiten.DrawImageOptions) {
			o.GeoM.Translate(g.camera.X, g.camera.Y)
		})
	}

	// Draw Colliders
	for _, collider := range g.colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(g.camera.X),
			float32(collider.Min.Y)+float32(g.camera.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			true,
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

	// tilemapImg, _, err := ebitenutil.NewImageFromFile("./assets/images/TilesetFloor.png")
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
		// tilemapImg:  tilemapImg,
		camera: NewCamera(0, 0),
		colliders: []image.Rectangle{
			image.Rect(120, 120, 136, 136),
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

// vim: ts=4
