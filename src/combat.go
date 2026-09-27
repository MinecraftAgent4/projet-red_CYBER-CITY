package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Attack(a Entity, t Entity, mod int) {
	damage := (a.Atk() * mod) - t.Def()
	if damage < 1 {
		damage = 1
	}
	t.Dmg(damage)
}

func randomGang() []Monster {
	count := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(3) + 1
	enemies := make([]Monster, 0, count)
	for i := 0; i < count; i++ {
		mob := initGoblin()
		mob.name = fmt.Sprintf("Gangster %d", i+1)
		enemies = append(enemies, mob)
	}
	return enemies
}

func combat(perso *Character, monstre Monster) bool {
	return combatGroupe(perso, []Monster{monstre})
}

func combatGroupe(perso *Character, monstres []Monster) bool {
	if len(monstres) == 0 {
		return false
	}

	fmt.Println("\n========================================")
	fmt.Println("                 COMBAT")
	fmt.Println("========================================")
	fmt.Printf("Ennemis : %d\n", len(monstres))

	turn := 1
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for perso.Hp() > 0 && aliveMonsters(monstres) {
		printCombatState(perso, monstres, turn)

		action, ok := playerTurn(perso, monstres)
		if !ok {
			continue
		}
		if action == "flee" {
			fmt.Println("\nTu fuis le combat.")
			return false
		}
		if action == "invalid" {
			continue
		}

		if !aliveMonsters(monstres) {
			break
		}

		// Chaque ennemi vivant joue une fois.
		for i := range monstres {
			if perso.Hp() <= 0 {
				break
			}
			if monstres[i].Hp() <= 0 {
				continue
			}

			fmt.Printf("\n%s attaque !\n", monstres[i].Name())
			if strings.HasPrefix(monstres[i].Name(), "Boss") {
				boss1Pattern(turn, &monstres[i], perso)
			} else {
				damageMod := 1
				if rng.Intn(4) == 3 {
					damageMod = 2
					fmt.Println("Attaque puissante !")
				}
				Attack(&monstres[i], perso, damageMod)
				fmt.Printf("Tu perds des PV : %d/%d restants.\n", perso.Hp(), perso.MaxHp())
			}
		}

		turn++
	}

	if perso.Hp() <= 0 {
		fmt.Println("\n========================================")
		fmt.Println("                 DEFAITE")
		fmt.Println("========================================")
		perso.isDead()
		return false
	}

	reward := 20 * len(monstres)
	xp := 50 * len(monstres)
	perso.money += reward
	perso.exp_joueur += xp

	fmt.Println("\n========================================")
	fmt.Println("                VICTOIRE")
	fmt.Println("========================================")
	fmt.Printf("Tu as vaincu tous les ennemis.\n")
	fmt.Printf("Recompenses : +%d po, +%d XP.\n", reward, xp)

	for perso.exp_joueur >= perso.exp_required {
		perso.lvlUp()
		fmt.Printf("Niveau %d atteint ! PV max : %d.\n", perso.niveau, perso.MaxHp())
	}

	return true
}

func printCombatState(perso *Character, monstres []Monster, turn int) {
	fmt.Println("\n----------------------------------------")
	fmt.Printf("TOUR %d\n", turn)
	fmt.Printf("%s [%s] - PV %d/%d - Armure %d\n",
		perso.Name(), perso.classe, perso.Hp(), perso.MaxHp(), perso.Def())
	fmt.Println("----------------------------------------")

	for i := range monstres {
		state := "VIVANT"
		if monstres[i].Hp() <= 0 {
			state = "VAINCU"
		}
		fmt.Printf("%d - %s : %d/%d PV | ATK %d | DEF %d | %s\n",
			i+1,
			monstres[i].Name(),
			monstres[i].Hp(),
			monstres[i].MaxHp(),
			monstres[i].Atk(),
			monstres[i].Def(),
			state,
		)
	}
	fmt.Println("----------------------------------------")
}

