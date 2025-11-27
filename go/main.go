/*
================================================================================
PIXEL ART WITH EBITENGINE (GO)
================================================================================

Ebitengine is a dead-simple 2D game engine for Go. It's been used in commercial
games like "Bear's Restaurant" and "Fishing Paradiso" (2M+ downloads).

This demo shows how to:
  - Define pixel art as string arrays (like your Python version)
  - Render pixels as colored rectangles
  - Animate sprites with palette swapping
  - Handle keyboard input
  - Export frames as PNG/GIF

HOW EBITENGINE WORKS:
--------------------
Ebitengine runs a game loop at 60 TPS (ticks per second):

    for {
        game.Update()   // Update game state (logic)
        game.Draw()     // Render to screen
        // Wait for next frame...
    }

You implement the ebiten.Game interface with 3 methods:
  - Update() error              - Game logic (called 60x/sec)
  - Draw(screen *ebiten.Image)  - Rendering (called every frame)
  - Layout(w, h int) (int, int) - Return logical screen size

RUN THIS:
---------
    go mod init pixelart
    go get github.com/hajimehoshi/ebiten/v2
    go run main.go

CONTROLS:
---------
    1, 2, 3, 4  - Switch palettes
    SPACE       - Toggle animation
    S           - Save screenshot
    ESC         - Quit

================================================================================
*/

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ============================================================================
// PIXEL ART DATA
// ============================================================================
// Same format as your Python version - strings where each char is a palette index

var MARIO = []string{
	"0000022222000000", // Row 0: Top of hat
	"0000222222222000", // Row 1: Hat brim
	"0000333113100000", // Row 2: Hair and face start
	"0003131113111000", // Row 3: Face details
	"0003133111311100", // Row 4: Face and ear
	"0003311113333000", // Row 5: Face bottom
	"0000011111110000", // Row 6: Neck area
	"0000332333300000", // Row 7: Shirt top
	"0003332332333000", // Row 8: Arms and body
	"0033332222333300", // Row 9: Body middle
	"0011321221231100", // Row 10: Belt area
	"0011122222211100", // Row 11: Overalls top
	"0011222222221100", // Row 12: Overalls bottom
	"0000222002220000", // Row 13: Legs gap
	"0003330000333000", // Row 14: Boots
	"0033330000333300", // Row 15: Boot soles
}

// Additional sprite for animation demo
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
	"0000222222220000", // Changed: leg forward
	"0000022200000000", // Changed
	"0000333300000000", // Changed
	"0003333000000000", // Changed
}

// ============================================================================
// COLOR PALETTES
// ============================================================================
// Each palette maps character -> RGBA color

type Palette map[byte]color.RGBA

var PALETTE_NES = Palette{
	'0': {148, 148, 255, 255}, // Sky blue (background)
	'1': {240, 172, 63, 255},  // Skin/peach
	'2': {185, 39, 22, 255},   // Red
	'3': {115, 103, 2, 255},   // Brown
}

var PALETTE_GAMEBOY = Palette{
	'0': {155, 188, 15, 255},  // Lightest green
	'1': {139, 172, 15, 255},  // Light green
	'2': {48, 98, 48, 255},    // Dark green
	'3': {15, 56, 15, 255},    // Darkest green
}

var PALETTE_GRAYSCALE = Palette{
	'0': {255, 255, 255, 255}, // White
	'1': {170, 170, 170, 255}, // Light gray
	'2': {85, 85, 85, 255},    // Dark gray
	'3': {0, 0, 0, 255},       // Black
}

var PALETTE_SYNTHWAVE = Palette{
	'0': {13, 2, 33, 255},     // Deep purple background
	'1': {255, 113, 206, 255}, // Hot pink
	'2': {1, 205, 254, 255},   // Cyan
	'3': {185, 103, 255, 255}, // Purple
}

// All palettes for cycling
var ALL_PALETTES = []Palette{
	PALETTE_NES,
	PALETTE_GAMEBOY,
	PALETTE_GRAYSCALE,
	PALETTE_SYNTHWAVE,
}

