package main

import (
	"flag"
	"github.com/jbaikge/ingress-maplinks/delaunay"
	"image"
	"image/color"
	"log"
	"os"

	_ "image/png"
)

var (
	borderColor = flag.Int("border", 16750848, "Portal border color (Default: #FF9900)")
	size        = flag.String("size", "16", "Portal diameter")

	portals = make([]image.Rectangle, 0, 512)
)

func main() {
	flag.Parse()
	log.SetFlags(log.Lmicroseconds)

	targetColor := color.RGBA{
		R: uint8((*borderColor & 0xFF0000) >> 16),
		G: uint8((*borderColor & 0x00FF00) >> 8),
		B: uint8((*borderColor & 0x0000FF)),
	}
	log.Printf("Target color: %s", targetColor)

	inFilename, outFilename := flag.Arg(0), flag.Arg(1)
	in, err := os.Open(inFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	imgIn, _, err := image.Decode(in)
	if err != nil {
		log.Fatal(err)
	}
	bounds := imgIn.Bounds()

	imgOut := NewSVG(bounds)
	imgOut.Elements = append(imgOut.Elements, Image{Width: bounds.Dx(), Height: bounds.Dy(), Href: inFilename})

	log.Printf("Searching image for portals")

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if colorEq(imgIn.At(x, y), targetColor) {
				addPortal(image.Pt(x, y))
			}
		}
	}

	log.Printf("Found %d portals", len(portals))

	verticies := make([]image.Point, len(portals))
	for i, p := range portals {
		//imgOut.AddPortal(p)
		verticies[i] = image.Pt(p.Min.X+p.Dx()/2, p.Min.Y+p.Dy()/2)
	}

	log.Println("Triangulating fields")
	fields := delaunay.Triangulate(verticies)
	log.Printf("Found %d fields", len(fields))

	log.Printf("Drawing links")
	drawn := make(map[string]map[string]bool, len(fields)/3)
	drawnCount := 0
	for i := range fields {
		for j, edges := 0, fields[i].Edges(); j < len(edges); j += 2 {
			a, b := edges[j], edges[j+1]
			aS, bS := a.String(), b.String()
			y, ok := drawn[aS]
			if !ok {
				drawn[aS] = make(map[string]bool, len(fields)/3)
				y = drawn[aS]
			}
			if y[bS] {
				continue
			}
			imgOut.AddLink(*a, *b)
			drawn[aS][bS] = true
			drawnCount++
		}
	}
	log.Printf("Drew %d links", drawnCount)

	log.Printf("Saving to %s", outFilename)
	out, err := os.OpenFile(outFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	imgOut.WriteTo(out)
	log.Printf("Done")
}

func addPortal(p image.Point) {
	for _, portal := range portals {
		if p.In(portal) {
			return
		}
	}
	portals = append(portals, image.Rect(p.X-7, p.Y, p.X+9, p.Y+16))
}

func colorEq(a, b color.Color) bool {
	aR, aG, aB, _ := a.RGBA()
	bR, bG, bB, _ := b.RGBA()
	return aR == bR && aG == bG && aB == bB
}
