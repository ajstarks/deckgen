// Package generate performs slide deck generation
package generate

import (
	"fmt"
	"io"
)

const (
	circlefmt   = `<ellipse xp="%.2f" yp="%.2f" wp="%.2f" hr="%.2f" opacity="%.2f" color="%s"/>`
	squarefmt   = `<rect xp="%.2f" yp="%.2f" wp="%.2f" hr="%.2f" opacity="%.2f" color="%s"/>`
	ellipsefmt  = `<ellipse xp="%.2f" yp="%.2f" wp="%.2f" hp="%.2f" opacity="%.2f" color="%s"/>`
	rectfmt     = `<rect xp="%.2f" yp="%.2f" wp="%.2f" hp="%.2f" opacity="%.2f" color="%s"/>`
	arcfmt      = `<arc xp="%.2f" yp="%.2f" wp="%.2f" hp="%.2f" sp="%.2f" a1="%.2f" a2="%.2f" opacity="%.2f" color="%s"/>`
	linefmt     = `<line xp1="%.2f" yp1="%.2f" xp2="%.2f" yp2="%.2f" sp="%.2f" opacity="%.2f" color="%s"/>`
	curvefmt    = `<curve xp1="%.2f" yp1="%.2f" xp2="%.2f" yp2="%.2f" xp3="%.2f" yp3="%.2f" sp="%.2f" opacity="%.2f" color="%s"/>`
	polygonfmt  = `<polygon xc="%s" yc="%s" opacity="%.2f" color="%s"/>`
	polylinefmt = `<polyline xc="%s" yc="%s" sp="%.2f" opacity="%.2f" color="%s"/>`
	textfmt     = `<text xp="%.2f" yp="%.2f" sp="%.2f" align="%s" wp="%.2f" font="%s" opacity="%.2f" color="%s" type="%s">%s</text>`
	textlinkfmt = `<text xp="%.2f" yp="%.2f" sp="%.2f" align="%s" wp="%.2f" font="%s" opacity="%.2f" color="%s" type="%s" link="%s">%s</text>`
	textrotfmt  = `<text xp="%.2f" yp="%.2f" sp="%.2f" align="%s" wp="%.2f" font="%s" opacity="%.2f" color="%s" type="%s" link="%s" rotation="%.2f">%s</text>`
	imagefmt    = `<image xp="%.2f" yp="%.2f" width="%d" height="%d" name="%s" link="%s"/>`
	listfmt     = `<list type="%s" xp="%.2f" yp="%.2f" sp="%.2f" lp="%.2f" wp="%.2f" font="%s" color="%s">`
	lifmt       = `<li>%s</li>`
	closelist   = `</list>`
	slidefmt    = `<slide>`
	slidebg     = `<slide bg="%s">`
	slidebgfg   = `<slide bg="%s" fg="%s">`
	closeslide  = `</slide>`
	deckfmt     = `<deck><canvas width="%d" height="%d"/>`
	closedeck   = `</deck>`
)

// deckmarkup defines the structure of a presentation deck
// The size of the canvas, and series of slides
type Deck struct {
	Title       string  `xml:"title"`
	Creator     string  `xml:"creator"`
	Subject     string  `xml:"subject"`
	Publisher   string  `xml:"publisher"`
	Description string  `xml:"description"`
	Date        string  `xml:"date"`
	Canvas      canvas  `xml:"canvas"`
	Slide       []Slide `xml:"slide"`
}

type canvas struct {
	Width  int `xml:"width,attr"`
	Height int `xml:"height,attr"`
}

// Slide is the structure of an individual slide within a deck
// <slide bg="black" fg="rgb(255,255,255)" duration="2s" note="hello, world">
// <slide gradcolor1="black" gradcolor2="white" gp="20" duration="2s" note="wassup">
type Slide struct {
	Bg          string     `xml:"bg,attr"`
	Fg          string     `xml:"fg,attr"`
	Gradcolor1  string     `xml:"gradcolor1,attr"`
	Gradcolor2  string     `xml:"gradcolor2,attr"`
	GradPercent float64    `xml:"gp,attr"`
	Duration    string     `xml:"duration,attr"`
	Note        string     `xml:"note"`
	List        []List     `xml:"list"`
	Text        []Text     `xml:"text"`
	Image       []Image    `xml:"image"`
	Ellipse     []Ellipse  `xml:"ellipse"`
	Line        []Line     `xml:"line"`
	Rect        []Rect     `xml:"rect"`
	Curve       []Curve    `xml:"curve"`
	Arc         []Arc      `xml:"arc"`
	Polygon     []Polygon  `xml:"polygon"`
	Polyline    []Polyline `xml:"polyline"`
}

