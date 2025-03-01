package delaunay

import (
	"fmt"
	"iter"
	"math"
)

type Triangulation[P Point] struct {
	Points       []P
	ConvexHull   []P
	Triangles    []int
	Halfedges    []int
	NumTriangles int
}

// Triangulate returns a Delaunay triangulation of the provided points.
func Triangulate[S ~[]P, P Point](points S) (*Triangulation[P], error) {
	t := newTriangulator(points)
	err := t.triangulate()
	return &Triangulation[P]{
		points,
		t.convexHull(),
		t.triangles,
		t.halfedges,
		len(t.triangles) / 3,
	}, err
}

// TriangleIdx returns the points of the triangle i, where 0 <= i <=
// t.NumTriangles.
func (t *Triangulation[P]) Triangle(i int) [3]P {
	i *= 3
	return [3]P{
		t.Points[t.Triangles[i]],
		t.Points[t.Triangles[i+1]],
		t.Points[t.Triangles[i+2]],
	}
}

// TriangleIdx returns the point indices of the triangle i, where 0 <= i <=
// t.NumTriangles.
func (t *Triangulation[P]) TriangleIdx(i int) [3]int {
	i *= 3
	return [3]int{t.Triangles[i], t.Triangles[i+1], t.Triangles[i+2]}
}

// Edges returns an iterator that yields each edge at most once. If hull is
// false the outer "hull" edges are skipped. Otherwise all edges will be
// iterated.
func (t *Triangulation[P]) Edges(hull bool) iter.Seq2[P, P] {
	return func(yield func(P, P) bool) {
		for i, j := range t.EdgesIdx(hull) {
			if !yield(t.Points[i], t.Points[j]) {
				return
			}
		}
	}
}

func (t *Triangulation[P]) EdgesIdx(hull bool) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		call := func(i, j int) bool {
			if t.Halfedges[i] < 0 {
				return true
			}
			p, q := t.Triangles[i], t.Triangles[j]
			if p < q {
				if !yield(p, q) {
					return false
				}
			}
			return true
		}
		if hull {
			call = func(i, j int) bool {
				p, q := t.Triangles[i], t.Triangles[j]
				if p < q || t.Halfedges[i] < 0 {
					if !yield(p, q) {
						return false
					}
				}
				return true
			}
		}
		l := len(t.Triangles)
		for i := 0; i < l; i += 3 {
			if !call(i, i+1) {
				return
			}
			if !call(i+1, i+2) {
				return
			}
			if !call(i+2, i) {
				return
			}
		}
	}
}

func (t *Triangulation[P]) area() float64 {
	var result float64
	points := t.Points
	ts := t.Triangles
	for i := 0; i < len(ts); i += 3 {
		p0 := points[ts[i+0]]
		p1 := points[ts[i+1]]
		p2 := points[ts[i+2]]
		result += area(p0, p1, p2)
	}
	return result / 2
}

// Validate performs several sanity checks on the Triangulation to check for
// potential errors. Returns nil if no issues were found. You normally
// shouldn't need to call this function but it can be useful for debugging.
func (t *Triangulation[P]) Validate() error {
	// verify halfedges
	for i1, i2 := range t.Halfedges {
		if i1 != -1 && t.Halfedges[i1] != i2 {
			return fmt.Errorf("invalid halfedge connection")
		}
		if i2 != -1 && t.Halfedges[i2] != i1 {
			return fmt.Errorf("invalid halfedge connection")
		}
	}

	// verify convex hull area vs sum of triangle areas
	hull1 := t.ConvexHull
	hull2 := ConvexHull(t.Points)
	area1 := polygonArea(hull1)
	area2 := polygonArea(hull2)
	area3 := t.area()
	if math.Abs(area1-area2) > 1e-9 || math.Abs(area1-area3) > 1e-9 {
		return fmt.Errorf("hull areas disagree: %f, %f, %f", area1, area2, area3)
	}

	// verify convex hull perimeter
	perimeter1 := polygonPerimeter(hull1)
	perimeter2 := polygonPerimeter(hull2)
	if math.Abs(perimeter1-perimeter2) > 1e-9 {
		return fmt.Errorf("hull perimeters disagree: %f, %f", perimeter1, perimeter2)
	}

	return nil
}
