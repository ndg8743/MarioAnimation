"""
================================================================================
PIXEL ART ANIMATION WITH MANIM
================================================================================

This script creates a beautiful animated visualization of pixel art using Manim,
the mathematical animation engine (made famous by 3Blue1Brown).

HOW IT WORKS:
-------------
1. We define pixel art as simple string arrays (like your original Mario)
2. Each character maps to a color in a palette
3. Manim converts each "pixel" into a Square object
4. We animate the squares appearing, transforming, and moving

MANIM BASICS:
-------------
- Scene: A container for all animations (like a movie scene)
- Mobject: Any "Mathematical Object" that can be displayed/animated
- VGroup: A group of Mobjects that move together
- Animation: Things like FadeIn, Transform, Create, etc.

RUN THIS WITH:
    manim -pql pixel_art_manim.py PixelArtScene
    
FLAGS:
    -p  = preview (auto-open when done)
    -ql = quality low (faster render, 480p)
    -qm = quality medium (720p)
    -qh = quality high (1080p)
    -q4k = 4K quality

================================================================================
"""

from manim import *
import numpy as np


# ============================================================================
# PIXEL ART DATA
# ============================================================================
# Each string represents a row of pixels
# Each character maps to a color in the palette

MARIO = [
    "0000022222000000",  # Row 0: Top of hat
    "0000222222222000",  # Row 1: Hat brim
    "0000333113100000",  # Row 2: Hair and face
    "0003131113111000",  # Row 3: Face details
    "0003133111311100",  # Row 4: Face and ear
    "0003311113333000",  # Row 5: Face bottom
    "0000011111110000",  # Row 6: Neck area
    "0000332333300000",  # Row 7: Shirt top
    "0003332332333000",  # Row 8: Arms and body
    "0033332222333300",  # Row 9: Body middle
    "0011321221231100",  # Row 10: Belt area
    "0011122222211100",  # Row 11: Overalls
    "0011222222221100",  # Row 12: Overalls
    "0000222002220000",  # Row 13: Legs gap
    "0003330000333000",  # Row 14: Boots
    "0033330000333300",  # Row 15: Boot soles
]

# Color palettes - RGB tuples converted to Manim hex colors
PALETTE_MARIO = {
    "0": "#9494FF",  # Sky blue (background/transparent)
    "1": "#F0AC3F",  # Skin tone (peach/yellow)
    "2": "#B92716",  # Red (hat, shirt)
    "3": "#736702",  # Brown (hair, shoes)
}

PALETTE_GRAYSCALE = {
    "0": "#FFFFFF",
    "1": "#AAAAAA", 
    "2": "#555555",
    "3": "#000000",
}

PALETTE_GAMEBOY = {
    "0": "#9BBC0F",  # Lightest green
    "1": "#8BAC0F",  # Light green
    "2": "#306230",  # Dark green
    "3": "#0F380F",  # Darkest green
}


# ============================================================================
# HELPER FUNCTIONS
# ============================================================================

def create_pixel_grid(
    pixel_data: list[str],
    palette: dict[str, str],
    pixel_size: float = 0.3,
    gap: float = 0.02
) -> VGroup:
    """
    Convert string-based pixel art into a Manim VGroup of squares.
    
    Args:
        pixel_data: List of strings where each char is a palette key
        palette: Dict mapping characters to hex color strings
        pixel_size: Size of each pixel square
        gap: Space between pixels (0 for no gap)
    
    Returns:
        VGroup containing all pixel squares, centered at origin
    
    How it works:
        1. Iterate through each row (y) and column (x)
        2. Create a Square for each pixel
        3. Color it based on the palette
        4. Position it in a grid layout
        5. Group everything together
    """
    pixels = VGroup()
    
    height = len(pixel_data)
    width = len(pixel_data[0]) if pixel_data else 0
    
    # Calculate total dimensions for centering
    total_width = width * (pixel_size + gap)
    total_height = height * (pixel_size + gap)
    
    for y, row in enumerate(pixel_data):
        for x, char in enumerate(row):
            # Skip if character not in palette (could be used for transparency)
            if char not in palette:
                continue
            
            # Create the pixel as a square
            pixel = Square(
                side_length=pixel_size,
                fill_color=palette[char],
                fill_opacity=1.0,
                stroke_width=0.5,
                stroke_color=BLACK,
                stroke_opacity=0.3
            )
            
            # Position the pixel
            # Note: Manim's y-axis goes UP, but our data goes DOWN
            # So we flip the y coordinate
            pixel.move_to([
                x * (pixel_size + gap) - total_width / 2 + pixel_size / 2,
                -y * (pixel_size + gap) + total_height / 2 - pixel_size / 2,
                0
            ])
            
            # Store metadata for later use in animations
            pixel.pixel_x = x
            pixel.pixel_y = y
            pixel.pixel_char = char
            
            pixels.add(pixel)
    
    return pixels