var PALETTE_NAMES = []string{
	"NES Classic",
	"Game Boy",
	"Grayscale",
	"Synthwave",
}

// ============================================================================
// PIXEL SPRITE
// ============================================================================

// PixelSprite holds pixel art data and can render to an ebiten.Image
type PixelSprite struct {
	Data      []string       // The string array defining pixels
	Palette   Palette        // Color mapping
	PixelSize int            // Size of each "pixel" in actual pixels
	Image     *ebiten.Image  // Cached rendered image
}

// NewPixelSprite creates a sprite from string data
func NewPixelSprite(data []string, palette Palette, pixelSize int) *PixelSprite {
	s := &PixelSprite{
		Data:      data,
		Palette:   palette,
		PixelSize: pixelSize,
	}
	s.Render()
	return s
}

// Render converts the string data to an ebiten.Image
// This is called once when created, and again when palette changes
func (s *PixelSprite) Render() {
	if len(s.Data) == 0 {
		return
	}

	height := len(s.Data)
	width := len(s.Data[0])

	// Create image with dimensions: (width * pixelSize) x (height * pixelSize)
	s.Image = ebiten.NewImage(width*s.PixelSize, height*s.PixelSize)

	// Draw each pixel as a filled rectangle
	for y, row := range s.Data {
		for x, char := range row {
			// Look up color in palette
			col, exists := s.Palette[byte(char)]
			if !exists {
				continue // Skip unknown characters
			}

			// Calculate rectangle position
			rectX := x * s.PixelSize
			rectY := y * s.PixelSize

			// Draw filled rectangle for this pixel
			// We create a tiny 1x1 image and scale it
			pixel := ebiten.NewImage(s.PixelSize, s.PixelSize)
			pixel.Fill(col)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(rectX), float64(rectY))
			s.Image.DrawImage(pixel, op)
		}
	}
}

// SetPalette changes the color palette and re-renders
func (s *PixelSprite) SetPalette(p Palette) {
	s.Palette = p
	s.Render()
}

// Width returns the rendered width in pixels
func (s *PixelSprite) Width() int {
	if s.Image == nil {
		return 0
	}
	return s.Image.Bounds().Dx()
}

// Height returns the rendered height in pixels
func (s *PixelSprite) Height() int {
	if s.Image == nil {
		return 0
	}
	return s.Image.Bounds().Dy()
}

// Draw renders the sprite to a target image at position (x, y)
func (s *PixelSprite) Draw(target *ebiten.Image, x, y float64) {
	if s.Image == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	target.DrawImage(s.Image, op)
}

// DrawWithOptions renders with custom options (rotation, scale, etc.)
func (s *PixelSprite) DrawWithOptions(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	if s.Image == nil {
		return
	}
	target.DrawImage(s.Image, op)
}

// ============================================================================
// ANIMATED SPRITE
// ============================================================================

// AnimatedSprite handles frame-by-frame animation
type AnimatedSprite struct {
	Frames       []*PixelSprite
	CurrentFrame int
	FrameTime    float64 // Seconds per frame
	Elapsed      float64 // Time accumulator
	Playing      bool
}

// NewAnimatedSprite creates an animation from multiple sprites
func NewAnimatedSprite(frames []*PixelSprite, fps float64) *AnimatedSprite {
	return &AnimatedSprite{
		Frames:    frames,
		FrameTime: 1.0 / fps,
		Playing:   true,
	}
}

// Update advances the animation (call every tick)
func (a *AnimatedSprite) Update(dt float64) {
	if !a.Playing || len(a.Frames) == 0 {
		return
	}

	a.Elapsed += dt
	if a.Elapsed >= a.FrameTime {
		a.Elapsed -= a.FrameTime
		a.CurrentFrame = (a.CurrentFrame + 1) % len(a.Frames)
	}
}

// Current returns the current frame's sprite
func (a *AnimatedSprite) Current() *PixelSprite {
	if len(a.Frames) == 0 {
		return nil
	}
	return a.Frames[a.CurrentFrame]
}

