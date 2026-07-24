package randomnick

// 50 x 20 = 1000 组内置英文昵称(数据库无数据时使用)
var (
	builtinAdjectives = []string{
		"Swift", "Calm", "Bold", "Bright", "Cool", "Dark", "Epic", "Fast", "Free", "Gold",
		"Happy", "Iron", "Jade", "Keen", "Lucky", "Magic", "Noble", "Quick", "Royal", "Sharp",
		"Silent", "Smart", "Soft", "Star", "True", "Wild", "Wise", "Brave", "Clear", "Deep",
		"Fair", "Grand", "High", "Kind", "Light", "Neat", "Pure", "Rich", "Safe", "Warm",
		"Young", "Zen", "Ace", "Blue", "Cyber", "Dream", "Echo", "Frost", "Glow", "Hyper",
	}
	builtinNouns = []string{
		"Fox", "Wolf", "Bear", "Hawk", "Lion", "Tiger", "Eagle", "Owl", "Moon", "Sun",
		"Cloud", "Storm", "River", "Ocean", "Peak", "Flame", "Spark", "Knight", "Ranger", "Spirit",
	}
)

func buildBuiltinEnglishNicknames() []string {
	out := make([]string, 0, len(builtinAdjectives)*len(builtinNouns))
	for _, adj := range builtinAdjectives {
		for _, noun := range builtinNouns {
			out = append(out, adj+noun)
		}
	}
	return out
}
