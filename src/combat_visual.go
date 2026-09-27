package main

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// VisualCombat est le combat temps réel affiché dans la fenêtre Raylib.
// Il remplace le combat dans le terminal quand on joue depuis la carte.
type VisualCombat struct {
	mobs        []Monster
	boss        bool
	target      int
	menuMode    string
	pending     int
	hack        int
	inventory   int
	phase       string
	phaseTimer  float32
	hitDone     bool
	enemyIndex  int
	turn        int
	message     string
	result      int // 0 en cours, 1 victoire, -1 défaite, -2 fuite
	rewardGiven bool
	playerAnim  *PlayerAnimations
	facingLeft  bool
	flash       float32
	damageText  string
	damageTimer float32
	shake       float32
}

func newVisualCombat(perso *Character, boss bool, playerAnim *PlayerAnimations) *VisualCombat {
	c := &VisualCombat{
		boss:       boss,
		target:     0,
		menuMode:   "actions",
		phase:      "player",
		playerAnim: playerAnim,
		message:    "À toi de jouer !",
	}

	if boss {
		c.mobs = []Monster{initBoss1()}
	} else {
		count := rand.Intn(3) + 1
		c.mobs = make([]Monster, count)
		for i := range c.mobs {
			c.mobs[i] = initGoblin()
			c.mobs[i].name = fmt.Sprintf("GOBLIN %d", i+1)
		}
	}

	if len(c.mobs) > 0 {
		c.target = 0
	}
	return c
}

func (c *VisualCombat) aliveCount() int {
	n := 0
	for i := range c.mobs {
		if c.mobs[i].Hp() > 0 {
			n++
		}
	}
	return n
}

func (c *VisualCombat) aliveTarget() int {
	if c.target >= 0 && c.target < len(c.mobs) && c.mobs[c.target].Hp() > 0 {
		return c.target
	}
	for i := range c.mobs {
		if c.mobs[i].Hp() > 0 {
			return i
		}
	}
	return -1
}

func (c *VisualCombat) setMessage(s string) {
	c.message = s
}

func (c *VisualCombat) startAttackAnimation() {
	c.phase = "playerAttack"
	c.phaseTimer = 0
	c.hitDone = false
	if c.playerAnim == nil {
		return
	}
	var anim *Animation
	switch c.pending {
	case 2:
		anim = c.playerAnim.attack2
	case 3:
		anim = c.playerAnim.attack3
	default:
		anim = c.playerAnim.attack1
	}
	if anim != nil {
		anim.Reset()
	}
}

func (c *VisualCombat) performPlayerAttack(perso *Character, target int, attackType int) {
	if target < 0 || target >= len(c.mobs) || c.mobs[target].Hp() <= 0 {
		c.menuMode = "actions"
		return
	}

	damage := 0
	switch perso.classe {
	case "merc":
		if attackType == 1 {
			damage = perso.Atk() - c.mobs[target].Def()
			c.setMessage("Attaque au corps à corps !")
		} else {
			damage = weaponDamage(perso.weapon2) - c.mobs[target].Def()
			c.setMessage("Tir au pistolet !")
		}
	case "cyberpsycho":
		if attackType == 1 {
			damage = weaponDamage(perso.weapon) - c.mobs[target].Def()
			c.setMessage("Tir gauche !")
		} else if attackType == 2 {
			damage = weaponDamage(perso.weapon2) - c.mobs[target].Def()
			c.setMessage("Tir droit !")
		} else {
			damage = weaponDamage(perso.weapon) - c.mobs[target].Def()
			if damage < 1 {
				damage = 1
			}
			c.mobs[target].Dmg(damage)
			c.damageText = fmt.Sprintf("-%d", damage)
			c.damageTimer = 0.65
			if c.mobs[target].Hp() > 0 {
				second := weaponDamage(perso.weapon2) - c.mobs[target].Def()
				if second < 1 {
					second = 1
				}
				c.mobs[target].Dmg(second)
				damage += second
			}
			c.setMessage("DOUBLE TIR !")
			c.flash = 0.18
			c.shake = 0.12
			c.hitDone = true
			return
		}
	}

	if damage < 1 {
		damage = 1
	}
	c.mobs[target].Dmg(damage)
	c.damageText = fmt.Sprintf("-%d", damage)
	c.damageTimer = 0.65
	c.flash = 0.18
	c.shake = 0.12
}

