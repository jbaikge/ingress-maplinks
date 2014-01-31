package delaunay

import (
	"image"
	"testing"
)

func TestDelaunay(t *testing.T) {
	tests := []struct {
		Points    []image.Point
		Triangles []*Triangle
	}{
		{
			[]image.Point{image.Pt(3, 0), image.Pt(5, 5), image.Pt(0, 2)},
			[]*Triangle{
				NewTriangle(&image.Point{0, 2}, &image.Point{3, 0}, &image.Point{5, 5}),
			},
		},
		{
			[]image.Point{image.Pt(3, 0), image.Pt(5, 5), image.Pt(0, 2), image.Pt(-1, -4)},
			[]*Triangle{
				NewTriangle(&image.Point{-1, -4}, &image.Point{0, 2}, &image.Point{3, 0}),
				NewTriangle(&image.Point{0, 2}, &image.Point{3, 0}, &image.Point{5, 5}),
			},
		},
		{
			[]image.Point{image.Pt(3, 0), image.Pt(5, 5), image.Pt(0, 2), image.Pt(-1, -4), image.Pt(1, 4)},
			[]*Triangle{
				NewTriangle(&image.Point{-1, -4}, &image.Point{0, 2}, &image.Point{3, 0}),
				NewTriangle(&image.Point{0, 2}, &image.Point{1, 4}, &image.Point{3, 0}),
				NewTriangle(&image.Point{1, 4}, &image.Point{3, 0}, &image.Point{5, 5}),
			},
		},
	}
	for _, test := range tests {
		out := Triangulate(test.Points)
		if out == nil {
			t.Fatal("Too few points? %v", test.Points)
		}
		if len(out) != len(test.Triangles) {
			t.Errorf("Improper number of triangles found: %d, expected %d", len(out), len(test.Triangles))
			continue
		}
		for i := range out {
			if !out[i].A.Eq(*test.Triangles[i].A) || !out[i].B.Eq(*test.Triangles[i].B) || !out[i].C.Eq(*test.Triangles[i].C) {
				t.Errorf("Improper vertex at %d", i)
				t.Errorf("    Expect: %s", test.Triangles[i])
				t.Errorf("    Got:    %s", out[i])
			}
			t.Logf("T: %s", out[i])
		}
	}
}
