package util

import (
	"math/rand"
	"time"
)

type Random struct {
	random *rand.Rand
}

func NewRandom() *Random {
	return &Random{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *Random) RandomNum(max int) int {
	if max > 0 {
		return g.random.Intn(max)
	}
	return 0
}

func (g *Random) RandomInt63n(max int64) int64 {
	if max > 0 {
		return g.random.Int63n(max)
	}
	return 0
}

func (g *Random) RandomInt63() int64 {
	return g.random.Int63()
}

func (g *Random) RandomFloat() float64 {
	return g.random.Float64()
}

func (g *Random) RandomBool() bool {
	return g.RandomNum(2) == 0
}

func (g *Random) RandomRange(n int, m int) int {
	if n == m {
		return n
	}

	if m < n {
		n, m = m, n
	}

	return n + g.RandomNum(m-n)
}
