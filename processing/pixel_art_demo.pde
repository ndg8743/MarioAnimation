/*
================================================================================
PIXEL ART WITH PROCESSING
================================================================================

Processing is a flexible software sketchbook and language for learning how to 
code within the context of the visual arts. It's been used for everything from
teaching programming to creating interactive installations and generative art.

Website: https://processing.org/
Download: https://processing.org/download/

HOW PROCESSING WORKS:
--------------------
Processing uses two main functions:
  - setup()  : Called once at the start
  - draw()   : Called repeatedly (60 fps by default)

This is very similar to game engines like Pygame, Ebitengine, etc.

TO RUN THIS:
-----------
1. Download Processing from https://processing.org/download/
2. Open Processing IDE
3. Copy this entire code
4. Click the Play button (or press Ctrl+R / Cmd+R)

CONTROLS:
---------
    1, 2, 3, 4  - Switch palettes
    SPACE       - Toggle animation
    W           - Toggle wave effect
    S           - Save screenshot
    R           - Toggle rotation

================================================================================
*/

// ============================================================================
// PIXEL ART DATA
// ============================================================================
// Same format as Python/Go versions - strings where each char maps to a color

String[] MARIO = {
  "0000022222000000",  // Row 0: Top of hat
  "0000222222222000",  // Row 1: Hat brim
  "0000333113100000",  // Row 2: Hair and face start
  "0003131113111000",  // Row 3: Face details
  "0003133111311100",  // Row 4: Face and ear
  "0003311113333000",  // Row 5: Face bottom
  "0000011111110000",  // Row 6: Neck area
  "0000332333300000",  // Row 7: Shirt top
  "0003332332333000",  // Row 8: Arms and body
  "0033332222333300",  // Row 9: Body middle
  "0011321221231100",  // Row 10: Belt area
  "0011122222211100",  // Row 11: Overalls top
  "0011222222221100",  // Row 12: Overalls bottom
  "0000222002220000",  // Row 13: Legs gap
  "0003330000333000",  // Row 14: Boots
  "0033330000333300"   // Row 15: Boot soles
};

// Walking frame (leg forward)
String[] MARIO_WALK = {
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
  "0000222222220000",  // Changed
  "0000022200000000",  // Changed
  "0000333300000000",  // Changed
  "0003333000000000"   // Changed
};

// ============================================================================
// COLOR PALETTES
// ============================================================================
// Processing uses color() function: color(r, g, b) or color(r, g, b, a)

// Palette arrays: index 0='0', 1='1', 2='2', 3='3'
color[] PALETTE_NES = {
  color(148, 148, 255),  // 0: Sky blue (background)
  color(240, 172, 63),   // 1: Skin/peach
  color(185, 39, 22),    // 2: Red
  color(115, 103, 2)     // 3: Brown
};

color[] PALETTE_GAMEBOY = {
  color(155, 188, 15),   // 0: Lightest green
  color(139, 172, 15),   // 1: Light green
  color(48, 98, 48),     // 2: Dark green
  color(15, 56, 15)      // 3: Darkest green
};

color[] PALETTE_GRAYSCALE = {
  color(255, 255, 255),  // 0: White
  color(170, 170, 170),  // 1: Light gray
  color(85, 85, 85),     // 2: Dark gray
  color(0, 0, 0)         // 3: Black
};

color[] PALETTE_SYNTHWAVE = {
  color(13, 2, 33),      // 0: Deep purple
  color(255, 113, 206),  // 1: Hot pink
  color(1, 205, 254),    // 2: Cyan
  color(185, 103, 255)   // 3: Purple
};

// All palettes in an array for easy switching
color[][] ALL_PALETTES = {
  PALETTE_NES,
  PALETTE_GAMEBOY,
  PALETTE_GRAYSCALE,
  PALETTE_SYNTHWAVE
};

String[] PALETTE_NAMES = {
  "NES Classic",
  "Game Boy",
  "Grayscale", 
  "Synthwave"
};

// ============================================================================
// GLOBAL STATE VARIABLES
// ============================================================================

int currentPalette = 0;      // Which palette is active
int pixelSize = 16;          // Size of each "pixel" square
boolean animating = false;   // Animation on/off
boolean waveEffect = false;  // Wave distortion on/off
boolean rotating = false;    // Rotation on/off

int currentFrame = 0;        // Current animation frame
int frameTimer = 0;          // Timer for animation
int animationSpeed = 8;      // Frames between animation updates

float waveTime = 0;          // Time accumulator for wave effect
float rotationAngle = 0;     // Current rotation angle

// ============================================================================
// SETUP - Called once at start
// ============================================================================

void setup() {
  // Create window: 400x400 pixels
  size(400, 400);
  
  // Set frame rate (default is 60)
  frameRate(60);
  
  // Use noSmooth() for crisp pixel edges
  // This prevents anti-aliasing which would blur our pixels
  noSmooth();
  
  // No outline on rectangles by default
  noStroke();
  
  println("=== PIXEL ART DEMO ===");
  println("Controls:");
  println("  1-4   : Switch palettes");
  println("  SPACE : Toggle animation");
  println("  W     : Toggle wave effect");
  println("  R     : Toggle rotation");
  println("  S     : Save screenshot");
}

// ============================================================================
// DRAW - Called every frame (60 times per second)
// ============================================================================

