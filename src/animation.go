package main

import (
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Animation gère une spritesheet horizontale découpée en frames.
type Animation struct {
	texture     rl.Texture2D
	frameWidth  int32
	frameHeight int32
	frameCount  int
	current     int
	timer       float32
	frameTime   float32
	loop        bool
}

func NewAnimation(path string, frameWidth, frameHeight, frameCount int, frameTime float32) *Animation {
	return &Animation{
		texture:     rl.LoadTexture(path),
		frameWidth:  int32(frameWidth),
		frameHeight: int32(frameHeight),
		frameCount:  frameCount,
		frameTime:   frameTime,
		loop:        true,
	}
}

func (a *Animation) Update(delta float32) {
	if a == nil || a.frameCount <= 1 {
		return
	}
	if a.frameTime <= 0 {
		a.frameTime = 0.1
	}
	a.timer += delta
	for a.timer >= a.frameTime {
		a.timer -= a.frameTime
		a.current++
		if a.current >= a.frameCount {
			if a.loop {
				a.current = 0
			} else {
				a.current = a.frameCount - 1
			}
		}
	}
}

func (a *Animation) Reset() {
	if a == nil {
		return
	}
	a.current = 0
	a.timer = 0
}

func (a *Animation) Draw(x, y, width, height float32) {
	a.draw(x, y, width, height, false)
}

func (a *Animation) DrawFlipped(x, y, width, height float32) {
	a.draw(x, y, width, height, true)
}

func (a *Animation) draw(x, y, width, height float32, flipped bool) {
	if a == nil || a.texture.ID == 0 {
		return
	}

	sourceWidth := float32(a.frameWidth)
	if flipped {
		sourceWidth = -sourceWidth
	}

	source := rl.Rectangle{
		X:      float32(a.current) * float32(a.frameWidth),
		Y:      0,
		Width:  sourceWidth,
		Height: float32(a.frameHeight),
	}
	destination := rl.Rectangle{X: x, Y: y, Width: width, Height: height}
	rl.DrawTexturePro(a.texture, source, destination, rl.Vector2{}, 0, rl.White)
}

func (a *Animation) Unload() {
	if a != nil && a.texture.ID != 0 {
		rl.UnloadTexture(a.texture)
		a.texture = rl.Texture2D{}
	}
}

type PlayerAnimations struct {
	idle    *Animation
	walk    *Animation
	jump    *Animation
	attack1 *Animation
	attack2 *Animation
	attack3 *Animation
}

func initPlayerAnimations(classe string) PlayerAnimations {
	folder := filepath.Join("asset valide", classe)

	var prefix string
	switch classe {
	case "netrunner":
		prefix = "Punk_"
	case "merc":
		prefix = "Biker_"
	case "cyberpsycho":
		prefix = "Cyborg_"
	default:
		return PlayerAnimations{}
	}

	return PlayerAnimations{
		idle:    NewAnimation(filepath.Join(folder, prefix+"idle.png"), 48, 48, 4, 0.12),
		walk:    NewAnimation(filepath.Join(folder, "Walk.png"), 48, 48, 6, 0.09),
		jump:    NewAnimation(filepath.Join(folder, prefix+"jump.png"), 48, 48, 4, 0.12),
		attack1: NewAnimation(filepath.Join(folder, prefix+"attack1.png"), 48, 48, 6, 0.07),
		attack2: NewAnimation(filepath.Join(folder, prefix+"attack2.png"), 48, 48, 8, 0.07),
		attack3: NewAnimation(filepath.Join(folder, prefix+"attack3.png"), 48, 48, 8, 0.07),
	}
}

func (p *PlayerAnimations) Update(delta float32) {
	if p == nil {
		return
	}
	if p.idle != nil {
		p.idle.Update(delta)
	}
	if p.walk != nil {
		p.walk.Update(delta)
	}
	if p.jump != nil {
		p.jump.Update(delta)
	}
	if p.attack1 != nil {
		p.attack1.Update(delta)
	}
	if p.attack2 != nil {
		p.attack2.Update(delta)
	}
	if p.attack3 != nil {
		p.attack3.Update(delta)
	}
}

func (p *PlayerAnimations) Unload() {
	if p == nil {
		return
	}
	p.idle.Unload()
	p.walk.Unload()
	p.jump.Unload()
	p.attack1.Unload()
	p.attack2.Unload()
	p.attack3.Unload()
}

// initPNJAnimations charge les sprites Idle des trois PNJ de la carte.
func initPNJAnimations() map[string]*Animation {
	base := "asset valide"
	return map[string]*Animation{
		"MARCHAND":       NewAnimation(filepath.Join(base, "marchand", "Idle.png"), 48, 48, 4, 0.25),
		"CHARCUDOC":      NewAnimation(filepath.Join(base, "charcudoc", "Idle.png"), 48, 48, 6, 0.25),
		"GUERRE DE GANG": NewAnimation(filepath.Join(base, "guerre de gang", "Idle.png"), 72, 72, 4, 0.25),
	}
}

func unloadPNJAnimations(anims map[string]*Animation) {
	for _, animation := range anims {
		animation.Unload()
	}
}
