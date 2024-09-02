package main

import (
	"encoding/json"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type UniformTilesetJSON struct {
	Path string `json:"image"`
}

type UniformTileset struct {
	img *ebiten.Image
	gid int
}

func (u *UniformTileset) Img(id int) *ebiten.Image {
	id -= u.gid
	// Convert to position in the tileset
	srcX := id % 22
	srcY := id / 22
	// Convert to pixel coordinates
	srcX *= 16
	srcY *= 16

	return u.img.SubImage(
		image.Rect(
			srcX, srcY, srcX+16, srcY+16,
		),
	).(*ebiten.Image)
}

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

type DynTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

type DynTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d *DynTileset) Img(id int) *ebiten.Image {
	id -= d.gid
	return d.imgs[id]
}

// HACK: Denna borde göras på ett annat sätt än att kolla path?!
func NewTileset(path string, gid int) (Tileset, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if strings.Contains(path, "buildings") {
		var dynTilesetJSON DynTilesetJSON
		err = json.Unmarshal(content, &dynTilesetJSON)
		if err != nil {
			return nil, err
		}
		dynTileset := DynTileset{}
		dynTileset.gid = gid
		dynTileset.imgs = make([]*ebiten.Image, 0)

		for _, tileJSON := range dynTilesetJSON.Tiles {

			tileJSONPath := tileJSON.Path
			tileJSONPath = strings.Trim(tileJSONPath, "../")
			tileJSONPath = filepath.Join("assets/", tileJSONPath)

			img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
			if err != nil {
				return nil, err
			}
			dynTileset.imgs = append(dynTileset.imgs, img)
		}

		return &dynTileset, nil
	} else {

		var uniformTilesetJSON UniformTilesetJSON
		err = json.Unmarshal(content, &uniformTilesetJSON)
		if err != nil {
			return nil, err
		}
		uniformTileset := UniformTileset{}

		uniformTilesetJSONPath := uniformTilesetJSON.Path
		uniformTilesetJSONPath = strings.Trim(uniformTilesetJSONPath, "../")
		uniformTilesetJSONPath = filepath.Join("assets/", uniformTilesetJSONPath)
		img, _, err := ebitenutil.NewImageFromFile(uniformTilesetJSONPath)
		if err != nil {
			return nil, err
		}
		uniformTileset.img = img
		uniformTileset.gid = gid
		return &uniformTileset, nil
	}
}

// vim: ts=4
