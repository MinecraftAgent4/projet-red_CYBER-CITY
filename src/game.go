package main

import (
	"fmt"
	"strings"
	"unicode"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type mapBlock struct {
	rect  rl.Rectangle
	color rl.Color
}

type mapCharacter struct {
	name      string
	rect      rl.Rectangle
	color     rl.Color
	animation *Animation
}


func game() {
	rl.InitWindow(1440, 810, "RedProject")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()

	initVirtualRenderer()
	defer unloadVirtualRenderer()
	initUI()
	initWorldArt()
	defer unloadUI()
	defer unloadWorldArt()

	marketObjets := initMarket()
	charcudocObjets := initCharcudoc()
	pnjAnimations := initPNJAnimations()
	defer unloadPNJAnimations(pnjAnimations)

	perso := initCharacter(
		"Joueur",
		"",
		1,
		100,
		100,
		[]Object{
			Resource{nom: "Ferraille", quantité: 10, quantité_max: 99},
			Resource{nom: "Composants", quantité: 4, quantité_max: 99},
			Resource{nom: "Poudre", quantité: 5, quantité_max: 99},
		},
		100,
	)

	screen := "classe"
	nameInput := ""
	selectedClass := ""

	blocks := []mapBlock{
		{rect: rl.Rectangle{X: 0, Y: 450, Width: 960, Height: 90}, color: rl.DarkGreen},
		{rect: rl.Rectangle{X: 80, Y: 360, Width: 180, Height: 24}, color: rl.Green},
		{rect: rl.Rectangle{X: 350, Y: 290, Width: 180, Height: 24}, color: rl.Green},
		{rect: rl.Rectangle{X: 630, Y: 360, Width: 180, Height: 24}, color: rl.Green},
		{rect: rl.Rectangle{X: 820, Y: 220, Width: 100, Height: 24}, color: rl.Green},
	}

	characters := []mapCharacter{
		{name: "MARCHAND", rect: rl.Rectangle{X: 180, Y: 308, Width: 48, Height: 48}, color: rl.Gold, animation: pnjAnimations["MARCHAND"]},
		{name: "CHARCUDOC", rect: rl.Rectangle{X: 411, Y: 238, Width: 48, Height: 48}, color: rl.Purple, animation: pnjAnimations["CHARCUDOC"]},
		{name: "GUERRE DE GANG", rect: rl.Rectangle{X: 644, Y: 288, Width: 72, Height: 72}, color: rl.Red, animation: pnjAnimations["GUERRE DE GANG"]},
	}

	player := rl.Rectangle{X: 40, Y: 408, Width: 32, Height: 42}
	playerSpeed := float32(220)
	playerFacingLeft := false
	playerAnims := PlayerAnimations{}
	defer playerAnims.Unload()
	var visualCombat *VisualCombat
	velocityY := float32(0)
	onGround := true
	message := "Explore la map — M pour ouvrir le menu"

	for !rl.WindowShouldClose() {
		delta := rl.GetFrameTime()

		for _, character := range characters {
			if character.animation != nil {
				character.animation.Update(delta)
			}
		}
		playerAnims.Update(delta)
		if screen == "classe" {
			if rl.IsKeyPressed(rl.KeyOne) {
				selectedClass = "netrunner"
				perso.classe = selectedClass
				playerAnims.Unload()
				playerAnims = initPlayerAnimations(selectedClass)
				perso.pv_total = 80
				perso.pv_actuelle = 80
				perso.weapon = Ranged{nom: "Pistolet neural", dmg: 8}
				perso.weapon2 = Ranged{nom: "Pistolet-mitrailleur", dmg: 10}
				screen = "nom"
			}
			if rl.IsKeyPressed(rl.KeyTwo) {
				selectedClass = "merc"
				perso.classe = selectedClass
				playerAnims.Unload()
				playerAnims = initPlayerAnimations(selectedClass)
				perso.pv_total = 100
				perso.pv_actuelle = 100
				perso.weapon = Melee{nom: "Lame cyber", dmg: 12}
				perso.weapon2 = Ranged{nom: "Pistolet lourd", dmg: 9}
				screen = "nom"
			}
			if rl.IsKeyPressed(rl.KeyThree) {
				selectedClass = "cyberpsycho"
				perso.classe = selectedClass
				playerAnims.Unload()
				playerAnims = initPlayerAnimations(selectedClass)
				perso.pv_total = 120
				perso.pv_actuelle = 120
				perso.weapon = Ranged{nom: "Canon gauche", dmg: 14}
				perso.weapon2 = Ranged{nom: "Canon droit", dmg: 16}
				screen = "nom"
			}
			drawClassScreen(selectedClass)
			continue
		}

		if screen == "nom" {
			nameInput = readNameInput(nameInput)

			if rl.IsKeyPressed(rl.KeyEnter) {
				if validName(nameInput) {
					perso.nom = formatName(nameInput)
					screen = "map"
					message = "Bienvenue " + perso.nom + " ! — M pour ouvrir le menu"
				}
			}

			beginVirtualDrawing(rl.NewColor(28, 37, 61, 255))
			drawPanel("CRÉATION DU PERSONNAGE")
			rl.DrawText("Choisis ton nom :", 90, 150, 24, rl.RayWhite)
			rl.DrawRectangle(90, 200, 780, 55, rl.DarkGray)
			rl.DrawRectangleLines(90, 200, 780, 55, rl.Gold)
			rl.DrawText(nameInput+"_",
				110, 215, 24, rl.RayWhite)
			rl.DrawText("Entrée : valider | Retour arrière : effacer",
				90, 300, 18, rl.LightGray)
			rl.DrawText("Classe : "+selectedClass, 90, 340, 20, rl.Gold)
			endVirtualDrawing()
			continue
		}

		if screen == "combat" {
			if visualCombat == nil {
				screen = "gang"
				continue
			}
			visualCombat.update(&perso, delta)
			visualCombat.draw(&perso)
			if visualCombat.result != 0 {
				if visualCombat.phaseTimer > 1.4 {
					if visualCombat.result == -1 {
						perso.isDead()
					}
					if visualCombat.result == 1 {
						message = "Victoire ! Retour à la carte."
					} else if visualCombat.result == -2 {
						message = "Tu as fui le combat."
					} else {
						message = "Tu as été vaincu."
					}
					visualCombat = nil
					screen = "map"
				}
			}
			continue
		}

		if screen != "map" {
			handlePanelInput(&screen, &perso, &marketObjets, &charcudocObjets, &message, &visualCombat, &playerAnims)
			drawGamePanel(screen, &perso, &marketObjets, &charcudocObjets, message)
			continue
		}

		if rl.IsKeyPressed(rl.KeyM) || rl.IsKeyPressed(rl.KeySemicolon) {
			screen = "menu"
			continue
		}
		direction := float32(0)
		if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
			direction = 1
		}
		if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
			direction = -1
		}
		if direction < 0 {
			playerFacingLeft = true
		} else if direction > 0 {
			playerFacingLeft = false
		}

		if rl.IsKeyPressed(rl.KeySpace) && onGround {
			velocityY = -520
			onGround = false
		}

		// Gravité
		velocityY += 1000 * delta

		// Ancienne position
		previousY := player.Y

		// Déplacement vertical
		player.Y += velocityY * delta
		onGround = false

		// Collision avec le dessus des plateformes
		if velocityY >= 0 {
			for _, block := range blocks {
				// Le joueur doit être au-dessus de la plateforme
				// et traverser son niveau pendant cette frame.
				previousBottom := previousY + player.Height
				currentBottom := player.Y + player.Height

				if previousBottom <= block.rect.Y &&
					currentBottom >= block.rect.Y &&
					player.X+player.Width > block.rect.X &&
					player.X < block.rect.X+block.rect.Width {

					player.Y = block.rect.Y - player.Height
					velocityY = 0
					onGround = true
					break
				}
			}
		}

		moveHorizontal(&player, direction*playerSpeed*delta)

		if player.X < 0 {
			player.X = 0
		}
		if player.X+player.Width > 960 {
			player.X = 960 - player.Width
		}
		if player.Y > 540 {
			player.X = 40
			player.Y = 408
			velocityY = 0
			onGround = true
		}

		message = "Explore la map — M pour ouvrir le menu"

		for _, character := range characters {
			if rl.CheckCollisionRecs(player, character.rect) {
				message = "Appuie sur E pour parler au " + character.name
				if rl.IsKeyPressed(rl.KeyE) {
					switch character.name {
					case "MARCHAND":
						screen = "marchand"
					case "CHARCUDOC":
						screen = "charcudoc"
					case "GUERRE DE GANG":
						screen = "gang"
					}
				}
			}
		}

		drawMap(&player, blocks, characters, message, &perso, &playerAnims, playerFacingLeft, direction != 0, onGround)
	}

	fmt.Println("Au revoir !")
}