func weaponDamage(w Arme) int {
	if w == nil {
		return 1
	}
	return w.Dmg()
}

func (c *VisualCombat) performHack(perso *Character, target int, hack int) {
	if target < 0 || target >= len(c.mobs) || c.mobs[target].Hp() <= 0 {
		return
	}

	var base int
	switch hack {
	case 1:
		base = 20
		c.setMessage("SURCHARGE : le système ennemi surchauffe !")
	case 2:
		base = 10
		c.setMessage("VIRUS : infection du système !")
	case 3:
		base = 12
		c.setMessage("COURT-CIRCUIT : décharge électrique !")
	}

	damage := base - c.mobs[target].Def()
	if damage < 1 {
		damage = 1
	}
	c.mobs[target].Dmg(damage)
	if hack == 2 && c.mobs[target].Hp() > 0 {
		// Le poison est visuel et non bloquant : pas de time.Sleep pendant un combat Raylib.
		poisonDamage := 5
		c.mobs[target].Dmg(poisonDamage)
		damage += poisonDamage
		c.setMessage("VIRUS : dégâts directs + poison !")
	}
	c.damageText = fmt.Sprintf("-%d", damage)
	c.damageTimer = 0.65
	c.flash = 0.18
	c.shake = 0.12
	_ = perso
}

func (c *VisualCombat) useInventory(perso *Character, index int) bool {
	if index < 0 || index >= len(perso.inventaire) {
		return false
	}

	object := perso.inventaire[index]
	switch item := object.(type) {
	case Armure:
		perso.equipArmor(item)
		c.setMessage(item.Nom() + " equipee.")
		return true
	case Item:
		switch item.nom {
		case "Steampack de basse qualité":
			old := perso.Hp()
			perso.AddPV(10)
			removeInventoryIndex(perso, index)
			c.setMessage(fmt.Sprintf("Steampack : +%d PV.", perso.Hp()-old))
			return true
		case "Steampack":
			old := perso.Hp()
			perso.AddPV(25)
			removeInventoryIndex(perso, index)
			c.setMessage(fmt.Sprintf("Steampack : +%d PV.", perso.Hp()-old))
			return true
		case "Grenade à fragmentation":
			target := c.aliveTarget()
			if target < 0 {
				return false
			}
			damage := 25 - c.mobs[target].Def()
			if damage < 1 {
				damage = 1
			}
			c.mobs[target].Dmg(damage)
			removeInventoryIndex(perso, index)
			c.damageText = fmt.Sprintf("-%d", damage)
			c.damageTimer = 0.65
			c.flash = 0.18
			c.shake = 0.12
			c.setMessage("GRENADE !")
			return true
		}
	}

	c.setMessage("Cet objet n'est pas utilisable pendant un combat.")
	return false
}

func (c *VisualCombat) update(perso *Character, delta float32) {
	if c == nil || c.result != 0 {
		return
	}

	if c.flash > 0 {
		c.flash -= delta
	}
	if c.shake > 0 {
		c.shake -= delta
	}
	if c.damageTimer > 0 {
		c.damageTimer -= delta
	}

	switch c.phase {
	case "player":
		c.updatePlayerInput(perso)
	case "playerAttack":
		c.phaseTimer += delta
		if !c.hitDone && c.phaseTimer >= 0.20 {
			c.hitDone = true
			target := c.aliveTarget()
			if target >= 0 {
				if c.hack > 0 {
					c.performHack(perso, target, c.hack)
				} else {
					c.performPlayerAttack(perso, target, c.pending)
				}
			}
		}
		if c.phaseTimer >= 0.55 {
			if c.aliveCount() == 0 {
				c.finishVictory(perso)
			} else {
				c.phase = "enemy"
				c.phaseTimer = 0
				c.enemyIndex = 0
			}
		}
	case "enemy":
		c.updateEnemyTurn(perso, delta)
	case "result":
		c.phaseTimer += delta
		if c.phaseTimer > 1.4 {
			if c.result == 1 || c.result == -1 || c.result == -2 {
				return
			}
		}
	}
}