// Draw renders the current frame
func (a *AnimatedSprite) Draw(target *ebiten.Image, x, y float64) {
	if frame := a.Current(); frame != nil {
		frame.Draw(target, x, y)
	}
}

// SetPalette changes palette for all frames
func (a *AnimatedSprite) SetPalette(p Palette) {
	for _, frame := range a.Frames {
		frame.SetPalette(p)
	}
}

// ============================================================================
// GAME STATE
// ============================================================================

// Game implements ebiten.Game interface
type Game struct {
	// Sprites
	mario     *PixelSprite
	marioAnim *AnimatedSprite

	// State
	currentPalette int
	animating      bool
	waveEffect     bool
	waveTime       float64

	// Position for wave effect
	spriteX float64
	spriteY float64

	// Screen dimensions
	screenWidth  int
	screenHeight int

	// For delta time calculation
	lastUpdate time.Time
}

// NewGame initializes the game
func NewGame() *Game {
	g := &Game{
		currentPalette: 0,
		animating:      false,
		screenWidth:    320,
		screenHeight:   240,
		lastUpdate:     time.Now(),
	}

	// Create main sprite
	g.mario = NewPixelSprite(MARIO, PALETTE_NES, 8)

	// Create animation frames
	frame1 := NewPixelSprite(MARIO, PALETTE_NES, 8)
	frame2 := NewPixelSprite(MARIO_WALK, PALETTE_NES, 8)
	g.marioAnim = NewAnimatedSprite([]*PixelSprite{frame1, frame2}, 4)

	// Center sprite
	g.spriteX = float64(g.screenWidth-g.mario.Width()) / 2
	g.spriteY = float64(g.screenHeight-g.mario.Height()) / 2

	return g
}

// Update handles game logic (called 60 times per second)
func (g *Game) Update() error {
	// Calculate delta time
	now := time.Now()
	dt := now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	// ========================================
	// INPUT HANDLING
	// ========================================

	// Number keys 1-4: Switch palettes
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.switchPalette(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.switchPalette(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.switchPalette(2)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.switchPalette(3)
	}

	// Space: Toggle animation
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.animating = !g.animating
		g.marioAnim.Playing = g.animating
	}

	// W: Toggle wave effect
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.waveEffect = !g.waveEffect
		g.waveTime = 0
	}

	// S: Save screenshot
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.saveScreenshot()
	}

	// G: Save as GIF (animation)
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		g.saveGIF()
	}

	// ESC: Quit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// ========================================
	// UPDATE ANIMATIONS
	// ========================================

	if g.animating {
		g.marioAnim.Update(dt)
	}

	if g.waveEffect {
		g.waveTime += dt
	}

	return nil
}

// switchPalette changes the color palette
func (g *Game) switchPalette(index int) {
	if index >= 0 && index < len(ALL_PALETTES) {
		g.currentPalette = index
		g.mario.SetPalette(ALL_PALETTES[index])
		g.marioAnim.SetPalette(ALL_PALETTES[index])
	}
}