// CommonAttr are the common attributes for text and list
type CommonAttr struct {
	Xp          float64 `xml:"xp,attr"`         // X coordinate
	Yp          float64 `xml:"yp,attr"`         // Y coordinate
	Sp          float64 `xml:"sp,attr"`         // size
	Lp          float64 `xml:"lp,attr"`         // linespacing (leading) percentage
	Rotation    float64 `xml:"rotation,attr"`   // Rotation (0-360 degrees)
	Type        string  `xml:"type,attr"`       // type: block, plain, code, number, bullet
	Align       string  `xml:"align,attr"`      // alignment: center, end, begin
	Color       string  `xml:"color,attr"`      // item color
	Gradcolor1  string  `xml:"gradcolor1,attr"` // gradient color 1
	Gradcolor2  string  `xml:"gradcolor2,attr"` // gradient color 2
	GradPercent float64 `xml:"gp,attr"`         // gradient percentage
	Opacity     float64 `xml:"opacity,attr"`    // opacity percentage
	Font        string  `xml:"font,attr"`       // font type: i.e. sans, serif, mono
	Link        string  `xml:"link,attr"`       // reference to other content (i.e. http:// or mailto:)
}

// Dimension describes a graphics object with width and height
type Dimension struct {
	CommonAttr
	Wp float64 `xml:"wp,attr"` // width percentage
	Hp float64 `xml:"hp,attr"` // height percentage
	Hr float64 `xml:"hr,attr"` // height relative percentage
	Hw float64 `xml:"hw,attr"` // height by width
}

// ListItem describes a list item
// <list xp="20" yp="70" sp="1.5">
//
//	<li>canvas<li>
//	<li>slide</li>
//
// </list>
type ListItem struct {
	Color    string  `xml:"color,attr"`
	Opacity  float64 `xml:"opacity,attr"`
	Font     string  `xml:"font,attr"`
	ListText string  `xml:",chardata"`
}

// List describes the list element
type List struct {
	CommonAttr
	Wp float64    `xml:"wp,attr"`
	Li []ListItem `xml:"li"`
}

// Text describes the text element
type Text struct {
	CommonAttr
	Wp    float64 `xml:"wp,attr"`
	File  string  `xml:"file,attr"`
	Tdata string  `xml:",chardata"`
}

// Image describes an image
// <image xp="20" yp="30" width="256" height="256" scale="50" name="picture.png" caption="Pretty picture"/>
type Image struct {
	CommonAttr
	Width     int     `xml:"width,attr"`     // image width
	Height    int     `xml:"height,attr"`    // image height
	Scale     float64 `xml:"scale,attr"`     // image scale percentage
	Autoscale string  `xml:"autoscale,attr"` // scale the image to the canvas
	Name      string  `xml:"name,attr"`      // image file name
	Caption   string  `xml:"caption,attr"`   // image caption
}

// Ellipse describes a rectangle with x,y,w,h
// <ellipse xp="45"  yp="10" wp="4" hr="75" color="rgb(0,127,0)"/>
type Ellipse struct {
	Dimension
}

// Rect describes a rectangle with x,y,w,h
// <rect xp="35"  yp="10" wp="4" hp="3"/>
type Rect struct {
	Dimension
}

// Line defines a straight line
// <line xp1="20" yp1="10" xp2="30" yp2="10"/>
type Line struct {
	Xp1     float64 `xml:"xp1,attr"`     // begin x coordinate
	Yp1     float64 `xml:"yp1,attr"`     // begin y coordinate
	Xp2     float64 `xml:"xp2,attr"`     // end x coordinate
	Yp2     float64 `xml:"yp2,attr"`     // end y coordinate
	Sp      float64 `xml:"sp,attr"`      // line thickness
	Color   string  `xml:"color,attr"`   // line color
	Opacity float64 `xml:"opacity,attr"` // line opacity (1-100)
}