func (c *VisualCombat) updatePlayerInput(perso *Character) {
	switch c.menuMode {
	case "actions":
		switch perso.classe {
		case "merc":
			if rl.IsKeyPressed(rl.KeyOne) {
				c.pending = 1
				c.hack = 0
				c.selectTargetOrAttack()
			} else if rl.IsKeyPressed(rl.KeyTwo) {
				c.pending = 2
				c.hack = 0
				c.selectTargetOrAttack()
			} else if rl.IsKeyPressed(rl.KeyThree) {
				c.menuMode = "inventory"
				c.inventory = 0
			} else if rl.IsKeyPressed(rl.KeyFour) {
				c.result = -2
				c.phase = "result"
				c.phaseTimer = 0
				c.setMessage("Tu prends la fuite.")
			}
		case "netrunner":
			if rl.IsKeyPressed(rl.KeyOne) {
				c.menuMode = "hacks"
			} else if rl.IsKeyPressed(rl.KeyTwo) {
				c.menuMode = "inventory"
				c.inventory = 0
			} else if rl.IsKeyPressed(rl.KeyThree) {
				c.result = -2
				c.phase = "result"
				c.phaseTimer = 0
				c.setMessage("Tu prends la fuite.")
			}
		case "cyberpsycho":
			if rl.IsKeyPressed(rl.KeyOne) || rl.IsKeyPressed(rl.KeyTwo) || rl.IsKeyPressed(rl.KeyThree) {
				c.pending = 1
				if rl.IsKeyPressed(rl.KeyTwo) {
					c.pending = 2
				}
				if rl.IsKeyPressed(rl.KeyThree) {
					c.pending = 3
				}
				c.hack = 0
				c.selectTargetOrAttack()
			} else if rl.IsKeyPressed(rl.KeyFour) {
				c.menuMode = "inventory"
				c.inventory = 0
			} else if rl.IsKeyPressed(rl.KeyFive) {
				c.result = -2
				c.phase = "result"
				c.phaseTimer = 0
				c.setMessage("Tu prends la fuite.")
			}
		}

	case "hacks":
		if rl.IsKeyPressed(rl.KeyOne) || rl.IsKeyPressed(rl.KeyTwo) || rl.IsKeyPressed(rl.KeyThree) {
			if rl.IsKeyPressed(rl.KeyOne) {
				c.hack = 1
			} else if rl.IsKeyPressed(rl.KeyTwo) {
				c.hack = 2
			} else {
				c.hack = 3
			}
			c.pending = 1
			c.selectTargetOrAttack()
		} else if rl.IsKeyPressed(rl.KeyFour) || rl.IsKeyPressed(rl.KeyEscape) {
			c.menuMode = "actions"
		}

	case "target":
		c.updateTargetInput()

	case "inventory":
		c.updateInventoryInput(perso)
	}
}

func (c *VisualCombat) selectTargetOrAttack() {
	if c.aliveCount() <= 1 {
		c.menuMode = "actions"
		c.startAttackAnimation()
		return
	}
	c.target = c.aliveTarget()
	c.menuMode = "target"
	c.setMessage("Choisis une cible : 1, 2 ou 3, puis Entrée.")
}

func (c *VisualCombat) updateTargetInput() {
	if rl.IsKeyPressed(rl.KeyLeft) {
		c.target--
		if c.target < 0 {
			c.target = len(c.mobs) - 1
		}
		c.target = c.aliveTargetNear(c.target, -1)
	}
	if rl.IsKeyPressed(rl.KeyRight) {
		c.target++
		if c.target >= len(c.mobs) {
			c.target = 0
		}
		c.target = c.aliveTargetNear(c.target, 1)
	}
	if rl.IsKeyPressed(rl.KeyOne) {
		c.target = 0
	}
	if rl.IsKeyPressed(rl.KeyTwo) && len(c.mobs) > 1 {
		c.target = 1
	}
	if rl.IsKeyPressed(rl.KeyThree) && len(c.mobs) > 2 {
		c.target = 2
	}
	if rl.IsKeyPressed(rl.KeyEnter) {
		if c.target >= 0 && c.target < len(c.mobs) && c.mobs[c.target].Hp() > 0 {
			c.menuMode = "actions"
			c.startAttackAnimation()
		}
	}
	if rl.IsKeyPressed(rl.KeyEscape) {
		c.menuMode = "actions"
	}
}