def create_pixel_grid_animated(
    pixel_data: list[str],
    palette: dict[str, str],
    pixel_size: float = 0.3
) -> list[list[Square]]:
    """
    Create pixel grid as 2D array for more control over animations.
    Returns grid[y][x] for easy row/column access.
    """
    height = len(pixel_data)
    width = len(pixel_data[0]) if pixel_data else 0
    
    total_width = width * pixel_size
    total_height = height * pixel_size
    
    grid = []
    
    for y, row in enumerate(pixel_data):
        grid_row = []
        for x, char in enumerate(row):
            pixel = Square(
                side_length=pixel_size,
                fill_color=palette.get(char, "#000000"),
                fill_opacity=1.0,
                stroke_width=0,
            )
            pixel.move_to([
                x * pixel_size - total_width / 2 + pixel_size / 2,
                -y * pixel_size + total_height / 2 - pixel_size / 2,
                0
            ])
            grid_row.append(pixel)
        grid.append(grid_row)
    
    return grid


# ============================================================================
# MAIN ANIMATION SCENE
# ============================================================================

class PixelArtScene(Scene):
    """
    Main scene demonstrating various pixel art animations.
    
    Scene structure:
        1. Title card
        2. Pixel-by-pixel reveal animation
        3. Palette swap effect
        4. Zoom and showcase
        5. Wave distortion effect
        6. Finale
    """
    
    def construct(self):
        """
        The construct() method is where all animation happens.
        This is automatically called by Manim when rendering.
        """
        
        # ================================================
        # PART 1: TITLE CARD
        # ================================================
        
        title = Text("Pixel Art with Manim", font_size=48)
        subtitle = Text("Building sprites one square at a time", font_size=24)
        subtitle.next_to(title, DOWN)
        
        # Animate title appearing
        self.play(Write(title), run_time=1.5)
        self.play(FadeIn(subtitle, shift=UP * 0.3), run_time=0.8)
        self.wait(1)
        
        # Fade out title
        self.play(FadeOut(title), FadeOut(subtitle))
        
        # ================================================
        # PART 2: PIXEL-BY-PIXEL REVEAL
        # ================================================
        
        # Create the pixel grid
        mario = create_pixel_grid(MARIO, PALETTE_MARIO, pixel_size=0.35, gap=0.01)
        
        # Add a label
        label = Text("16×16 Pixel Art", font_size=28)
        label.to_edge(UP)
        self.play(Write(label))
        
        # Animate pixels appearing one by one (spiral pattern)
        # First, sort pixels by distance from center for cool effect
        center = mario.get_center()
        sorted_pixels = sorted(
            mario,
            key=lambda p: np.linalg.norm(p.get_center() - center)
        )
        
        # Create animations for each pixel
        # Using LaggedStart for a wave-like appearance
        self.play(
            LaggedStart(
                *[GrowFromCenter(p) for p in sorted_pixels],
                lag_ratio=0.02,  # Delay between each pixel
                run_time=3
            )
        )
        
        self.wait(1)
        
        # ================================================
        # PART 3: PALETTE SWAP ANIMATION
        # ================================================
        
        # Create same sprite with different palette
        mario_gameboy = create_pixel_grid(MARIO, PALETTE_GAMEBOY, pixel_size=0.35, gap=0.01)
        mario_grayscale = create_pixel_grid(MARIO, PALETTE_GRAYSCALE, pixel_size=0.35, gap=0.01)
        
        # Update label
        new_label = Text("Palette Swap: Game Boy Style", font_size=28)
        new_label.to_edge(UP)
        
        # Transform to Game Boy palette
        self.play(
            Transform(mario, mario_gameboy),
            Transform(label, new_label),
            run_time=1.5
        )
        self.wait(0.8)
        
        # Transform to grayscale
        new_label2 = Text("Palette Swap: Grayscale", font_size=28)
        new_label2.to_edge(UP)
        
        self.play(
            Transform(mario, mario_grayscale),
            Transform(label, new_label2),
            run_time=1.5
        )
        self.wait(0.8)
        
        # Back to original colors
        mario_original = create_pixel_grid(MARIO, PALETTE_MARIO, pixel_size=0.35, gap=0.01)
        new_label3 = Text("Back to Original NES Colors", font_size=28)
        new_label3.to_edge(UP)
        
        self.play(
            Transform(mario, mario_original),
            Transform(label, new_label3),
            run_time=1.5
        )
        self.wait(1)
        
        # ================================================
        # PART 4: SCALE AND ROTATE
        # ================================================
        
        new_label4 = Text("Pixel Art Scales Perfectly!", font_size=28)
        new_label4.to_edge(UP)
        
        self.play(
            mario.animate.scale(1.8),
            Transform(label, new_label4),
            run_time=1.5
        )
        self.wait(0.5)
        
        # Rotate
        self.play(
            Rotate(mario, angle=TAU, run_time=2)  # TAU = 2π = full rotation
        )
        self.wait(0.5)
        
        # Scale back
        self.play(mario.animate.scale(1/1.8), run_time=0.8)
        
        # ================================================
        # PART 5: WAVE DISTORTION EFFECT
        # ================================================
        
        new_label5 = Text("Wave Animation Effect", font_size=28)
        new_label5.to_edge(UP)
        self.play(Transform(label, new_label5))
        
        # Store original positions
        original_positions = [p.get_center().copy() for p in mario]
        
        # Animate a wave passing through
        def wave_update(mob, dt, time_tracker):
            """Apply sine wave distortion to pixels"""
            t = time_tracker[0]
            time_tracker[0] += dt
            
            for i, pixel in enumerate(mob):
                orig_pos = original_positions[i]
                # Sine wave based on x position and time
                wave_offset = 0.15 * np.sin(orig_pos[0] * 3 + t * 4)
                pixel.move_to([orig_pos[0], orig_pos[1] + wave_offset, 0])
        
        # Run wave animation for a few seconds
        time_tracker = [0]
        mario.add_updater(lambda m, dt: wave_update(m, dt, time_tracker))
        self.wait(3)
        mario.remove_updater(lambda m, dt: wave_update(m, dt, time_tracker))
        
        # Reset positions
        for i, pixel in enumerate(mario):
            pixel.move_to(original_positions[i])
        
        # ================================================
        # PART 6: FINALE - EXPLODE AND REFORM
        # ================================================
        
        new_label6 = Text("Thanks for Watching!", font_size=36)
        new_label6.to_edge(UP)
        self.play(Transform(label, new_label6))
        
        # Explode pixels outward
        self.play(
            *[
                p.animate.shift(
                    (p.get_center() - mario.get_center()) * 3  # Move away from center
                ).set_opacity(0)
                for p in mario
            ],
            run_time=1.5
        )
        
        # Recreate and implode back
        mario_final = create_pixel_grid(MARIO, PALETTE_MARIO, pixel_size=0.35, gap=0.01)
        
        # Start pixels scattered
        for p in mario_final:
            original = p.get_center().copy()
            p.shift((original - mario_final.get_center()) * 3)
            p.set_opacity(0)
            p.original_pos = original
        
        # Animate back together
        self.play(
            *[
                p.animate.move_to(p.original_pos).set_opacity(1)
                for p in mario_final
            ],
            run_time=1.5
        )
        
        self.wait(2)
        
        # Final fade out
        self.play(
            FadeOut(mario_final),
            FadeOut(label),
            run_time=1
        )