// Curve defines a quadratic Bezier curve
// The begining, ending, and control points are required:
// <curve xp1="60" yp1="10" xp2="75" yp2="20" xp3="70" yp3="10" />
type Curve struct {
	Xp1     float64 `xml:"xp1,attr"`
	Yp1     float64 `xml:"yp1,attr"`
	Xp2     float64 `xml:"xp2,attr"`
	Yp2     float64 `xml:"yp2,attr"`
	Xp3     float64 `xml:"xp3,attr"`
	Yp3     float64 `xml:"yp3,attr"`
	Sp      float64 `xml:"sp,attr"`
	Color   string  `xml:"color,attr"`
	Opacity float64 `xml:"opacity,attr"`
}

// Arc defines an elliptical arc
// the arc is defined by a beginning and ending angle in percentages
// <arc xp="55"  yp="10" wp="4" hr="75" a1="0" a2="180"/>
type Arc struct {
	Dimension
	A1      float64 `xml:"a1,attr"`
	A2      float64 `xml:"a2,attr"`
	Sp      float64 `xml:"sp,attr"`
	Opacity float64 `xml:"opacity,attr"`
}

// Polygon defines a polygon, x and y coordinates are specified by
// strings of space-separated percentages:
// <polygon xc="10 20 30" yc="30 40 50"/>
type Polygon struct {
	XC      string  `xml:"xc,attr"`
	YC      string  `xml:"yc,attr"`
	Color   string  `xml:"color,attr"`
	Opacity float64 `xml:"opacity,attr"`
}

// Polyline defines a polyline, x and y coordinates are specified by
// strings of space-separated percentages:
// <polyline xc="10 20 30" yc="30 40 50"/>
type Polyline struct {
	XC      string  `xml:"xc,attr"`
	YC      string  `xml:"yc,attr"`
	Sp      float64 `xml:"sp,attr"` // line thickness
	Color   string  `xml:"color,attr"`
	Opacity float64 `xml:"opacity,attr"`
}

// DeckGen is the generated deck structure.
type DeckGen struct {
	width, height int
	dest          io.Writer
}

// NewSlides initializes he generated deck structure.
func NewSlides(where io.Writer, w, h int) *DeckGen {
	return &DeckGen{dest: where, width: w, height: h}
}

// StartDeck begins a slide
func (p *DeckGen) StartDeck() {
	fmt.Fprintf(p.dest, deckfmt, p.width, p.height)
}

// EndDeck ends a slide.
func (p *DeckGen) EndDeck() {
	fmt.Fprintln(p.dest, closedeck)
}

// StartSlide begins a slide.
func (p *DeckGen) StartSlide(colors ...string) {
	switch len(colors) {
	case 1:
		fmt.Fprintf(p.dest, slidebg, colors[0])
	case 2:
		fmt.Fprintf(p.dest, slidebgfg, colors[0], colors[1])
	default:
		fmt.Fprintln(p.dest, slidefmt)
	}
}

// EndSlide ends a slide.
func (p *DeckGen) EndSlide() {
	fmt.Fprintln(p.dest, closeslide)
}

// square makes square markup from the rect structure.
func (p *DeckGen) square(r Rect) {
	fmt.Fprintf(p.dest, squarefmt, r.Xp, r.Yp, r.Wp, r.Hr, r.Opacity, r.Color)
}

// circle makes square markup from the ellipse structure.
func (p *DeckGen) circle(e Ellipse) {
	fmt.Fprintf(p.dest, circlefmt, e.Xp, e.Yp, e.Wp, e.Hr, e.Opacity, e.Color)
}

// ellipse makes ellipse markup from the ellipse structure.
func (p *DeckGen) ellipse(e Ellipse) {
	fmt.Fprintf(p.dest, ellipsefmt, e.Xp, e.Yp, e.Wp, e.Hp, e.Opacity, e.Color)
}

// rect makes rect markup rom the rect structure.
func (p *DeckGen) rect(r Rect) {
	fmt.Fprintf(p.dest, rectfmt, r.Xp, r.Yp, r.Wp, r.Hp, r.Opacity, r.Color)
}

