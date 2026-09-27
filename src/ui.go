package main

import (
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UI struct {
	font        rl.Font
	frame       rl.Texture2D
	button      rl.Texture2D
	buttonHover rl.Texture2D
	logo        rl.Texture2D
	health      rl.Texture2D
	energy      rl.Texture2D
	loaded      bool
}

var ui UI

func initUI() {
	base := filepath.Join("asset valide", "Free-GUI-for-Cyberpunk-Pixel-Art1")

	ui.font = rl.LoadFont(filepath.Join(base, "10 Font", "CyberpunkCraftpixPixel.otf"))
	ui.frame = rl.LoadTexture(filepath.Join(base, "1 Frames", "FrameMap.png"))
	ui.button = rl.LoadTexture(filepath.Join(base, "6 Buttons", "Button_01.png"))
	ui.buttonHover = rl.LoadTexture(filepath.Join(base, "6 Buttons", "Button_02.png"))
	ui.logo = rl.LoadTexture(filepath.Join(base, "5 Logo", "Logo1.png"))
	ui.health = rl.LoadTexture(filepath.Join(base, "2 Bars", "HealthBar1.png"))
	ui.energy = rl.LoadTexture(filepath.Join(base, "2 Bars", "EnergyBar1.png"))

	rl.SetTextureFilter(ui.frame, rl.FilterPoint)
	rl.SetTextureFilter(ui.button, rl.FilterPoint)
	rl.SetTextureFilter(ui.buttonHover, rl.FilterPoint)
	rl.SetTextureFilter(ui.logo, rl.FilterPoint)
	rl.SetTextureFilter(ui.health, rl.FilterPoint)
	rl.SetTextureFilter(ui.energy, rl.FilterPoint)

	ui.loaded = true
}

func unloadUI() {
	if !ui.loaded {
		return
	}
	rl.UnloadFont(ui.font)
	rl.UnloadTexture(ui.frame)
	rl.UnloadTexture(ui.button)
	rl.UnloadTexture(ui.buttonHover)
	rl.UnloadTexture(ui.logo)
	rl.UnloadTexture(ui.health)
	rl.UnloadTexture(ui.energy)
	ui.loaded = false
}

func uiText(text string, x, y, size int, color rl.Color) {
	text = sansAccents(text)
	if ui.loaded && ui.font.Texture.ID != 0 {
		rl.DrawTextEx(ui.font, text, rl.Vector2{X: float32(x), Y: float32(y)}, float32(size), 1, color)
		return
	}
	rl.DrawText(text, int32(x), int32(y), int32(size), color)
}

func sansAccents(text string) string {
	r := strings.NewReplacer(
		"à", "a", "â", "a", "ä", "a", "À", "A", "Â", "A", "Ä", "A",
		"ç", "c", "Ç", "C",
		"é", "e", "è", "e", "ê", "e", "ë", "e", "É", "E", "È", "E", "Ê", "E", "Ë", "E",
		"î", "i", "ï", "i", "Î", "I", "Ï", "I",
		"ô", "o", "ö", "o", "Ô", "O", "Ö", "O",
		"ù", "u", "û", "u", "ü", "u", "Ù", "U", "Û", "U", "Ü", "U",
		"ÿ", "y", "Ÿ", "Y",
		"œ", "oe", "Œ", "OE",
		"æ", "ae", "Æ", "AE",
		"⚠", "!", "•", "-", "→", "->", "←", "<-", "↑", "^", "↓", "v",
	)
	return r.Replace(text)
}

func drawUIFrame(rect rl.Rectangle) {
	if !ui.loaded || ui.frame.ID == 0 {
		rl.DrawRectangleRec(rect, rl.NewColor(18, 20, 35, 245))
		rl.DrawRectangleLinesEx(rect, 2, rl.Gold)
		return
	}

	// Le FrameMap est un élément décoratif carré. On le répète dans le panneau
	// pour garder le rendu pixel-art au lieu d'étirer toute l'image.
	tile := float32(72)
	for y := rect.Y; y < rect.Y+rect.Height; y += tile {
		for x := rect.X; x < rect.X+rect.Width; x += tile {
			dst := rl.Rectangle{X: x, Y: y, Width: tile, Height: tile}
			src := rl.Rectangle{X: 0, Y: 0, Width: float32(ui.frame.Width), Height: float32(ui.frame.Height)}
			rl.DrawTexturePro(ui.frame, src, dst, rl.Vector2{}, 0, rl.White)
		}
	}

	rl.DrawRectangleRec(rl.Rectangle{X: rect.X + 8, Y: rect.Y + 8, Width: rect.Width - 16, Height: rect.Height - 16}, rl.NewColor(10, 13, 25, 235))
}

func drawUIPanel(title string) {
	rect := rl.Rectangle{X: 42, Y: 30, Width: 876, Height: 460}
	drawUIFrame(rect)

	// Bandeau supérieur inspiré des cadres du pack.
	rl.DrawRectangle(58, 48, 844, 48, rl.NewColor(40, 0, 35, 245))
	rl.DrawRectangleLinesEx(rl.Rectangle{X: 58, Y: 48, Width: 844, Height: 48}, 2, rl.Gold)
	uiText(title, 78, 58, 27, rl.Gold)
}

func drawUIButton(text string, rect rl.Rectangle, selected bool) {
	texture := ui.button
	if selected {
		texture = ui.buttonHover
	}

	if ui.loaded && texture.ID != 0 {
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(texture.Width), Height: float32(texture.Height)}
		rl.DrawTexturePro(texture, src, rect, rl.Vector2{}, 0, rl.White)
	} else {
		rl.DrawRectangleRec(rect, rl.NewColor(75, 0, 65, 255))
		rl.DrawRectangleLinesEx(rect, 2, rl.Gold)
	}

	uiText(text, int(rect.X+16), int(rect.Y+7), 20, rl.RayWhite)
}

