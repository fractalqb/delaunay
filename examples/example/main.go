package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/fogleman/delaunay"
	"github.com/fogleman/delaunay/internal"
	"github.com/fogleman/gg"
	"github.com/fogleman/poissondisc"
)

const (
	W = 2048
	H = 2048
	N = 5000
)

func generatePoints() []internal.Point {
	s := math.Sqrt(float64(N) * 1.618)
	points := poissondisc.Sample(-s, -s, s, s, 1, 32, nil)
	sort.Slice(points, func(i, j int) bool {
		p1 := points[i]
		p2 := points[j]
		d1 := math.Hypot(p1.X, p1.Y)
		d2 := math.Hypot(p2.X, p2.Y)
		return d1 < d2
	})
	points = points[:N]
	result := make([]internal.Point, len(points))
	for i, p := range points {
		result[i].Px = p.X
		result[i].Py = p.Y
	}
	return result
}

func main() {
	// generate points
	points := generatePoints()

	// triangulate
	start := time.Now()
	triangulation, err := delaunay.Triangulate(points)
	elapsed := time.Since(start)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(elapsed)
	fmt.Println(len(triangulation.Triangles) / 3)

	// compute point bounds for rendering
	pmin := points[0]
	pmax := points[0]
	for _, p := range points {
		pmin.Px = min(pmin.X(), p.X())
		pmin.Py = min(pmin.Y(), p.Y())
		pmax.Px = max(pmax.X(), p.X())
		pmax.Py = max(pmax.Y(), p.Y())
	}

	size := internal.Point{Px: pmax.X() - pmin.X(), Py: pmax.Y() - pmin.Y()}
	center := internal.Point{Px: pmin.X() + size.X()/2, Py: pmin.Y() + size.Y()/2}
	scale := math.Min(W/size.X(), H/size.Y()) * 0.9

	// render points and edges
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	dc.Translate(W/2, H/2)
	dc.Scale(scale, scale)
	dc.Translate(-center.X(), -center.Y())

	for p, q := range triangulation.Edges(false) {
		dc.DrawLine(p.X(), p.Y(), q.X(), q.Y())
	}
	dc.Stroke()

	for _, p := range points {
		dc.DrawPoint(p.X(), p.Y(), 5)
	}
	dc.Fill()

	for _, p := range triangulation.ConvexHull {
		dc.LineTo(p.X(), p.Y())
	}
	dc.ClosePath()
	dc.SetLineWidth(5)
	dc.Stroke()

	dc.SavePNG("out.png")
}
