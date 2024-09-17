package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"rpg-tutorial/animations"
	"rpg-tutorial/entities"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	camera      *Camera
	colliders   []image.Rectangle
}

func (g *Game) Update() error {

	g.player.Update()
	g.player.AdjustForColliders(g.colliders)

	// HACK: Få enemies att följa mig

	for _, enemy := range g.enemies {
		enemy.Update()
		enemy.AdjustForColliders(g.colliders)
	}

	// HACK:	Testar att få Potion
	for _, potion := range g.potions {
		if potion.HealingPower > 0 && g.player.Rect().Overlaps(potion.Rect()) {
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

	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, 16)
	enemySpriteSheet := spritesheet.NewSpriteSheet(4, 7, 16)
	potionSpriteSheet := spritesheet.NewSpriteSheet(1, 1, 16)

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img:         playerImg,
				X:           100,
				Y:           100,
				SpriteSheet: playerSpriteSheet,
			},
			Health: 0,
			Animations: map[entities.PlayerDirection]*animations.Animation{
				entities.Down:  animations.NewAnimation(0, 12, 4, 10.0),
				entities.Up:    animations.NewAnimation(1, 13, 4, 10.0),
				entities.Left:  animations.NewAnimation(2, 14, 4, 10.0),
				entities.Right: animations.NewAnimation(3, 15, 4, 10.0),
			},
		},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img:         skeletonImg,
					X:           200,
					Y:           200,
					SpriteSheet: enemySpriteSheet,
				},
				CanFollow: true,
			},
			{
				Sprite: &entities.Sprite{
					Img:         skeletonImg,
					X:           300,
					Y:           100,
					SpriteSheet: enemySpriteSheet,
				},
				CanFollow: false,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img:         potionImg,
					X:           100,
					Y:           200,
					SpriteSheet: potionSpriteSheet,
				},
				HealingPower: 10,
			},
			{
				Sprite: &entities.Sprite{
					Img:         potionImg,
					X:           50,
					Y:           300,
					SpriteSheet: potionSpriteSheet,
				},
				HealingPower: 40,
			},
			{
				Sprite: &entities.Sprite{
					Img:         potionImg,
					X:           310,
					Y:           150,
					SpriteSheet: potionSpriteSheet,
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

	//HACK: Måste sätta denna relation i efterhand eftersom de skapas i samma struct ovan.
	for _, enemy := range game.enemies {
		enemy.Player = game.player
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

// vim: ts=4
