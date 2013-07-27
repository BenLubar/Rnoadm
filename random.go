package main

type RandomSource uint64

func (s *RandomSource) Int63() int64 {
	const a = 6364136223846793005
	const c = 1442695040888963407
	const sixtythree = 1<<63 - 1

	// https://en.wikipedia.org/wiki/Linear_congruential_generator
	*s = a**s + c

	return int64(*s & sixtythree)
}

func (s *RandomSource) Seed(n int64) {
	*s = RandomSource(n)

	// obfuscate the seed a bit
	for i := 0; i < 1000; i++ {
		s.Int63()
	}
}
