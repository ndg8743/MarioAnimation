/*
================================================================================
PIXEL ART PNG GENERATOR (HEADLESS GO)
================================================================================

This version doesn't need Ebitengine - it uses pure Go standard library
to generate PNG images from pixel art data.

Perfect for:
  - Generating sprites without a display
  - Server-side sprite generation
  - CI/CD pipelines
  - Learning how PNG encoding works in Go

RUN:
    go run headless.go

OUTPUT:
    Creates mario_nes.png, mario_gameboy.png, etc.

================================================================================
*/

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
)

// ============================================================================
// PIXEL ART DATA
// ============================================================================

var MARIO = []string{
	"0000022222000000",
	"0000222222222000",
	"0000333113100000",
	"0003131113111000",
	"0003133111311100",
	"0003311113333000",
	"0000011111110000",
	"0000332333300000",
	"0003332332333000",
	"0033332222333300",
	"0011321221231100",
	"0011122222211100",
	"0011222222221100",
	"0000222002220000",
	"0003330000333000",
	"0033330000333300",
}

var MARIO_WALK = []string{
	"0000022222000000",
	"0000222222222000",
	"0000333113100000",
	"0003131113111000",
	"0003133111311100",
	"0003311113333000",
	"0000011111110000",
	"0000332333300000",
	"0003332332333000",
	"0033332222333300",
	"0011321221231100",
	"0011122222211100",
	"0000222222220000",
	"0000022200000000",
	"0000333300000000",
	"0003333000000000",
}

// ============================================================================
// PALETTES
// ============================================================================

type Palette map[byte]color.RGBA

var PALETTE_NES = Palette{
	'0': {148, 148, 255, 255},
	'1': {240, 172, 63, 255},
	'2': {185, 39, 22, 255},
	'3': {115, 103, 2, 255},
}

var PALETTE_GAMEBOY = Palette{
	'0': {155, 188, 15, 255},
	'1': {139, 172, 15, 255},
	'2': {48, 98, 48, 255},
	'3': {15, 56, 15, 255},
}

var PALETTE_GRAYSCALE = Palette{
	'0': {255, 255, 255, 255},
	'1': {170, 170, 170, 255},
	'2': {85, 85, 85, 255},
	'3': {0, 0, 0, 255},
}

var PALETTE_SYNTHWAVE = Palette{
	'0': {13, 2, 33, 255},
	'1': {255, 113, 206, 255},
	'2': {1, 205, 254, 255},
	'3': {185, 103, 255, 255},
}

// ============================================================================
// RENDERING FUNCTIONS
// ============================================================================

// RenderPixelArt converts string pixel data to an image.RGBA
func RenderPixelArt(data []string, palette Palette, scale int) *image.RGBA {
	if len(data) == 0 {
		return nil
	}

	height := len(data)
	width := len(data[0])

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, width*scale, height*scale))

	// Fill pixels
	for y, row := range data {
		for x, char := range row {
			col, exists := palette[byte(char)]
			if !exists {
				continue
			}

			// Fill the scaled pixel area
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.SetRGBA(x*scale+dx, y*scale+dy, col)
				}
			}
		}
	}

	return img
}

// SavePNG saves an image as PNG
func SavePNG(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// SaveGIF saves multiple frames as an animated GIF
func SaveGIF(frames []*image.RGBA, delays []int, filename string) error {
	if len(frames) == 0 {
		return fmt.Errorf("no frames provided")
	}

	// Build palette from first frame's colors
	palette := make([]color.Color, 0, 256)
	colorSet := make(map[color.RGBA]bool)

	for _, frame := range frames {
		bounds := frame.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := frame.RGBAAt(x, y)
				if !colorSet[c] {
					colorSet[c] = true
					palette = append(palette, c)
				}
			}
		}
	}

	// Convert frames to paletted images
	var images []*image.Paletted
	for _, frame := range frames {
		bounds := frame.Bounds()
		paletted := image.NewPaletted(bounds, palette)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				paletted.Set(x, y, frame.At(x, y))
			}
		}
		images = append(images, paletted)
	}

	// Write GIF
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}

// ============================================================================
// MAIN
// ============================================================================

func main() {
	fmt.Println("🎮 Pixel Art Generator (Go)")
	fmt.Println("===========================")

	scale := 8 // Each pixel becomes 8x8

	// Generate sprites with different palettes
	palettes := map[string]Palette{
		"nes":       PALETTE_NES,
		"gameboy":   PALETTE_GAMEBOY,
		"grayscale": PALETTE_GRAYSCALE,
		"synthwave": PALETTE_SYNTHWAVE,
	}

	for name, palette := range palettes {
		img := RenderPixelArt(MARIO, palette, scale)
		filename := fmt.Sprintf("mario_%s.png", name)

		if err := SavePNG(img, filename); err != nil {
			fmt.Printf("❌ Error saving %s: %v\n", filename, err)
		} else {
			fmt.Printf("✅ Saved: %s\n", filename)
		}
	}

	// Generate animated GIF
	fmt.Println("\n🎬 Generating animation...")

	frame1 := RenderPixelArt(MARIO, PALETTE_NES, scale)
	frame2 := RenderPixelArt(MARIO_WALK, PALETTE_NES, scale)

	frames := []*image.RGBA{frame1, frame2, frame1, frame2}
	delays := []int{25, 25, 25, 25} // 250ms each (in centiseconds)

	if err := SaveGIF(frames, delays, "mario_walk.gif"); err != nil {
		fmt.Printf("❌ Error saving GIF: %v\n", err)
	} else {
		fmt.Println("✅ Saved: mario_walk.gif")
	}

	fmt.Println("\n🎉 Done!")
}
