package hero

type CraftItemType uint64

const (
	CraftLongHandle CraftItemType = iota
	CraftShortHandle
	CraftGrip
	CraftPommelCross
	CraftHilt
	CraftMaceHead
	CraftChain
	CraftHatchetHead
	CraftPickaxeHead
	CraftAxeHead
	CraftHammerHead
	CraftSpearHead
	CraftHalberdBlade
	CraftKnifeBlade
	CraftDaggerBlade
	CraftShortBlade
	CraftLongBlade
	CraftBroadBlade
	CraftThinBlade
	CraftHandleChain
	CraftMace
	CraftFlail
	CraftHatchet
	CraftPickaxe
	CraftBattleaxe
	CraftPoleaxe
	CraftWarhammer
	CraftSpear
	CraftHalberd
	CraftKnife
	CraftDagger
	CraftShortsword
	CraftLongsword
	CraftBroadsword
	CraftRapier
	CraftSwordChain
	CraftSwordChucks

	craftItemCount
)

type craftingIngredientType uint64

const (
	craftIngOtherIngredient craftingIngredientType = 0
	craftIngWood            craftingIngredientType = 1 << (63 - iota)
	craftIngStone
	craftIngMetal
)

type craftingIngredient struct {
	kind   craftingIngredientType
	volume uint64
	other  CraftItemType
}

var craftItemData = [craftItemCount]struct {
	name        string
	ingredientA *craftingIngredient
	ingredientB *craftingIngredient

	equip     bool
	equipSlot EquipSlot
	equipKind uint64
}{
	CraftLongHandle: {
		name: "long handle",
		ingredientA: &craftingIngredient{
			kind:   craftIngWood,
			volume: 4800,
		},
	},
	CraftShortHandle: {
		name: "short handle",
		ingredientA: &craftingIngredient{
			kind:   craftIngWood,
			volume: 1800,
		},
	},
	CraftGrip: {
		name: "grip",
		ingredientA: &craftingIngredient{
			kind:   craftIngWood | craftIngMetal,
			volume: 500,
		},
	},
	CraftPommelCross: {
		name: "pommel & cross",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 500,
		},
	},
	CraftHilt: {
		name: "hilt",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftGrip,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftPommelCross,
		},
	},
	CraftMaceHead: {
		name: "mace head",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 5000,
		},
	},
	CraftChain: {
		name: "chain",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 1000,
		},
	},
	CraftHatchetHead: {
		name: "hatchet head",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 600,
		},
	},
	CraftPickaxeHead: {
		name: "pickaxe head",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 600,
		},
	},
	CraftAxeHead: {
		name: "axe head",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 1500,
		},
	},
	CraftHammerHead: {
		name: "hammer head",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 4000,
		},
	},
	CraftSpearHead: {
		name: "spear head",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 500,
		},
	},
	CraftHalberdBlade: {
		name: "halberd blade",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 2500,
		},
	},
	CraftKnifeBlade: {
		name: "knife blade",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 4000,
		},
	},
	CraftDaggerBlade: {
		name: "dagger blade",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 4000,
		},
	},
	CraftShortBlade: {
		name: "short blade",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 6000,
		},
	},
	CraftLongBlade: {
		name: "long blade",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 8000,
		},
	},
	CraftBroadBlade: {
		name: "broad blade",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 10000,
		},
	},
	CraftThinBlade: {
		name: "thin blade",
		ingredientA: &craftingIngredient{
			kind:   craftIngMetal,
			volume: 3000,
		},
	},
	CraftHandleChain: {
		name: "handle & chain",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftChain,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftShortHandle,
		},
	},
	CraftMace: {
		name: "mace",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftMaceHead,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftShortHandle,
		},

		equip: true,
	},
	CraftFlail: {
		name: "flail",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftMaceHead,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftHandleChain,
		},

		equip: true,
	},
	CraftHatchet: {
		name: "hatchet",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftHatchetHead,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftShortHandle,
		},

		equip: true,
	},
	CraftPickaxe: {
		name: "pickaxe",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftPickaxeHead,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftShortHandle,
		},

		equip: true,
	},
	CraftBattleaxe: {
		name: "battleaxe",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftAxeHead,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftShortHandle,
		},

		equip: true,
	},
	CraftPoleaxe: {
		name: "poleaxe",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftAxeHead,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftLongHandle,
		},

		equip: true,
	},
	CraftWarhammer: {
		name: "warhammer",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftHammerHead,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftShortHandle,
		},

		equip: true,
	},
	CraftSpear: {
		name: "spear",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftSpearHead,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftLongHandle,
		},

		equip: true,
	},
	CraftHalberd: {
		name: "halberd",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftHalberdBlade,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftLongHandle,
		},

		equip: true,
	},
	CraftKnife: {
		name: "knife",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftKnifeBlade,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftGrip,
		},

		equip: true,
	},
	CraftDagger: {
		name: "dagger",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftDaggerBlade,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftHilt,
		},

		equip: true,
	},
	CraftShortsword: {
		name: "shortsword",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftShortBlade,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftHilt,
		},

		equip: true,
	},
	CraftLongsword: {
		name: "longsword",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftLongBlade,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftHilt,
		},

		equip: true,
	},
	CraftBroadsword: {
		name: "broadsword",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftBroadBlade,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftHilt,
		},

		equip: true,
	},
	CraftRapier: {
		name: "rapier",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftThinBlade,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftHilt,
		},

		equip: true,
	},
	CraftSwordChain: {
		name: "sword & chain",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftShortBlade,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftChain,
		},
	},
	CraftSwordChucks: {
		name: "sword-chucks",
		ingredientA: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftShortBlade,
		},
		ingredientB: &craftingIngredient{
			kind:  craftIngOtherIngredient,
			other: CraftSwordChain,
		},

		equip: true,
	},
}