func handlePanelInput(screen *string, perso *Character, market *Market, charcudoc *Market, message *string, visualCombat **VisualCombat, playerAnims *PlayerAnimations) {
	if rl.IsKeyPressed(rl.KeyEscape) {

		*screen = "map"
		return
	}

	if *screen == "menu" {
		switch {
		case rl.IsKeyPressed(rl.KeyOne):
			*screen = "perso"
		case rl.IsKeyPressed(rl.KeyTwo):
			*screen = "inventaire"
		case rl.IsKeyPressed(rl.KeyThree):
			rl.CloseWindow()
		}
		return
	}

	if *screen == "marchand" {
		if rl.IsKeyPressed(rl.KeyOne) {
			_, result := market.buy(perso, 1)
			*message = result
			return
		}
		if rl.IsKeyPressed(rl.KeyTwo) {
			_, result := market.buy(perso, 2)
			*message = result
			return
		}
		if rl.IsKeyPressed(rl.KeyThree) {
			_, result := market.buy(perso, 3)
			*message = result
			return
		}
		if rl.IsKeyPressed(rl.KeyFour) {
			*screen = "map"
		}
		return
	}

	if *screen == "charcudoc" {
		if rl.IsKeyPressed(rl.KeyOne) {
			_, result := charcudoc.buy(perso, 1)
			*message = result
			return
		}
		if rl.IsKeyPressed(rl.KeyTwo) {
			_, result := charcudoc.buy(perso, 2)
			*message = result
			return
		}
		if rl.IsKeyPressed(rl.KeyThree) {
			_, result := charcudoc.buy(perso, 3)
			*message = result
			return
		}
		if rl.IsKeyPressed(rl.KeyFour) {
			*screen = "map"
		}
		return
	}

	if *screen == "gang" {
		if rl.IsKeyPressed(rl.KeyOne) {
			// Le combat est maintenant entièrement visuel dans la fenêtre Raylib.
			*visualCombat = newVisualCombat(perso, false, playerAnims)
			*screen = "combat"
			return
		}
		if rl.IsKeyPressed(rl.KeyTwo) {
			*visualCombat = newVisualCombat(perso, true, playerAnims)
			*screen = "combat"
			return
		}
		if rl.IsKeyPressed(rl.KeyFour) {
			*screen = "map"
		}
		return
	}
}

func validName(name string) bool {
	if name == "" {
		return false
	}

	for _, r := range name {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func formatName(name string) string {
	if name == "" {
		return ""
	}
	runes := []rune(strings.ToLower(name))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func itoa(value int) string {
	return fmt.Sprintf("%d", value)
}

func moveHorizontal(player *rl.Rectangle, movement float32) {
	player.X += movement
}

func touchesBlockSide(player rl.Rectangle, block rl.Rectangle) bool {
	return player.Y < block.Y && player.Y+player.Height > block.Y+4
}

func overlapsX(first rl.Rectangle, second rl.Rectangle) bool {
	return first.X < second.X+second.Width && first.X+first.Width > second.X
}