void draw() {
  // Get current palette
  color[] palette = ALL_PALETTES[currentPalette];
  
  // Fill background with palette's background color
  background(palette[0]);
  
  // Update animation
  if (animating) {
    frameTimer++;
    if (frameTimer >= animationSpeed) {
      frameTimer = 0;
      currentFrame = (currentFrame + 1) % 2;  // Toggle between 0 and 1
    }
  }
  
  // Update wave time
  if (waveEffect) {
    waveTime += 0.05;
  }
  
  // Update rotation
  if (rotating) {
    rotationAngle += 0.02;
  }
  
  // Choose which sprite data to use
  String[] spriteData = (animating && currentFrame == 1) ? MARIO_WALK : MARIO;
  
  // Calculate center position
  int spriteWidth = spriteData[0].length() * pixelSize;
  int spriteHeight = spriteData.length * pixelSize;
  float centerX = width / 2;
  float centerY = height / 2;
  
  // Draw the sprite
  // We use pushMatrix/popMatrix to isolate transformations
  pushMatrix();
  
  // Move to center of screen
  translate(centerX, centerY);
  
  // Apply rotation if enabled
  if (rotating) {
    rotate(rotationAngle);
  }
  
  // Draw each pixel
  for (int y = 0; y < spriteData.length; y++) {
    String row = spriteData[y];
    
    for (int x = 0; x < row.length(); x++) {
      // Get character at this position
      char c = row.charAt(x);
      
      // Convert char to palette index (0-3)
      int colorIndex = c - '0';  // '0'->0, '1'->1, etc.
      
      // Get color from palette
      color pixelColor = palette[colorIndex];
      
      // Calculate pixel position (centered on origin)
      float px = (x - spriteData[0].length() / 2.0) * pixelSize;
      float py = (y - spriteData.length / 2.0) * pixelSize;
      
      // Apply wave effect if enabled
      if (waveEffect) {
        // Sine wave offset based on x position and time
        float waveOffset = sin(px * 0.05 + waveTime * 4) * 8;
        py += waveOffset;
      }
      
      // Set fill color and draw rectangle
      fill(pixelColor);
      rect(px, py, pixelSize, pixelSize);
    }
  }
  
  popMatrix();
  
  // Draw UI text
  drawUI();
}

// ============================================================================
// DRAW UI - Shows info and controls
// ============================================================================

void drawUI() {
  // Semi-transparent black bar at top
  fill(0, 0, 0, 150);
  noStroke();
  rect(0, 0, width, 30);
  
  // White text
  fill(255);
  textSize(12);
  textAlign(LEFT, TOP);
  
  // Palette name
  text("Palette: " + PALETTE_NAMES[currentPalette] + " (1-4)", 10, 8);
  
  // Status indicators
  String status = "";
  if (animating) status += "[ANIM] ";
  if (waveEffect) status += "[WAVE] ";
  if (rotating) status += "[ROTATE] ";
  
  textAlign(RIGHT, TOP);
  text(status, width - 10, 8);
  
  // Bottom bar with controls
  fill(0, 0, 0, 150);
  rect(0, height - 25, width, 25);
  
  fill(255);
  textAlign(CENTER, BOTTOM);
  text("SPACE=anim  W=wave  R=rotate  S=save", width/2, height - 5);
}

// ============================================================================
// KEY PRESSED - Handle keyboard input
// ============================================================================

void keyPressed() {
  // Number keys 1-4: Switch palettes
  if (key == '1') currentPalette = 0;
  if (key == '2') currentPalette = 1;
  if (key == '3') currentPalette = 2;
  if (key == '4') currentPalette = 3;
  
  // Space: Toggle animation
  if (key == ' ') {
    animating = !animating;
    println("Animation: " + (animating ? "ON" : "OFF"));
  }
  
  // W: Toggle wave effect
  if (key == 'w' || key == 'W') {
    waveEffect = !waveEffect;
    waveTime = 0;
    println("Wave effect: " + (waveEffect ? "ON" : "OFF"));
  }
  
  // R: Toggle rotation
  if (key == 'r' || key == 'R') {
    rotating = !rotating;
    if (!rotating) rotationAngle = 0;  // Reset angle when stopping
    println("Rotation: " + (rotating ? "ON" : "OFF"));
  }
  
  // S: Save screenshot
  if (key == 's' || key == 'S') {
    String filename = "mario_" + PALETTE_NAMES[currentPalette].toLowerCase().replace(" ", "_") + "_" + millis() + ".png";
    saveFrame(filename);
    println("Saved: " + filename);
  }
  
  // +/-: Adjust pixel size
  if (key == '+' || key == '=') {
    pixelSize = min(pixelSize + 2, 32);
    println("Pixel size: " + pixelSize);
  }
  if (key == '-' || key == '_') {
    pixelSize = max(pixelSize - 2, 4);
    println("Pixel size: " + pixelSize);
  }
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// Convert pixel art to PImage for more advanced manipulation
PImage spriteToImage(String[] data, color[] palette) {
  int w = data[0].length();
  int h = data.length;
  
  PImage img = createImage(w, h, RGB);
  img.loadPixels();
  
  for (int y = 0; y < h; y++) {
    for (int x = 0; x < w; x++) {
      char c = data[y].charAt(x);
      int colorIndex = c - '0';
      img.pixels[y * w + x] = palette[colorIndex];
    }
  }
  
  img.updatePixels();
  return img;
}

// Draw sprite from PImage with scaling
void drawSprite(PImage img, float x, float y, float scale) {
  // Disable smoothing for crisp pixel scaling
  noSmooth();
  image(img, x, y, img.width * scale, img.height * scale);
}
