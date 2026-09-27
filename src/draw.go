package main

import (
	"fmt"
	"strings"
	"unicode"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	virtualWidth  int32 = 960
	virtualHeight int32 = 540
	windowWidth   int32 = 1440
	windowHeight  int32 = 810
)

var virtualTarget rl.RenderTexture2D

func initVirtualRenderer() {
	virtualTarget = rl.LoadRenderTexture(virtualWidth, virtualHeight)
	rl.SetTextureFilter(virtualTarget.Texture, rl.FilterPoint)
}

func unloadVirtualRenderer() {
	if virtualTarget.ID != 0 {
		rl.UnloadRenderTexture(virtualTarget)
		virtualTarget = rl.RenderTexture2D{}
	}
}

func beginVirtualDrawing(background rl.Color) {
	rl.BeginTextureMode(virtualTarget)
	rl.ClearBackground(background)
}

func endVirtualDrawing() {
	rl.EndTextureMode()
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	source := rl.Rectangle{X: 0, Y: 0, Width: float32(virtualWidth), Height: -float32(virtualHeight)}
	destination := rl.Rectangle{X: 0, Y: 0, Width: float32(windowWidth), Height: float32(windowHeight)}
	rl.DrawTexturePro(virtualTarget.Texture, source, destination, rl.Vector2{}, 0, rl.White)
	rl.EndDrawing()
}

func drawMap(player *rl.Rectangle, blocks []mapBlock, characters []mapCharacter, message string, perso *Character, playerAnims *PlayerAnimations, facingLeft bool, moving bool, onGround bool) {
	beginVirtualDrawing(rl.NewColor(8, 10, 20, 255))

	// Décor du niveau : backgrounds, plateformes en tiles, machines et overlay animé.
	drawWorldBackground()

	for _, block := range blocks {
		drawPlatformTexture(block.rect)
	}

	drawWorldDecor()

	for _, character := range characters {
		// La collision reste active, mais aucune hitbox n'est dessinée.
		if character.animation != nil {
			character.animation.Draw(character.rect.X, character.rect.Y, character.rect.Width, character.rect.Height)
		} else {
			rl.DrawRectangleRec(character.rect, character.color)
		}

		uiText(character.name, int(character.rect.X-20), int(character.rect.Y-24), 12, rl.RayWhite)
	}

	// Joueur : animations propres à sa classe.
	if playerAnims != nil && playerAnims.idle != nil && playerAnims.idle.texture.ID != 0 {
		drawX := player.X - 8
		drawY := player.Y - 6
		animation := playerAnims.idle
		if !onGround {
			animation = playerAnims.jump
		} else if moving {
			animation = playerAnims.walk
		}
		if animation != nil {
			if facingLeft {
				animation.DrawFlipped(drawX, drawY, 48, 48)
			} else {
				animation.Draw(drawX, drawY, 48, 48)
			}
		}
	} else {
		rl.DrawRectangleRec(*player, rl.SkyBlue)
	}

	// Effet de pluie/lumière au-dessus de toute la scène.
	drawWorldOverlay()

	// HUD basé sur les assets du pack.
	drawUIHealthBar(24, 20, 180, perso.pv_actuelle, perso.pv_total)
	uiText(perso.nom+" | "+perso.classe, 220, 18, 18, rl.Gold)
	uiText("PV "+itoa(perso.pv_actuelle)+"/"+itoa(perso.pv_total), 220, 43, 15, rl.RayWhite)

	if message != "" {
		drawUIPrompt(message, 24, 82)
	}

	drawUIButton("M", rl.Rectangle{X: 870, Y: 18, Width: 58, Height: 38}, false)
	drawUIHint("A/D ou ←/→ : déplacement    Espace : saut    M : menu", 24, 510)

	endVirtualDrawing()
}

func drawClassScreen(selected string) {
	beginVirtualDrawing(rl.NewColor(8, 10, 20, 255))
	drawUIPanel("CRÉATION DU PERSONNAGE")

	uiText("Choisis une classe", 90, 125, 26, rl.RayWhite)
	drawUIButton("1 - NETRUNNER     80 PV", rl.Rectangle{X: 90, Y: 180, Width: 330, Height: 42}, selected == "netrunner")
	drawUIButton("2 - MERC         100 PV", rl.Rectangle{X: 90, Y: 235, Width: 330, Height: 42}, selected == "merc")
	drawUIButton("3 - CYBERPSYCHO  120 PV", rl.Rectangle{X: 90, Y: 290, Width: 330, Height: 42}, selected == "cyberpsycho")

	if selected != "" {
		uiText("Classe sélectionnée : "+selected, 90, 370, 20, rl.Gold)
	}
	drawUIHint("Appuie sur 1, 2 ou 3", 90, 425)
	endVirtualDrawing()
}