// Draw renders the game (called every frame)
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with background color from current palette
	bgColor := ALL_PALETTES[g.currentPalette]['0']
	screen.Fill(bgColor)

	// ========================================
	// DRAW SPRITE
	// ========================================

	var sprite *PixelSprite
	if g.animating {
		sprite = g.marioAnim.Current()
	} else {
		sprite = g.mario
	}

	if sprite != nil {
		op := &ebiten.DrawImageOptions{}

		if g.waveEffect {
			// Apply wave distortion effect
			// Move sprite up/down based on sine wave
			waveOffset := math.Sin(g.waveTime*4) * 10
			op.GeoM.Translate(g.spriteX, g.spriteY+waveOffset)

			// Also rotate slightly
			// First translate to center, rotate, translate back
			cx := float64(sprite.Width()) / 2
			cy := float64(sprite.Height()) / 2
			op.GeoM.Reset()
			op.GeoM.Translate(-cx, -cy)
			op.GeoM.Rotate(math.Sin(g.waveTime*2) * 0.1) // Small rotation
			op.GeoM.Translate(cx, cy)
			op.GeoM.Translate(g.spriteX, g.spriteY+waveOffset)
		} else {
			op.GeoM.Translate(g.spriteX, g.spriteY)
		}

		sprite.DrawWithOptions(screen, op)
	}

	// ========================================
	// DRAW UI TEXT
	// ========================================

	// Debug info at top
	info := fmt.Sprintf("Palette: %s (1-4 to switch)", PALETTE_NAMES[g.currentPalette])
	ebitenutil.DebugPrint(screen, info)

	// Controls at bottom
	controls := "SPACE=animate W=wave S=screenshot G=gif ESC=quit"
	ebitenutil.DebugPrintAt(screen, controls, 10, g.screenHeight-20)

	// Status
	status := ""
	if g.animating {
		status += "[ANIMATING] "
	}
	if g.waveEffect {
		status += "[WAVE] "
	}
	if status != "" {
		ebitenutil.DebugPrintAt(screen, status, 10, 20)
	}
}

// Layout returns the logical screen size
// Ebitengine will scale this to fit the window
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

// ============================================================================
// EXPORT FUNCTIONS
// ============================================================================

// saveScreenshot saves the current sprite as PNG
func (g *Game) saveScreenshot() {
	filename := fmt.Sprintf("mario_%s.png", time.Now().Format("20060102_150405"))

	// Get current sprite image
	var img *ebiten.Image
	if g.animating && g.marioAnim.Current() != nil {
		img = g.marioAnim.Current().Image
	} else {
		img = g.mario.Image
	}

	if img == nil {
		log.Println("No image to save")
		return
	}

	// Convert to standard image
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	// Save to file
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return
	}
	defer f.Close()

	if err := png.Encode(f, rgba); err != nil {
		log.Printf("Failed to encode PNG: %v", err)
		return
	}

	log.Printf("Saved: %s", filename)
}

// saveGIF saves the animation as a GIF file
func (g *Game) saveGIF() {
	filename := fmt.Sprintf("mario_anim_%s.gif", time.Now().Format("20060102_150405"))

	// Collect frames
	var images []*image.Paletted
	var delays []int

	// Get palette colors for GIF
	gifPalette := []color.Color{
		ALL_PALETTES[g.currentPalette]['0'],
		ALL_PALETTES[g.currentPalette]['1'],
		ALL_PALETTES[g.currentPalette]['2'],
		ALL_PALETTES[g.currentPalette]['3'],
	}

	for _, frame := range g.marioAnim.Frames {
		if frame.Image == nil {
			continue
		}

		bounds := frame.Image.Bounds()
		paletted := image.NewPaletted(bounds, gifPalette)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				paletted.Set(x, y, frame.Image.At(x, y))
			}
		}

		images = append(images, paletted)
		delays = append(delays, 25) // 250ms per frame (in centiseconds)
	}

	if len(images) == 0 {
		log.Println("No frames to save")
		return
	}

	f, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return
	}
	defer f.Close()

	err = gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
	if err != nil {
		log.Printf("Failed to encode GIF: %v", err)
		return
	}

	log.Printf("Saved: %s", filename)
}

// ============================================================================
// MAIN
// ============================================================================

func main() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════════╗
║              PIXEL ART WITH EBITENGINE (GO)                   ║
╠═══════════════════════════════════════════════════════════════╣
║  Controls:                                                    ║
║    1, 2, 3, 4  - Switch color palettes                        ║
║    SPACE       - Toggle walking animation                     ║
║    W           - Toggle wave distortion effect                ║
║    S           - Save screenshot as PNG                       ║
║    G           - Save animation as GIF                        ║
║    ESC         - Quit                                         ║
╚═══════════════════════════════════════════════════════════════╝
`)

	game := NewGame()

	// Configure window
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Pixel Art - Ebitengine")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Run the game!
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