// line makes line markup from the deck line structure.
func (p *DeckGen) line(l Line) {
	fmt.Fprintf(p.dest, linefmt, l.Xp1, l.Yp1, l.Xp2, l.Yp2, l.Sp, l.Opacity, l.Color)
}

// curve makes curve markup from the curve structure.
func (p *DeckGen) curve(c Curve) {
	fmt.Fprintf(p.dest, curvefmt, c.Xp1, c.Yp1, c.Xp2, c.Yp2, c.Xp3, c.Yp3, c.Sp, c.Opacity, c.Color)
}

// arc makes arc markup from the arc structure.
func (p *DeckGen) arc(a Arc) {
	fmt.Fprintf(p.dest, arcfmt, a.Xp, a.Yp, a.Wp, a.Hp, a.Sp, a.A1, a.A2, a.Opacity, a.Color)
}

// polygon makes polygon markup from the polygon structure.
func (p *DeckGen) polygon(poly Polygon) {
	fmt.Fprintf(p.dest, polygonfmt, poly.XC, poly.YC, poly.Opacity, poly.Color)
}

// polyline makes polyline markup from the polyline structure.
func (p *DeckGen) polyline(poly Polyline) {
	fmt.Fprintf(p.dest, polylinefmt, poly.XC, poly.YC, poly.Sp, poly.Opacity, poly.Color)
}

// text makes text markup from the deck text structure.
func (p *DeckGen) text(t Text) {
	fmt.Fprintf(p.dest, textfmt, t.Xp, t.Yp, t.Sp, t.Align, t.Wp, t.Font, t.Opacity, t.Color, t.Type, t.Tdata)
}

// textlink makes text markup from the deck text structure, including a link
func (p *DeckGen) textlink(t Text) {
	fmt.Fprintf(p.dest, textlinkfmt, t.Xp, t.Yp, t.Sp, t.Align, t.Wp, t.Font, t.Opacity, t.Color, t.Type, t.Link, t.Tdata)
}

// textrotate makes text markup from the deck text structure, including a link
func (p *DeckGen) textrotate(t Text) {
	fmt.Fprintf(p.dest, textrotfmt, t.Xp, t.Yp, t.Sp, t.Align, t.Wp, t.Font, t.Opacity, t.Color, t.Type, t.Link, t.Rotation, t.Tdata)
}

// image makes image markup from the deck image structure.
func (p *DeckGen) image(pic Image) {
	fmt.Fprintf(p.dest, imagefmt, pic.Xp, pic.Yp, pic.Width, pic.Height, pic.Name, pic.Link)
}

// list makes markup from the list deck structure.
func (p *DeckGen) list(l List, items []string, ltype, font, color string) {
	fmt.Fprintf(p.dest, listfmt, ltype, l.Xp, l.Yp, l.Sp, l.Lp, l.Wp, l.Font, l.Color)
	for _, s := range items {
		fmt.Fprintf(p.dest, lifmt, s)
	}
	fmt.Fprintln(p.dest, closelist)
}

// Text places plain text aligned at (x,y), with specified font, size and color. Opacity is optional
func (p *DeckGen) Text(x, y float64, s, font string, size float64, color string, opacity ...float64) {
	t := Text{}
	t.Xp = x
	t.Yp = y
	t.Sp = size
	t.Font = font
	t.Color = color
	t.Tdata = s
	if len(opacity) > 0 {
		t.Opacity = opacity[0]
	} else {
		t.Opacity = 100
	}
	p.text(t)
}

// TextMid places centered text aligned at (x,y), with specified font, size and color. Opacity is optional.
func (p *DeckGen) TextMid(x, y float64, s, font string, size float64, color string, opacity ...float64) {
	t := Text{}
	t.Xp = x
	t.Yp = y
	t.Sp = size
	t.Font = font
	t.Tdata = s
	t.Color = color
	t.Align = "center"
	if len(opacity) > 0 {
		t.Opacity = opacity[0]
	} else {
		t.Opacity = 100
	}
	p.text(t)
}

