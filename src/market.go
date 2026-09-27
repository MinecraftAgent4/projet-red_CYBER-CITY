package main

type Market struct {
	liste_offres []Trade
}

type Trade struct {
	result      Object
	price       int
	ingredients []Resource
	repeatable  bool
}

func (m *Market) addObject(object Object, price int, ingredients []Resource, isRepeatable bool) {
	m.liste_offres = append(m.liste_offres, Trade{result: object, price: price, ingredients: ingredients, repeatable: isRepeatable})
}

func initMarket() Market {
	var market Market
	market.addObject(Potion{nom: "Steampack de basse qualité", effect: func(target Entity) { target.AddPV(20) }}, 12, []Resource{}, false)
	market.addObject(Potion{nom: "Steampack", effect: func(target Entity) { target.AddPV(40) }}, 28, []Resource{}, false)
	market.addObject(Item{nom: "Grenade à fragmentation"}, 40, []Resource{}, false)
	market.addObject(Item{nom: "Sacoche"}, 30, []Resource{}, true)
	return market
}

func initCharcudoc() Market {
	var charcudoc Market
	charcudoc.addObject(Casque{nom: "Interface neurale", defense: 5}, 0, []Resource{Resource{nom: "Composants", quantité: 3}, Resource{nom: "Poudre", quantité: 1}}, true)
	charcudoc.addObject(Plastron{nom: "Optiques cybernétiques", defense: 10}, 0, []Resource{Resource{nom: "Composants", quantité: 4}, Resource{nom: "Ferraille", quantité: 2}}, true)
	charcudoc.addObject(Bottes{nom: "Réflexes augmentés", defense: 15}, 0, []Resource{Resource{nom: "Composants", quantité: 5}, Resource{nom: "Poudre", quantité: 2}}, true)
	return charcudoc
}

func (m Market) buy(character *Character, choice int) (bool, string) {
	if choice < 1 || choice > len(m.liste_offres) {
		return false, "Choix invalide."
	}

	trade := m.liste_offres[choice-1]
	if !character.inventoryLimit() {
		return false, "Inventaire plein."
	}
	if character.money < trade.price {
		return false, "Pas assez de pièces d'or."
	}
	for _, ingredient := range trade.ingredients {
		if character.resourceQuantity(ingredient.nom) < ingredient.quantité {
			return false, "Ressources insuffisantes."
		}
	}
	for _, ingredient := range trade.ingredients {
		character.removeResource(ingredient.nom, ingredient.quantité)
	}
	character.money -= trade.price
	if trade.result.Nom() != "Sacoche" {
		character.GiveItem(trade.result)
	}
	return true, trade.result.Nom() + " acheté."
}
