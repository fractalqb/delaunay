package internal

import (
	"math"
	"math/rand/v2"
)

type Point struct {
	Px, Py float64
}

func (p Point) X() float64             { return p.Px }
func (p Point) Y() float64             { return p.Py }
func (p Point) New(x, y float64) Point { return Point{x, y} }

func Uniform(n int, rnd *rand.Rand) []Point {
	points := make([]Point, n)
	for i := range points {
		x := rnd.Float64()
		y := rnd.Float64()
		points[i] = Point{Px: x, Py: y}
	}
	return points
}

func Normal(n int, rnd *rand.Rand) []Point {
	points := make([]Point, n)
	for i := range points {
		x := rnd.NormFloat64()
		y := rnd.NormFloat64()
		points[i] = Point{Px: x, Py: y}
	}
	return points
}

func Grid(n int, rnd *rand.Rand) []Point {
	side := int(math.Floor(math.Sqrt(float64(n))))
	n = side * side
	points := make([]Point, 0, n)
	for y := 0; y < side; y++ {
		for x := 0; x < side; x++ {
			p := Point{Px: float64(x), Py: float64(y)}
			points = append(points, p)
		}
	}
	return points
}

func Circle(n int, rnd *rand.Rand) []Point {
	points := make([]Point, n)
	for i := range points {
		t := float64(i) / float64(n)
		x := math.Cos(t)
		y := math.Sin(t)
		points[i] = Point{Px: x, Py: y}
	}
	return points
}
