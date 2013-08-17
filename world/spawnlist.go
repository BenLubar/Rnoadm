package world

var spawnFuncs []func(string) Visible

func RegisterSpawnFunc(f func(string) Visible) {
	spawnFuncs = append(spawnFuncs, f)
}

func Spawn(s string) Visible {
	for _, f := range spawnFuncs {
		if v := f(s); v != nil {
			return v
		}
	}
	return nil
}
