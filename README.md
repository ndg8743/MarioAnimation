# Mario Animation

A collection of pixel art animation projects featuring Mario sprites created with Python and Pygame.

![Mario Animation](pixel_art_animation.gif)

## Files

- **`mario_pygame.py`** - Simple Mario sprite display using Pygame
- **`pixel_art_engine.py`** - Full-featured pixel art engine with animated sprites (Mario walking animation and Goomba)
- **`pixel_art_manim.py`** - Manim-based animation (generates the GIF above)

## Requirements

```bash
pip install pygame
```

For the Manim animation:
```bash
pip install manim
```

## Usage

### Simple Mario Display
```bash
python mario_pygame.py
```

### Pixel Art Engine Demo
```bash
python pixel_art_engine.py
```

Controls:
- **ESC** - Quit
- **S** - Save screenshot

### Generate Animation GIF
```bash
manim -pql pixel_art_manim.py
```

## Features

- 16x16 pixel art sprites
- NES-style color palette
- Animated walking sprites
- Sprite sheet generation
- PNG export functionality

