package main

import "fmt"

const (
	reset  = "\x1b[0m"
	bold   = "\x1b[1m"
	italic = "\x1b[3m"

	red   = "\x1b[38;5;160m"
	gold  = "\x1b[38;5;220m"
	gray  = "\x1b[38;5;244m"
	white = "\x1b[97m"
	green = "\x1b[38;5;114m"
)

var menu = gold + "╔════════════════════════════════════╗" + reset + "\n" +
	gold + "║" + reset + "       " + bold + gold + "▓▓▓  LE QUARTIER  ▓▓▓" + reset + "        " + gold + "║" + reset + "\n" +
	gold + "╠════════════════════════════════════╣" + reset + "\n" +
	gold + "║" + reset + "                                    " + gold + "║" + reset + "\n" +
	gold + "║" + reset + "   " + red + bold + "1." + reset + " " + white + "Perso" + reset + "                         " + gold + "║" + reset + "\n" +
	gold + "║" + reset + "   " + red + bold + "2." + reset + " " + white + "Marchand" + reset + "                      " + gold + "║" + reset + "\n" +
	gold + "║" + reset + "   " + red + bold + "3." + reset + " " + white + "Charcudoc" + reset + "                     " + gold + "║" + reset + "\n" +
	gold + "║" + reset + "   " + red + bold + "4." + reset + " " + white + "Guerre de gang" + reset + "                " + gold + "║" + reset + "\n" +
	gold + "║" + reset + "   " + red + bold + "5." + reset + " " + white + "Inventaire" + reset + "                    " + gold + "║" + reset + "\n" +
	gold + "║" + reset + "   " + red + bold + "6." + reset + " " + white + "EXIT" + reset + "                          " + gold + "║" + reset + "\n" +
	gold + "║" + reset + "                                    " + gold + "║" + reset + "\n" +
	gold + "║" + reset + gray + italic + "  › choisis un numéro pour agir" + reset + "     " + gold + "║" + reset + "\n" +
	gold + "╚════════════════════════════════════╝" + reset + "\n"

func marketMenu(market Market, character Character) string {
	return gold + "╔═════════════════════════════════════════╗" + reset + "\n" +
		gold + "║" + reset + "            " + bold + gold + "▓▓▓  MARCHAND  ▓▓▓" + reset + "           " + gold + "║" + reset + "\n" +
		gold + "╠═════════════════════════════════════════╣" + reset + "\n" +
		gold + "║" + reset + "                                         " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + white + "Money :" + reset + " " + green + bold + fmt.Sprintf("%d", character.money) + "po" + reset + "                         " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + red + bold + "1." + reset + " " + white + "Steampack basse qualité" + reset + " - " + green + bold + "12 po" + reset + "    " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + red + bold + "2." + reset + " " + white + "Steampack" + reset + " - " + green + bold + "28 po" + reset + "                  " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + red + bold + "3." + reset + " " + white + "Grenade à fragmentation" + reset + " - " + green + bold + "40 po" + reset + "    " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + red + bold + "4." + reset + " " + white + "Sacoche" + reset + " - " + green + bold + "30 po" + reset + "                    " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + red + bold + "5." + reset + " " + white + "Retour" + reset + "                             " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "                                         " + gold + "║" + reset + "\n" +
		gold + "║" + reset + gray + italic + "  › choisis un numéro pour acheter" + reset + "       " + gold + "║" + reset + "\n" +
		gold + "╚═════════════════════════════════════════╝" + reset + "\n"

}

func charcudocMenu(market Market, character Character) string {
	return gold + "╔═══════════════════════════════════════════════════╗" + reset + "\n" +
		gold + "║" + reset + "                " + bold + gold + "▓▓▓  CHARCUDOC  ▓▓▓" + reset + "                " + gold + "║" + reset + "\n" +
		gold + "╠═══════════════════════════════════════════════════╣" + reset + "\n" +
		gold + "║" + reset + "   " + red + bold + "IMPLANTS" + reset + "                                        " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   1. Interface neurale       " + green + bold + "3 Comp., 1 Poudre" + reset + "    " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   2. Optiques cybernétiques " + green + bold + "4 Comp., 2 Ferraille" + reset + "  " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   3. Réflexes augmentés     " + green + bold + "5 Comp., 2 Poudre" + reset + "     " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   4. Retour                                       " + gold + "║" + reset + "\n" +
		gold + "╚═══════════════════════════════════════════════════╝" + reset + "\n"
}

func charactinfoMenu(character Character) string {
	casque := "Aucun"
	plastron := "Aucun"
	bottes := "Aucunes"
	if character.equipe.helmet != nil {
		casque = character.equipe.helmet.Nom()
	}
	if character.equipe.torso != nil {
		plastron = character.equipe.torso.Nom()
	}
	if character.equipe.boots != nil {
		bottes = character.equipe.boots.Nom()
	}
	return gold + "╔═════════════════════════════════════════╗" + reset + "\n" +
		gold + "║" + reset + "         " + bold + gold + "▓▓▓  INFOS JOUEUR  ▓▓▓" + reset + "          " + gold + "║" + reset + "\n" +
		gold + "╠═════════════════════════════════════════╣" + reset + "\n" +
		gold + "║" + reset + "                                         " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + white + bold + "Nom :" + reset + " " + white + fmt.Sprintf("%s", character.nom) + reset + "                              " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + white + bold + "Classe :" + reset + " " + white + fmt.Sprintf("%s", character.classe) + reset + "                    " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + white + bold + "Niveau :" + reset + " " + white + fmt.Sprintf("%d", character.niveau) + reset + "                            " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + white + bold + "PV :" + reset + " " + white + fmt.Sprintf("%d / %d", character.pv_actuelle, character.pv_total) + reset + "                          " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + white + bold + "Experience :" + reset + " " + white + fmt.Sprintf("%d / %d", character.exp_joueur, character.exp_required) + reset + "                    " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + white + bold + "Argent :" + reset + " " + white + fmt.Sprintf("%d", character.money) + reset + "                          " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + white + bold + "Armure :" + reset + " " + white + fmt.Sprintf("%d", character.armorValue()) + reset + "                            " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + white + bold + "Casque :" + reset + " " + white + fmt.Sprintf(casque) + "                        " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + white + bold + "Plastron :" + reset + " " + white + fmt.Sprintf(plastron) + reset + "                      " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "   " + white + bold + "Bottes :" + reset + " " + white + fmt.Sprintf(bottes) + "                      " + gold + "║" + reset + "\n" +
		gold + "║" + reset + "                                         " + gold + "║" + reset + "\n" +
		gold + "╚═════════════════════════════════════════╝" + reset + "\n"

}
