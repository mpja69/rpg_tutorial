package main

import (
	"encoding/json"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

type TilemapLayerJSON struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

type TilemapJSON struct {
	Layers   []TilemapLayerJSON `json:"layers"`
	Tilesets []map[string]any   `json:"tilesets"`
}

func (t *TilemapJSON) GenTilesets() ([]Tileset, error) {
	tilesets := make([]Tileset, 0)
	for _, tilesetData := range t.Tilesets {
		tilesetPath := path.Join("assets/maps", tilesetData["source"].(string))
		tileset, err := NewTileset(tilesetPath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, err
		}
		tilesets = append(tilesets, tileset)
	}
	return tilesets, nil
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var tilemapJSON TilemapJSON
	err = json.Unmarshal(content, &tilemapJSON)
	if err != nil {
		return nil, err
	}
	return &tilemapJSON, nil
}

func (tl *TilemapLayerJSON) Draw(screen *ebiten.Image, tileset Tileset, camera *Camera) {
	opts := ebiten.DrawImageOptions{}
	for index, id := range tl.Data {

		if id == 0 {
			continue
		}
		x := index % tl.Width
		y := index / tl.Width
		x *= 16
		y *= 16

		img := tileset.Img(id)

		// Move according to the current tile
		opts.GeoM.Translate(float64(x), float64(y))
		// Move according to the tiles anchor point (Top-Left instead of Bottom-Left)
		opts.GeoM.Translate(0.0, -float64(img.Bounds().Dy()+16))

		// Move according to the Camera
		opts.GeoM.Translate(camera.X, camera.Y)

		screen.DrawImage(img, &opts)
		opts.GeoM.Reset()
	}

}

// vim: ts=4