func drawGamePanel(screen string, perso *Character, market *Market, charcudoc *Market, message string) {
	beginVirtualDrawing(rl.NewColor(8, 10, 20, 255))

	switch screen {
	case "menu":
		drawMenuPanel()
	case "perso":
		drawCharacterPanel(perso)
	case "inventaire":
		drawInventoryPanel(perso)
	case "marchand":
		drawMarketPanel("MARCHAND", market)
	case "charcudoc":
		drawMarketPanel("CHARCUDOC", charcudoc)
	case "gang":
		drawGangPanel()
	}

	if message != "" && screen != "menu" && screen != "gang" {
		uiText(message, 70, 450, 16, rl.Gold)
	}
	drawUIHint("Échap : retour", 760, 505)
	endVirtualDrawing()
}

func drawMenuPanel() {
	drawUIPanel("LE QUARTIER")
	drawUIButton("1. PERSONNAGE", rl.Rectangle{X: 90, Y: 130, Width: 330, Height: 45}, false)
	drawUIButton("2. INVENTAIRE", rl.Rectangle{X: 90, Y: 190, Width: 330, Height: 45}, false)
	drawUIButton("3. EXIT", rl.Rectangle{X: 90, Y: 250, Width: 330, Height: 45}, false)
	drawUIHint("Les services sont disponibles auprès des PNJ sur la carte.", 90, 340)
}

func drawCharacterPanel(perso *Character) {
	drawUIPanel("PERSONNAGE")
	uiText("Nom : "+perso.nom, 90, 125, 22, rl.RayWhite)
	uiText("Classe : "+perso.classe, 90, 165, 22, rl.RayWhite)
	uiText("Niveau : "+itoa(perso.niveau), 90, 205, 22, rl.RayWhite)
	uiText("Argent : "+itoa(perso.money)+" po", 90, 245, 22, rl.Gold)
	uiText("Emplacements : "+itoa(len(perso.inventaire))+" / "+itoa(perso.maxslots), 90, 285, 22, rl.RayWhite)
	uiText("PV", 90, 330, 18, rl.RayWhite)
	drawUIHealthBar(130, 328, 350, perso.pv_actuelle, perso.pv_total)
	uiText(itoa(perso.pv_actuelle)+" / "+itoa(perso.pv_total), 495, 328, 18, rl.RayWhite)
}

func drawInventoryPanel(perso *Character) {
	drawUIPanel("INVENTAIRE")

	if len(perso.inventaire) == 0 {
		uiText("Inventaire vide.", 90, 145, 24, rl.RayWhite)
		return
	}

	for i, object := range perso.inventaire {
		if i >= 8 {
			break
		}
		drawUIButton(fmt.Sprintf("%d. %s", i+1, object.Nom()), rl.Rectangle{X: 80, Y: float32(120 + i*42), Width: 520, Height: 36}, false)
	}
}

func drawMarketPanel(title string, market *Market) {
	drawUIPanel(title)

	for i, trade := range market.liste_offres {
		line := fmt.Sprintf("%d. %s", i+1, trade.result.Nom())
		if trade.price > 0 {
			line += " - " + itoa(trade.price) + " po"
		}
		if len(trade.ingredients) > 0 {
			ingredients := make([]string, 0, len(trade.ingredients))
			for _, ingredient := range trade.ingredients {
				ingredients = append(ingredients, fmt.Sprintf("%d %s", ingredient.quantité, ingredient.nom))
			}
			line += " [" + strings.Join(ingredients, ", ") + "]"
		}
		drawUIButton(line, rl.Rectangle{X: 70, Y: float32(115 + i*52), Width: 780, Height: 40}, false)
	}
	drawUIButton("4. RETOUR À LA CARTE", rl.Rectangle{X: 70, Y: 355, Width: 300, Height: 40}, false)
}

func drawGangPanel() {
	drawUIPanel("GUERRE DE GANG")
	uiText("Le gang est prêt à se battre.", 90, 120, 24, rl.RayWhite)
	drawUIButton("1. COMBAT ALÉATOIRE (1 À 3 MOBS)", rl.Rectangle{X: 90, Y: 175, Width: 500, Height: 42}, false)
	drawUIButton("2. COMBAT CONTRE LE BOSS", rl.Rectangle{X: 90, Y: 230, Width: 500, Height: 42}, false)
	drawUIButton("4. RETOUR À LA CARTE", rl.Rectangle{X: 90, Y: 285, Width: 300, Height: 42}, false)
	uiText("Le nombre de mobs est tiré au hasard à chaque combat.", 90, 355, 17, rl.LightGray)
}

func drawPanel(title string) {
	drawUIPanel(title)
}

func readNameInput(current string) string {
	if rl.IsKeyPressed(rl.KeyBackspace) && len(current) > 0 {
		runes := []rune(current)
		current = string(runes[:len(runes)-1])
	}
	for {
		char := rl.GetCharPressed()
		if char == 0 {
			break
		}
		r := rune(char)
		if unicode.IsLetter(r) && len([]rune(current)) < 18 {
			current += string(r)
		}
	}
	return current
}