# ============================================================================
# BONUS: ROW-BY-ROW REVEAL SCENE
# ============================================================================

class RowByRowScene(Scene):
    """
    Alternative scene showing row-by-row scanline reveal effect.
    Like an old CRT TV drawing the image.
    """
    
    def construct(self):
        title = Text("Scanline Reveal Effect", font_size=36)
        title.to_edge(UP)
        self.add(title)
        
        # Create grid as 2D array for row access
        grid = create_pixel_grid_animated(MARIO, PALETTE_MARIO, pixel_size=0.4)
        
        # Flatten for adding to scene, but animate row by row
        for row in grid:
            # Animate each row appearing left to right
            self.play(
                LaggedStart(
                    *[FadeIn(pixel, scale=0.5) for pixel in row],
                    lag_ratio=0.05,
                    run_time=0.3
                )
            )
        
        self.wait(2)


# ============================================================================
# BONUS: TYPING EFFECT SCENE
# ============================================================================

class TypingEffectScene(Scene):
    """
    Shows the pixel art data being "typed" while the sprite builds.
    Educational scene showing the data structure.
    """
    
    def construct(self):
        # Code display on left
        code_title = Text("Pixel Data:", font_size=24)
        code_title.to_corner(UL)
        self.add(code_title)
        
        # Sprite display on right
        sprite_area = Rectangle(
            width=5, height=5,
            stroke_color=WHITE,
            stroke_width=2
        )
        sprite_area.to_edge(RIGHT)
        self.add(sprite_area)
        
        all_pixels = VGroup()
        
        # Type out each row and build sprite simultaneously
        for i, row_data in enumerate(MARIO[:8]):  # First 8 rows for time
            # Create code text
            code_line = Text(
                f'"{row_data}"',
                font_size=16,
                font="Monospace"
            )
            code_line.next_to(code_title, DOWN, buff=0.3 + i * 0.35)
            code_line.align_to(code_title, LEFT)
            
            # Create pixel row
            row_pixels = VGroup()
            for x, char in enumerate(row_data):
                pixel = Square(
                    side_length=0.25,
                    fill_color=PALETTE_MARIO.get(char, "#000000"),
                    fill_opacity=1,
                    stroke_width=0
                )
                pixel.move_to(sprite_area.get_center() + [
                    (x - 8) * 0.25,
                    (4 - i) * 0.25,
                    0
                ])
                row_pixels.add(pixel)
            
            # Animate both together
            self.play(
                Write(code_line, run_time=0.5),
                LaggedStart(
                    *[GrowFromCenter(p) for p in row_pixels],
                    lag_ratio=0.02,
                    run_time=0.5
                )
            )
            all_pixels.add(*row_pixels)
        
        self.wait(2)


# ============================================================================
# RUN INSTRUCTIONS
# ============================================================================

if __name__ == "__main__":
    print("""
    ╔═══════════════════════════════════════════════════════════════╗
    ║              PIXEL ART MANIM ANIMATION                        ║
    ╠═══════════════════════════════════════════════════════════════╣
    ║                                                               ║
    ║  To render the main animation:                                ║
    ║    manim -pql pixel_art_manim.py PixelArtScene                ║
    ║                                                               ║
    ║  To render the scanline effect:                               ║
    ║    manim -pql pixel_art_manim.py RowByRowScene                ║
    ║                                                               ║
    ║  To render the typing effect:                                 ║
    ║    manim -pql pixel_art_manim.py TypingEffectScene            ║
    ║                                                               ║
    ║  Quality options:                                             ║
    ║    -ql  = 480p (fast)                                         ║
    ║    -qm  = 720p (medium)                                       ║
    ║    -qh  = 1080p (slow)                                        ║
    ║    -qk  = 4K (very slow)                                      ║
    ║                                                               ║
    ║  Add -p to auto-preview when done                             ║
    ║  Add --format=gif for GIF output                              ║
    ║                                                               ║
    ╚═══════════════════════════════════════════════════════════════╝
    """)
