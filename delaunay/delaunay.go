package delaunay

import (
	"fmt"
	"image"
	"sort"
)

type byX []image.Point

type Triangle struct {
	A, B, C *image.Point
	R, X, Y float64
}

func NewTriangle(a, b, c *image.Point) (t *Triangle) {
	t = &Triangle{A: a, B: b, C: c}
	A := b.X - a.X
	B := b.Y - a.Y
	C := c.X - a.X
	D := c.Y - a.Y
	E := A*(a.X+b.X) + B*(a.Y+b.Y)
	F := C*(a.X+c.X) + D*(a.Y+c.Y)
	G := 2 * (A*(c.Y-b.Y) - B*(c.X-b.X))

	if abs(G) == 0 {
		minx := float64(min(a.X, min(b.X, c.X)))
		miny := float64(min(a.Y, min(b.Y, c.Y)))
		dx := (float64(max(a.X, max(b.X, c.X))) - minx) / 2
		dy := (float64(max(a.Y, max(b.Y, c.Y))) - miny) / 2

		t.X = minx + dx
		t.Y = miny + dy
		t.R = dx*dx + dy*dy
	} else {
		t.X = float64(D*E-B*F) / float64(G)
		t.Y = float64(A*F-C*E) / float64(G)
		dx := t.X - float64(a.X)
		dy := t.Y - float64(a.Y)
		t.R = dx*dx + dy*dy
	}
	return
}

func (t *Triangle) Edges() (e []*image.Point) {
	order := func(a, b *image.Point) (m, n *image.Point) {
		if (a.X < b.X) || (a.X == b.X && a.Y < b.Y) {
			return a, b
		}
		return b, a
	}
	e = make([]*image.Point, 6)
	e[0], e[1] = order(t.A, t.B)
	e[2], e[3] = order(t.B, t.C)
	e[4], e[5] = order(t.C, t.A)
	return
}

func (t *Triangle) String() string {
	return fmt.Sprintf("A%8s B%8s C%8s X%8.3f Y%8.3f R%8.3f", t.A, t.B, t.C, t.X, t.Y, t.R)
}

var _ sort.Interface = byX{}

func (by byX) Len() int           { return len(by) }
func (by byX) Swap(i, j int)      { by[i], by[j] = by[j], by[i] }
func (by byX) Less(i, j int) bool { return by[i].X > by[j].X }

func Triangulate(verticies []image.Point) []*Triangle {
	if len(verticies) < 3 {
		return nil
	}
	sort.Sort(byX(verticies))

	xmin := verticies[len(verticies)-1].X
	xmax := verticies[0].X
	ymin := verticies[len(verticies)-1].Y
	ymax := ymin

	for _, v := range verticies {
		if v.Y < ymin {
			ymin = v.Y
		}
		if v.Y > ymax {
			ymax = v.Y
		}
	}

	dx := xmax - xmin
	dy := ymax - ymin
	dmax := max(dx, dy)
	xmid := (xmax + xmin) / 2
	ymid := (ymax + ymin) / 2
	st := NewTriangle(
		&image.Point{X: xmid - 20*dmax, Y: ymid - dmax},
		&image.Point{X: xmid, Y: ymax + 20*dmax},
		&image.Point{X: xmid + 20*dmax, Y: ymid - dmax},
	)
	open := make([]*Triangle, 1, len(verticies))
	open[0] = st
	closed := make([]*Triangle, 0, len(verticies))
	edges := make([]*image.Point, 0, len(verticies)*6)

	for i := len(verticies) - 1; i >= 0; i-- {
		edges = edges[:0]
		for j := len(open) - 1; j >= 0; j-- {
			dx := float64(verticies[i].X) - open[j].X
			if dx > 0 && dx*dx > open[j].R {
				closed = append(closed, open[j])
				open = append(open[:j], open[j+1:]...)
				continue
			}
			dy := float64(verticies[i].Y) - open[j].Y
			if dx*dx+dy*dy > open[j].R {
				continue
			}
			edges = append(edges, open[j].Edges()...)
			open = append(open[:j], open[j+1:]...)
		}
		edges = dedupe(edges)
		for j := len(edges) - 2; j >= 0; j -= 2 {
			open = append(open, NewTriangle(edges[j], edges[j+1], &verticies[i]))
		}
	}
	closed = append(closed, open...)
	for i := len(closed) - 1; i >= 0; i-- {
		c := closed[i]
		if c.A == st.A || c.A == st.B || c.A == st.C ||
			c.B == st.A || c.B == st.B || c.B == st.C ||
			c.C == st.A || c.C == st.B || c.C == st.C {
			closed = append(closed[:i], closed[i+1:]...)
		}
	}
	return closed
}

func dedupe(edges []*image.Point) []*image.Point {
OuterLoop:
	for j := len(edges) - 2; j >= 0; j -= 2 {
		a, b := edges[j], edges[j+1]
		for i := j - 2; i >= 0; i -= 2 {
			m, n := edges[i], edges[i+1]
			if (a == m && b == n) || (a == n && b == m) {
				edges = append(edges[:j], edges[j+2:]...)
				edges = append(edges[:i], edges[i+2:]...)
				j -= 2
				continue OuterLoop
			}
		}
	}
	return edges
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
