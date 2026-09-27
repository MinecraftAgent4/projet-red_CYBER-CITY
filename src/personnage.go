package main

import (
	"fmt"
)

type Character struct {
	nom          string
	classe       string
	niveau       int
	pv_total     int
	pv_actuelle  int
	inventaire   []Object
	maxslots     int
	money        int
	equipe       Equipement
	weapon       Arme
	weapon2      Arme
	spellbook    SpellBook
	exp_required int
	exp_joueur   int
}

type Equipement struct {
	helmet Armure
	torso  Armure
	boots  Armure
}

func (c Character) armorValue() int {
	total := 0

	if c.equipe.helmet != nil {
		total += c.equipe.helmet.Defense()
	}
	if c.equipe.torso != nil {
		total += c.equipe.torso.Defense()
	}
	if c.equipe.boots != nil {
		total += c.equipe.boots.Defense()
	}

	return total
}

func (c *Character) equipArmor(piece Armure) {
	emplacement_inv := -1
	for i, x := range c.inventaire {
		if x == piece {
			emplacement_inv = i
			break
		}
	}

	if emplacement_inv == -1 {
		return
	}

	c.inventaire = append(
		c.inventaire[:emplacement_inv],
		c.inventaire[emplacement_inv+1:]...,
	)

	switch armure := piece.(type) {
	case Casque:
		if c.equipe.helmet != nil {
			c.inventaire = append(c.inventaire, c.equipe.helmet)
		}
		c.equipe.helmet = armure

	case Plastron:
		if c.equipe.torso != nil {
			c.inventaire = append(c.inventaire, c.equipe.torso)
		}
		c.equipe.torso = armure

	case Bottes:
		if c.equipe.boots != nil {
			c.inventaire = append(c.inventaire, c.equipe.boots)
		}
		c.equipe.boots = armure
	}
}

func (c *Character) GiveItem(object Object) {
	if c.inventoryLimit() {
		c.inventaire = append(c.inventaire, object)
	}
}

func (c *Character) AddInventory(object Object) {
	c.GiveItem(object)
}

func (c Character) resourceQuantity(name string) int {
	quantity := 0
	for _, object := range c.inventaire {
		resource, ok := object.(Resource)
		if ok && resource.nom == name {
			quantity += resource.quantité
		}
	}
	return quantity
}

func (c *Character) removeResource(name string, quantity int) {
	for index := 0; index < len(c.inventaire) && quantity > 0; index++ {
		resource, ok := c.inventaire[index].(Resource)
		if !ok || resource.nom != name {
			continue
		}
		if resource.quantité <= quantity {
			quantity -= resource.quantité
			c.inventaire = append(c.inventaire[:index], c.inventaire[index+1:]...)
			index--
			continue
		}
		resource.quantité -= quantity
		c.inventaire[index] = resource
		quantity = 0
	}
}

func (c *Character) TakePot(p Potion) {
	for i, objet := range c.inventaire {
		if objet.Nom() == p.nom {
			if p.effect != nil {
				p.effect(c)
			}

			c.inventaire = append(c.inventaire[:i], c.inventaire[i+1:]...)
			return
		}
	}
}

func (c *Character) accessInventory() {
	for i, objet := range c.inventaire {
		fmt.Println(i+1, "- "+objet.Nom())
	}
	var objetChoisi int
	fmt.Scanln(&objetChoisi)
	objetChoisi--
	if objetChoisi < 0 || objetChoisi >= len(c.inventaire) {
		return
	}
	objett := c.inventaire[objetChoisi]
	if armure, ok := objett.(Armure); ok {
		c.equipArmor(armure)
		return
	}
	if potion, ok := objett.(Potion); ok {
		c.TakePot(potion)
	} else {
		fmt.Println("L'objet sélectionné n'est pas utilisable")
	}
}

func (c *Character) isDead() {
	if c.pv_actuelle <= 0 {
		c.pv_actuelle = c.pv_total / 2

		if c.pv_actuelle < 1 {
			c.pv_actuelle = 1
		}

		fmt.Printf("Vous récupérez %d PV.\n", c.pv_actuelle)
	}
}

func (c *Character) inventoryLimit() bool {
	if len(c.inventaire) == c.maxslots {
		return false
	}
	return true
}

func (c *Character) upgradeInventorySlot() {
	c.maxslots += 10
}

func (c *Character) lvlUp() {
	c.niveau += 1
	c.exp_joueur -= c.exp_required
	c.exp_required *= 2
	c.pv_actuelle += 20
	c.pv_total += 20
}

func removeInventoryIndex(perso *Character, index int) {
	if perso == nil {
		return
	}

	if index < 0 || index >= len(perso.inventaire) {
		return
	}

	perso.inventaire = append(
		perso.inventaire[:index],
		perso.inventaire[index+1:]...,
	)
}
