//special thanks to John Cox, who optimized this code to run faster and generate smaller GIFs.

package main

import (
	"C"
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
)

// ImagesToGIF() takes a slice of images and uses them to generate an animated GIF
// with the name "filename.out.gif" where filename is an input parameter.
func ImagesToGIF(imglist []image.Image, filename string) {

	// get ready to write images to files
	w, err := os.Create(filename + ".out.gif")

	if err != nil {
		fmt.Println("Sorry: couldn't create the file!")
		os.Exit(1)
	}

	defer w.Close()
	var g gif.GIF
	g.Delay = make([]int, len(imglist))
	g.Image = make([]*image.Paletted, len(imglist))
	g.LoopCount = 100

	for i := range imglist {
		g.Image[i] = ImageToPaletted(imglist[i])
		g.Delay[i] = 10
	}

	gif.EncodeAll(w, &g)
}

// ImageToPaletted converts an image to an image.Paletted with 256 colors.
// It is used by a subroutine by process() to generate an animated GIF.
func ImageToPalettedVersion1(img image.Image) *image.Paletted {
	pm, ok := img.(*image.Paletted)
	if !ok {
		b := img.Bounds()
		pm = image.NewPaletted(b, palette.WebSafe)
		draw.Draw(pm, pm.Bounds(), img, image.Point{}, draw.Src)
	}
	return pm
}

var mapOfColorIndices map[color.Color]uint8

func init() {
	mapOfColorIndices = make(map[color.Color]uint8)
}

func ImageToPaletted(img image.Image) *image.Paletted {
	pm, ok := img.(*image.Paletted)
	if !ok {
		b := img.Bounds()
		pm = image.NewPaletted(b, palette.WebSafe)
		var prevC color.Color = nil
		var idx uint8
		var ok bool
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				c := img.At(x, y)
				if c != prevC {
					if idx, ok = mapOfColorIndices[c]; !ok {
						idx = uint8(pm.Palette.Index(c))
						mapOfColorIndices[c] = idx
					}
					prevC = c
				}
				i := pm.PixOffset(x, y)
				pm.Pix[i] = idx
			}
		}
	}
	return pm
}

type Canvas struct {
	gc     *draw2dimg.GraphicContext
	img    image.Image
	width  int // both width and height are in pixels
	height int
}

func (c *Canvas) GetImage() image.Image {
	return c.img
}

// Create a new canvas
func CreateNewCanvas(w, h int) Canvas {
	i := image.NewRGBA(image.Rect(0, 0, w, h))
	gc := draw2dimg.NewGraphicContext(i)

	gc.SetStrokeColor(image.Black)
	gc.SetFillColor(image.White)
	// fill the background
	gc.Clear()
	gc.SetFillColor(image.Black)

	return Canvas{gc, i, w, h}
}

// Create a new Paletted canvas
func CreateNewPalettedCanvas(w, h int, cp color.Palette) Canvas {
	if cp == nil {
		cp = palette.WebSafe
	}
	i := image.NewPaletted(image.Rect(0, 0, w, h), cp)
	gc := draw2dimg.NewGraphicContext(i)

	gc.SetStrokeColor(image.Black)
	gc.SetFillColor(image.White)
	// fill the background
	gc.Clear()
	gc.SetFillColor(image.Black)

	return Canvas{gc, i, w, h}
}

// Create a new color
func MakeColor(r, g, b uint8) color.Color {
	return &color.RGBA{r, g, b, 255}
}

// Move the current point to (x,y)
func (c *Canvas) MoveTo(x, y float64) {
	c.gc.MoveTo(x, y)
}

// Draw a line from the current point to (x,y), and set the current point to (x,y)
func (c *Canvas) LineTo(x, y float64) {
	c.gc.LineTo(x, y)
}

// Draw an arc from the current point to (x, y)
// Can be used to easily draw a circle or an ellipse
func (c *Canvas) ArcTo(x, y, radiusX, radiusY, degStart, degEnd float64) {
	c.gc.ArcTo(x, y, radiusX, radiusY, degStart, degEnd)
}

// Set the line color
func (c *Canvas) SetStrokeColor(col color.Color) {
	c.gc.SetStrokeColor(col)
}

// Set the fill color
func (c *Canvas) SetFillColor(col color.Color) {
	c.gc.SetFillColor(col)
}

// Set the line width
func (c *Canvas) SetLineWidth(w float64) {
	c.gc.SetLineWidth(w)
}

// Actually draw the lines you've set up with LineTo
func (c *Canvas) Stroke() {
	c.gc.Stroke()
}

// Fill the area inside the lines you've set up with LineTo
func (c *Canvas) FillStroke() {
	c.gc.FillStroke()
}

// Fill the area inside the lines you've set up with LineTo, but don't
// draw the lines
func (c *Canvas) Fill() {
	c.gc.Fill()
}

// Fill the whole canvas with the fill color
func (c *Canvas) Clear() {
	c.gc.Clear()
}

// Fill the given rectangle with the fill color
func (c *Canvas) ClearRect(x1, y1, x2, y2 int) {
	c.gc.ClearRect(x1, y1, x2, y2)
}

// Draws an empty circle
// Fill the given circle with the fill color
// Stroke() each time to avoid connected circles
func (c *Canvas) Circle(cx, cy, r float64) {
	c.gc.ArcTo(cx, cy, r, r, 0, -math.Pi*2)
	c.gc.Close()
}

// Draws an empty ellipse
// Fill the given ellipse with the fill color
// Stroke() each time to avoid connected ellipses
func (c *Canvas) Ellipse(cx, cy, rx, ry float64) {
	c.gc.ArcTo(cx, cy, rx, ry, 0, -math.Pi*2)
	c.gc.Close()
}

// Save the current canvas to a PNG file
func (c *Canvas) SaveToPNG(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, c.img)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Wrote %s OK.\n", filename)
}

// Return the width of the canvas
func (c *Canvas) Width() int {
	return c.width
}

// Return the height of the canvas
func (c *Canvas) Height() int {
	return c.height
}