// TextEnd places right-justified text aligned at (x,y), with specified font, size and color. Opacity is optional.
func (p *DeckGen) TextEnd(x, y float64, s, font string, size float64, color string, opacity ...float64) {
	t := Text{}
	t.Xp = x
	t.Yp = y
	t.Sp = size
	t.Font = font
	t.Tdata = s
	t.Color = color
	t.Align = "right"
	if len(opacity) > 0 {
		t.Opacity = opacity[0]
	} else {
		t.Opacity = 100
	}
	p.text(t)
}

// TextBlock makes a block of text aligned at (x,y), wrapped at margin; with specified font, size and color. Opacity is optional.
func (p *DeckGen) TextBlock(x, y float64, s, font string, size, margin float64, color string, opacity ...float64) {
	t := Text{}
	t.Xp = x
	t.Yp = y
	t.Sp = size
	t.Font = font
	t.Wp = margin
	t.Tdata = s
	t.Color = color
	t.Type = "block"
	if len(opacity) > 0 {
		t.Opacity = opacity[0]
	} else {
		t.Opacity = 100
	}
	p.text(t)
}

// TextLink places text aligned at (x,y) with a link
func (p *DeckGen) TextLink(x, y float64, s, link, font string, size float64, color string, opacity ...float64) {
	t := Text{}
	t.Xp = x
	t.Yp = y
	t.Sp = size
	t.Font = font
	t.Tdata = s
	t.Color = color
	t.Link = link
	t.Type = "plain"
	if len(opacity) > 0 {
		t.Opacity = opacity[0]
	} else {
		t.Opacity = 100
	}
	p.textlink(t)
}

// TextRotate places rotated text
func (p *DeckGen) TextRotate(x, y float64, s, link, font string, rotation, size float64, color string, opacity ...float64) {
	t := Text{}
	t.Xp = x
	t.Yp = y
	t.Sp = size
	t.Font = font
	t.Tdata = s
	t.Color = color
	t.Link = link
	t.Rotation = rotation
	t.Type = "plain"
	if len(opacity) > 0 {
		t.Opacity = opacity[0]
	} else {
		t.Opacity = 100
	}
	p.textrotate(t)
}

// Code makes a code block at (x,y), with specified size and color (opacity is optional),
// on a light gray background with the specified margin width.
func (p *DeckGen) Code(x, y float64, s string, size, margin float64, color string, opacity ...float64) {
	t := Text{}
	t.Xp = x
	t.Yp = y
	t.Sp = size
	t.Wp = margin
	t.Tdata = s
	t.Color = color
	t.Type = "code"
	if len(opacity) > 0 {
		t.Opacity = opacity[0]
	} else {
		t.Opacity = 100
	}
	p.text(t)
}

// List makes a plain, bullet, or plain list with the specified font, size and color, with optional spacing
func (p *DeckGen) List(x, y, size, spacing, wrap float64, items []string, ltype, font, color string) {
	l := List{}
	l.Xp = x
	l.Yp = y
	l.Sp = size
	l.Lp = spacing
	l.Wp = wrap
	l.Font = font
	l.Color = color
	p.list(l, items, ltype, font, color)
}

// Square makes a square, centered at (x,y), with width w, at the specified color and optional opacity.
func (p *DeckGen) Square(x, y, w float64, color string, opacity ...float64) {
	r := Rect{}
	r.Xp = x
	r.Yp = y
	r.Wp = w
	r.Hr = 100
	r.Color = color
	if len(opacity) > 0 {
		r.Opacity = opacity[0]
	} else {
		r.Opacity = 100
	}
	p.square(r)
}

// Circle makes a circle, centered at (x,y) with width w, at the specified color and optional opacity.
func (p *DeckGen) Circle(x, y, w float64, color string, opacity ...float64) {
	e := Ellipse{}
	e.Xp = x
	e.Yp = y
	e.Wp = w
	e.Hr = 100
	e.Color = color
	if len(opacity) > 0 {
		e.Opacity = opacity[0]
	} else {
		e.Opacity = 100
	}
	p.circle(e)
}