func playerTurn(perso *Character, monstres []Monster) (string, bool) {
	for {
		fmt.Println("\nActions :")

		switch perso.classe {
		case "merc":
			fmt.Println("1 - Attaque au corps a corps")
			fmt.Println("2 - Tirer au pistolet")
			fmt.Println("3 - Se soigner")
			fmt.Println("4 - Inventaire")
			fmt.Println("5 - Fuir")
		case "netrunner":
			fmt.Println("1 - Hacks")
			fmt.Println("2 - Se soigner")
			fmt.Println("3 - Inventaire")
			fmt.Println("4 - Fuir")
		case "cyberpsycho":
			fmt.Println("1 - Tir gauche")
			fmt.Println("2 - Tir droit")
			fmt.Println("3 - Double tir")
			fmt.Println("4 - Se soigner")
			fmt.Println("5 - Inventaire")
			fmt.Println("6 - Fuir")
		}

		fmt.Print("Choix : ")
		choice := readInt()

		switch perso.classe {
		case "merc":
			switch choice {
			case 1:
				target := chooseTarget(monstres)
				if target < 0 {
					return "invalid", false
				}
				fmt.Printf("\n%s frappe avec %s !\n", perso.Name(), perso.weapon.Nom())
				before := monstres[target].Hp()
				Attack(perso, &monstres[target], 1)
				fmt.Printf("Degats infliges : %d.\n", before-monstres[target].Hp())
				return "attack", true
			case 2:
				target := chooseTarget(monstres)
				if target < 0 {
					return "invalid", false
				}
				fmt.Printf("\n%s tire avec %s !\n", perso.Name(), perso.weapon2.Nom())
				old := perso.weapon
				perso.weapon = perso.weapon2
				before := monstres[target].Hp()
				Attack(perso, &monstres[target], 1)
				perso.weapon = old
				fmt.Printf("Degats infliges : %d.\n", before-monstres[target].Hp())
				return "attack", true
			case 3:
				heal(perso, 15)
				return "heal", true
			case 4:
				openCombatInventory(perso)
				return "item", true
			case 5:
				return "flee", true
			}

		case "netrunner":
			switch choice {
			case 1:
				used := netrunnerAttack(perso, monstres)
				if used {
					return "attack", true
				}
				continue
			case 2:
				heal(perso, 15)
				return "heal", true
			case 3:
				openCombatInventory(perso)
				return "item", true
			case 4:
				return "flee", true
			}

		case "cyberpsycho":
			switch choice {
			case 1, 2:
				target := chooseTarget(monstres)
				if target < 0 {
					return "invalid", false
				}
				old := perso.weapon
				if choice == 2 {
					perso.weapon = perso.weapon2
				}
				before := monstres[target].Hp()
				fmt.Printf("\nTir avec %s !\n", perso.weapon.Nom())
				Attack(perso, &monstres[target], 1)
				perso.weapon = old
				fmt.Printf("Degats infliges : %d.\n", before-monstres[target].Hp())
				return "attack", true
			case 3:
				target := chooseTarget(monstres)
				if target < 0 {
					return "invalid", false
				}
				fmt.Println("\nDOUBLE TIR !")
				before := monstres[target].Hp()
				Attack(perso, &monstres[target], 1)
				if monstres[target].Hp() > 0 {
					old := perso.weapon
					perso.weapon = perso.weapon2
					Attack(perso, &monstres[target], 1)
					perso.weapon = old
				}
				fmt.Printf("Degats infliges : %d.\n", before-monstres[target].Hp())
				return "attack", true
			case 4:
				heal(perso, 15)
				return "heal", true
			case 5:
				openCombatInventory(perso)
				return "item", true
			case 6:
				return "flee", true
			}
		}

		fmt.Println("Choix invalide.")
	}
}

func chooseTarget(monstres []Monster) int {
	alive := 0
	for _, m := range monstres {
		if m.Hp() > 0 {
			alive++
		}
	}
	if alive <= 1 {
		for i := range monstres {
			if monstres[i].Hp() > 0 {
				return i
			}
		}
		return -1
	}

	for {
		fmt.Print("Choisis la cible : ")
		choice := readInt() - 1
		if choice >= 0 && choice < len(monstres) && monstres[choice].Hp() > 0 {
			return choice
		}
		fmt.Println("Cible invalide.")
	}
}

