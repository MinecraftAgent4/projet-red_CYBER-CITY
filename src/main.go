package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputReader = bufio.NewReader(os.Stdin)

func initCharacter(nom string, classe string, niveau int, pvTotal int, pvActuelle int, inventaire []Object, money int) Character {
	return Character{
		nom:          nom,
		classe:       classe,
		niveau:       niveau,
		pv_total:     pvTotal,
		pv_actuelle:  pvActuelle,
		inventaire:   inventaire,
		maxslots:     10,
		money:        money,
		weapon:       nil,
		weapon2:      nil,
		spellbook:    SpellBook{},
		exp_required: 100,
		exp_joueur:   0,
	}
}

func Isformated(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	runes := []rune(strings.ToLower(s))
	for _, r := range runes {
		if r < 'a' || r > 'z' {
			return ""
		}
	}
	runes[0] -= 'a' - 'A'
	return string(runes)
}

func saisiePrenom() string {
	for {
		fmt.Print("Entre ton prenom (lettres uniquement) : ")
		saisie := readLine()
		nom := Isformated(saisie)
		if nom != "" {
			return nom
		}
		fmt.Println("Prenom invalide.")
	}
}

func création_perso() Character {
	fmt.Println("\n========================================")
	fmt.Println("          CREATION DU PERSONNAGE")
	fmt.Println("========================================")

	var classe string
	for {
		fmt.Println("\nChoisis une classe :")
		fmt.Println("1 - Netrunner   (80 PV)")
		fmt.Println("2 - Mercenaire  (100 PV)")
		fmt.Println("3 - Cyberpsycho  (120 PV)")
		fmt.Print("Choix : ")

		switch readInt() {
		case 1:
			classe = "netrunner"
		case 2:
			classe = "merc"
		case 3:
			classe = "cyberpsycho"
		default:
			fmt.Println("Choix invalide.")
			continue
		}
		break
	}

	pv := 100
	switch classe {
	case "netrunner":
		pv = 80
	case "cyberpsycho":
		pv = 120
	}

	perso := initCharacter(
		saisiePrenom(), classe, 1, pv, pv,
		[]Object{
			Resource{nom: "Ferraille", quantité: 10, quantité_max: 99},
			Resource{nom: "Composants", quantité: 4, quantité_max: 99},
			Resource{nom: "Poudre", quantité: 5, quantité_max: 99},
		},
		100,
	)

	switch classe {
	case "merc":
		perso.weapon = Melee{nom: "Epee de mercenaire", dmg: 12}
		perso.weapon2 = Ranged{nom: "Pistolet", dmg: 8}
	case "netrunner":
		perso.spellbook = createSpellBook(Spell{nom: "Cyberdeck"})
	case "cyberpsycho":
		perso.weapon = Ranged{nom: "Gros calibre gauche", dmg: 15}
		perso.weapon2 = Ranged{nom: "Gros calibre droit", dmg: 15}
	}

	return perso
}

func readLine() string {
	line, err := inputReader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return ""
	}
	return strings.TrimSpace(line)
}

func readInt() int {
	value, err := strconv.Atoi(readLine())
	if err != nil {
		return -1
	}
	return value
}

func pause() {
	fmt.Print("\nAppuie sur Entree pour continuer...")
	readLine()
}

func Cli() {
	perso := création_perso()
	marketObjets := initMarket()
	charcudocObjets := initCharcudoc()

	for perso.Hp() > 0 {
		fmt.Println()
		fmt.Print(menu)
		fmt.Print("Choix : ")

		switch readInt() {
		case 1:
			fmt.Print(charactinfoMenu(perso))
			pause()

		case 2:
			openMarket(&perso, &marketObjets)

		case 3:
			openCharcudoc(&perso, &charcudocObjets)

		case 4:
			openGang(&perso)

		case 5:
			openInventory(&perso)

		case 6:
			fmt.Println("\nFin de partie. A bientot !")
			return

		default:
			fmt.Println("Choix invalide : entre 1 et 6.")
		}
	}

	fmt.Println("\n========================================")
	fmt.Println("              FIN DE PARTIE")
	fmt.Println("========================================")
}

func openMarket(perso *Character, market *Market) {
	for {
		fmt.Println(marketMenu(*market, *perso))
		fmt.Print("Choix : ")
		choice := readInt()

		if choice == 5 {
			return
		}
		if choice < 1 || choice > len(market.liste_offres) {
			fmt.Println("Choix invalide.")
			continue
		}

		if choice == 4 {
			if perso.maxslots >= 40 {
				fmt.Println("La sacoche est deja au maximum.")
				continue
			}
			if perso.money < 30 {
				fmt.Println("Pas assez d'argent.")
				continue
			}
			perso.money -= 30
			perso.upgradeInventorySlot()
			fmt.Printf("Sacoche amelioree ! Capacite : %d objets.\n", perso.maxslots)
			continue
		}

		_, message := market.buy(perso, choice)
		fmt.Println(message)
	}
}

func openCharcudoc(perso *Character, market *Market) {
	for {
		fmt.Println(charcudocMenu(*market, *perso))
		fmt.Print("Choix : ")
		choice := readInt()
		if choice == 4 {
			return
		}
		if choice < 1 || choice > len(market.liste_offres) {
			fmt.Println("Choix invalide.")
			continue
		}
		_, message := market.buy(perso, choice)
		fmt.Println(message)
	}
}

func openInventory(perso *Character) {
	for {
		fmt.Println("\n========================================")
		fmt.Println("               INVENTAIRE")
		fmt.Printf("Cases : %d/%d\n", len(perso.inventaire), perso.maxslots)
		fmt.Println("========================================")

		if len(perso.inventaire) == 0 {
			fmt.Println("Inventaire vide.")
		} else {
			for i, object := range perso.inventaire {
				fmt.Printf("%2d - %s\n", i+1, object.Nom())
			}
		}

		fmt.Println("\n0 - Retour")
		fmt.Println("Choisis un objet pour l'utiliser/equiper.")
		fmt.Print("Choix : ")
		choice := readInt()
		if choice == 0 {
			return
		}
		if choice < 1 || choice > len(perso.inventaire) {
			fmt.Println("Choix invalide.")
			continue
		}

		useInventoryObject(perso, choice-1)
	}
}

func useInventoryObject(perso *Character, index int) {
	object := perso.inventaire[index]

	switch item := object.(type) {
	case Armure:
		perso.equipArmor(item)
		fmt.Printf("%s equipe.\n", item.Nom())
	case Potion:
		before := perso.Hp()
		perso.TakePot(item)
		fmt.Printf("%s utilisee : +%d PV.\n", item.Nom(), perso.Hp()-before)
	case Item:
		fmt.Printf("%s ne peut pas etre utilise ici.\n", item.Nom())
	default:
		fmt.Println("Cet objet n'est pas utilisable ici.")
	}
}

func openGang(perso *Character) {
	for {
		fmt.Println("\n========================================")
		fmt.Println("             GUERRE DE GANG")
		fmt.Println("========================================")
		fmt.Println("1 - Patrouille : 1 a 3 ennemis")
		fmt.Println("2 - Affronter le boss")
		fmt.Println("3 - Retour au quartier")
		fmt.Print("Choix : ")

		switch readInt() {
		case 1:
			enemies := randomGang()
			won := combatGroupe(perso, enemies)
			if !won {
				return
			}
		case 2:
			won := combatGroupe(perso, []Monster{initBoss1()})
			if !won {
				return
			}
		case 3:
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}

func main() {
	game()
}
