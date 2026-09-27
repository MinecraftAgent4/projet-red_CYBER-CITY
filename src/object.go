package main

import "fmt"

type Object interface {
	Nom() string
}

type Arme interface {
	Nom() string
	Dmg() int
	DmgType() string
}

type Melee struct {
	nom string
	dmg int
}

type Ranged struct {
	nom string
	dmg int
}

func (i Melee) Dmg() int         { return i.dmg }
func (i Ranged) Dmg() int        { return i.dmg }
func (i Melee) DmgType() string  { return "melee" }
func (i Ranged) DmgType() string { return "ranged" }
func (i Melee) Nom() string      { return i.nom }
func (i Ranged) Nom() string     { return i.nom }

type Armure interface {
	Nom() string
	Defense() int
}

type Casque struct {
	nom     string
	defense int
}

type Plastron struct {
	nom     string
	defense int
}

type Bottes struct {
	nom     string
	defense int
}

func (a Casque) Nom() string    { return a.nom }
func (a Plastron) Nom() string  { return a.nom }
func (a Bottes) Nom() string    { return a.nom }
func (a Casque) Defense() int   { return a.defense }
func (a Plastron) Defense() int { return a.defense }
func (a Bottes) Defense() int   { return a.defense }

type Item struct {
	nom string
}

type Spell struct {
	nom    string
	effect func()
}

type Potion struct {
	nom    string
	effect func(target Entity)
}

type SpellBook struct {
	spell Spell
}

type Resource struct {
	nom          string
	quantité     int
	quantité_max int
}

func (i Potion) Nom() string    { return i.nom }
func (i Item) Nom() string      { return i.nom }
func (i SpellBook) Nom() string { return "Livre de Sort : " + i.spell.nom }
func (i Resource) Nom() string  { return fmt.Sprintf("%s: %d/%d", i.nom, i.quantité, i.quantité_max) }
