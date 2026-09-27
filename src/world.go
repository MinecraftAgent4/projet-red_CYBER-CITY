package main

import (
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type WorldArt struct {
	bg      [5]rl.Texture2D
	overlay [30]rl.Texture2D
	tiles   map[int]rl.Texture2D
	deco    [8]rl.Texture2D
	loaded  bool
}

var worldArt WorldArt

func initWorldArt() {
	root := "asset valide"
	bgRoot := filepath.Join(root, "Power-Station-Free-Tileset-Pixel-Art", "2 Background", "Night")
	overlayRoot := filepath.Join(root, "free-cyberpunk-overlay-effects-for-platformer-game", "Overlay")
	tileRoot := filepath.Join(root, "Power-Station-Free-Tileset-Pixel-Art", "1 Tiles")
	decoRoot := filepath.Join(root, "Power-Station-Free-Tileset-Pixel-Art", "3 Objects", "2 Decoration")

	for i := 0; i < 5; i++ {
		worldArt.bg[i] = rl.LoadTexture(filepath.Join(bgRoot, itoa(i+1)+".png"))
		rl.SetTextureFilter(worldArt.bg[i], rl.FilterPoint)
	}

	for i := 0; i < 30; i++ {
		worldArt.overlay[i] = rl.LoadTexture(filepath.Join(overlayRoot, itoa(i+1)+".png"))
		rl.SetTextureFilter(worldArt.overlay[i], rl.FilterBilinear)
	}

	worldArt.tiles = make(map[int]rl.Texture2D)
	// 46 = surface/floor sombre, 03 = bordure de plateforme, 61 = dalle métallique.
	for _, id := range []int{3, 46, 61} {
		worldArt.tiles[id] = rl.LoadTexture(filepath.Join(tileRoot, "Tile_"+format2(id)+".png"))
		rl.SetTextureFilter(worldArt.tiles[id], rl.FilterPoint)
	}

	// Décor : machines, cuves, consoles et structure métallique.
	for i, id := range []int{26, 24, 27, 25, 16, 14, 20, 19} {
		worldArt.deco[i] = rl.LoadTexture(filepath.Join(decoRoot, itoa(id)+".png"))
		rl.SetTextureFilter(worldArt.deco[i], rl.FilterPoint)
	}

	worldArt.loaded = true
}

func unloadWorldArt() {
	if !worldArt.loaded {
		return
	}
	for i := range worldArt.bg {
		rl.UnloadTexture(worldArt.bg[i])
	}
	for i := range worldArt.overlay {
		rl.UnloadTexture(worldArt.overlay[i])
	}
	for _, tex := range worldArt.tiles {
		rl.UnloadTexture(tex)
	}
	for i := range worldArt.deco {
		rl.UnloadTexture(worldArt.deco[i])
	}
	worldArt.loaded = false
}

func format2(n int) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}

func drawWorldBackground() {
	if !worldArt.loaded {
		rl.ClearBackground(rl.NewColor(8, 10, 20, 255))
		return
	}

	dst := rl.Rectangle{X: 0, Y: 0, Width: 960, Height: 540}
	for i := 0; i < len(worldArt.bg); i++ {
		tex := worldArt.bg[i]
		if tex.ID == 0 {
			continue
		}
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(tex.Width), Height: float32(tex.Height)}
		rl.DrawTexturePro(tex, src, dst, rl.Vector2{}, 0, rl.White)
	}
}

func drawPlatformTexture(rect rl.Rectangle) {
	if !worldArt.loaded {
		rl.DrawRectangleRec(rect, rl.DarkGreen)
		return
	}

	floor := worldArt.tiles[46]
	top := worldArt.tiles[3]
	if floor.ID == 0 || top.ID == 0 {
		rl.DrawRectangleRec(rect, rl.NewColor(20, 23, 35, 255))
		return
	}

	// Remplissage avec la tuile de sol.
	for x := rect.X; x < rect.X+rect.Width; x += 32 {
		w := float32(32)
		if x+w > rect.X+rect.Width {
			w = rect.X + rect.Width - x
		}
		src := rl.Rectangle{X: 0, Y: 0, Width: 32, Height: 32}
		dst := rl.Rectangle{X: x, Y: rect.Y, Width: w, Height: rect.Height}
		rl.DrawTexturePro(floor, src, dst, rl.Vector2{}, 0, rl.White)
	}

	// Bordure métallique en haut.
	for x := rect.X; x < rect.X+rect.Width; x += 32 {
		w := float32(32)
		if x+w > rect.X+rect.Width {
			w = rect.X + rect.Width - x
		}
		src := rl.Rectangle{X: 0, Y: 0, Width: 32, Height: 32}
		dst := rl.Rectangle{X: x, Y: rect.Y - 5, Width: w, Height: 32}
		rl.DrawTexturePro(top, src, dst, rl.Vector2{}, 0, rl.White)
	}
}

func drawWorldDecor() {
	if !worldArt.loaded {
		return
	}

	drawWorldObjectOnGround(worldArt.deco[0], 140, 355, 1.35)
	drawWorldObjectOnGround(worldArt.deco[1], 285, 445, 1.1)
	drawWorldObjectOnGround(worldArt.deco[2], 500, 445, 1.15)
	drawWorldObjectOnGround(worldArt.deco[3], 745, 355, 1.05)
	drawWorldObject(worldArt.deco[4], 400, 250, 1.4)
	drawWorldObject(worldArt.deco[5], 830, 183, 1.25)
	drawWorldObjectOnGround(worldArt.deco[6], 870, 445, 1.3)
	drawWorldObjectOnGround(worldArt.deco[7], 690, 445, 1.25)
}

func drawWorldObjectOnGround(tex rl.Texture2D, x, groundY, scale float32) {
	if tex.ID == 0 {
		return
	}

	height := float32(tex.Height) * scale
	drawWorldObject(tex, x, groundY-height, scale)
}

func drawWorldObject(tex rl.Texture2D, x, y, scale float32) {
	if tex.ID == 0 {
		return
	}
	src := rl.Rectangle{X: 0, Y: 0, Width: float32(tex.Width), Height: float32(tex.Height)}
	dst := rl.Rectangle{X: x, Y: y, Width: float32(tex.Width) * scale, Height: float32(tex.Height) * scale}
	rl.DrawTexturePro(tex, src, dst, rl.Vector2{}, 0, rl.White)
}

func drawWorldOverlay() {
	if !worldArt.loaded {
		return
	}
	// Les overlays forment une animation très légère de pluie/lumière/scanlines.
	frame := int(rl.GetTime()) % len(worldArt.overlay)
	tex := worldArt.overlay[frame]
	if tex.ID == 0 {
		return
	}
	src := rl.Rectangle{X: 0, Y: 0, Width: float32(tex.Width), Height: float32(tex.Height)}
	dst := rl.Rectangle{X: 0, Y: 0, Width: 960, Height: 540}
	rl.DrawTexturePro(tex, src, dst, rl.Vector2{}, 0, rl.NewColor(255, 255, 255, 42))
}