func (c *VisualCombat) aliveTargetNear(start, direction int) int {
	if len(c.mobs) == 0 {
		return -1
	}
	index := start
	for i := 0; i < len(c.mobs); i++ {
		if index < 0 {
			index = len(c.mobs) - 1
		}
		if index >= len(c.mobs) {
			index = 0
		}
		if c.mobs[index].Hp() > 0 {
			return index
		}
		index += direction
	}
	return -1
}

func (c *VisualCombat) updateInventoryInput(perso *Character) {
	if rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyFour) {
		c.menuMode = "actions"
		return
	}
	if rl.IsKeyPressed(rl.KeyUp) {
		c.inventory--
		if c.inventory < 0 {
			c.inventory = len(perso.inventaire) - 1
		}
	}
	if rl.IsKeyPressed(rl.KeyDown) {
		c.inventory++
		if c.inventory >= len(perso.inventaire) {
			c.inventory = 0
		}
	}
	if rl.IsKeyPressed(rl.KeyEnter) && len(perso.inventaire) > 0 {
		if c.useInventory(perso, c.inventory) {
			c.menuMode = "actions"
			c.phase = "enemy"
			c.phaseTimer = 0
			c.enemyIndex = 0
		}
	}
}

func (c *VisualCombat) updateEnemyTurn(perso *Character, delta float32) {
	c.phaseTimer += delta
	if c.enemyIndex >= len(c.mobs) {
		c.turn++
		c.phase = "player"
		c.phaseTimer = 0
		c.menuMode = "actions"
		c.setMessage("À toi de jouer !")
		return
	}

	mob := &c.mobs[c.enemyIndex]
	if mob.Hp() <= 0 {
		c.enemyIndex++
		c.phaseTimer = 0
		return
	}

	if c.phaseTimer < 0.25 {
		return
	}

	if c.phaseTimer < 0.55 {
		damage := mob.Atk() - perso.Def()
		if c.boss {
			if c.turn > 0 && c.turn%4 == 0 {
				mob.AddPV(3)
				c.setMessage("Le BOSS se regenere : +3 PV !")
				c.enemyIndex++
				c.phaseTimer = 0
				return
			}
			damage = mob.Atk() - perso.Def()
		}
		if damage < 1 {
			damage = 1
		}
		perso.Dmg(damage)
		c.damageText = fmt.Sprintf("-%d", damage)
		c.damageTimer = 0.65
		c.flash = 0.18
		c.shake = 0.12
		c.setMessage(mob.Name() + " attaque !")
		return
	}

	c.enemyIndex++
	c.phaseTimer = 0
	if perso.Hp() <= 0 {
		c.result = -1
		c.phase = "result"
		c.phaseTimer = 0
		c.setMessage("DÉFAITE")
	}
}

func (c *VisualCombat) finishVictory(perso *Character) {
	if c.rewardGiven {
		return
	}
	c.rewardGiven = true
	c.result = 1
	c.phase = "result"
	c.phaseTimer = 0

	totalMoney := 0
	for i := range c.mobs {
		if c.mobs[i].Hp() <= 0 {
			reward := c.mobs[i].moneyReward
			if reward <= 0 {
				reward = 20
			}
			totalMoney += reward
			c.mobs[i].exp_reward(perso)
		}
	}
	perso.money += totalMoney
	c.setMessage(fmt.Sprintf("VICTOIRE ! +%d po", totalMoney))
}

