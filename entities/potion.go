package entities

type Potion struct {
	*Sprite
	HealingPower  uint
	CloseToPlayer bool
}