// Rect makes a rectangle, centered at (x,y), with (w,h) dimensions, at the specified color and optional opacity.
func (p *DeckGen) Rect(x, y, w, h float64, color string, opacity ...float64) {
	r := Rect{}
	r.Xp = x
	r.Yp = y
	r.Wp = w
	r.Hp = h
	r.Color = color
	if len(opacity) > 0 {
		r.Opacity = opacity[0]
	} else {
		r.Opacity = 100
	}
	p.rect(r)
}

// Ellipse makes a ellipse graphic, centered at (x,y), with (w,h) dimensions, at the specified color and optional opacity.
func (p *DeckGen) Ellipse(x, y, w, h float64, color string, opacity ...float64) {
	e := Ellipse{}
	e.Xp = x
	e.Yp = y
	e.Wp = w
	e.Hp = h
	e.Color = color
	if len(opacity) > 0 {
		e.Opacity = opacity[0]
	} else {
		e.Opacity = 100
	}
	p.ellipse(e)
}

// Line makes a line from (x1,y1) to (x2, y2), with the specified color with optional opacity; thickness is size.
func (p *DeckGen) Line(x1, y1, x2, y2, size float64, color string, opacity ...float64) {
	l := Line{Xp1: x1, Xp2: x2, Yp1: y1, Yp2: y2, Sp: size, Color: color}
	if len(opacity) > 0 {
		l.Opacity = opacity[0]
	} else {
		l.Opacity = 100
	}
	p.line(l)
}

// Arc makes an arc centered at (x,y), with specified color (with optional opacity),
// with dimensions (w,h), between angle a1 and a2 (specified in degrees).
func (p *DeckGen) Arc(x, y, w, h, size, a1, a2 float64, color string, opacity ...float64) {
	a := Arc{A1: a1, A2: a2}
	a.Xp = x
	a.Yp = y
	a.Wp = w
	a.Hp = h
	a.Sp = size
	a.Color = color
	if len(opacity) > 0 {
		a.Opacity = opacity[0]
	} else {
		a.Opacity = 100
	}
	p.arc(a)
}

// Curve makes a Bezier curve between (x1, y2) and (x3, y3), with control points at (x2, y2), thickness is specified by size.
func (p *DeckGen) Curve(x1, y1, x2, y2, x3, y3, size float64, color string, opacity ...float64) {
	c := Curve{Xp1: x1, Xp2: x2, Xp3: x3, Yp1: y1, Yp2: y2, Yp3: y3, Sp: size, Color: color}
	if len(opacity) > 0 {
		c.Opacity = opacity[0]
	} else {
		c.Opacity = 100
	}
	p.curve(c)
}

// Polygon makes a polygon with the specified color (with optional opacity), with coordinates in x and y slices.
func (p *DeckGen) Polygon(x, y []float64, color string, opacity ...float64) {
	xc, yc := Polycoord(x, y)
	poly := Polygon{XC: xc, YC: yc, Color: color}
	if len(opacity) > 0 {
		poly.Opacity = opacity[0]
	}
	p.polygon(poly)
}

// Polyline makes a polyline with the specified color and thickness (with optional opacity), with coordinates in x and y slices.
func (p *DeckGen) Polyline(x, y []float64, size float64, color string, opacity ...float64) {
	xc, yc := Polycoord(x, y)
	poly := Polyline{XC: xc, YC: yc, Sp: size, Color: color}
	if len(opacity) > 0 {
		poly.Opacity = opacity[0]
	}
	p.polyline(poly)
}

// Polycoord converts slices of coordinates to strings.
func Polycoord(px, py []float64) (string, string) {
	var xc, yc string
	np := len(px)
	if np < 3 || len(py) != np {
		return xc, yc
	}
	for i := 0; i < np-1; i++ {
		xc += fmt.Sprintf("%.2f ", px[i])
		yc += fmt.Sprintf("%.2f ", py[i])
	}
	xc += fmt.Sprintf("%.2f", px[np-1])
	yc += fmt.Sprintf("%.2f", py[np-1])
	return xc, yc
}

// Image places the named image centered at (x, y), with dimensions of (w, h).
func (p *DeckGen) Image(x, y float64, w, h int, name, link string) {
	i := Image{Width: w, Height: h, Name: name}
	i.Xp = x
	i.Yp = y
	i.CommonAttr.Link = link
	p.image(i)
}
