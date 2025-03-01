package delaunay

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/fogleman/delaunay/internal"
)

func Benchmark_triangulator(b *testing.B) {
	dists := map[string]dist{
		"uniform": internal.Uniform,
		"normal":  internal.Normal,
		"grid":    internal.Grid,
		"circle":  internal.Circle,
	}
	for name, dist := range dists {
		for n := 10; n <= 10000; n *= 10 {
			b.Run(
				fmt.Sprintf("%d × %s", n, name),
				func(b *testing.B) { benchmark(b, dist, n) },
			)
		}
	}
}

type dist func(n int, rnd *rand.Rand) []internal.Point

func benchmark(b *testing.B, f dist, n int) {
	rnd := rand.New(&rand.PCG{})
	points := f(n, rnd)
	for b.Loop() {
		if _, err := Triangulate(points); err != nil {
			b.Fatal(err)
		}
	}
}