func openCombatInventory(perso *Character) {
	usable := false
	fmt.Println("\n--- Inventaire de combat ---")
	for i, object := range perso.inventaire {
		if _, ok := object.(Potion); ok {
			fmt.Printf("%d - %s\n", i+1, object.Nom())
			usable = true
		}
	}
	if !usable {
		fmt.Println("Aucune potion disponible.")
		return
	}
	fmt.Println("0 - Annuler")
	fmt.Print("Choix : ")
	choice := readInt()
	if choice == 0 {
		return
	}
	if choice < 1 || choice > len(perso.inventaire) {
		fmt.Println("Choix invalide.")
		return
	}
	potion, ok := perso.inventaire[choice-1].(Potion)
	if !ok {
		fmt.Println("Cet objet n'est pas une potion.")
		return
	}
	before := perso.Hp()
	perso.TakePot(potion)
	fmt.Printf("%s utilisee : +%d PV.\n", potion.Nom(), perso.Hp()-before)
}

func heal(perso *Character, amount int) {
	before := perso.Hp()
	perso.AddPV(amount)
	fmt.Printf("Soin : +%d PV (%d/%d).\n", perso.Hp()-before, perso.Hp(), perso.MaxHp())
}

func aliveMonsters(monstres []Monster) bool {
	for _, monster := range monstres {
		if monster.Hp() > 0 {
			return true
		}
	}
	return false
}

func (m *Monster) exp_reward(c *Character) {
	if m.pv_actuelle > 0 {
		return
	}
	c.exp_joueur += 50
	for c.exp_joueur >= c.exp_required {
		c.lvlUp()
		fmt.Printf("Niveau %d atteint !\n", c.niveau)
	}
}

func mercAttack(perso *Character, monstre *Monster, choix int) {
	if choix == 1 {
		Attack(perso, monstre, 1)
	} else if choix == 2 {
		old := perso.weapon
		perso.weapon = perso.weapon2
		Attack(perso, monstre, 1)
		perso.weapon = old
	}
}

func cyberpsychoAttack(perso *Character, monstre *Monster, choix int) {
	switch choix {
	case 1:
		Attack(perso, monstre, 1)
	case 2:
		old := perso.weapon
		perso.weapon = perso.weapon2
		Attack(perso, monstre, 1)
		perso.weapon = old
	case 3:
		Attack(perso, monstre, 1)
		if monstre.Hp() > 0 {
			old := perso.weapon
			perso.weapon = perso.weapon2
			Attack(perso, monstre, 1)
			perso.weapon = old
		}
	}
}

func netrunnerAttack(perso *Character, monstres []Monster) bool {
	fmt.Println("\n===== HACKS =====")
	fmt.Println("1 - Surcharge : 20 degats")
	fmt.Println("2 - Virus : 10 degats + poison")
	fmt.Println("3 - Court-circuit : 12 degats")
	fmt.Println("4 - Retour")
	fmt.Print("Choix : ")

	choice := readInt()
	if choice == 4 {
		return false
	}

	target := chooseTarget(monstres)
	if target < 0 {
		return false
	}

	var base int
	switch choice {
	case 1:
		base = 20
		fmt.Println("Surcharge du systeme !")
	case 2:
		base = 10
		fmt.Println("Virus injecte !")
	case 3:
		base = 12
		fmt.Println("Court-circuit !")
	default:
		fmt.Println("Choix invalide.")
		return false
	}

	damage := base - monstres[target].Def()
	if damage < 1 {
		damage = 1
	}
	monstres[target].Dmg(damage)
	fmt.Printf("Le hack inflige %d degats.\n", damage)

	if choice == 2 && monstres[target].Hp() > 0 {
		poisonDamage := 5
		monstres[target].Dmg(poisonDamage)
		fmt.Printf("Le virus inflige encore %d degats.\n", poisonDamage)
	}
	return true
}

func boss1Pattern(turn int, m *Monster, t *Character) {
	if turn%4 == 0 {
		fmt.Println("Le boss regenere 10 PV !")
		m.AddPV(10)
		return
	}
	mod := 1
	if turn%4 == 3 {
		mod = 2
	}
	Attack(m, t, mod)
	fmt.Printf("Tu perds des PV : %d/%d restants.\n", t.Hp(), t.MaxHp())
}
