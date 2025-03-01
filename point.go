package delaunay

import (
	"math"
)

type Point interface {
	X() float64
	Y() float64
}

func squaredDistance[P Point](a, b P) float64 {
	dx := a.X() - b.X()
	dy := a.Y() - b.Y()
	return dx*dx + dy*dy
}

func distance[P Point](a, b P) float64 {
	return math.Hypot(a.X()-b.X(), a.Y()-b.Y())
}

func sub[P Point](a, b P) pt { return pt{a.X() - b.X(), a.Y() - b.Y()} }

type pt struct{ x, y float64 }

func (p pt) X() float64 { return p.x }
func (p pt) Y() float64 { return p.y }