func drawUIHealthBar(x, y, width int, current, max int) {
	if max <= 0 {
		max = 1
	}
	ratio := float32(current) / float32(max)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	if ui.loaded && ui.health.ID != 0 {
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(ui.health.Width), Height: float32(ui.health.Height)}
		rl.DrawTexturePro(ui.health, src, rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(width), Height: 18}, rl.Vector2{}, 0, rl.White)
		if ratio < 1 {
			rl.DrawRectangle(int32(x+int(float32(width)*ratio)), int32(y+2), int32(int(float32(width)*(1-ratio))), 14, rl.NewColor(10, 13, 25, 210))
		}
	} else {
		rl.DrawRectangle(int32(x), int32(y), int32(width), 18, rl.DarkGray)
		rl.DrawRectangle(int32(x), int32(y), int32(int(float32(width)*ratio)), 18, rl.Red)
	}
}

func drawUIEnergyBar(x, y, width int, current, max int) {
	if max <= 0 {
		max = 1
	}
	ratio := float32(current) / float32(max)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	if ui.loaded && ui.energy.ID != 0 {
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(ui.energy.Width), Height: float32(ui.energy.Height)}
		rl.DrawTexturePro(ui.energy, src, rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(width), Height: 18}, rl.Vector2{}, 0, rl.White)
		if ratio < 1 {
			rl.DrawRectangle(int32(x+int(float32(width)*ratio)), int32(y+2), int32(int(float32(width)*(1-ratio))), 14, rl.NewColor(10, 13, 25, 210))
		}
	} else {
		rl.DrawRectangle(int32(x), int32(y), int32(width), 18, rl.DarkGray)
		rl.DrawRectangle(int32(x), int32(y), int32(width), 18, rl.SkyBlue)
	}
}

func drawUIHint(text string, x, y int) {
	uiText(text, x, y, 16, rl.LightGray)
}

func drawUIPrompt(text string, x, y int) {
	pulse := 0.75 + 0.25*float32((rl.GetTime()-float64(int(rl.GetTime()))))
	size := 18 + int(pulse*2)
	uiText(text, x, y, size, rl.Gold)
}