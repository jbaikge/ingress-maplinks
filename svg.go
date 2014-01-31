package main

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"io"
)

type SVG struct {
	XMLName  xml.Name `xml:"svg"`
	Version  string   `xml:"version,attr"`
	Height   int      `xml:"height,attr"`
	Width    int      `xml:"width,attr"`
	XMLNS    string   `xml:"xmlns,attr"`
	Xlink    string   `xml:"xmlns:xlink,attr"`
	Elements []Element
}

func NewSVG(r image.Rectangle) *SVG {
	return &SVG{
		Version:  "1.1",
		Height:   r.Dy(),
		Width:    r.Dx(),
		Elements: make([]Element, 0, 1024),
		XMLNS:    "http://www.w3.org/2000/svg",
		Xlink:    "http://www.w3.org/1999/xlink",
	}
}

func (s *SVG) AddLink(a, b image.Point) {
	l := new(Line)
	l.SetColor(color.RGBA{255, 0, 0, 255})
	l.Link(a, b)
	s.Elements = append(s.Elements, l)
}

func (s *SVG) AddPortal(r image.Rectangle) {
	c := new(Circle)
	c.SetColor(color.Black)
	c.SetSize(r)
	s.Elements = append(s.Elements, c)
}

func (s *SVG) WriteTo(w io.Writer) (n int64, err error) {
	enc := xml.NewEncoder(w)
	enc.Indent("", "\t")
	err = enc.Encode(s)
	return
}

type Element interface {
}

////////////////////////////////////////////////////////////////////////////////

type Line struct {
	XMLName xml.Name `xml:"line"`
	X1      int      `xml:"x1,attr"`
	X2      int      `xml:"x2,attr"`
	Y1      int      `xml:"y1,attr"`
	Y2      int      `xml:"y2,attr"`
	Style   string   `xml:"style,attr"`
}

func (e *Line) Link(a, b image.Point) {
	e.X1 = a.X
	e.Y1 = a.Y
	e.X2 = b.X
	e.Y2 = b.Y
}

func (e *Line) SetColor(c color.Color) {
	r, g, b, _ := c.RGBA()
	e.Style = fmt.Sprintf("stroke:rgb(%d,%d,%d);stroke-width:1.5", r>>8, g>>8, b>>8)
}

////////////////////////////////////////////////////////////////////////////////

type Circle struct {
	XMLName xml.Name `xml:"circle"`
	X       int      `xml:"cx,attr"`
	Y       int      `xml:"cy,attr"`
	Radius  int      `xml:"r,attr"`
	Style   string   `xml:"style,attr"`
}

func (e *Circle) SetColor(c color.Color) {
	r, g, b, _ := c.RGBA()
	e.Style = fmt.Sprintf("stroke:rgb(%d,%d,%d);stroke-width:1;opacity:0.5;", r>>8, g>>8, b>>8)
}

func (e *Circle) SetSize(r image.Rectangle) {
	e.Radius = r.Dx() / 2
	e.X = r.Min.X + e.Radius
	e.Y = r.Min.Y + e.Radius
}

////////////////////////////////////////////////////////////////////////////////

type Image struct {
	XMLName xml.Name `xml:"image"`
	Href    string   `xml:"xlink:href,attr"`
	X       int      `xml:"x,attr"`
	Y       int      `xml:"y,attr"`
	Width   int      `xml:"width,attr"`
	Height  int      `xml:"height,attr"`
}