func (c *VisualCombat) draw(perso *Character) {
	beginVirtualDrawing(rl.NewColor(8, 10, 20, 255))

	// Fond de l'arène.
	rl.DrawRectangle(0, 0, 960, 540, rl.NewColor(16, 18, 35, 255))
	for x := int32(0); x < 960; x += 48 {
		rl.DrawLine(x, 150, x, 500, rl.NewColor(27, 30, 52, 255))
	}
	for y := int32(150); y < 500; y += 48 {
		rl.DrawLine(0, y, 960, y, rl.NewColor(27, 30, 52, 255))
	}

	title := "COMBAT"
	if c.boss {
		title = "COMBAT CONTRE LE BOSS"
	}
	uiText(title, 42, 28, 26, rl.Gold)
	uiText(fmt.Sprintf("TOUR %d", c.turn+1), 800, 32, 18, rl.LightGray)

	// Joueur.
	drawUIHealthBar(45, 72, 300, perso.Hp(), perso.MaxHp())
	uiText(perso.nom+" | "+perso.classe, 45, 98, 18, rl.RayWhite)
	drawCombatPlayer(c, 120, 190)

	// Ennemis.
	for i := range c.mobs {
		x := float32(500 + i*135)
		if len(c.mobs) == 1 {
			x = 650
		}
		drawCombatEnemy(&c.mobs[i], i == c.target && c.menuMode == "target", x, 215, c.boss, c.phaseTimer)
	}

	// Message.
	if c.message != "" {
		drawUIPrompt(c.message, 42, 145)
	}

	// Menus.
	drawCombatMenu(c, perso)

	if c.damageTimer > 0 && c.damageText != "" {
		uiText(c.damageText, 650, 180, 24, rl.Red)
	}

	if c.result == 1 {
		drawCombatResult("VICTOIRE", rl.Gold)
	} else if c.result == -1 {
		drawCombatResult("DÉFAITE", rl.Red)
	} else if c.result == -2 {
		drawCombatResult("FUITE", rl.LightGray)
	}

	drawUIHint("Combat visuel • touches 1-5 • ↑/↓ inventaire • ←/→ cible • Entrée valider", 42, 510)
	endVirtualDrawing()
}

func drawCombatPlayer(c *VisualCombat, x, y float32) {
	if c.playerAnim != nil {
		anim := c.playerAnim.idle
		if c.phase == "playerAttack" {
			switch c.pending {
			case 2:
				anim = c.playerAnim.attack2
			case 3:
				anim = c.playerAnim.attack3
			default:
				anim = c.playerAnim.attack1
			}
		}
		if anim != nil && anim.texture.ID != 0 {
			anim.Draw(x, y, 112, 112)
			return
		}
	}
	rl.DrawRectangle(int32(x+30), int32(y+10), 48, 88, rl.SkyBlue)
}

func drawCombatEnemy(mob *Monster, selected bool, x, y float32, boss bool, timer float32) {
	if mob.Hp() <= 0 {
		uiText("X", int(x+38), int(y+35), 40, rl.DarkGray)
		return
	}

	bob := float32(0)
	if timer > 0 {
		bob = float32((int(timer*8) % 2) * 3)
	}
	if boss {
		// Boss animé sans nouveau fichier : respiration + oscillation.
		rl.DrawRectangle(int32(x+20), int32(y+10+bob), 82, 120, rl.NewColor(90, 25, 70, 255))
		rl.DrawRectangle(int32(x+35), int32(y-12+bob), 52, 35, rl.NewColor(130, 35, 80, 255))
		rl.DrawRectangle(int32(x+46), int32(y-2+bob), 8, 8, rl.Gold)
		rl.DrawRectangle(int32(x+68), int32(y-2+bob), 8, 8, rl.Gold)
		uiText("BOSS", int(x+30), int(y+140), 18, rl.Red)
	} else {
		// Petit gobelin pixel-art animé : flottement + yeux.
		rl.DrawRectangle(int32(x+35), int32(y+35+bob), 55, 70, rl.NewColor(60, 115, 85, 255))
		rl.DrawRectangle(int32(x+25), int32(y+20+bob), 75, 32, rl.NewColor(80, 145, 95, 255))
		rl.DrawRectangle(int32(x+38), int32(y+30+bob), 8, 8, rl.Red)
		rl.DrawRectangle(int32(x+78), int32(y+30+bob), 8, 8, rl.Red)
		uiText(fmt.Sprintf("%d", int(mob.Hp())), int(x+48), int(y+118), 15, rl.RayWhite)
	}

	drawUIHealthBar(int(x), int(y+160), 115, mob.Hp(), mob.MaxHp())
	uiText(mob.Name(), int(x), int(y+182), 14, rl.RayWhite)

	if selected {
		rl.DrawRectangleLinesEx(rl.Rectangle{X: x - 8, Y: y - 22, Width: 130, Height: 220}, 3, rl.Gold)
		uiText("CIBLE", int(x+25), int(y-45), 15, rl.Gold)
	}
}

