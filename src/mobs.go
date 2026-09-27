package main

type Monster struct {
	name        string
	pv_total    int
	pv_actuelle int
	weapon      Arme
	armor       Armure
	moneyReward int
}

type Entity interface {
	Name() string
	MaxHp() int
	Hp() int
	Atk() int
	Def() int
	AddPV(n int)
	Dmg(n int)
}

func (mob *Monster) Name() string { return mob.name }
func (mob *Monster) MaxHp() int   { return mob.pv_total }
func (mob *Monster) Hp() int      { return mob.pv_actuelle }
func (mob *Monster) Atk() int     { return mob.weapon.Dmg() }
func (mob *Monster) Def() int {
	if mob.armor == nil {
		return 0
	}
	return mob.armor.Defense()
}
func (mob *Monster) Dmg(n int) {
	mob.pv_actuelle -= n
	if mob.pv_actuelle < 0 {
		mob.pv_actuelle = 0
	}
}
func (mob *Monster) AddPV(n int) {
	mob.pv_actuelle += n
	if mob.pv_actuelle > mob.pv_total {
		mob.pv_actuelle = mob.pv_total
	}
}

func (guy *Character) Name() string { return guy.nom }
func (guy *Character) MaxHp() int   { return guy.pv_total }
func (guy *Character) Hp() int      { return guy.pv_actuelle }
func (guy *Character) Atk() int {
	if guy.weapon == nil {
		return 1
	}
	return guy.weapon.Dmg()
}
func (guy *Character) Def() int { return guy.armorValue() }
func (guy *Character) Dmg(n int) {
	guy.pv_actuelle -= n
	if guy.pv_actuelle < 0 {
		guy.pv_actuelle = 0
	}
}
func (guy *Character) AddPV(n int) {
	guy.pv_actuelle += n
	if guy.pv_actuelle > guy.pv_total {
		guy.pv_actuelle = guy.pv_total
	}
}

func initMonster(nom string, hp int, arme Arme, armure Armure) Monster {
	return Monster{name: nom, pv_total: hp, pv_actuelle: hp, weapon: arme, armor: armure}
}

func initGoblin() Monster {
	return initMonster("Gangster", 40, Melee{"Dague", 5}, nil)
}

func initBoss1() Monster {
	return initMonster("Boss de gang", 150, Melee{"Lame du boss", 12}, Casque{"Casque blindé", 4})
}
