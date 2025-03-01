package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"reflect"
	"runtime"
	"time"

	"github.com/fogleman/delaunay"
	"github.com/fogleman/delaunay/internal"
)

type dist func(n int, rnd *rand.Rand) []internal.Point

func getFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func test(f dist, n int) {
	rnd := rand.New(&rand.PCG{})
	points := f(n, rnd)
	start := time.Now()
	_, err := delaunay.Triangulate(points)
	elapsed := time.Since(start)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n, elapsed)
}

func main() {
	dists := []dist{
		internal.Uniform,
		internal.Normal,
		internal.Grid,
		internal.Circle,
	}
	for _, f := range dists {
		fmt.Println(getFunctionName(f))
		for n := 10; n <= 1000000; n *= 10 {
			test(f, n)
		}
		fmt.Println()
	}
}