func drawCombatMenu(c *VisualCombat, perso *Character) {
	panel := rl.Rectangle{X: 42, Y: 355, Width: 876, Height: 135}
	drawUIFrame(panel)

	switch c.menuMode {
	case "actions":
		switch perso.classe {
		case "merc":
			drawUIButton("1. CAC", rl.Rectangle{X: 65, Y: 375, Width: 180, Height: 38}, false)
			drawUIButton("2. PISTOLET", rl.Rectangle{X: 255, Y: 375, Width: 180, Height: 38}, false)
			drawUIButton("3. INVENTAIRE", rl.Rectangle{X: 445, Y: 375, Width: 180, Height: 38}, false)
			drawUIButton("4. FUIR", rl.Rectangle{X: 635, Y: 375, Width: 180, Height: 38}, false)
		case "netrunner":
			drawUIButton("1. HACKS", rl.Rectangle{X: 100, Y: 375, Width: 220, Height: 38}, false)
			drawUIButton("2. INVENTAIRE", rl.Rectangle{X: 350, Y: 375, Width: 220, Height: 38}, false)
			drawUIButton("3. FUIR", rl.Rectangle{X: 600, Y: 375, Width: 180, Height: 38}, false)
		case "cyberpsycho":
			drawUIButton("1. GAUCHE", rl.Rectangle{X: 55, Y: 375, Width: 160, Height: 38}, false)
			drawUIButton("2. DROIT", rl.Rectangle{X: 225, Y: 375, Width: 160, Height: 38}, false)
			drawUIButton("3. DOUBLE", rl.Rectangle{X: 395, Y: 375, Width: 160, Height: 38}, false)
			drawUIButton("4. INV", rl.Rectangle{X: 565, Y: 375, Width: 130, Height: 38}, false)
			drawUIButton("5. FUIR", rl.Rectangle{X: 705, Y: 375, Width: 130, Height: 38}, false)
		}
	case "hacks":
		drawUIButton("1. SURCHARGE", rl.Rectangle{X: 70, Y: 375, Width: 220, Height: 38}, false)
		drawUIButton("2. VIRUS", rl.Rectangle{X: 310, Y: 375, Width: 220, Height: 38}, false)
		drawUIButton("3. COURT-CIRCUIT", rl.Rectangle{X: 550, Y: 375, Width: 220, Height: 38}, false)
		uiText("4. Retour", 75, 430, 16, rl.LightGray)
	case "target":
		uiText("CHOISIR UNE CIBLE", 70, 375, 20, rl.Gold)
		uiText("1/2/3 ou ←/→ pour sélectionner • Entrée pour attaquer • Échap pour annuler", 70, 410, 16, rl.RayWhite)
	case "inventory":
		uiText("INVENTAIRE DE COMBAT", 70, 365, 20, rl.Gold)
		if len(perso.inventaire) == 0 {
			uiText("Inventaire vide", 70, 400, 17, rl.LightGray)
		} else {
			start := c.inventory
			if start < 0 {
				start = 0
			}
			if start >= len(perso.inventaire) {
				start = len(perso.inventaire) - 1
			}
			for i := 0; i < 4 && start+i < len(perso.inventaire); i++ {
				index := start + i
				prefix := "  "
				if index == c.inventory {
					prefix = "> "
				}
				uiText(prefix+perso.inventaire[index].Nom(), 70, 395+i*22, 14, rl.RayWhite)
			}
			uiText("↑/↓ sélectionner • Entrée utiliser • Échap retour", 600, 430, 13, rl.LightGray)
		}
	}
}

func drawCombatResult(text string, color rl.Color) {
	rl.DrawRectangle(270, 175, 420, 120, rl.NewColor(8, 10, 20, 235))
	rl.DrawRectangleLinesEx(rl.Rectangle{X: 270, Y: 175, Width: 420, Height: 120}, 3, color)
	uiText(text, 390, 205, 38, color)
	uiText("Retour à la carte...", 385, 255, 16, rl.RayWhite)
}
